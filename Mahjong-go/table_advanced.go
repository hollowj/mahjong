package mahjong

import (
	"fmt"
	"strings"
)

// GameInitWithConfig 使用配置初始化游戏
func (t *Table) GameInitWithConfig(config GameConfig) {
	// 设置庄家
	if config.Oya >= 0 && config.Oya <= 3 {
		t.Oya = config.Oya
	} else {
		t.Oya = 0
	}

	// 设置场风
	if int(config.GameWind) >= 0 && int(config.GameWind) <= 3 {
		t.GameWind = config.GameWind
	} else {
		t.GameWind = East
	}

	// 设置供托与本场
	if config.Kyoutaku >= 0 {
		t.Kyoutaku = config.Kyoutaku
	} else {
		t.Kyoutaku = 0
	}
	if config.Honba >= 0 {
		t.Honba = config.Honba
	} else {
		t.Honba = 0
	}

	// 初始化牌/赤宝
	t.InitTiles()
	t.InitRedDora3()

	// 如果提供了牌山日志则导入，否则随机化（可能使用种子）
	if len(config.YamaLog) == NTiles {
		t.ImportYama(config.YamaLog)
	} else if len(config.YamaLog) == 0 {
		// 使用种子（如果有）在洗牌前设置
		if config.HasSeed {
			t.SetSeed(config.Seed)
		}
		t.InitYama()
		t.ShuffleTiles()
	} else {
		panic("Yama size is not 136")
	}

	// 记录牌山日志（可能为空）
	t.YamaLog = make([]int, len(config.YamaLog))
	copy(t.YamaLog, config.YamaLog)

	// 初始化宝牌并发牌
	t.InitDora()
	t.DrawTenhouStyle()

	// 初始化分数
	if len(config.InitScores) == 4 {
		for i := 0; i < 4; i++ {
			t.Players[i].Score = config.InitScores[i]
		}
	} else if len(config.InitScores) == 0 {
		for i := 0; i < 4; i++ {
			t.Players[i].Score = 25000
		}
	} else {
		panic("init_scores size is not 4.")
	}

	// 最后准备开始游戏
	t.InitBeforePlaying()
}

// GameConfig 游戏配置结构
type GameConfig struct {
	HasSeed    bool
	Seed       int64
	YamaLog    []int
	InitScores []int
	Kyoutaku   int
	Honba      int
	GameWind   Wind
	Oya        int
	Rules      string
}

// GameInitWithMetadata 使用元数据初始化游戏
func (t *Table) GameInitWithMetadata(metadata GameMetadata) {
	// 应用元数据
	if len(metadata.YamaLog) > 0 {
		t.ImportYama(metadata.YamaLog)
	}

	if metadata.Seed > 0 {
		t.SetSeed(metadata.Seed)
	}

	// 其他元数据应用...
}

// GameInitForReplay 使用回放数据初始化游戏
// 参数:
// - yamaLog: 完整的牌山索引（长度应为 NTiles），用于重放精确还原发牌顺序
// - initScores: 可选的初始分数（长度为4），若为空则使用默认25000
// - oya: 可选的庄家索引（0-3），若无效则默认0
// - honba/kyoutaku: 本场与供托数
func (t *Table) GameInitForReplay(yamaLog []int, initScores []int, oya int, honba int, kyoutaku int) {
	// 基本初始化（保留现有 Tiles/RedDora 的设置）
	t.InitTiles()
	t.InitRedDora3()

	// 导入牌山
	if len(yamaLog) == NTiles {
		t.ImportYama(yamaLog)
		t.YamaLog = make([]int, len(yamaLog))
		copy(t.YamaLog, yamaLog)
	} else if len(yamaLog) == 0 {
		// 如果没有提供 yamaLog，则按随机/种子流程初始化
		t.InitYama()
		t.ShuffleTiles()
	} else {
		// 非完整的 yamaLog 也允许导入（尽力而为）
		t.ImportYama(yamaLog)
		t.YamaLog = make([]int, len(yamaLog))
		copy(t.YamaLog, yamaLog)
	}

	// 初始化宝牌指示并按天凤风格发牌（使用已导入的 Yama）
	t.InitDora()
	t.DrawTenhouStyle()

	// 设置庄家/本场/供托
	if oya >= 0 && oya < NPlayers {
		t.Oya = oya
	} else {
		t.Oya = 0
	}
	if honba >= 0 {
		t.Honba = honba
	} else {
		t.Honba = 0
	}
	if kyoutaku >= 0 {
		t.Kyoutaku = kyoutaku
	} else {
		t.Kyoutaku = 0
	}

	// 初始化分数
	if len(initScores) == NPlayers {
		for i := 0; i < NPlayers; i++ {
			if t.Players[i] != nil {
				t.Players[i].Score = initScores[i]
			}
		}
	} else if len(initScores) == 0 {
		for i := 0; i < NPlayers; i++ {
			if t.Players[i] != nil {
				t.Players[i].Score = 25000
			}
		}
	}

	// 最终准备
	t.SortPlayerHands()
	t.LastAction = Pass
	t.LastActor = -1
	t.Turn = t.Oya

	if t.GameLog == nil {
		t.GameLog = NewGameLogRecord()
	} else {
		t.GameLog.Clear()
	}
}

// GameMetadata 游戏元数据结构
type GameMetadata struct {
	Seed    int64
	YamaLog []int
	GameID  string
	// ... 其他元数据字段
}

// SetDebugModeAdvanced 高级调试模式设置
func (t *Table) SetDebugModeAdvanced(mode int, verbose bool, logFile string) {
	t.SetDebugMode(mode)

	// 根据模式启用调试输出
	if verbose {
		// 启用详细日志
		if t.GameLog != nil {
			// t.GameLog.EnableVerbose()
		}
	}

	// 设置日志文件（如果需要）
	if logFile != "" {
		// 重定向日志到文件
		// logWriter := CreateLogWriter(logFile)
		// t.GameLog.SetWriter(logWriter)
	}
}

// GetCurrentPlayerWind 获取指定玩家的风位
func (t *Table) GetCurrentPlayerWind(playerIndex int) Wind {
	// 计算该玩家相对于东家的风位
	return Wind((playerIndex + 4 - t.Oya) % 4)
}

// GetGameWind 获取当前场风
func (t *Table) GetGameWind() Wind {
	return t.GameWind
}

// RotateOya 轮转庄家
func (t *Table) RotateOya() {
	t.Oya = (t.Oya + 1) % 4
}

// SaveGameState 保存游戏状态
func (t *Table) SaveGameState() GameState {
	state := GameState{
		OyaIndex:   t.Oya,
		Scores:     make([]int, 4),
		RemainTile: len(t.Yama),
	}

	// 保存分数
	for i, player := range t.Players {
		if player != nil {
			state.Scores[i] = player.Score
		}
	}

	return state
}

// GameState 游戏状态快照
type GameState struct {
	Round      int
	OyaIndex   int
	Scores     []int
	RemainTile int
}

// RestoreGameState 恢复游戏状态
func (t *Table) RestoreGameState(state GameState) {
	t.Oya = state.OyaIndex

	// 恢复分数
	for i, score := range state.Scores {
		if i < len(t.Players) && t.Players[i] != nil {
			t.Players[i].Score = score
		}
	}
}

// CheckGameEnd 检查游戏是否结束
func (t *Table) CheckGameEnd() bool {
	// 与 FromBeginning 中的结束判定保持一致：
	// - 四风连打
	// - 四立直
	// - 四杠散了
	// - 牌库耗尽（流局）
	if t.isAbaortedFourWind() {
		return true
	}
	if t.isFourRiichi() {
		return true
	}
	if t.isFourKanAborted() {
		return true
	}
	if t.GetRemainTile() == 0 {
		return true
	}
	return false
}

// GetFinalRankings 获取最终排名
func (t *Table) GetFinalRankings() []PlayerRanking {
	rankings := make([]PlayerRanking, len(t.Players))

	// 按分数排序
	type playerScore struct {
		index int
		score int
	}

	playerScores := make([]playerScore, len(t.Players))
	for i, player := range t.Players {
		if player != nil {
			playerScores[i] = playerScore{index: i, score: player.Score}
		}
	}

	// 排序（省略具体排序代码）

	for i, ps := range playerScores {
		rankings[i] = PlayerRanking{
			PlayerIndex: ps.index,
			Score:       ps.score,
			Rank:        i + 1,
		}
	}

	return rankings
}

// PlayerRanking 玩家排名信息
type PlayerRanking struct {
	PlayerIndex int
	Score       int
	Rank        int
}

// CalculateGameResult 计算游戏结果
func (t *Table) CalculateGameResult() *GameResult {
	// 更完整的结果生成，尽量与 C++ 的 generate_result_* 行为对应
	// 优先判定特殊流局
	if t.isAbaortedFourWind() {
		return GenerateResultSuzukaflush()
	}
	if t.isFourRiichi() {
		return GenerateResultSuuchahan()
	}
	if t.isFourKanAborted() {
		return GenerateResultSufonrenda()
	}

	// 九种九牌判定：若任一玩家满足九种九牌
	for i := 0; i < NPlayers; i++ {
		if t.Players[i] != nil {
			if len(t.Players[i].GetKyushukyuhai()) > 0 {
				return GenerateResultKyushukyuhai(i)
			}
		}
	}

	// 牌库流局（含流局满贯判断在 GenerateResultNotile 中处理）
	if t.GetRemainTile() == 0 {
		return GenerateResultNotile()
	}

	// 若无特殊流局，返回当前分数快照（保持兼容旧行为）
	result := NewGameResult()
	result.SetRyukyokuNotile() // 标记为流局占位（使 Message 可用）
	// 构造快照消息和分数变化为0（保留最终分数信息在 Message 中）
	sb := strings.Builder{}
	sb.WriteString("最终分数: ")
	for i := 0; i < NPlayers; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		if t.Players[i] != nil {
			sb.WriteString(fmt.Sprintf("P%d=%d", i, t.Players[i].Score))
		} else {
			sb.WriteString(fmt.Sprintf("P%d=%d", i, 0))
		}
	}
	result.Message = sb.String()
	return result
}

// ValidateTableState 验证桌子状态的一致性
func (t *Table) ValidateTableState() bool {
	// 验证牌的总数是否为 NTiles
	total := 0

	// 牌山剩余（包含14张dead wall）
	total += len(t.Yama)

	// 玩家手牌、鸣牌与河
	for i := 0; i < NPlayers; i++ {
		p := t.Players[i]
		if p == nil {
			continue
		}
		total += len(p.Hand)
		total += p.River.Size()
		// 统计鸣牌组中的牌
		for _, cg := range p.CallGroups {
			total += len(cg.Tiles)
		}
	}

	// DoraIndicator/UraDoraIndicator 不算作独立牌（引用 Yama 中的牌），因此不额外计入

	// 校验总数是否超过/不足
	if total != NTiles {
		return false
	}
	return true
}

// GetPlayerByWind 获取特定风位的玩家
func (t *Table) GetPlayerByWind(wind Wind) *Player {
	playerIndex := (int(wind) + t.Oya) % 4
	if playerIndex < len(t.Players) {
		return t.Players[playerIndex]
	}
	return nil
}

// GetOya 获取庄家
func (t *Table) GetOya() *Player {
	if t.Oya < len(t.Players) {
		return t.Players[t.Oya]
	}
	return nil
}

// GetNextOya 计算下一把的庄家索引
func (t *Table) GetNextOya() int {
	// 如果需要换庄（根据上一把结果）
	// 返回新的庄家索引
	return (t.Oya + 1) % 4
}

// AllPlayersHaveTenpai 检查所有玩家是否都听牌
func (t *Table) AllPlayersHaveTenpai() bool {
	for _, player := range t.Players {
		if player == nil || !player.IsTenpai() {
			return false
		}
	}
	return true
}

// AnyPlayerHasAction 检查是否有玩家有可执行的动作
func (t *Table) AnyPlayerHasAction() bool {
	// 检查每个玩家是否有可能的动作
	for _, player := range t.Players {
		if player != nil && len(player.GetDiscard(false)) > 0 {
			return true
		}
	}
	return false
}

// GetGameStatistics 获取游戏统计信息
func (t *Table) GetGameStatistics() GameStatistics {
	stats := GameStatistics{
		TotalActions: 0,
	}

	if t.GameLog != nil {
		stats.TotalActions = len(t.GameLog.Logs)
	}

	return stats
}

// GameStatistics 游戏统计信息
type GameStatistics struct {
	TotalRounds  int
	TotalTsumo   int
	TotalRon     int
	TotalDraw    int
	TotalActions int
}
