package mahjong

import (
	"fmt"
	"sort"
	"strings"
)

// PaipuReplayer 用于重放和模拟麻将游戏
type PaipuReplayer struct {
	Table         *Table        // 游戏桌子
	Paipu         []BaseAction  // 行动日志
	SelectionLog  []int         // 选择日志
	ResultLog     []*GameResult // 结果日志
	CurrentRound  int           // 当前回合
	GameState     int           // 游戏状态（0=进行中，1=完成）
	DebugMode     bool          // 调试模式
	RoundStartIdx int           // 当前回合的起始索引
}

// NewPaipuReplayer 创建一个新的重放器
func NewPaipuReplayer(table *Table) *PaipuReplayer {
	return &PaipuReplayer{
		Table:        table,
		Paipu:        make([]BaseAction, 0),
		SelectionLog: make([]int, 0),
		ResultLog:    make([]*GameResult, 0),
		CurrentRound: 0,
		GameState:    0,
	}
}

// Init 初始化重放器
func (pr *PaipuReplayer) Init() {
	pr.Table.InitBeforePlaying()
	pr.Table.DrawTenhouStyle()
	pr.Table.SortPlayerHands()
	pr.CurrentRound = 0
	pr.GameState = 0
	pr.RoundStartIdx = 0
	pr.Paipu = make([]BaseAction, 0)
	pr.SelectionLog = make([]int, 0)
	pr.ResultLog = make([]*GameResult, 0)
}

// InitWithYama 用给定的牌山初始化
func (pr *PaipuReplayer) InitWithYama(yamaLog []int) {
	pr.Table.InitBeforePlaying()
	pr.Table.ImportYama(yamaLog)
	pr.Table.DrawTenhouStyle()
	pr.Table.SortPlayerHands()
	pr.CurrentRound = 0
	pr.GameState = 0
}

// GetSelfActions 获取当前玩家可以执行的自主行动
func (pr *PaipuReplayer) GetSelfActions(playerIdx int) []*SelfAction {
	player := pr.Table.Players[playerIdx]
	actions := make([]*SelfAction, 0)

	// 获取所有可能的行动
	actions = append(actions, player.GetDiscard(false)...)
	actions = append(actions, player.GetAnkan()...)
	actions = append(actions, player.GetKakan()...)

	if player.IsTenpai() {
		actions = append(actions, player.GetTsumo(pr.Table)...)
	}

	if player.IsMenzen() && !player.IsRiichi() {
		actions = append(actions, player.GetRiichi()...)
	}

	if pr.CurrentRound == 0 {
		actions = append(actions, player.GetKyushukyuhai()...)
	}

	return actions
}

// GetResponseActions 获取当前玩家对他人行动的反应行动
func (pr *PaipuReplayer) GetResponseActions(playerIdx int, action BaseAction, tile *Tile) []*ResponseAction {
	player := pr.Table.Players[playerIdx]
	actions := make([]*ResponseAction, 0)

	switch action {
	case Discard:
		// 可以选择吃、碰、杠、荣和
		actions = append(actions, player.GetChi(tile)...)
		ponActions := player.GetPon(tile)
		if ponActions != nil {
			actions = append(actions, ponActions...)
		}
		kanActions := player.GetKan(tile)
		if kanActions != nil {
			actions = append(actions, kanActions...)
		}
		ronActions := player.GetRon(pr.Table, tile)
		if ronActions != nil {
			actions = append(actions, ronActions...)
		}
	case AnKan:
		// 可以选择抢暗杠或荣和
		chanAnkanActions := player.GetChanAnkan(tile)
		if chanAnkanActions != nil {
			actions = append(actions, chanAnkanActions...)
		}
		chankanActions := player.GetChankan(tile)
		if chankanActions != nil {
			actions = append(actions, chankanActions...)
		}
	case Kan:
		// 可以选择抢杠或荣和
		chankanActions := player.GetChankan(tile)
		if chankanActions != nil {
			actions = append(actions, chankanActions...)
		}
		ronActions := player.GetRon(pr.Table, tile)
		if ronActions != nil {
			actions = append(actions, ronActions...)
		}
	}

	return actions
}

// MakeSelection 玩家选择一个行动
func (pr *PaipuReplayer) MakeSelection(playerIdx int, actionIdx int) bool {
	actions := pr.GetSelfActions(playerIdx)
	if actionIdx >= 0 && actionIdx < len(actions) {
		pr.SelectionLog = append(pr.SelectionLog, actionIdx)
		pr.ExecuteAction(playerIdx, &actions[actionIdx].Action)
		return true
	}
	return false
}

// MakeSelectionFromAction 直接选择一个具体的行动
func (pr *PaipuReplayer) MakeSelectionFromAction(playerIdx int, action *SelfAction) bool {
	pr.ExecuteAction(playerIdx, &action.Action)
	return true
}

// ExecuteAction 执行一个行动
func (pr *PaipuReplayer) ExecuteAction(playerIdx int, action *Action) {
	player := pr.Table.Players[playerIdx]

	switch action.GetAction() {
	case Discard:
		// 弃牌
		if len(action.GetCorrespondTiles()) > 0 {
			tile := action.GetCorrespondTiles()[0]
			// 添加到河中
			riverTile := RiverTile{
				Tile:     tile,
				Number:   pr.CurrentRound,
				Riichi:   player.IsRiichi(),
				Remain:   true,
				FromHand: true,
			}
			player.River.PushBack(riverTile)
			player.RemoveFromHand(tile)

			// 更新振听状态
			if player.IsRiichi() && player.IsTenpai() {
				if BaseTileInSlice(tile.Tile, player.AtariTiles) {
					player.Minogashi()
				}
			}

			// 记录弃牌日志
			if pr.Table != nil && pr.Table.GameLog != nil {
				actionLog := LogDiscardFromHand
				if player.IsRiichi() {
					actionLog = LogRiichiDiscardFromHand
				}
				pr.Table.GameLog.AddActionLog(playerIdx, -1, actionLog, tile, nil)
			}
		}

	case AnKan, Kan:
		// 杠牌
		if len(action.GetCorrespondTiles()) > 0 {
			tile := action.GetCorrespondTiles()[0]
			if action.GetAction() == AnKan {
				player.ExecuteAnkan(tile.Tile)
			} else {
				player.ExecuteKakan(tile)
			}
			pr.Table.NewDora()
			// 摸一张牌
			pr.Table.DrawRinshan(playerIdx)

			// 记录杠动作日志
			if pr.Table != nil && pr.Table.GameLog != nil {
				if action.GetAction() == AnKan {
					pr.Table.GameLog.AddActionLog(playerIdx, -1, LogAnKan, nil, nil)
				} else {
					pr.Table.GameLog.AddActionLog(playerIdx, -1, LogKan, nil, nil)
				}
			}
		}

	case Tsumo:
		// 自摸胡牌
		if player.IsTenpai() {
			counter := &ScoreCounter{}
			baseTiles := ConvertTilesToBaseTiles(player.Hand)
			isSevenPair := IsSevenPairPattern(baseTiles)
			result := counter.CalculateScore(pr.Table, player, baseTiles, player.CallGroups, baseTiles[len(baseTiles)-1], isSevenPair)
			if result != nil {
				gameResult := GenerateResultTsumo(playerIdx, result)
				pr.ResultLog = append(pr.ResultLog, gameResult)
				var players [4]*Player
				for i := 0; i < 4; i++ {
					players[i] = pr.Table.Players[i]
				}
				gameResult.ApplyScoreChanges(players)

				// 记录自摸与分数日志
				if pr.Table != nil && pr.Table.GameLog != nil {
					// 尝试取胜利牌（手牌最后一张）
					var winTile *Tile = nil
					if len(player.Hand) > 0 {
						winTile = player.Hand[len(player.Hand)-1]
					}
					pr.Table.GameLog.AddActionLog(playerIdx, -1, LogTsumo, winTile, nil)
					// 添加分数快照
					scores := [NPlayers]int{}
					for i := 0; i < NPlayers; i++ {
						scores[i] = pr.Table.Players[i].Score
					}
					pr.Table.GameLog.AddScoreLog(scores)
				}
				pr.GameState = 1
			}
		}

	case Ron:
		// 荣和
		// 这需要在响应行动中处理
		break

	case Riichi:
		// 立直
		player.Riichi = true
		player.Ippatsu = true
		break

	case Kyushukyuhai:
		// 九种九牌
		gameResult := GenerateResultKyushukyuhai(playerIdx)
		pr.ResultLog = append(pr.ResultLog, gameResult)
		pr.GameState = 1
		// 记录九种九牌日志
		if pr.Table != nil && pr.Table.GameLog != nil {
			pr.Table.GameLog.AddActionLog(playerIdx, -1, LogKyushukyuhai, nil, nil)
		}
		break
	}

	pr.Paipu = append(pr.Paipu, action.GetAction())
}

// SimulateToCompletion 模拟游戏直到完成
func (pr *PaipuReplayer) SimulateToCompletion() *GameResult {
	for pr.GameState == 0 && len(pr.Table.Yama) > 14 {
		currentPlayer := pr.Table.Players[pr.Table.Turn]

		// 摸一张牌（如果不是初始状态）
		if pr.CurrentRound > 0 {
			pr.Table.DrawNormal(pr.Table.Turn)
			currentPlayer.UpdateAtariTiles()
		}

		pr.CurrentRound++

		// 获取可能的行动
		selfActions := pr.GetSelfActions(pr.Table.Turn)
		if len(selfActions) > 0 {
			// 选择第一个行动（弃牌）
			selectedAction := selfActions[0]
			pr.ExecuteAction(pr.Table.Turn, &selectedAction.Action)
		}

		// 推进回合
		pr.Table.NextTurn((pr.Table.Turn + 1) % 4)
	}

	// 如果游戏还没有结束，产生流局
	if pr.GameState == 0 {
		gameResult := GenerateResultNotile()
		pr.ResultLog = append(pr.ResultLog, gameResult)
		pr.GameState = 1
		return gameResult
	}

	if len(pr.ResultLog) > 0 {
		return pr.ResultLog[len(pr.ResultLog)-1]
	}

	return nil
}

// PrintGameState 打印当前游戏状态
func (pr *PaipuReplayer) PrintGameState() {
	fmt.Println(pr.Table.String())
	fmt.Printf("当前回合: %d\n", pr.CurrentRound)
	fmt.Printf("游戏状态: %d\n", pr.GameState)
}

// GetGameLog 获取游戏日志
func (pr *PaipuReplayer) GetGameLog() string {
	sb := strings.Builder{}
	sb.WriteString("=== 麻将游戏日志 ===\n")
	sb.WriteString(fmt.Sprintf("总回合数: %d\n", pr.CurrentRound))
	sb.WriteString(fmt.Sprintf("行动数: %d\n", len(pr.Paipu)))
	sb.WriteString("\n== 行动序列 ==\n")

	for i, action := range pr.Paipu {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i, BaseActionToString(action)))
	}

	sb.WriteString("\n== 结果 ==\n")
	for i, result := range pr.ResultLog {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i, result.String()))
	}

	sb.WriteString("\n== 最终分数 ==\n")
	for i := 0; i < 4; i++ {
		sb.WriteString(fmt.Sprintf("玩家 %d: %d\n", i, pr.Table.Players[i].Score))
	}

	return sb.String()
}

// ExportYamaLog 导出牌山日志
func (pr *PaipuReplayer) ExportYamaLog() string {
	return pr.Table.ExportYama()
}

// GetRoundInfo 获取指定回合的信息
func (pr *PaipuReplayer) GetRoundInfo(roundNum int) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("=== 回合 %d ===\n", roundNum))

	// 遍历该回合的行动
	for i, action := range pr.Paipu {
		roundIndex := i / 4
		playerIdx := i % 4
		if roundIndex == roundNum {
			sb.WriteString(fmt.Sprintf("玩家 %d: %s\n", playerIdx, BaseActionToString(action)))
		}
	}

	return sb.String()
}

// 辅助函数

// ConvertTilesToBaseTiles 将Tile数组转换为BaseTile数组
func ConvertTilesToBaseTiles(tiles []*Tile) []BaseTile {
	result := make([]BaseTile, len(tiles))
	for i, tile := range tiles {
		result[i] = tile.Tile
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}

// BaseTileInSlice 检查BaseTile是否在切片中
func BaseTileInSlice(tile BaseTile, slice []BaseTile) bool {
	for _, t := range slice {
		if t == tile {
			return true
		}
	}
	return false
}

// IsSevenPairPattern 判断是否为七对子形式
func IsSevenPairPattern(tiles []BaseTile) bool {
	if len(tiles) != 14 {
		return false
	}
	// 检查是否为7对对子
	counts := make(map[BaseTile]int)
	for _, tile := range tiles {
		counts[tile]++
	}
	pairCount := 0
	for _, count := range counts {
		if count == 2 {
			pairCount++
		} else if count != 0 {
			return false
		}
	}
	return pairCount == 7
}

// FindInTiles 在Tile数组中查找第一个匹配的索引
func FindInTiles(tiles []*Tile, target *Tile) int {
	for i, tile := range tiles {
		if tile.ID == target.ID {
			return i
		}
	}
	return -1
}

// GetNCopies 获取n张相同的牌
func GetNCopies(tiles []*Tile, baseTile BaseTile, n int) []*Tile {
	result := make([]*Tile, 0, n)
	for _, tile := range tiles {
		if tile.Tile == baseTile && len(result) < n {
			result = append(result, tile)
		}
	}
	return result
}
