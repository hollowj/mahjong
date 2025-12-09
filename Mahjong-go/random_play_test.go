package mahjong

import (
	"testing"
	"time"
)

// TestRandomPlay 简化版随机自测：执行若干局，验证不会 panic 且能结束
func TestRandomPlay(t *testing.T) {
	games := 5
	start := time.Now()
	for i := 0; i < games; i++ {
		table := NewTable()
		pr := NewPaipuReplayer(table)
		pr.Init()
		res := pr.SimulateToCompletion()
		if res == nil {
			t.Fatalf("SimulateToCompletion returned nil result on game %d", i)
		}
	}
	dur := time.Since(start)
	t.Logf("%d random plays passed, duration=%v", games, dur)
}
