package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"github.com/suapapa/croquis-king/internal/lobby"
	"github.com/suapapa/croquis-king/internal/timer"
	"github.com/suapapa/croquis-king/internal/ws"
)

func TestFullSessionFlowHTTPAndWS(t *testing.T) {
	const (
		drawDuration = 150 * time.Millisecond
		tickInterval = 50 * time.Millisecond
	)

	store := lobby.NewMemoryStore()
	lobbySync := ws.NewSnapshotSync(store)
	wsHandler := ws.NewHandler(lobbySync, nil)
	scheduler := timer.NewScheduler(store, lobbySync, tickInterval)
	router := newTestRouter(store, drawDuration, nil, wsHandler, lobbySync)

	server := httptest.NewServer(router)
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go scheduler.Run(ctx)

	httpClient := &http.Client{Timeout: 5 * time.Second}

	created := createLobbyViaHTTP(t, httpClient, server.URL)
	conn := dialLobbyWSFromServer(t, server, created.ID)
	defer conn.Close()

	waitingSnap := readSnapshotFromConn(t, conn)
	if waitingSnap.Phase != lobby.PhaseWaiting {
		t.Fatalf("initial Phase = %q, want WAITING", waitingSnap.Phase)
	}

	photos := []lobby.Photo{
		{PixabayID: 1, LargeImageURL: "https://cdn.example/1.jpg", PreviewURL: "https://cdn.example/1s.jpg"},
		{PixabayID: 2, LargeImageURL: "https://cdn.example/2.jpg", PreviewURL: "https://cdn.example/2s.jpg"},
	}
	putPhotosViaHTTP(t, httpClient, server.URL, created, photos)
	selectingSnap := readSnapshotFromConn(t, conn)
	if selectingSnap.Phase != lobby.PhaseSelecting {
		t.Fatalf("Phase after photos = %q, want SELECTING", selectingSnap.Phase)
	}
	if selectingSnap.SelectedCount != len(photos) {
		t.Fatalf("SelectedCount = %d, want %d", selectingSnap.SelectedCount, len(photos))
	}

	postAdminViaHTTP(t, httpClient, server.URL, created, "/ready")
	readySnap := readSnapshotFromConn(t, conn)
	if readySnap.Phase != lobby.PhaseReady {
		t.Fatalf("Phase after ready = %q, want READY", readySnap.Phase)
	}
	if readySnap.TotalRounds != len(photos) {
		t.Fatalf("TotalRounds = %d, want %d", readySnap.TotalRounds, len(photos))
	}

	postAdminViaHTTP(t, httpClient, server.URL, created, "/start")
	drawingSnap := readSnapshotFromConn(t, conn)
	if drawingSnap.Phase != lobby.PhaseDrawing {
		t.Fatalf("Phase after start = %q, want DRAWING", drawingSnap.Phase)
	}
	if drawingSnap.CurrentRound != 1 {
		t.Fatalf("CurrentRound = %d, want 1", drawingSnap.CurrentRound)
	}
	if drawingSnap.CurrentPhoto == nil {
		t.Fatal("CurrentPhoto = nil, want photo in DRAWING")
	}

	betweenSnap := waitForSnapshotPhase(t, conn, lobby.PhaseBetweenRounds, 2*time.Second)
	if betweenSnap.CurrentPhoto != nil {
		t.Fatal("CurrentPhoto should be hidden in BETWEEN_ROUNDS")
	}

	postAdminViaHTTP(t, httpClient, server.URL, created, "/next")
	secondDrawingSnap := readSnapshotFromConn(t, conn)
	if secondDrawingSnap.Phase != lobby.PhaseDrawing {
		t.Fatalf("Phase after next = %q, want DRAWING", secondDrawingSnap.Phase)
	}
	if secondDrawingSnap.CurrentRound != 2 {
		t.Fatalf("CurrentRound = %d, want 2", secondDrawingSnap.CurrentRound)
	}

	betweenSnap = waitForSnapshotPhase(t, conn, lobby.PhaseBetweenRounds, 2*time.Second)
	if betweenSnap.Phase != lobby.PhaseBetweenRounds {
		t.Fatalf("Phase after timer = %q, want BETWEEN_ROUNDS", betweenSnap.Phase)
	}

	postAdminViaHTTP(t, httpClient, server.URL, created, "/next")
	finishedSnap := readSnapshotFromConn(t, conn)
	if finishedSnap.Phase != lobby.PhaseFinished {
		t.Fatalf("Phase after final next = %q, want FINISHED", finishedSnap.Phase)
	}
	if finishedSnap.CurrentPhoto != nil {
		t.Fatal("CurrentPhoto should be hidden in FINISHED")
	}
}

func createLobbyViaHTTP(t *testing.T, client *http.Client, baseURL string) createLobbyResponse {
	t.Helper()

	resp, err := client.Post(baseURL+"/api/lobbies", "application/json", nil)
	if err != nil {
		t.Fatalf("POST /api/lobbies error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("POST /api/lobbies status = %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	var created createLobbyResponse
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if created.ID == "" || created.AdminToken == "" {
		t.Fatal("create response missing id or admin_token")
	}
	return created
}

func putPhotosViaHTTP(t *testing.T, client *http.Client, baseURL string, created createLobbyResponse, photos []lobby.Photo) {
	t.Helper()

	body, err := json.Marshal(setPhotosRequest{Photos: photos})
	if err != nil {
		t.Fatalf("marshal photos: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, baseURL+"/api/lobbies/"+created.ID+"/photos", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("NewRequest() error = %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(lobby.AdminTokenHeader, created.AdminToken)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("PUT photos error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("PUT photos status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func postAdminViaHTTP(t *testing.T, client *http.Client, baseURL string, created createLobbyResponse, suffix string) {
	t.Helper()

	req, err := http.NewRequest(http.MethodPost, baseURL+"/api/lobbies/"+created.ID+suffix, nil)
	if err != nil {
		t.Fatalf("NewRequest() error = %v", err)
	}
	req.Header.Set(lobby.AdminTokenHeader, created.AdminToken)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("POST %s error = %v", suffix, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("POST %s status = %d, want %d", suffix, resp.StatusCode, http.StatusOK)
	}
}

func dialLobbyWSFromServer(t *testing.T, server *httptest.Server, lobbyID string) *websocket.Conn {
	t.Helper()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/lobby/" + lobbyID
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Dial() error = %v", err)
	}
	return conn
}

func waitForSnapshotPhase(t *testing.T, conn *websocket.Conn, phase lobby.LobbyPhase, timeout time.Duration) lobby.LobbySnapshot {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			break
		}
		readWait := 200 * time.Millisecond
		if remaining < readWait {
			readWait = remaining
		}

		snapshot, ok := tryReadSnapshotFromConn(conn, readWait)
		if ok && snapshot.Phase == phase {
			return snapshot
		}
	}
	t.Fatalf("timed out waiting for phase %q", phase)
	return lobby.LobbySnapshot{}
}

func tryReadSnapshotFromConn(conn *websocket.Conn, timeout time.Duration) (lobby.LobbySnapshot, bool) {
	_ = conn.SetReadDeadline(time.Now().Add(timeout))
	_, message, err := conn.ReadMessage()
	if err != nil {
		return lobby.LobbySnapshot{}, false
	}

	var envelope struct {
		Type    string          `json:"type"`
		Payload json.RawMessage `json:"payload"`
	}
	if err := json.Unmarshal(message, &envelope); err != nil {
		return lobby.LobbySnapshot{}, false
	}

	var snapshot lobby.LobbySnapshot
	if err := json.Unmarshal(envelope.Payload, &snapshot); err != nil {
		return lobby.LobbySnapshot{}, false
	}
	return snapshot, true
}
