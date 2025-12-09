package mahjong

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	NBaseTiles = 34             // 基础牌种类数
	NTiles     = NBaseTiles * 4 // 总牌数(136张)
	NPlayers   = 4              // 玩家数
)

// Table 表示麻将桌子，管理整个游戏的状态
type Table struct {
	// 牌和宝牌相关
	Tiles            [NTiles]*Tile // 所有牌的数组
	NActiveDora      int           // 翻开的宝牌指示牌数量
	DoraIndicator    []*Tile       // 宝牌指示牌
	UraDoraIndicator []*Tile       // 里宝牌指示牌
	Yama             []*Tile       // 牌山（剩余的牌）

	// 玩家和游戏状态
	Players    [NPlayers]*Player // 四个玩家
	Turn       int               // 当前玩家的回合
	LastAction BaseAction        // 上一个行动
	GameWind   Wind              // 场风（东、南、西、北）
	Oya        int               // 庄家索引
	Honba      int               // 本场数
	Kyoutaku   int               // 供托数

	// 随机数生成器
	Rand    *rand.Rand // 随机数生成器
	UseSeed bool       // 是否使用种子
	Seed    int64      // 随机种子

	// 调试信息
	YamaLog      []int          // 牌山日志
	SelectionLog []int          // 选择日志
	GameLog      *GameLogRecord // 游戏日志记录器
	LastActor    int            // 上一个执行动作的玩家索引（用于响应阶段）
}

// NewTable 创建一个新的Table实例
func NewTable() *Table {
	table := &Table{
		Turn:         0,
		LastAction:   Pass,
		GameWind:     East,
		Oya:          0,
		Honba:        0,
		Kyoutaku:     0,
		Rand:         rand.New(rand.NewSource(time.Now().UnixNano())),
		NActiveDora:  1,
		YamaLog:      make([]int, 0),
		SelectionLog: make([]int, 0),
		GameLog:      NewGameLogRecord(),
		LastActor:    -1,
	}

	// 初始化玩家
	for i := 0; i < NPlayers; i++ {
		table.Players[i] = NewPlayer(Wind(i), i == 0)
	}

	return table
}

// NewDora 翻出新的宝牌
func (t *Table) NewDora() {
	t.NActiveDora++
}

// GetDora 获取所有已翻开的宝牌
func (t *Table) GetDora() []BaseTile {
	doras := make([]BaseTile, 0, t.NActiveDora)
	for i := 0; i < t.NActiveDora && i < len(t.DoraIndicator); i++ {
		doras = append(doras, GetDoraNext(t.DoraIndicator[i].Tile))
	}
	return doras
}

// GetUraDora 获取所有里宝牌
func (t *Table) GetUraDora() []BaseTile {
	uras := make([]BaseTile, 0, t.NActiveDora)
	for i := 0; i < t.NActiveDora && i < len(t.UraDoraIndicator); i++ {
		uras = append(uras, GetDoraNext(t.UraDoraIndicator[i].Tile))
	}
	return uras
}

// GetRemainKanTile 获取剩余的杠牌数
func (t *Table) GetRemainKanTile() int {
	if len(t.DoraIndicator) == 0 {
		return 0
	}
	for i, tile := range t.Yama {
		if tile.ID == t.DoraIndicator[0].ID {
			return i - 1
		}
	}
	return 0
}

// GetRemainTile 获取剩余的牌数
func (t *Table) GetRemainTile() int {
	return len(t.Yama) - 14
}

// InitTiles 初始化所有牌
func (t *Table) InitTiles() {
	tileID := 0
	for baseTile := BaseTile(0); baseTile < NBaseTiles; baseTile++ {
		for i := 0; i < 4; i++ {
			t.Tiles[tileID] = &Tile{
				Tile:    baseTile,
				RedDora: false,
				ID:      tileID,
			}
			tileID++
		}
	}
}

// InitRedDora3 初始化3张赤宝牌
func (t *Table) InitRedDora3() {
	for color := 0; color < 3; color++ {
		baseTile := BaseTile(_5m + BaseTile(color*9))
		idx := t.Rand.Intn(4)
		t.Tiles[int(baseTile)*4+idx].RedDora = true
	}
}

// ShuffleTiles 洗牌
func (t *Table) ShuffleTiles() {
	tiles := make([]*Tile, NTiles)
	copy(tiles, t.Tiles[:])

	for i := NTiles - 1; i > 0; i-- {
		j := t.Rand.Intn(i + 1)
		tiles[i], tiles[j] = tiles[j], tiles[i]
	}

	copy(t.Tiles[:], tiles)
}

// InitYama 初始化牌山
func (t *Table) InitYama() {
	t.Yama = make([]*Tile, 0, NTiles)

	for i := 0; i < NTiles; i++ {
		t.Yama = append(t.Yama, t.Tiles[i])
	}
}

// InitDora 初始化宝牌指示牌 - 初始化5组宝牌和里宝牌
func (t *Table) InitDora() {
	t.DoraIndicator = make([]*Tile, 0, 5)
	t.UraDoraIndicator = make([]*Tile, 0, 5)

	// 宝牌指示牌位置: Yama[5], Yama[7], Yama[9], Yama[11], Yama[13]
	// 里宝牌指示牌位置: Yama[4], Yama[6], Yama[8], Yama[10], Yama[12]
	doraPositions := []int{5, 7, 9, 11, 13}
	uraPositions := []int{4, 6, 8, 10, 12}

	for _, pos := range doraPositions {
		if len(t.Yama) > pos {
			t.DoraIndicator = append(t.DoraIndicator, t.Yama[len(t.Yama)-pos-1])
		}
	}

	for _, pos := range uraPositions {
		if len(t.Yama) > pos {
			t.UraDoraIndicator = append(t.UraDoraIndicator, t.Yama[len(t.Yama)-pos-1])
		}
	}

	t.NActiveDora = 1 // 初始只翻1张
}

// InitBeforePlaying 在开始游戏前的初始化
func (t *Table) InitBeforePlaying() {
	t.InitTiles()
	t.InitRedDora3()
	t.ShuffleTiles()
	t.InitYama()
	t.InitDora()
}

// ImportYama 导入预定义的牌山（用于重放）
func (t *Table) ImportYama(yamaLog []int) {
	t.Yama = make([]*Tile, 0, len(yamaLog))
	for _, idx := range yamaLog {
		if idx >= 0 && idx < NTiles {
			t.Yama = append(t.Yama, t.Tiles[idx])
		}
	}
}

// ExportYama 导出牌山（用于日志记录）
func (t *Table) ExportYama() string {
	var export string
	for _, tile := range t.Yama {
		export += fmt.Sprintf("%d,", tile.ID)
	}
	return export
}

// SetSeed 设置随机种子
func (t *Table) SetSeed(seed int64) {
	t.Seed = seed
	t.UseSeed = true
	t.Rand = rand.New(rand.NewSource(seed))
}

// DrawTenhouStyle 按照天凤风格摸牌
func (t *Table) DrawTenhouStyle() {
	for round := 0; round < 3; round++ {
		for i := 0; i < NPlayers; i++ {
			if len(t.Yama) > 0 {
				tile := t.Yama[len(t.Yama)-1]
				t.Yama = t.Yama[:len(t.Yama)-1]
				idx := (t.Oya + i) % NPlayers
				t.Players[idx].Hand = append(t.Players[idx].Hand, tile)
				t.YamaLog = append(t.YamaLog, tile.ID)
			}
		}
	}
	if len(t.Yama) > 0 {
		tile := t.Yama[len(t.Yama)-1]
		t.Yama = t.Yama[:len(t.Yama)-1]
		t.Players[t.Oya].Hand = append(t.Players[t.Oya].Hand, tile)
		t.YamaLog = append(t.YamaLog, tile.ID)
	}
}

// DrawNormal 按照标准方式摸牌
func (t *Table) DrawNormal(playerIndex int) {
	if len(t.Yama) > 0 {
		tile := t.Yama[len(t.Yama)-1]
		t.Yama = t.Yama[:len(t.Yama)-1]
		t.Players[playerIndex].Hand = append(t.Players[playerIndex].Hand, tile)
		t.YamaLog = append(t.YamaLog, tile.ID)
		if t.GameLog != nil {
			t.GameLog.AddActionLog(playerIndex, -1, LogDrawNormal, tile, nil)
		}
	}
}

// DrawNormalNoRecord 摸牌但不记录（用于初始化）
func (t *Table) DrawNormalNoRecord(playerIndex int) {
	if len(t.Yama) > 0 {
		tile := t.Yama[len(t.Yama)-1]
		t.Yama = t.Yama[:len(t.Yama)-1]
		t.Players[playerIndex].Hand = append(t.Players[playerIndex].Hand, tile)
	}
}

// DrawNNormal 摸n张牌，不记录
func (t *Table) DrawNNormal(playerIndex int, nTiles int) {
	for i := 0; i < nTiles; i++ {
		t.DrawNormalNoRecord(playerIndex)
	}
}

// DrawRinshan 从岭上摸牌
func (t *Table) DrawRinshan(playerIndex int) {
	if len(t.Yama) > 14 {
		tile := t.Yama[0]
		t.Yama = t.Yama[1:]
		t.Players[playerIndex].Hand = append(t.Players[playerIndex].Hand, tile)
		if t.GameLog != nil {
			t.GameLog.AddActionLog(playerIndex, -1, LogDrawRinshan, tile, nil)
		}
	}
}

// NextTurn 推进到下一个回合
func (t *Table) NextTurn(nextTurn int) {
	t.Turn = nextTurn
}

// SetDebugMode 设置调试模式
// mode: 0 = no debug, 1 = debug by buffer, 2 = debug by stdout
func (t *Table) SetDebugMode(mode int) {
	if mode < 0 {
		mode = 0
	}
	if mode > 2 {
		mode = 2
	}
	// mode: 0 = no debug, 1 = debug by buffer, 2 = debug by stdout
	switch mode {
	case 0:
		// 关闭日志记录
		t.GameLog = nil
	case 1:
		// 使用内存缓冲记录日志
		if t.GameLog == nil {
			t.GameLog = NewGameLogRecord()
		}
	case 2:
		// 也使用内存记录（可扩展为同时输出到 stdout）
		if t.GameLog == nil {
			t.GameLog = NewGameLogRecord()
		}
		// TODO: 将日志输出到 stdout/文件（当前先保留内存记录）
	}
}

// SortPlayerHands 对所有玩家的手牌进行排序
func (t *Table) SortPlayerHands() {
	for i := 0; i < NPlayers; i++ {
		t.Players[i].SortHand()
	}
}

// GetPhase 获取当前游戏阶段
// 返回值与 C++ PhaseEnum 对应: 0-3为自动作, 4-7为响应, 8-11为抢杠响应, 12-15为抢暗杠响应, 16为游戏结束
func (t *Table) GetPhase() int {
	// If game over
	if t.GetRemainTile() <= 0 {
		return 16 // 游戏结束
	}

	// If the last action was a discard, enter response phase for next players
	if t.LastAction == Discard {
		// response phases 4-7, map to next player (turn+1 ..)
		next := (t.Turn + 1) % NPlayers
		return 4 + next
	}

	// If last action was a kan, some players may have chankan opportunity
	if t.LastAction == Kan || t.LastAction == AnKan || t.LastAction == KaKan {
		next := (t.Turn + 1) % NPlayers
		return 8 + next
	}

	// Default: autonomous action phase for current player (0-3)
	return t.Turn
}

// MakeSelection 根据选择执行游戏逻辑
func (t *Table) MakeSelection(selection int) bool {
	phase := t.GetPhase()

	// helper: action priority map (higher is higher priority)
	actionPriority := func(a BaseAction) int {
		switch a {
		case Ron:
			return 4
		case Kan:
			return 3
		case Pon:
			return 2
		case Chi:
			return 1
		default:
			return 0
		}
	}

	// 自动作阶段（0-3）
	if phase >= 0 && phase <= 3 {
		playerIdx := t.Turn
		player := t.Players[playerIdx]
		if player == nil {
			return false
		}

		// 构造可能的自主行动（与 PaipuReplayer.GetSelfActions 保持一致）
		actions := make([]*SelfAction, 0)
		actions = append(actions, player.GetDiscard(false)...)
		actions = append(actions, player.GetAnkan()...)
		actions = append(actions, player.GetKakan()...)

		if player.IsTenpai() {
			actions = append(actions, player.GetTsumo(t)...)
		}

		if player.IsMenzen() && !player.IsRiichi() {
			actions = append(actions, player.GetRiichi()...)
		}

		if player.FirstRound {
			actions = append(actions, player.GetKyushukyuhai()...)
		}

		if selection < 0 || selection >= len(actions) {
			return false
		}

		act := actions[selection]

		// record actor
		t.LastActor = playerIdx

		switch act.GetAction() {
		case Discard:
			tile := act.GetCorrespondTiles()[0]
			// 添加到河中并从手中移除
			rt := RiverTile{Tile: tile, Number: 0, Riichi: player.IsRiichi(), Remain: true, FromHand: true}
			player.River.PushBack(rt)
			player.RemoveFromHand(tile)

			// 更新振听状态
			if player.IsRiichi() && player.IsTenpai() {
				if BaseTileInSlice(tile.Tile, player.AtariTiles) {
					player.Minogashi()
				}
			}

			// 记录弃牌日志
			if t.GameLog != nil {
				actionLog := LogDiscardFromHand
				if player.IsRiichi() {
					actionLog = LogRiichiDiscardFromHand
				}
				t.GameLog.AddActionLog(playerIdx, -1, actionLog, tile, nil)
			}

			t.LastAction = Discard
			// 不立即推进回合，等待响应阶段处理（由上层控制器或接下来的 MakeSelection 响应阶段触发）
			return true

		case AnKan:
			// 暗杠（使用第一个对应牌的基牌）
			tiles := act.GetCorrespondTiles()
			if len(tiles) == 0 {
				return false
			}
			base := tiles[0].Tile
			player.ExecuteAnkan(base)
			t.NewDora()
			t.DrawRinshan(playerIdx)
			if t.GameLog != nil {
				t.GameLog.AddActionLog(playerIdx, -1, LogAnKan, nil, nil)
			}
			t.LastAction = AnKan
			return true

		case KaKan:
			tiles := act.GetCorrespondTiles()
			if len(tiles) == 0 {
				return false
			}
			tile := tiles[0]
			player.ExecuteKakan(tile)
			t.NewDora()
			t.DrawRinshan(playerIdx)
			if t.GameLog != nil {
				t.GameLog.AddActionLog(playerIdx, -1, LogKan, nil, nil)
			}
			t.LastAction = KaKan
			return true

		case Tsumo:
			if player.IsTenpai() {
				counter := &ScoreCounter{}
				baseTiles := ConvertTilesToBaseTiles(player.Hand)
				isSevenPair := IsSevenPairPattern(baseTiles)
				result := counter.CalculateScore(t, player, baseTiles, player.CallGroups, baseTiles[len(baseTiles)-1], isSevenPair)
				if result != nil {
					gameResult := GenerateResultTsumo(playerIdx, result)
					// 应用分数变化
					var playersArr [4]*Player
					for i := 0; i < 4; i++ {
						playersArr[i] = t.Players[i]
					}
					gameResult.ApplyScoreChanges(playersArr)

					if t.GameLog != nil {
						var winTile *Tile = nil
						if len(player.Hand) > 0 {
							winTile = player.Hand[len(player.Hand)-1]
						}
						t.GameLog.AddActionLog(playerIdx, -1, LogTsumo, winTile, nil)
						scores := [NPlayers]int{}
						for i := 0; i < NPlayers; i++ {
							scores[i] = t.Players[i].Score
						}
						t.GameLog.AddScoreLog(scores)
					}
				}
			}
			return true

		case Riichi:
			player.Riichi = true
			player.Ippatsu = true
			if t.GameLog != nil {
				t.GameLog.AddActionLog(playerIdx, -1, LogRiichiSuccess, nil, nil)
			}
			return true

		case Kyushukyuhai:
			gameResult := GenerateResultKyushukyuhai(playerIdx)
			var playersArr [4]*Player
			for i := 0; i < 4; i++ {
				playersArr[i] = t.Players[i]
			}
			gameResult.ApplyScoreChanges(playersArr)
			if t.GameLog != nil {
				t.GameLog.AddActionLog(playerIdx, -1, LogKyushukyuhai, nil, nil)
			}
			return true
		default:
			return false
		}
	}

	// 响应阶段处理: 处理他人弃牌后的响应（吃/碰/杠/荣和/抢杠）
	// 处理范围：4-7 为常规响应（对弃牌），8-11 为抢杠(chankan)，12-15 为抢暗杠(chan ankan)
	if phase >= 4 && phase <= 15 {
		if t.LastActor < 0 || t.LastActor >= NPlayers {
			return false
		}

		// 获取被弃的牌（LastActor 最后丢出的牌）
		la := t.Players[t.LastActor]
		if la == nil || la.River.Size() == 0 {
			return false
		}
		discarded := la.River.River[la.River.Size()-1].Tile

		// 收集每个玩家可用的响应，并选出优先级最高的动作
		type resp struct {
			player int
			action BaseAction
		}
		responses := make([]resp, 0)

		// 按座位顺序从弃牌者下一家开始收集
		for d := 1; d <= 3; d++ {
			idx := (t.LastActor + d) % NPlayers
			p := t.Players[idx]
			if p == nil {
				continue
			}

			// 生成响应动作列表（针对弃牌）
			candidateActions := make([]BaseAction, 0)
			// 荣和
			ron := p.GetRon(t, discarded)
			if ron != nil && len(ron) > 0 {
				candidateActions = append(candidateActions, Ron)
			}
			// 抢杠/抢暗杠相关会在chankan/chanankan阶段处理（这里忽略）
			// 杠/碰/吃
			kan := p.GetKan(discarded)
			if kan != nil && len(kan) > 0 {
				candidateActions = append(candidateActions, Kan)
			}
			pon := p.GetPon(discarded)
			if pon != nil && len(pon) > 0 {
				candidateActions = append(candidateActions, Pon)
			}
			chi := p.GetChi(discarded)
			if chi != nil && len(chi) > 0 {
				candidateActions = append(candidateActions, Chi)
			}

			// 选择该玩家优先级最高的响应
			best := Pass
			bestPri := 0
			for _, a := range candidateActions {
				pri := actionPriority(a)
				if pri > bestPri {
					best = a
					bestPri = pri
				}
			}
			responses = append(responses, resp{player: idx, action: best})
		}

		// 计算最终最高优先级
		finalAction := Pass
		finalPri := 0
		for _, r := range responses {
			pri := actionPriority(r.action)
			if pri > finalPri {
				finalPri = pri
				finalAction = r.action
			}
		}

		if finalAction == Pass {
			// 所有人都放弃 -> 将弃牌标记为不可取，并推进到下一家摸牌
			la.River.SetNotRemain()
			// 清除 last action 后推进回合
			t.LastAction = Pass
			t.NextTurn((t.LastActor + 1) % NPlayers)
			return true
		}

		// 在有最终动作的情况下，找到优先的玩家（同优先级则靠近弃牌者的玩家优先）
		winner := -1
		for d := 1; d <= 3; d++ {
			idx := (t.LastActor + d) % NPlayers
			for _, r := range responses {
				if r.player == idx && r.action == finalAction {
					winner = idx
					break
				}
			}
			if winner >= 0 {
				break
			}
		}
		if winner < 0 {
			return false
		}

		// 执行最终响应
		switch finalAction {
		case Ron:
			// 胡牌：winner 胡 LastActor 的弃牌
			winnerP := t.Players[winner]
			if winnerP == nil {
				return false
			}
			// 构造正确的牌组用于计分
			baseTiles := ConvertTilesToBaseTiles(winnerP.Hand)
			// 将弃牌作为最后一张
			baseTiles = append(baseTiles, discarded.Tile)
			isSeven := IsSevenPairPattern(baseTiles)
			counter := &ScoreCounter{}
			result := counter.CalculateScore(t, winnerP, baseTiles, winnerP.CallGroups, discarded.Tile, isSeven)
			var gameResult *GameResult
			if result != nil {
				gameResult = GenerateResultRon(winner, t.LastActor, result)
				var playersArr [4]*Player
				for i := 0; i < 4; i++ {
					playersArr[i] = t.Players[i]
				}
				gameResult.ApplyScoreChanges(playersArr)
			} else {
				// 没有役（异常情况），返回 false
				return false
			}
			// 记录日志
			if t.GameLog != nil {
				t.GameLog.AddActionLog(winner, t.LastActor, LogRon, discarded, nil)
				scores := [NPlayers]int{}
				for i := 0; i < NPlayers; i++ {
					scores[i] = t.Players[i].Score
				}
				t.GameLog.AddScoreLog(scores)
			}
			t.LastAction = Ron
			return true

		case Chi, Pon, Kan:
			// 执行鸣牌（吃/碰/大明杠）
			caller := t.Players[winner]
			if caller == nil {
				return false
			}
			// 从被弃牌家的河中删除该牌（标记为已取走）
			la.River.SetNotRemain()
			// 执行鸣牌
			caller.ExecuteNaki(discarded, finalAction)
			// 记录日志
			if t.GameLog != nil {
				var logAction LogAction
				switch finalAction {
				case Chi:
					logAction = LogChi
				case Pon:
					logAction = LogPon
				case Kan:
					logAction = LogKan
				}
				t.GameLog.AddActionLog(winner, t.LastActor, logAction, discarded, nil)
			}
			// 消除第一巡与一发
			for i := 0; i < NPlayers; i++ {
				if t.Players[i] != nil {
					t.Players[i].FirstRound = false
					t.Players[i].Ippatsu = false
				}
			}
			// 设置当前回合为鸣牌者
			t.Turn = winner
			t.LastAction = finalAction
			return true
		}
	}

	return false
}

// FromBeginning 游戏主循环的开始
// 处理流局判定、摸牌、生成行动列表等
func (t *Table) FromBeginning() {
	// 检查游戏是否已结束
	// TODO: 添加阶段检查

	// 检查四风连打（四个玩家都弃了相同的字牌）
	if t.isAbaortedFourWind() {
		// 游戏结束
		return
	}

	// 检查四立直
	if t.isFourRiichi() {
		// 游戏结束
		return
	}

	// 检查四杠散了
	if t.isFourKanAborted() {
		// 游戏结束
		return
	}

	// 检查是否没有牌了（流局）
	if t.GetRemainTile() == 0 {
		// 游戏结束
		return
	}

	// 整理所有玩家的手牌
	t.SortPlayerHands()

	// 根据上一个行动决定是否摸牌
	if t.LastAction == Kan || t.LastAction == AnKan || t.LastAction == KaKan {
		// 杠后从岭上摸牌
		t.DrawRinshan(t.Turn)
	} else if t.LastAction != Chi && t.LastAction != Pon {
		// 吃碰后不摸牌，其他时候摸牌
		t.DrawNormal(t.Turn)
	}

	// 更新听牌
	if t.Players[t.Turn] != nil {
		t.Players[t.Turn].UpdateAtariTiles()
	}

	// 生成本回合的可能行动列表（基础实现，供外部调用者查询）
	// 包括弃牌、自摸、暗杠、加杠、立直、九种九牌
	// 具体选择由调用方（如 PaipuReplayer 或 UI）执行 MakeSelection
	// 这里只做简单的可行性检查并记录（SelectionLog用于调试）
	player := t.Players[t.Turn]
	if player != nil {
		// 清理之前的选择日志
		t.SelectionLog = t.SelectionLog[:0]
		// 如果有暗杠选项，则记下
		ankan := player.GetAnkan()
		if len(ankan) > 0 {
			t.SelectionLog = append(t.SelectionLog, int(AnKan))
		}
		// 加杠
		kakan := player.GetKakan()
		if len(kakan) > 0 {
			t.SelectionLog = append(t.SelectionLog, int(KaKan))
		}
		// 自摸
		if player.IsTenpai() {
			t.SelectionLog = append(t.SelectionLog, int(Tsumo))
		}
		// 立直
		if player.IsMenzen() && !player.IsRiichi() {
			t.SelectionLog = append(t.SelectionLog, int(Riichi))
		}
		// 九种九牌
		if player.FirstRound {
			ky := player.GetKyushukyuhai()
			if len(ky) > 0 {
				t.SelectionLog = append(t.SelectionLog, int(Kyushukyuhai))
			}
		}
	}
}

// 辅助方法

// isAbaortedFourWind 判断是否为四风连打流局
func (t *Table) isAbaortedFourWind() bool {
	if len(t.Players[0].River.River) == 1 &&
		len(t.Players[1].River.River) == 1 &&
		len(t.Players[2].River.River) == 1 &&
		len(t.Players[3].River.River) == 1 &&
		len(t.Players[0].CallGroups) == 0 &&
		len(t.Players[1].CallGroups) == 0 &&
		len(t.Players[2].CallGroups) == 0 &&
		len(t.Players[3].CallGroups) == 0 {

		// 检查四个玩家的第一张弃牌是否都是相同的字牌
		t0 := t.Players[0].River.River[0].Tile.Tile
		t1 := t.Players[1].River.River[0].Tile.Tile
		t2 := t.Players[2].River.River[0].Tile.Tile
		t3 := t.Players[3].River.River[0].Tile.Tile

		return t0 == t1 && t0 == t2 && t0 == t3 &&
			t0 >= _1z && t0 <= _4z
	}
	return false
}

// isFourRiichi 判断是否为四立直
func (t *Table) isFourRiichi() bool {
	richiCount := 0
	for i := 0; i < NPlayers; i++ {
		if t.Players[i].IsRiichi() {
			richiCount++
		}
	}
	return richiCount == NPlayers
}

// isFourKanAborted 判断是否为四杠散了
func (t *Table) isFourKanAborted() bool {
	if t.GetRemainKanTile() != 0 {
		return false
	}

	kanCount := 0
	for i := 0; i < NPlayers; i++ {
		hasKan := false
		for _, group := range t.Players[i].CallGroups {
			if group.Type == Kantsu {
				hasKan = true
				break
			}
		}
		if hasKan {
			kanCount++
		}
	}

	// 2个或以上玩家杠过
	return kanCount >= 2
}

// String 返回Table的字符串表示
func (t *Table) String() string {
	str := fmt.Sprintf("场风: %d, 庄家: %d, 本场: %d, 供托: %d\n",
		t.GameWind, t.Oya, t.Honba, t.Kyoutaku)
	str += fmt.Sprintf("牌山剩余: %d\n", len(t.Yama))
	for i := 0; i < NPlayers; i++ {
		str += fmt.Sprintf("玩家 %d: 点数=%d, 手牌=%s\n",
			i, t.Players[i].Score, t.Players[i].HandToString())
	}
	return str
}
