package mahjong

import (
	"fmt"
	"strings"
)

// ResultType 表示游戏结果的类型
type ResultType int

const (
	// 胡牌结果
	RonAgari      ResultType = iota // 荣和
	TsumoAgari                      // 自摸
	NagashiMangan                   // 流局满贯

	// 流局类型
	RyukyokuNotile       // 牌山无牌（牌库流局）
	RyukyokuSuzukaflush  // 四风连打
	RyukyokuSufonrenda   // 四杠散了
	RyukyokuSuuchahan    // 四家立直
	RyukyokuKyushukyuhai // 九种九牌
	RyukyokuFourWind     // 四风连打（复合）

	// 其他
	NoResult // 游戏正常进行中
)

// GameResult 表示一局麻将的结果
type GameResult struct {
	Type         ResultType          // 结果类型
	WinnerIdx    int                 // 胜者索引（-1表示无胜者）
	LoserIdx     int                 // 失败者索引（-1表示无失败者）
	Score        *ScoreCounterResult // 计分结果
	ScoreChanges [4]int              // 四个玩家的分数变化
	Yakus        []Yaku              // 成立的役
	Fan          int                 // 番数
	Fu           int                 // 符数
	Message      string              // 结果信息

	// 用于重新导出结果
	FanCount   int  // 番数（对于非胡牌结果）
	IsYakuman  bool // 是否为役满
	ResultTime int  // 结果产生的时间戳
}

// NewGameResult 创建一个新的GameResult实例
func NewGameResult() *GameResult {
	return &GameResult{
		Type:         NoResult,
		WinnerIdx:    -1,
		LoserIdx:     -1,
		ScoreChanges: [4]int{0, 0, 0, 0},
		Yakus:        make([]Yaku, 0),
	}
}

// SetRonAgari 设置为荣和结果
func (r *GameResult) SetRonAgari(winnerIdx int, loserIdx int, score *ScoreCounterResult) {
	r.Type = RonAgari
	r.WinnerIdx = winnerIdx
	r.LoserIdx = loserIdx
	r.Score = score

	// 计算分数变化
	for i := 0; i < 4; i++ {
		r.ScoreChanges[i] = 0
	}
	r.ScoreChanges[winnerIdx] = score.RonScore
	r.ScoreChanges[loserIdx] = -score.RonScore

	// 复制役
	r.Yakus = make([]Yaku, len(score.Yakus))
	copy(r.Yakus, score.Yakus)
	r.Fan = score.Fan
	r.Fu = score.Fu

	r.Message = fmt.Sprintf("玩家%d荣和玩家%d的牌，%d番%d符，得分%d",
		winnerIdx, loserIdx, score.Fan, score.Fu, score.RonScore)
}

// SetTsumoAgari 设置为自摸结果
func (r *GameResult) SetTsumoAgari(winnerIdx int, score *ScoreCounterResult) {
	r.Type = TsumoAgari
	r.WinnerIdx = winnerIdx
	r.LoserIdx = -1
	r.Score = score

	// 计算分数变化
	r.ScoreChanges = [4]int{0, 0, 0, 0}

	if winnerIdx == 0 { // 庄家自摸
		r.ScoreChanges[1] = -score.TsumoScore[0]
		r.ScoreChanges[2] = -score.TsumoScore[1]
		r.ScoreChanges[3] = -score.TsumoScore[2]
		r.ScoreChanges[0] = score.TsumoScore[0]*3 - (r.ScoreChanges[1] + r.ScoreChanges[2] + r.ScoreChanges[3])
	} else { // 子家自摸
		r.ScoreChanges[0] = -score.TsumoScore[0]
		for i := 1; i < 4; i++ {
			if i == winnerIdx {
				continue
			}
			if i == 0 {
				r.ScoreChanges[i] = -score.TsumoScore[0]
			} else {
				r.ScoreChanges[i] = -score.TsumoScore[1]
			}
		}
		total := 0
		for i := 0; i < 4; i++ {
			if i != winnerIdx {
				total += (-r.ScoreChanges[i])
			}
		}
		r.ScoreChanges[winnerIdx] = total
	}

	r.Yakus = make([]Yaku, len(score.Yakus))
	copy(r.Yakus, score.Yakus)
	r.Fan = score.Fan
	r.Fu = score.Fu

	r.Message = fmt.Sprintf("玩家%d自摸，%d番%d符",
		winnerIdx, score.Fan, score.Fu)
}

// SetNagashiMangan 设置为流局满贯
func (r *GameResult) SetNagashiMangan(winnerIdx int) {
	r.Type = NagashiMangan
	r.WinnerIdx = winnerIdx
	r.LoserIdx = -1
	r.Fan = 5
	r.Fu = 30
	r.FanCount = 5

	// 计算满贯分数（8000点）
	for i := 0; i < 4; i++ {
		r.ScoreChanges[i] = 0
	}

	if winnerIdx == 0 { // 庄家
		r.ScoreChanges[0] = 4000
		r.ScoreChanges[1] = -4000
		r.ScoreChanges[2] = -4000
		r.ScoreChanges[3] = -4000
	} else { // 子家
		r.ScoreChanges[winnerIdx] = 8000
		r.ScoreChanges[0] = -4000
		for i := 1; i < 4; i++ {
			if i != winnerIdx {
				r.ScoreChanges[i] = -2000
			}
		}
	}

	r.Message = fmt.Sprintf("玩家%d流局满贯", winnerIdx)
}

// SetRyukyokuNotile 设置为牌库流局
func (r *GameResult) SetRyukyokuNotile() {
	r.Type = RyukyokuNotile
	r.WinnerIdx = -1
	r.LoserIdx = -1
	r.Message = "牌库流局（无牌可摸）"
}

// SetRyukyokuSuzukaflush 设置为四风连打
func (r *GameResult) SetRyukyokuSuzukaflush() {
	r.Type = RyukyokuSuzukaflush
	r.WinnerIdx = -1
	r.LoserIdx = -1
	r.Message = "四风连打流局"
}

// SetRyukyokuSufonrenda 设置为四杠散了
func (r *GameResult) SetRyukyokuSufonrenda() {
	r.Type = RyukyokuSufonrenda
	r.WinnerIdx = -1
	r.LoserIdx = -1
	r.Message = "四杠散了流局"
}

// SetRyukyokuSuuchahan 设置为四家立直
func (r *GameResult) SetRyukyokuSuuchahan() {
	r.Type = RyukyokuSuuchahan
	r.WinnerIdx = -1
	r.LoserIdx = -1
	r.Message = "四家立直流局"
}

// SetRyukyokuKyushukyuhai 设置为九种九牌
func (r *GameResult) SetRyukyokuKyushukyuhai(playerIdx int) {
	r.Type = RyukyokuKyushukyuhai
	r.WinnerIdx = playerIdx
	r.LoserIdx = -1
	r.Message = fmt.Sprintf("玩家%d九种九牌流局", playerIdx)
}

// ApplyScoreChanges 将分数变化应用到玩家
func (r *GameResult) ApplyScoreChanges(players [4]*Player) {
	for i := 0; i < 4; i++ {
		players[i].Score += r.ScoreChanges[i]
	}
}

// IsAgari 判断是否为胡牌结果
func (r *GameResult) IsAgari() bool {
	return r.Type == RonAgari || r.Type == TsumoAgari || r.Type == NagashiMangan
}

// IsRyukyoku 判断是否为流局
func (r *GameResult) IsRyukyoku() bool {
	return r.Type >= RyukyokuNotile
}

// String 返回结果的字符串表示
func (r *GameResult) String() string {
	sb := strings.Builder{}
	sb.WriteString(r.Message)
	sb.WriteString("\n")

	if r.IsAgari() {
		sb.WriteString(fmt.Sprintf("番数: %d, 符数: %d\n", r.Fan, r.Fu))
		sb.WriteString("役: ")
		for i, yaku := range r.Yakus {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(YakuToString(yaku))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("分数变化: ")
	for i := 0; i < 4; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		if r.ScoreChanges[i] > 0 {
			sb.WriteString(fmt.Sprintf("+%d", r.ScoreChanges[i]))
		} else {
			sb.WriteString(fmt.Sprintf("%d", r.ScoreChanges[i]))
		}
	}
	sb.WriteString("\n")

	return sb.String()
}

// GenerateResultNotile 生成牌库流局的结果
func GenerateResultNotile() *GameResult {
	result := NewGameResult()
	result.SetRyukyokuNotile()
	return result
}

// GenerateResultSuzukaflush 生成四风连打的结果
func GenerateResultSuzukaflush() *GameResult {
	result := NewGameResult()
	result.SetRyukyokuSuzukaflush()
	return result
}

// GenerateResultSufonrenda 生成四杠散了的结果
func GenerateResultSufonrenda() *GameResult {
	result := NewGameResult()
	result.SetRyukyokuSufonrenda()
	return result
}

// GenerateResultSuuchahan 生成四家立直的结果
func GenerateResultSuuchahan() *GameResult {
	result := NewGameResult()
	result.SetRyukyokuSuuchahan()
	return result
}

// GenerateResultKyushukyuhai 生成九种九牌的结果
func GenerateResultKyushukyuhai(playerIdx int) *GameResult {
	result := NewGameResult()
	result.SetRyukyokuKyushukyuhai(playerIdx)
	return result
}

// GenerateResultRon 生成荣和的结果
func GenerateResultRon(winnerIdx int, loserIdx int, score *ScoreCounterResult) *GameResult {
	result := NewGameResult()
	result.SetRonAgari(winnerIdx, loserIdx, score)
	return result
}

// GenerateResultTsumo 生成自摸的结果
func GenerateResultTsumo(winnerIdx int, score *ScoreCounterResult) *GameResult {
	result := NewGameResult()
	result.SetTsumoAgari(winnerIdx, score)
	return result
}

// GenerateResultNagashiMangan 生成流局满贯的结果
func GenerateResultNagashiMangan(playerIdx int) *GameResult {
	result := NewGameResult()
	result.SetNagashiMangan(playerIdx)
	return result
}
