package timer

import (
	"context"
	"log/slog"
	"time"
)

// DefaultTickInterval is how often the scheduler checks for expired draw timers.
const DefaultTickInterval = 1 * time.Second

// ExpiringStore expires draw timers in active lobbies.
type ExpiringStore interface {
	ExpireDrawingTimers(ctx context.Context, now time.Time) ([]string, error)
}

// Broadcaster pushes lobby snapshot updates to connected clients.
type Broadcaster interface {
	Broadcast(ctx context.Context, lobbyID string) error
}

// Scheduler periodically expires DRAWING lobbies and broadcasts snapshot updates.
type Scheduler struct {
	store       ExpiringStore
	broadcaster Broadcaster
	interval    time.Duration
	now         func() time.Time
}

// NewScheduler creates a draw timer scheduler.
func NewScheduler(store ExpiringStore, broadcaster Broadcaster, interval time.Duration) *Scheduler {
	if interval <= 0 {
		interval = DefaultTickInterval
	}

	return &Scheduler{
		store:       store,
		broadcaster: broadcaster,
		interval:    interval,
		now:         time.Now,
	}
}

// Run ticks until ctx is canceled.
func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.RunOnce(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.RunOnce(ctx)
		}
	}
}

// RunOnce checks for expired draw timers and broadcasts updates once.
func (s *Scheduler) RunOnce(ctx context.Context) {
	if s.store == nil || s.broadcaster == nil {
		return
	}

	expired, err := s.store.ExpireDrawingTimers(ctx, s.now())
	if err != nil {
		slog.Error("expire drawing timers", "err", err)
		return
	}

	for _, lobbyID := range expired {
		if err := s.broadcaster.Broadcast(ctx, lobbyID); err != nil {
			slog.Error("broadcast after timer expiry", "lobby_id", lobbyID, "err", err)
			continue
		}
		slog.Info("drawing timer expired", "lobby_id", lobbyID)
	}
}
