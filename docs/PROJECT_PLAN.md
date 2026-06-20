# Croquis King — 프로젝트 구상 및 구현 계획

> 크로키(5분 따라그리기) 모임을 위한 실시간 동기화 웹앱  
> **백엔드:** Go · **프론트엔드:** React  
> 구현 순서: **백엔드 WorkItem(001–017) 완료 후 → 프론트엔드 WorkItem(101–112)**

---

## 1. 배경 및 문제

매주 수요일 화상 미팅으로 모여 PixaBay에서 사진을 고른 뒤, 한 사람이 이미지 뷰어와 타이머를 화면 공유하며 사진당 5분씩 따라그리는 모임을 진행 중이다.

### 현재 불편함

| 문제 | 설명 |
|------|------|
| 낮은 화면 공유 해상도 | 화상 공유로 사진을 보면 선명도가 떨어짐 |
| 타이머 + 사진 동시 공유 | 타이머 창 때문에 사진을 크게 보기 어려움 |
| 시작 전 조기 그리기 | 공식 시작 전에 그리기 시작하는 참가자 발생 |
| 수동 워크플로 | PixaBay 사이트 접속 → 검색 → 다운로드 → 뷰어/타이머 수동 운영 |

### 목표

- 참가자 전원이 **같은 URL**에 접속해 **동일한 화면 변화**를 실시간으로 본다.
- 화상 공유 없이 **고해상도 사진**과 **서버 동기화 타이머**를 제공한다.
- PixaBay 검색·선택을 앱 안에서 처리한다.
- 사진 순서는 **시작 전까지 비공개**, **랜덤**으로 결정한다.
- 그리기 시간(5분) 동안만 사진을 **최대 크기**로 표시하고, 시간이 끝나면 숨긴다.

---

## 2. 핵심 개념

### 2.1 역할

| 역할 | 설명 |
|------|------|
| **관리자(Admin)** | 로비를 생성·운영. 사진 검색·선택, 세션 시작·다음 사진 진행 |
| **참가자(Participant)** | 로비 링크로 접속. 화면을 보기만 함 (초기 버전) |

관리자 권한은 로비 생성 시 발급되는 **관리자 토큰**으로 구분한다.  
참가자는 토큰 없이 링크만으로 접속 가능하다.

### 2.2 로비(Lobby)

하나의 그리기 세션 단위. 관리자가 로비를 열면 **공유 가능한 URL**이 생성된다.

```
https://{host}/lobby/{lobby_id}
```

같은 `lobby_id`에 접속한 모든 클라이언트는 WebSocket을 통해 **서버가 권위를 가진 상태(state)** 를 동기화받는다.

### 2.3 세션 상태 머신

서버가 상태를 단일 진실 공급원(Single Source of Truth)으로 관리한다.  
**타이머·사진 노출 여부는 모두 서버 시각 기준**으로 결정하여 조기 그리기를 방지한다.

```
                    ┌─────────────┐
                    │   WAITING   │  로비 대기 (참가자 접속)
                    └──────┬──────┘
                           │ 관리자: 사진 선택 시작
                           ▼
                    ┌─────────────┐
                    │  SELECTING  │  PixaBay 검색·사진 선택
                    └──────┬──────┘
                           │ 관리자: 선택 완료
                           ▼
                    ┌─────────────┐
                    │    READY    │  순서 셔플 완료, 사진 비공개
                    └──────┬──────┘
                           │ 관리자: 세션 시작
                           ▼
              ┌────────────────────────┐
              │       DRAWING          │  현재 사진 표시 + 5분 타이머
              └───────────┬────────────┘
                          │ 타이머 만료 (서버 자동)
                          ▼
              ┌────────────────────────┐
              │   BETWEEN_ROUNDS       │  사진 숨김, 다음 준비
              └───────────┬────────────┘
                          │ 관리자: 다음 사진 / 또는 자동 진행(옵션)
                          ▼
                    (마지막이 아니면 DRAWING 반복)
                          │
                          ▼
                    ┌─────────────┐
                    │  FINISHED   │  전체 세션 종료
                    └─────────────┘
```

**사진 노출 규칙**

- `READY`, `BETWEEN_ROUNDS`, `FINISHED`, `SELECTING`: 사진 **미표시**
- `DRAWING`: 현재 순번 사진만 **전체 화면 최대 크기**로 표시
- `photo_order`(셔플된 인덱스 배열)는 **관리자·참가자 모두** `DRAWING` 시작 전까지 API/WebSocket으로 전달하지 않음

---

## 3. 기술 아키텍처

### 3.1 전체 구조

```
┌──────────────┐     HTTP/WS      ┌─────────────────────────────┐
│ React SPA    │ ◄──────────────► │ Go Backend                  │
│ (frontend/)  │                  │  ├─ REST API                │
└──────────────┘                  │  ├─ WebSocket Hub           │
                                  │  ├─ Lobby State Machine     │
                                  │  ├─ Timer Scheduler         │
                                  │  └─ PixaBay Proxy Client    │
                                  └──────────────┬──────────────┘
                                                 │ HTTPS
                                                 ▼
                                  ┌─────────────────────────────┐
                                  │ PixaBay API                 │
                                  │ https://pixabay.com/api/    │
                                  └─────────────────────────────┘
```

### 3.2 실시간 동기화

- **WebSocket** (`/ws/lobby/{id}`): 로비 상태 변경, 타이머 잔여 시간, 참가자 수 등 브로드캐스트
- **REST API**: 로비 생성, PixaBay 검색, 관리자 액션(사진 선택·세션 시작 등)
- 타이머는 서버에서 `draw_ends_at`(Unix ms)을 저장하고, 클라이언트는 수신한 시각을 기준으로 UI만 렌더링  
  → 네트워크 지연이 있어도 **서버 만료 시점**이 기준

### 3.3 PixaBay API 연동

- API Key는 **백엔드 환경변수**에만 보관 (`PIXABAY_API_KEY`)
- 프론트엔드는 백엔드 REST를 호출; PixaBay에 직접 요청하지 않음
- 참고: [PixaBay API Docs](https://pixabay.com/api/docs/)

**초기 지원 파라미터**

| 파라미터 | 설명 | 기본값 |
|----------|------|--------|
| `q` | 검색어 (필수) | — |
| `order` | `popular` \| `latest` | `popular` |
| `page` | 페이지 | `1` |
| `per_page` | 페이지당 결과 (3–200) | `20` |
| `image_type` | `photo` 권장 | `photo` |
| `safesearch` | 전연령 필터 | `true` |

**응답에서 사용할 필드 (예시)**

- `id`, `pageURL`, `previewURL`, `webformatURL`, `largeImageURL`
- `imageWidth`, `imageHeight`, `views`, `downloads`, `likes`
- UI 썸네일: `previewURL` / 그리기용: `largeImageURL`

**PixaBay 이용 조건**: 검색 결과 표시 시 이미지 출처(PixaBay) 표기 필요 ([API 가이드라인](https://pixabay.com/api/docs/))

### 3.4 프로젝트 디렉터리 (목표 구조)

```
croquis-king/
├── cmd/
│   └── server/
│       └── main.go              # 진입점
├── internal/
│   ├── config/                  # 환경설정
│   ├── lobby/                   # 로비 도메인·상태머신·저장소
│   ├── pixabay/                 # PixaBay API 클라이언트
│   ├── timer/                   # 서버 타이머 스케줄러
│   ├── ws/                      # WebSocket 허브
│   └── http/                    # HTTP 핸들러·라우터
├── frontend/                    # React SPA (WorkItem 101 이후)
├── docs/
│   └── PROJECT_PLAN.md          # 본 문서
├── go.mod
└── README.md
```

### 3.5 주요 데이터 모델 (초안)

```go
// Lobby — 서버 내부 전체 상태
type Lobby struct {
    ID            string
    AdminToken    string        // 생성 시 1회 반환, 이후 헤더로 인증
    Phase         LobbyPhase
    SelectedPhotos []Photo      // 관리자가 고른 사진 목록
    PhotoOrder    []int          // SelectedPhotos 인덱스의 셔플 결과 (비공개)
    CurrentRound  int            // 0-based, PhotoOrder 내 위치
    DrawDuration  time.Duration  // 기본 5m
    DrawEndsAt    *time.Time     // DRAWING 중에만 설정
    CreatedAt     time.Time
}

type Photo struct {
    PixabayID     int
    PreviewURL    string
    LargeImageURL string
    PageURL       string
    Width         int
    Height        int
}

// LobbySnapshot — WebSocket/API로 클라이언트에 전달하는 공개 스냅샷
// Phase·DrawEndsAt·CurrentRound 등은 포함하되,
// PhotoOrder 및 아직 공개되지 않은 사진 URL은 Phase에 따라 마스킹
type LobbySnapshot struct {
    ID             string
    Phase          LobbyPhase
    ParticipantCount int
    SelectedCount  int            // 선택된 사진 수 (URL 없이)
    CurrentRound   int            // 1-based 표시용 (DRAWING/FINISHED에서만)
    TotalRounds    int
    DrawEndsAt     *time.Time     // DRAWING에서만
    CurrentPhoto   *Photo         // DRAWING에서만
    ServerTime     time.Time      // 클라이언트 시계 보정용
}
```

---

## 4. API 설계 (초안)

### 4.1 REST

| Method | Path | 권한 | 설명 |
|--------|------|------|------|
| `POST` | `/api/lobbies` | — | 로비 생성 → `{ id, admin_token, join_url }` |
| `GET` | `/api/lobbies/{id}` | — | 현재 스냅샷 조회 |
| `GET` | `/api/pixabay/search` | Admin* | PixaBay 검색 프록시 |
| `PUT` | `/api/lobbies/{id}/photos` | Admin | 선택 사진 목록 설정 |
| `POST` | `/api/lobbies/{id}/ready` | Admin | 선택 완료 → 셔플 → READY |
| `POST` | `/api/lobbies/{id}/start` | Admin | 첫 DRAWING 시작 |
| `POST` | `/api/lobbies/{id}/next` | Admin | BETWEEN → 다음 DRAWING |
| `POST` | `/api/lobbies/{id}/finish` | Admin | 세션 강제 종료 |

\* PixaBay 검색은 로비 AdminToken 헤더(`X-Admin-Token`) 또는 로비 ID + 토큰 검증

### 4.2 WebSocket

| Path | 설명 |
|------|------|
| `GET /ws/lobby/{id}` | 로비 구독. 연결 시 현재 스냅샷 전송, 이후 `snapshot` 이벤트 브로드캐스트 |

**메시지 예시**

```json
{
  "type": "snapshot",
  "payload": { /* LobbySnapshot */ }
}
```

---

## 5. 프론트엔드 화면 (개요)

백엔드 완료 후 구현. 화면별 요약만 기술한다.

| 화면 | Phase | 설명 |
|------|-------|------|
| 로비 대기 | `WAITING` | 참가자 수, "사진 고르는 중…" 안내 |
| 사진 선택 | `SELECTING` | 검색·정렬·썸네일 그리드·선택 (Admin만 조작) |
| 준비 | `READY` | "N장 준비됨. 시작을 기다리는 중" — 사진 미리보기 없음 |
| 그리기 | `DRAWING` | 상단 얇은 progress bar + 남은 시간, 중앙 최대 크기 사진 |
| 라운드 간 | `BETWEEN_ROUNDS` | 사진 숨김, "잠시 쉬어가요" |
| 종료 | `FINISHED` | 세션 완료 |

**DRAWING UI 상세**

- 상단: viewport 너비 100% 얇은 bar, `남은시간 / 5:00` 비율로 줄어듦
- 중앙: `object-fit: contain`으로 가능한 한 크게 표시
- 하단(선택): 현재 라운드 `3 / 5` 정도의 최소 정보만

---

## 6. WorkItem 목록

각 항목은 **인덱스 번호**로 구분한다.  
"**003 구현해줘**"처럼 번호만 지정하면 해당 WorkItem만 진행할 수 있다.

**범례:** `[BE]` 백엔드 · `[FE]` 프론트엔드 · `deps` = 선행 WorkItem

---

### Phase A — 백엔드 (001–017)

| Index | 제목 | 설명 | deps |
|-------|------|------|------|
| **001** | `[BE]` 프로젝트 스캐폴딩 | `cmd/server`, `internal/` 디렉터리 구조, `go mod`, Makefile 또는 `go run` 진입점 정리. 루트 `main.go`는 `cmd/server`로 이전 | — [DONE] |
| **002** | `[BE]` 설정·환경변수 | `PORT`, `PIXABAY_API_KEY`, `CORS_ORIGINS`, `DRAW_DURATION`(기본 5m) 등 envconfig 로딩 | 001 [DONE] |
| **003** | `[BE]` HTTP 서버 골격 | chi/gin 등 라우터, health check (`GET /health`), graceful shutdown | 002 |
| **004** | `[BE]` 로비 도메인 모델 | `Lobby`, `Photo`, `LobbyPhase`, `LobbySnapshot` 타입 및 Phase 전환 규칙 정의 | 001 |
| **005** | `[BE]` 로비 저장소 | 인메모리 `LobbyStore` (map + mutex). 생성·조회·스냅샷 변환. 추후 Redis 등으로 교체 가능하게 interface 분리 | 004 |
| **006** | `[BE]` 로비 생성·조회 API | `POST /api/lobbies`, `GET /api/lobbies/{id}`. AdminToken 생성·검증 헬퍼 | 003, 005 |
| **007** | `[BE]` WebSocket 허브 | `/ws/lobby/{id}` 업그레이드, 연결 등록/해제, 로비별 브로드캐스트 | 003, 005 |
| **008** | `[BE]` 실시간 스냅샷 동기화 | 로비 상태 변경 시 모든 WS 클라이언트에 `snapshot` push. 연결 직후 초기 스냅샷 전송 | 006, 007 |
| **009** | `[BE]` PixaBay API 클라이언트 | HTTP 클라이언트, 검색 요청/응답 파싱, API key 주입, 에러·rate limit 처리 | 002 |
| **010** | `[BE]` PixaBay 검색 API | `GET /api/pixabay/search?q&order&page&per_page`. AdminToken 검증, PixaBay 응답 정규화 | 006, 009 |
| **011** | `[BE]` 사진 선택 API | `PUT /api/lobbies/{id}/photos` — 선택 목록 저장, Phase를 `SELECTING`으로, 스냅샷 갱신 | 006, 008 |
| **012** | `[BE]` 사진 순서 랜덤 셔플 | `POST .../ready` — Fisher-Yates 셔플로 `PhotoOrder` 생성. `PhotoOrder`는 스냅샷에 **미포함** | 011 |
| **013** | `[BE]` 서버 권위 타이머 | `DrawEndsAt` 설정/초기화, `DRAWING` 진입 시 `now + DrawDuration`, 만료 판정 헬퍼 | 004, 005 |
| **014** | `[BE]` 세션 진행 API | `POST .../start`, `.../next`, `.../finish`. Phase 전환 + 스냅샷에서 현재 사진만 노출 | 012, 013 |
| **015** | `[BE]` 타이머 틱·자동 전환 | 백그라운드 goroutine/ticker: `DRAWING` 중 `DrawEndsAt` 경과 시 자동으로 `BETWEEN_ROUNDS` 전환 및 WS broadcast | 013, 014, 008 |
| **016** | `[BE]` CORS·Admin 미들웨어 | `X-Admin-Token` 검증 미들웨어, CORS 설정, 프론트 dev 서버 origin 허용 | 003, 006 |
| **017** | `[BE]` 백엔드 통합 테스트 | 로비 생성 → 사진 선택 → ready → start → 타이머 만료 → next → finish 시나리오 HTTP/WS 테스트 | 014, 015 |

---

### Phase B — 프론트엔드 (101–112)

> **선행 조건:** WorkItem **001–017** (백엔드) 완료

| Index | 제목 | 설명 | deps |
|-------|------|------|------|
| **101** | `[FE]` React 프로젝트 초기화 | Vite + React + TypeScript, `frontend/` 디렉터리, ESLint/Prettier, env (`VITE_API_BASE`) | 017 |
| **102** | `[FE]` 라우팅·로비 접속 | `/` (로비 생성), `/lobby/:id` (참가). AdminToken은 create 응답 후 sessionStorage 저장 | 101 |
| **103** | `[FE]` WebSocket 클라이언트 훅 | `useLobbySocket(id)` — snapshot 수신, reconnect, 서버 시각 offset 보정 | 102 |
| **104** | `[FE]` 로비 공통 레이아웃 | 참가자 수, Phase별 안내 문구, 로딩/연결 끊김 UI | 103 |
| **105** | `[FE]` Admin / Participant UI 분기 | AdminToken 유무에 따라 조작 버튼 노출/비노출 | 104 |
| **106** | `[FE]` PixaBay 검색 UI | 검색어 입력, `popular`/`latest` 정렬, 페이지네이션, 결과 그리드 (Admin, `SELECTING`) | 105, 010 |
| **107** | `[FE]` 사진 선택 UI | 다중 선택(기본 5장 권장), 선택 목록 확인, "선택 완료" → ready API | 106, 011 |
| **108** | `[FE]` READY 화면 | "N장 준비됨" — 썸네일·순서 미표시. Admin "시작" 버튼 | 107, 012 |
| **109** | `[FE]` DRAWING 화면 | 상단 progress bar + `mm:ss` 잔여 시간, 중앙 `largeImageURL` 최대 표시, PixaBay 출처 표기 | 108, 013 |
| **110** | `[FE]` BETWEEN / FINISHED 화면 | 사진 숨김, 라운드 간 대기, Admin "다음" 버튼, 종료 메시지 | 109, 014 |
| **111** | `[FE]` 반응형·전체화면 UX | 모바일/태블릿 대응, 불필요 chrome 최소화, 그리기에 방해되지 않는 UI | 109 |
| **112** | `[FE]` 프론트엔드 통합·빌드 | 프로덕션 빌드, Go static embed 또는 reverse proxy 연동 문서화 | 110, 111, 016 |

---

## 7. 구현 순서 요약

```
001 → 002 → 003
         ↓
004 → 005 → 006 → 007 → 008
         ↓
009 → 010
         ↓
011 → 012 → 013 → 014 → 015
         ↓
016 (병행 가능) → 017
         ↓
    [백엔드 완료]
         ↓
101 → 102 → 103 → 104 → 105
         ↓
106 → 107 → 108 → 109 → 110 → 111 → 112
```

**권장 첫 구현 묶음**

1. **MVP-Backend (001–008):** 로비 생성 + WebSocket 동기화만으로 "같은 화면을 본다" 검증
2. **MVP-PixaBay (009–011):** 검색·선택
3. **MVP-Session (012–015):** 셔플·타이머·자동 전환
4. **MVP-Frontend (101–109):** 실제 모임에 쓸 수 있는 최소 UI

---

## 8. 비기능 요구사항 (초기)

| 항목 | 방침 |
|------|------|
| 동시 로비 수 | 초기: 인메모리, 단일 인스턴스. 수십 명 규모면 충분 |
| 보안 | AdminToken은 UUID v4+, HTTPS 배포 권장 |
| PixaBay Rate Limit | 검색 debounce(프론트), 백엔드 simple cache(동일 query 30s) 고려 |
| 배포 | 단일 바이너리 + React static. Docker 선택 |

---

## 9. 향후 확장 (본 문서 범위外, 참고용)

- 참가자 닉네임·프레즌스
- `next` 자동 진행(countdown 후)
- 로비 영속화(Redis/DB)
- 그리기 결과물 업로드·갤러리
- `DRAW_DURATION` 로비별 커스터마izatio
- 다국어 UI

---

## 10. 환경변수 (예정)

| 변수 | 필수 | 설명 |
|------|------|------|
| `PORT` | | HTTP 포트 (default: `8080`) |
| `PIXABAY_API_KEY` | ✓ | PixaBay API 키 |
| `DRAW_DURATION` | | 예: `5m` |
| `CORS_ORIGINS` | | 예: `http://localhost:5173` |

---

*문서 버전: 2025-06-20 · croquis-king 초기 구상*
