package timer

import (
	"context"
	"testing"
	"time"
)

type mockExpiringStore struct {
	now      time.Time
	expired  []string
	calledAt time.Time
}

func (m *mockExpiringStore) ExpireDrawingTimers(_ context.Context, now time.Time) ([]string, error) {
	m.calledAt = now
	return append([]string(nil), m.expired...), nil
}

type mockBroadcaster struct {
	calls []string
}

func (m *mockBroadcaster) Broadcast(_ context.Context, lobbyID string) error {
	m.calls = append(m.calls, lobbyID)
	return nil
}

func TestSchedulerRunOnceBroadcastsExpiredLobbies(t *testing.T) {
	t.Parallel()

	fixedNow := time.Date(2026, 6, 20, 12, 5, 0, 0, time.UTC)
	store := &mockExpiringStore{
		now:     fixedNow,
		expired: []string{"lobby-1"},
	}
	broadcaster := &mockBroadcaster{}
	scheduler := NewScheduler(store, broadcaster, time.Second)
	scheduler.now = func() time.Time {
		return fixedNow
	}

	scheduler.RunOnce(context.Background())

	if !store.calledAt.Equal(fixedNow) {
		t.Fatalf("ExpireDrawingTimers now = %v, want %v", store.calledAt, fixedNow)
	}
	if len(broadcaster.calls) != 1 || broadcaster.calls[0] != "lobby-1" {
		t.Fatalf("Broadcast calls = %v, want [lobby-1]", broadcaster.calls)
	}
}

func TestSchedulerRunOnceNoBroadcastWhenNoneExpired(t *testing.T) {
	t.Parallel()

	store := &mockExpiringStore{}
	broadcaster := &mockBroadcaster{}
	scheduler := NewScheduler(store, broadcaster, time.Second)

	scheduler.RunOnce(context.Background())

	if len(broadcaster.calls) != 0 {
		t.Fatalf("Broadcast calls = %v, want none", broadcaster.calls)
	}
}

func TestSchedulerRunOnceSkipsWhenDependenciesNil(t *testing.T) {
	t.Parallel()

	scheduler := NewScheduler(nil, nil, time.Second)
	scheduler.RunOnce(context.Background())
}
