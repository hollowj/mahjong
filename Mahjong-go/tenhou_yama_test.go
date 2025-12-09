package mahjong

import "testing"

// TestTenhouYama 验证 DrawTenhouStyle 初始化后牌山尺寸和总牌数
func TestTenhouYama(t *testing.T) {
	table := NewTable()
	table.InitBeforePlaying()
	// DrawTenhouStyle 会把起始牌加入手牌（在 InitWithYama/Init 中也会被调用）
	table.DrawTenhouStyle()

	// 验证牌总数保持不变：牌山 + 所有玩家手牌 == NTiles
	total := len(table.Yama)
	for i := 0; i < NPlayers; i++ {
		total += len(table.Players[i].Hand)
	}
	if total != NTiles {
		t.Fatalf("expected total tiles %d, got %d (Yama=%d)", NTiles, total, len(table.Yama))
	}
}
