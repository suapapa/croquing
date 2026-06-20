# Croquis King — 구현 진도

> **자동 생성 문서** — 직접 수정하지 마세요.  
> 소스: [`workitems.json`](workitems.json) · 갱신: `make progress`

마지막 갱신: **2026-06-20 18:54 KST**

## Phase A — 백엔드 (001-017)

진행: **9 / 17** (52%)

| Index | 상태 | 제목 | deps | 산출물 검증 |
|-------|------|------|------|-------------|
| **001** | ✅ done | [BE] 프로젝트 스캐폴딩 | — | ✓ |
| **002** | ✅ done | [BE] 설정·환경변수 | 001 | ✓ |
| **003** | ✅ done | [BE] HTTP 서버 골격 | 002 | ✓ |
| **004** | ✅ done | [BE] 로비 도메인 모델 | 001 | ✓ |
| **005** | ✅ done | [BE] 로비 저장소 | 004 | ✓ |
| **006** | ✅ done | [BE] 로비 생성·조회 API | 003, 005 | ✓ |
| **007** | ✅ done | [BE] WebSocket 허브 | 003, 005 | ✓ |
| **008** | ⬜ pending | [BE] 실시간 스냅샷 동기화 | 006, 007 | — |
| **009** | ✅ done | [BE] PixaBay API 클라이언트 | 002 | ✓ |
| **010** | ✅ done | [BE] PixaBay 검색 API | 006, 009 | ✓ |
| **011** | ⬜ pending | [BE] 사진 선택 API | 006, 008 | — |
| **012** | ⬜ pending | [BE] 사진 순서 랜덤 셔플 | 011 | — |
| **013** | ⬜ pending | [BE] 서버 권위 타이머 | 004, 005 | — |
| **014** | ⬜ pending | [BE] 세션 진행 API | 012, 013 | — |
| **015** | ⬜ pending | [BE] 타이머 틱·자동 전환 | 013, 014, 008 | — |
| **016** | ⬜ pending | [BE] CORS·Admin 미들웨어 | 003, 006 | — |
| **017** | ⬜ pending | [BE] 백엔드 통합 테스트 | 014, 015 | — |

### 완료 항목 상세

#### 001 — [BE] 프로젝트 스캐폴딩

- 완료일: 2026-06-20
- 커밋: `bce5ae5`
- 산출물: cmd/server/main.go, internal/, Makefile, go.mod
- 메모: 루트 main.go 제거. make run / make build 동작.

#### 002 — [BE] 설정·환경변수

- 완료일: 2026-06-20
- 커밋: `bce5ae5`
- 산출물: internal/config/config.go
- 메모: caarlosh/envconfig 사용. cmd/server에서 config.Load() 호출.

#### 003 — [BE] HTTP 서버 골격

- 완료일: 2026-06-20
- 산출물: internal/http/router.go, internal/http/server.go, internal/http/router_test.go
- 메모: Gin 라우터, GET /health, SIGINT/SIGTERM graceful shutdown.

#### 004 — [BE] 로비 도메인 모델

- 완료일: 2026-06-20
- 산출물: internal/lobby/model.go, internal/lobby/phase.go, internal/lobby/errors.go, internal/lobby/auth.go, internal/lobby/lobby_test.go
- 메모: Phase 전환 규칙, Snapshot 마스킹(DRAWING에서만 사진·타이머 노출).

#### 005 — [BE] 로비 저장소

- 완료일: 2026-06-20
- 산출물: internal/lobby/store.go
- 메모: Store interface + MemoryStore. UUID id/admin token, clone on read.

#### 006 — [BE] 로비 생성·조회 API

- 완료일: 2026-06-20
- 산출물: internal/http/lobby_handlers.go, internal/http/lobby_handlers_test.go
- 메모: POST create → id/admin_token/join_url. GET snapshot. X-Admin-Token 검증 헬퍼.

#### 007 — [BE] WebSocket 허브

- 완료일: 2026-06-20
- 산출물: internal/ws/hub.go, internal/ws/client.go, internal/ws/handler.go, internal/ws/message.go, internal/ws/hub_test.go, internal/http/ws_handlers.go
- 메모: GET /ws/lobby/:id. Hub Register/Unregister/Broadcast, ping/pong, ClientCount.

#### 009 — [BE] PixaBay API 클라이언트

- 완료일: 2026-06-20
- 산출물: internal/pixabay/client.go, internal/pixabay/types.go, internal/pixabay/errors.go, internal/pixabay/client_test.go
- 메모: Search API 클라이언트. 429 ErrRateLimited, X-RateLimit-* 헤더 파싱, httptest 단위 테스트.

#### 010 — [BE] PixaBay 검색 API

- 완료일: 2026-06-20
- 산출물: internal/http/pixabay_handlers.go, internal/http/admin_auth.go, internal/http/pixabay_handlers_test.go
- 메모: GET /api/pixabay/search. lobby_id + X-Admin-Token 검증. snake_case JSON 정규화.

## Phase B — 프론트엔드 (101-112)

선행 조건: backend 완료 (001-017)

진행: **0 / 12** (0%)

| Index | 상태 | 제목 | deps | 산출물 검증 |
|-------|------|------|------|-------------|
| **101** | ⬜ pending | [FE] React 프로젝트 초기화 | 017 | — |
| **102** | ⬜ pending | [FE] 라우팅·로비 접속 | 101 | — |
| **103** | ⬜ pending | [FE] WebSocket 클라이언트 훅 | 102 | — |
| **104** | ⬜ pending | [FE] 로비 공통 레이아웃 | 103 | — |
| **105** | ⬜ pending | [FE] Admin / Participant UI 분기 | 104 | — |
| **106** | ⬜ pending | [FE] PixaBay 검색 UI | 105, 010 | — |
| **107** | ⬜ pending | [FE] 사진 선택 UI | 106, 011 | — |
| **108** | ⬜ pending | [FE] READY 화면 | 107, 012 | — |
| **109** | ⬜ pending | [FE] DRAWING 화면 | 108, 013 | — |
| **110** | ⬜ pending | [FE] BETWEEN / FINISHED 화면 | 109, 014 | — |
| **111** | ⬜ pending | [FE] 반응형·전체화면 UX | 109 | — |
| **112** | ⬜ pending | [FE] 프론트엔드 통합·빌드 | 110, 111, 016 | — |

## 전체 요약

- **전체:** 9 / 29 완료 (31%)
- **백엔드:** 9 / 17 (52%)
- **프론트엔드:** 0 / 12 (0%)

## 다음 작업 후보

- **008** — [BE] 실시간 스냅샷 동기화 (`pending`)
- **013** — [BE] 서버 권위 타이머 (`pending`)
- **016** — [BE] CORS·Admin 미들웨어 (`pending`)

## 진도 갱신 방법

1. [`workitems.json`](workitems.json)에서 해당 WorkItem의 `status`를 수정합니다.
   - `pending` · `in_progress` · `done` · `blocked`
2. 완료 시 `completed_at`, `commit`, `artifacts`, `notes`를 채웁니다.
3. 저장 후 `make progress`를 실행합니다.

상세 스펙·API·아키텍처는 [`../PROJECT_PLAN.md`](../PROJECT_PLAN.md)를 참고하세요.
