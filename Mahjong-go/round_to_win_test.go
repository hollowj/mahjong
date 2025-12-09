package mahjong

import "testing"

// TestLoadSyanten 验证 syanten.dat 可以加载且计算不会崩溃
func TestLoadSyanten(t *testing.T) {
	// 构造简单手牌：1m1m 2m3m4m 1p1p1p 2s3s4s
	hand := []*Tile{
		{Tile: _1m}, {Tile: _1m},
		{Tile: _2m}, {Tile: _3m}, {Tile: _4m},
		{Tile: _1p}, {Tile: _1p}, {Tile: _1p},
		{Tile: _2s}, {Tile: _3s}, {Tile: _4s},
		{Tile: _5s}, {Tile: _6s}, {Tile: _7s},
	}

	// 调用 CalculateRoundToWin（会加载 syanten.dat）
	res := CalculateRoundToWin(hand, 0)
	if res < 0 {
		t.Fatalf("CalculateRoundToWin returned negative: %d", res)
	}
}
