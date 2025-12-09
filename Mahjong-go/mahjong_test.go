package mahjong

import "testing"

// TestTileConstants 测试瓷牌常量
func TestTileConstants(t *testing.T) {
	if NBaseTiles != 34 {
		t.Errorf("Expected NBaseTiles=34, got %d", NBaseTiles)
	}
	if NTiles != 136 {
		t.Errorf("Expected NTiles=136, got %d", NTiles)
	}
	if NPlayers != 4 {
		t.Errorf("Expected NPlayers=4, got %d", NPlayers)
	}
}

// TestBaseTileValues 测试基础牌值
func TestBaseTileValues(t *testing.T) {
	if _1m != 0 {
		t.Error("_1m should be 0")
	}
	if _9m != 8 {
		t.Error("_9m should be 8")
	}
	if _1p != 9 {
		t.Error("_1p should be 9")
	}
	if _1s != 18 {
		t.Error("_1s should be 18")
	}
	if _1z != 27 {
		t.Error("_1z should be 27")
	}
}

// TestWindEnum 测试风向
func TestWindEnum(t *testing.T) {
	if East != 0 || South != 1 || West != 2 || North != 3 {
		t.Error("Wind values incorrect")
	}
}

// TestTileUtility 测试瓷牌工具函数
func TestTileUtility(t *testing.T) {
	if !Is1hai(_1m) {
		t.Error("_1m should be 1 tile")
	}
	if !Is9hai(_9m) {
		t.Error("_9m should be 9 tile")
	}
	if !IsYaochuhai(_1m) {
		t.Error("_1m should be yaochu tile")
	}
	if IsYaochuhai(_5m) {
		t.Error("_5m should not be yaochu tile")
	}
}

// TestTileBasics 测试瓷牌基础
func TestTileBasics(t *testing.T) {
	tile := &Tile{
		Tile:    _1m,
		RedDora: false,
		ID:      0,
	}

	str := tile.String()
	if str == "" {
		t.Error("Tile.String() should not be empty")
	}

	if BaseTileToString(_1m) == "" {
		t.Error("BaseTileToString should not be empty")
	}
}

// TestPlayerBasics 测试玩家基础
func TestPlayerBasics(t *testing.T) {
	player := &Player{
		Wind:   East,
		Oya:    true,
		Score:  30000,
	}

	if player.Wind != East {
		t.Error("Player wind should be East")
	}
	if player.Score != 30000 {
		t.Error("Player score should be 30000")
	}
}

// TestTableInitialization 测试桌子初始化
func TestTableInitialization(t *testing.T) {
	table := NewTable()
	table.InitTiles()

	if len(table.Tiles) != NTiles {
		t.Errorf("Expected %d tiles, got %d", NTiles, len(table.Tiles))
	}

	redCount := 0
	for _, tile := range table.Tiles {
		if tile != nil && tile.RedDora {
			redCount++
		}
	}

	// 初始化前红5牌数应该为0
	if redCount != 0 {
		t.Errorf("Before InitRedDora3, red dora count should be 0, got %d", redCount)
	}
}

// TestGameResultType 测试游戏结果类型
func TestGameResultType(t *testing.T) {
	result := &GameResult{
		Type:      RonAgari,
		WinnerIdx: 0,
		LoserIdx:  1,
	}

	if result.Type != RonAgari {
		t.Error("Result type should be RonAgari")
	}
	if !result.IsAgari() {
		t.Error("RonAgari should be agari")
	}
	if result.IsRyukyoku() {
		t.Error("RonAgari should not be ryukyoku")
	}
}

// TestActionType 测试动作类型
func TestActionType(t *testing.T) {
	if BaseActionToString(Discard) == "" {
		t.Error("Discard should have string representation")
	}
}

// TestYakuType 测试役型
func TestYakuType(t *testing.T) {
	if YakuToString(Tanyao) == "" {
		t.Error("Tanyao should have string representation")
	}

	if GetFanCount(Tanyao) <= 0 {
		t.Error("Tanyao should have positive fan count")
	}
}
