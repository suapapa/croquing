package lobby

import "testing"

func TestShuffleIndices(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		order, err := shuffleIndices(0)
		if err != nil {
			t.Fatalf("shuffleIndices() error = %v", err)
		}
		if order != nil {
			t.Fatalf("order = %v, want nil", order)
		}
	})

	t.Run("single", func(t *testing.T) {
		t.Parallel()

		order, err := shuffleIndices(1)
		if err != nil {
			t.Fatalf("shuffleIndices() error = %v", err)
		}
		if len(order) != 1 || order[0] != 0 {
			t.Fatalf("order = %v, want [0]", order)
		}
	})

	t.Run("permutation", func(t *testing.T) {
		t.Parallel()

		const n = 20
		order, err := shuffleIndices(n)
		if err != nil {
			t.Fatalf("shuffleIndices() error = %v", err)
		}
		if len(order) != n {
			t.Fatalf("len(order) = %d, want %d", len(order), n)
		}

		seen := make([]bool, n)
		for _, idx := range order {
			if idx < 0 || idx >= n {
				t.Fatalf("index out of range: %d", idx)
			}
			if seen[idx] {
				t.Fatalf("duplicate index: %d", idx)
			}
			seen[idx] = true
		}
	})
}
