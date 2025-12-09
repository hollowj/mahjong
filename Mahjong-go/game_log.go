package mahjong

// LogAction 表示游戏日志中记录的各种动作
type LogAction int

const (
	LogAnKan                  LogAction = iota // 暗杠
	LogPon                                     // 碰
	LogChi                                     // 吃
	LogKan                                     // 杠
	LogKaKan                                   // 加杠
	LogDiscardFromHand                         // 手切
	LogDiscardFromTsumo                        // 摸切
	LogRiichiDiscardFromHand                   // 宣告手切立
	LogRiichiDiscardFromTsumo                  // 宣告摸切立
	LogRiichiSuccess                           // 立直通过
	LogDrawNormal                              // 正常摸牌
	LogDrawRinshan                             // 摸岭上牌
	LogDoraReveal                              // 翻1张宝牌
	LogKyushukyuhai                            // 宣告九种九牌
	LogRon                                     // 宣告Ron
	LogTsumo                                   // 宣告自摸
	LogInvalidAction                           // 无效动作
)

// LogActionToString 将LogAction转换为字符串
func LogActionToString(action LogAction) string {
	names := []string{
		"AnKan",                  // LogAnKan
		"Pon",                    // LogPon
		"Chi",                    // LogChi
		"Kan",                    // LogKan
		"KaKan",                  // LogKaKan
		"DiscardFromHand",        // LogDiscardFromHand
		"DiscardFromTsumo",       // LogDiscardFromTsumo
		"RiichiDiscardFromHand",  // LogRiichiDiscardFromHand
		"RiichiDiscardFromTsumo", // LogRiichiDiscardFromTsumo
		"RiichiSuccess",          // LogRiichiSuccess
		"DrawNormal",             // LogDrawNormal
		"DrawRinshan",            // LogDrawRinshan
		"DoraReveal",             // LogDoraReveal
		"Kyushukyuhai",           // LogKyushukyuhai
		"Ron",                    // LogRon
		"Tsumo",                  // LogTsumo
		"InvalidAction",          // LogInvalidAction
	}

	if int(action) < len(names) {
		return names[action]
	}
	return "Unknown"
}

// BaseGameLog 基础游戏日志
// 记录游戏中发生的各种事件
type BaseGameLog struct {
	Player    int           // 主要玩家索引
	Player2   int           // 副手牌(用于部分动作)
	Action    LogAction     // 动作类型
	Tile      *Tile         // 相关的牌
	CallTiles []*Tile       // 面子牌组
	Scores    [NPlayers]int // 各玩家的分数
}

// NewBaseGameLogWithAction 使用动作创建日志
func NewBaseGameLogWithAction(p1, p2 int, action LogAction, tile *Tile, callGroup []*Tile) *BaseGameLog {
	log := &BaseGameLog{
		Player:    p1,
		Player2:   p2,
		Action:    action,
		Tile:      tile,
		CallTiles: make([]*Tile, len(callGroup)),
		Scores:    [NPlayers]int{},
	}

	// 复制callGroup
	copy(log.CallTiles, callGroup)

	return log
}

// NewBaseGameLogWithScores 使用分数创建日志(用于结算)
func NewBaseGameLogWithScores(scores [NPlayers]int) *BaseGameLog {
	log := &BaseGameLog{
		Player:    -1,
		Player2:   -1,
		Action:    LogInvalidAction,
		Tile:      nil,
		CallTiles: []*Tile{},
		Scores:    scores,
	}
	return log
}

// String 返回日志的字符串表示
func (g *BaseGameLog) String() string {
	result := ""

	if g.Action == LogInvalidAction && g.Player == -1 {
		// 这是分数日志
		result += "SCORE: "
		for i := 0; i < NPlayers; i++ {
			if i > 0 {
				result += ", "
			}
			result += "P" + string(rune('0'+i)) + "=" + string(rune(g.Scores[i]))
		}
	} else {
		// 这是动作日志
		result += "P" + string(rune('0'+g.Player)) + " "
		result += LogActionToString(g.Action)

		if g.Tile != nil {
			result += " " + g.Tile.String()
		}

		if len(g.CallTiles) > 0 {
			result += " ("
			for i, tile := range g.CallTiles {
				if i > 0 {
					result += ","
				}
				result += tile.String()
			}
			result += ")"
		}
	}

	return result
}

// GameLogRecord 游戏日志记录器
// 用于记录完整的游戏过程
type GameLogRecord struct {
	Logs []*BaseGameLog // 所有日志条目
}

// NewGameLogRecord 创建新的游戏日志记录器
func NewGameLogRecord() *GameLogRecord {
	return &GameLogRecord{
		Logs: make([]*BaseGameLog, 0),
	}
}

// AddLog 添加一条日志
func (g *GameLogRecord) AddLog(log *BaseGameLog) {
	g.Logs = append(g.Logs, log)
}

// AddActionLog 添加动作日志的便捷方法
func (g *GameLogRecord) AddActionLog(p1, p2 int, action LogAction, tile *Tile, callGroup []*Tile) {
	log := NewBaseGameLogWithAction(p1, p2, action, tile, callGroup)
	g.AddLog(log)
}

// AddScoreLog 添加分数日志的便捷方法
func (g *GameLogRecord) AddScoreLog(scores [NPlayers]int) {
	log := NewBaseGameLogWithScores(scores)
	g.AddLog(log)
}

// Clear 清空所有日志
func (g *GameLogRecord) Clear() {
	g.Logs = make([]*BaseGameLog, 0)
}

// GetLogCount 获取日志数量
func (g *GameLogRecord) GetLogCount() int {
	return len(g.Logs)
}

// GetLog 获取指定索引的日志
func (g *GameLogRecord) GetLog(index int) *BaseGameLog {
	if index >= 0 && index < len(g.Logs) {
		return g.Logs[index]
	}
	return nil
}

// String 返回所有日志的字符串表示
func (g *GameLogRecord) String() string {
	result := ""
	for i, log := range g.Logs {
		if i > 0 {
			result += "\n"
		}
		result += log.String()
	}
	return result
}
