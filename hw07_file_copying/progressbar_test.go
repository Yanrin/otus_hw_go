package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestProgressBarRun(t *testing.T) {
	t.Run("final 100%", func(t *testing.T) {
		N := 7
		pb := NewProgressBar(int64(N))

		for i := 1; i <= N; i++ {
			pb.Add(1)
			time.Sleep(100 * time.Millisecond)
		}
		pb.Finish()
		require.Equal(t, int64(N), pb.count)
		require.Equal(t, pb.current, pb.count)
	})

	t.Run("50% and final 100%", func(t *testing.T) {
		N := 8
		pb := NewProgressBar(int64(N))

		for i := 1; i <= N; i++ {
			pb.Add(1)
			time.Sleep(100 * time.Millisecond)
			if i == 4 {
				require.Equal(t, pb.current, int64(N)/2)
			}
		}
		pb.Finish()
		require.Equal(t, pb.current, pb.count)
	})

	t.Run("overloaded 200%", func(t *testing.T) {
		N := 3
		pb := NewProgressBar(int64(N))

		for i := 1; i <= 2*N; i++ {
			pb.Add(1)
			time.Sleep(100 * time.Millisecond)
		}
		pb.Finish()
		require.Equal(t, pb.current, 2*pb.count)
	})

	t.Run("60%", func(t *testing.T) {
		N := 5
		limit := N - 2
		pb := NewProgressBar(int64(N))

		for i := 1; i <= limit; i++ {
			pb.Add(1)
			time.Sleep(100 * time.Millisecond)
		}
		pb.Finish()
		require.Equal(t, pb.current, int64(limit))
	})
	t.Run("zero count", func(t *testing.T) {
		N := 5
		delta := 10

		pb := NewProgressBar(int64(0))
		for i := 1; i <= N; i++ {
			pb.Add(int64(delta))
			time.Sleep(100 * time.Millisecond)
		}
		pb.Finish()
		require.Equal(t, pb.current, int64(delta*N))
	})
	t.Run("negative count", func(t *testing.T) {
		N := -7
		delta := 10

		pb := NewProgressBar(int64(N))
		for i := 1; i <= -N; i++ {
			pb.Add(int64(delta))
			time.Sleep(100 * time.Millisecond)
		}
		pb.Finish()
		require.Equal(t, pb.current, -int64(delta*N))
	})
}
