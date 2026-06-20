package timer

import (
	"testing"
	"time"
)

func TestEndsAt(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)
	duration := 5 * time.Minute

	got := EndsAt(now, duration)
	want := now.Add(duration)
	if !got.Equal(want) {
		t.Fatalf("EndsAt() = %v, want %v", got, want)
	}
}

func TestNewDeadline(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)
	duration := 5 * time.Minute

	deadline := NewDeadline(now, duration)
	if deadline == nil {
		t.Fatal("NewDeadline() = nil, want timestamp")
	}
	if !deadline.Equal(now.Add(duration)) {
		t.Fatalf("NewDeadline() = %v, want %v", deadline, now.Add(duration))
	}
}

func TestIsExpired(t *testing.T) {
	t.Parallel()

	endsAt := time.Date(2026, 6, 20, 12, 5, 0, 0, time.UTC)

	tests := []struct {
		name string
		at   time.Time
		want bool
	}{
		{"before expiry", endsAt.Add(-time.Second), false},
		{"at expiry", endsAt, true},
		{"after expiry", endsAt.Add(time.Second), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := IsExpired(&endsAt, tt.at); got != tt.want {
				t.Fatalf("IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}

	if IsExpired(nil, endsAt) {
		t.Fatal("IsExpired(nil) = true, want false")
	}
}

func TestRemaining(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)
	endsAt := now.Add(5 * time.Minute)

	if got := Remaining(&endsAt, now); got != 5*time.Minute {
		t.Fatalf("Remaining() = %v, want 5m", got)
	}
	if got := Remaining(&endsAt, endsAt); got != 0 {
		t.Fatalf("Remaining(at expiry) = %v, want 0", got)
	}
	if got := Remaining(&endsAt, endsAt.Add(time.Second)); got != 0 {
		t.Fatalf("Remaining(after expiry) = %v, want 0", got)
	}
	if got := Remaining(nil, now); got != 0 {
		t.Fatalf("Remaining(nil) = %v, want 0", got)
	}
}
