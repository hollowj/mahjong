package mahjong

import (
	"fmt"
)

// ScoreCounterResult 表示计分的结果
type ScoreCounterResult struct {
	Fan        int    // 番数
	Fu         int    // 符数
	BaseScore  int    // 基础分
	TsumoScore [3]int // 自摸时的分数（非庄、非庄、非庄）
	RonScore   int    // 荣和时的分数
	IsYakuman  bool   // 是否为役满
	Yakus      []Yaku // 所有成立的役
}

// ScoreCounter 是麻将计分器
type ScoreCounter struct {
	Player      *Player     // 玩家
	Tiles       []BaseTile  // 牌（14张）
	CallGroups  []CallGroup // 鸣牌组
	WinTile     BaseTile    // 胡牌（第14张）
	IsSevenPair bool        // 是否为七对子形式
	Table       *Table      // 游戏桌（用于场风等信息）
}

// CalculateScore 计算分数
func (s *ScoreCounter) CalculateScore(
	table *Table,
	player *Player,
	tiles []BaseTile,
	callGroups []CallGroup,
	winTile BaseTile,
	isSevenPair bool,
) *ScoreCounterResult {
	// 将评分流程改为：枚举所有拆分 -> 对每个拆分枚举和牌所在的位置（含顺子的三种位置）
	// 对每个变体计算成立的役与对应符数，最后按 番->符->荣和分 选择最佳结果（与 C++ 一致的选择规则）。
	s.Player = player
	s.Table = table
	s.Tiles = tiles
	s.CallGroups = callGroups
	s.WinTile = winTile
	s.IsSevenPair = isSevenPair

	splitter := GetTileSplitter()
	completedList := splitter.GetAllCompletedTiles(s.Tiles)

	var candidates []*ScoreCounterResult

	for _, ct := range completedList {
		// 构造合并的组序列：雀头（如存在） -> 副露组 -> 手中面子
		combined := make([]TileGroup, 0)
		hasHead := false
		if len(ct.Head.Tiles) > 0 {
			combined = append(combined, ct.Head)
			hasHead = true
		}
		for _, cg := range callGroups {
			combined = append(combined, TileGroup{Type: cg.Type, Tiles: cg.Tiles})
		}
		for _, g := range ct.Body {
			combined = append(combined, g)
		}

		// body 在 combined 中的起始偏移
		bodyOffset := 0
		if hasHead {
			bodyOffset = 1
		}
		bodyOffset += len(callGroups)

		// 枚举 combined 中包含和牌的位置
		for i, g := range combined {
			switch g.Type {
			case Shuntsu:
				// 顺子可能在三种位置上和牌
				for pos := 0; pos < len(g.Tiles) && pos < 3; pos++ {
					if g.Tiles[pos] != winTile {
						continue
					}
					// 变体：ct + win 在 bodyIndex (若属于手牌部分) + 顺子位置 pos
					variant := s.evaluateVariant(&ct, callGroups, hasHead, bodyOffset, i, pos)
					if variant != nil {
						candidates = append(candidates, variant)
					}
				}
			default:
				// 刻子/对子/杠：只要包含和牌即可
				for _, t := range g.Tiles {
					if t == winTile {
						variant := s.evaluateVariant(&ct, callGroups, hasHead, bodyOffset, i, -1)
						if variant != nil {
							candidates = append(candidates, variant)
						}
						break
					}
				}
			}
		}
	}

	// 从候选结果中选择最佳
	best := s.GetBestResult(candidates)
	return best
}

// evaluateVariant 对单个拆分(ct)的某个和牌位置进行评估，返回对应的 ScoreCounterResult
// params:
// - ct: 拆分
// - callGroups: 副露组
// - hasHead: ct 是否有雀头
// - bodyOffset: combined 中 body 的起始偏移（用于将 combined 索引映射到 ct.Body）
// - combinedIndex: 和牌所在的 combined 索引
// - shuntsuPos: 若为顺子，和牌在顺子中的位置(0/1/2)，否则为 -1
func (s *ScoreCounter) evaluateVariant(ct *CompletedTiles, callGroups []CallGroup, hasHead bool, bodyOffset, combinedIndex, shuntsuPos int) *ScoreCounterResult {
	// 设置 ScoreCounter 上下文以便 CheckYaku 使用现有实现
	prevTiles := s.Tiles
	prevCallGroups := s.CallGroups
	prevWin := s.WinTile
	prevIsSeven := s.IsSevenPair
	defer func() {
		s.Tiles = prevTiles
		s.CallGroups = prevCallGroups
		s.WinTile = prevWin
		s.IsSevenPair = prevIsSeven
	}()

	// s.Tiles 已经是用于拆分的整副牌（含和牌），保持不变
	s.CallGroups = callGroups
	s.IsSevenPair = (len(ct.Body) == 7 && len(ct.Head.Tiles) == 0)

	// 识别 tsumo：若 winTile 在玩家手中，则认为是自摸
	tsumo := false
	if CountTile(s.Tiles, s.WinTile) > 0 {
		tsumo = true
	}

	// 收集成立的役（包含役满）
	yakus := make([]Yaku, 0)
	for y := Yaku(0); y < MaxYaku; y++ {
		if s.CheckYaku(y) {
			yakus = append(yakus, y)
		}
	}
	if len(yakus) == 0 {
		return nil
	}

	// 计算符数（基于当前拆分与和牌位置）
	fu := s.calculateFuForVariant(ct, callGroups, tsumo, hasHead, bodyOffset, combinedIndex, shuntsuPos)

	fan := s.CalculateFan(yakus)

	res := &ScoreCounterResult{
		Fan:   fan,
		Fu:    fu,
		Yakus: yakus,
	}
	if fan >= 13 {
		res.IsYakuman = true
		res.BaseScore = 8000
		res.RonScore = 12000
		res.TsumoScore = [3]int{4000, 4000, 4000}
	} else {
		res.BaseScore = s.CalculateBaseScore(fan, fu)
		res.RonScore = s.CalculateRonScore(res.BaseScore)
		res.TsumoScore = s.CalculateTsumoScore(res.BaseScore)
	}

	return res
}

// calculateFuForVariant 基于给定的 CompletedTiles 与和牌位置，计算符数（模仿 C++ 实现，包含听牌/边张/坎张判断）
func (s *ScoreCounter) calculateFuForVariant(ct *CompletedTiles, callGroups []CallGroup, tsumo bool, hasHead bool, bodyOffset, combinedIndex, shuntsuPos int) int {
	// 七对子固定25
	if len(ct.Body) == 7 && len(ct.Head.Tiles) == 0 {
		return 25
	}

	fu := 20

	// 雀头符（役牌对）：检查 ct.Head
	if len(ct.Head.Tiles) > 0 {
		head := ct.Head.Tiles[0]
		if s.Table != nil {
			if IsYakuhai(head, s.Table.GameWind, s.Player.Wind) {
				fu += 2
			}
		} else {
			if Is567z(head) {
				fu += 2
			}
		}
	}

	// 副露与杠的符
	for _, cg := range callGroups {
		switch cg.Type {
		case Shuntsu:
			// 吃不计符
		case Koutsu:
			if len(cg.Tiles) == 0 {
				continue
			}
			tile := cg.Tiles[0]
			if cg.IsOpen {
				if IsTerminalOrHonor(tile) {
					fu += 4
				} else {
					fu += 2
				}
			} else {
				if IsTerminalOrHonor(tile) {
					fu += 8
				} else {
					fu += 4
				}
			}
		case Kantsu:
			if len(cg.Tiles) == 0 {
				continue
			}
			tile := cg.Tiles[0]
			if cg.IsOpen {
				if IsTerminalOrHonor(tile) {
					fu += 16
				} else {
					fu += 8
				}
			} else {
				if IsTerminalOrHonor(tile) {
					fu += 32
				} else {
					fu += 16
				}
			}
		}
	}

	// 来自手中未鸣的面子（暗刻/暗杠）
	for _, g := range ct.Body {
		switch g.Type {
		case Koutsu:
			tile := g.Tiles[0]
			if IsTerminalOrHonor(tile) {
				fu += 8
			} else {
				fu += 4
			}
		case Kantsu:
			tile := g.Tiles[0]
			if IsTerminalOrHonor(tile) {
				fu += 32
			} else {
				fu += 16
			}
		}
	}

	// 听牌型判断（单骑/坎张/边张）——基于 combinedIndex 与 shuntsuPos
	// 先判断单骑（和在雀头）
	if len(ct.Head.Tiles) > 0 {
		if combinedIndex == 0 && len(ct.Head.Tiles) > 0 && ct.Head.Find(s.WinTile) >= 0 {
			// 单骑
			fu += 2
		}
	}

	// 若和在手中某个顺子，判断是否为坎张或边张
	// 判断 combinedIndex 是否对应到 ct.Body 中的某个组
	if combinedIndex >= bodyOffset && combinedIndex-bodyOffset < len(ct.Body) {
		bi := combinedIndex - bodyOffset
		g := ct.Body[bi]
		if g.Type == Shuntsu {
			// 找到和牌在该顺子的位置
			pos := shuntsuPos
			if pos < 0 {
				// 安全回退：尝试在组内查找
				for i, t := range g.Tiles {
					if t == s.WinTile {
						pos = i
						break
					}
				}
			}
			// 坎张：中张
			if pos == 1 {
				fu += 2
			}
			// 边张：1-2-3 且和为3 或 7-8-9 且和为1
			if pos == 2 {
				// 顺子的首张是 g.Tiles[0]
				if Is1hai(g.Tiles[0]) {
					fu += 2
				}
			}
			if pos == 0 {
				// 如果首张是 7（即 7-8-9），则 win 在首位为边张
				rank := int(g.Tiles[0]) % 9
				if rank == 6 { // 0-based: 6 -> 7
					fu += 2
				}
			}
		}
	}

	// 自摸/门清荣和
	if tsumo {
		fu += 2
	} else {
		if s.Player.IsMenzen() {
			fu += 10
		}
	}

	// 副露平和的特殊处理：若存在副露且非门清且符仍为20则变为30
	if !s.Player.IsMenzen() {
		// 若为荣和并非门清，且当前符为20，则按 C++ 逻辑设置为30
		if fu == 20 {
			fu = 30
		}
	}

	// 平和自摸 20符处理
	// 如果满足平和（Pinfu）并且自摸，则固定20符
	if s.CheckPinfu() && tsumo {
		fu = 20
	}

	if fu%10 != 0 {
		fu = ((fu / 10) + 1) * 10
	}

	return fu
}

// CalculateFan 计算番数
func (s *ScoreCounter) CalculateFan(yakus []Yaku) int {
	fanCount := 0
	yakumanCount := 0

	for _, yaku := range yakus {
		fanVal := GetFanCount(yaku)
		if fanVal >= 13 {
			yakumanCount++
		} else {
			fanCount += fanVal
		}
	}

	if yakumanCount > 0 {
		return 13 * yakumanCount
	}
	return fanCount
}

// CalculateFu 计算符数
func (s *ScoreCounter) CalculateFu() int {
	// 七对子固定25符
	if s.IsSevenPair {
		return 25
	}

	// 基础符
	baseFu := 20

	// 判断是否自摸（tsumo）：如果玩家的手牌中已包含胜利牌，则为自摸
	handBase := ConvertTilesToBaseTiles(s.Player.Hand)
	tsumo := CountTile(handBase, s.WinTile) > 0

	bestFu := 0

	// 遍历所有可能的拆牌（不同拆法可能导致不同的符数），选择最大的符数作为安全值
	splitter := GetTileSplitter()
	completedList := splitter.GetAllCompletedTiles(s.Tiles)
	if len(completedList) == 0 {
		// 兜底：返回最小的符（进位后）
		fu := baseFu
		if tsumo {
			fu += 2
		} else if s.Player.IsMenzen() {
			// 门清荣和 +10
			fu += 10
		}
		if fu%10 != 0 {
			fu = ((fu / 10) + 1) * 10
		}
		return fu
	}

	for _, ct := range completedList {
		fu := baseFu

		// 雀头计符：如果是役牌（场风/自风/三元牌）则 +2
		if len(ct.Head.Tiles) > 0 {
			head := ct.Head.Tiles[0]
			if s.Table != nil {
				if IsYakuhai(head, s.Table.GameWind, s.Player.Wind) {
					fu += 2
				}
			} else {
				// 保守判断：若为三元牌也计符
				if head == _5z || head == _6z || head == _7z {
					fu += 2
				}
			}
		}

		// 先统计副露（鸣牌/杠）带来的符 — 这些在 CallGroups 中
		for _, cg := range s.CallGroups {
			if cg.Type == Shuntsu {
				// 吃不计符
				continue
			}
			if cg.Type == Koutsu {
				if len(cg.Tiles) == 0 {
					continue
				}
				tile := cg.Tiles[0]
				if cg.IsOpen {
					// 明刻
					if IsTerminalOrHonor(tile) {
						fu += 4
					} else {
						fu += 2
					}
				} else {
					// 暗刻（暗杠在CallGroups也可能以Kantsu出现）
					if IsTerminalOrHonor(tile) {
						fu += 8
					} else {
						fu += 4
					}
				}
			}
			if cg.Type == Kantsu {
				if len(cg.Tiles) == 0 {
					continue
				}
				tile := cg.Tiles[0]
				if cg.IsOpen {
					// 明杠
					if IsTerminalOrHonor(tile) {
						fu += 16
					} else {
						fu += 8
					}
				} else {
					// 暗杠
					if IsTerminalOrHonor(tile) {
						fu += 32
					} else {
						fu += 16
					}
				}
			}
		}

		// 再统计手中未鸣的面子（ct.Body），这些通常为暗的刻/杠
		for _, g := range ct.Body {
			if g.Type == Koutsu {
				tile := g.Tiles[0]
				// 因为这些是来自手牌的组，视为暗刻
				if IsTerminalOrHonor(tile) {
					fu += 8
				} else {
					fu += 4
				}
			} else if g.Type == Kantsu {
				tile := g.Tiles[0]
				// 来自手牌的杠，视为暗杠
				if IsTerminalOrHonor(tile) {
					fu += 32
				} else {
					fu += 16
				}
			}
		}

		// 自摸或门清荣和
		if tsumo {
			fu += 2
		} else {
			if s.Player.IsMenzen() {
				// 门清荣和 +10 符
				fu += 10
			}
		}

		// 七对子在此已处理，上面不会命中

		// 进位到最近的10
		if fu%10 != 0 {
			fu = ((fu / 10) + 1) * 10
		}

		if fu > bestFu {
			bestFu = fu
		}
	}

	if bestFu == 0 {
		// 兜底
		bestFu = 20
		if tsumo {
			bestFu += 2
		} else if s.Player.IsMenzen() {
			bestFu += 10
		}
		if bestFu%10 != 0 {
			bestFu = ((bestFu / 10) + 1) * 10
		}
	}

	return bestFu
}

// CalculateBaseScore 计算基础分
func (s *ScoreCounter) CalculateBaseScore(fan int, fu int) int {
	if fan >= 13 {
		return 8000 // 役满
	} else if fan >= 11 {
		return 6000 // 三倍满
	} else if fan >= 8 {
		return 4000 // 倍满
	} else if fan >= 6 {
		return 3000 // 跳满
	} else if fan >= 5 {
		return 2000 // 满贯
	} else if fan == 4 {
		baseScore := fu * 1000
		if baseScore > 2000 {
			return 2000
		}
		return baseScore
	} else if fan == 3 {
		baseScore := fu * 500
		if baseScore > 2000 {
			return 2000
		}
		return baseScore
	} else if fan == 2 {
		baseScore := fu * 250
		if baseScore > 2000 {
			return 2000
		}
		return baseScore
	} else {
		baseScore := fu * 100
		if baseScore > 2000 {
			return 2000
		}
		return baseScore
	}
}

// CalculateRonScore 计算荣和分
func (s *ScoreCounter) CalculateRonScore(baseScore int) int {
	if s.Player.Oya {
		return baseScore * 6
	}
	return baseScore * 4
}

// CalculateTsumoScore 计算自摸分
func (s *ScoreCounter) CalculateTsumoScore(baseScore int) [3]int {
	var tsumoScore [3]int
	if s.Player.Oya {
		// 庄家自摸，每家都付2倍基础分
		tsumoScore[0] = baseScore * 2
		tsumoScore[1] = baseScore * 2
		tsumoScore[2] = baseScore * 2
	} else {
		// 子家自摸，庄家付2倍基础分，其他付1倍基础分
		tsumoScore[0] = baseScore * 2
		tsumoScore[1] = baseScore
		tsumoScore[2] = baseScore
	}
	return tsumoScore
}

// CheckYaku 检查是否满足某个役
func (s *ScoreCounter) CheckYaku(yaku Yaku) bool {
	switch yaku {
	case Tanyao:
		return s.CheckTanyao()
	case Pinfu:
		return s.CheckPinfu()
	case Iipeikou:
		return s.CheckIipeikou()
	case Ryanpeikou:
		return s.CheckRyanpeikou()
	case Yakuhai, YakuhaiWind, YakuhaiWhiteBoard, YakuhaiGreenBoard, YakuhaiRedBoard:
		return s.CheckYakuhai()
	case Honitsu:
		return s.CheckHonitsu()
	case Chinitsu:
		return s.CheckChinitsu()
	case Chanta:
		return s.CheckChanta()
	case Honchanta:
		return s.CheckHonchanta()
	case Ittsu:
		return s.CheckIkkitsuukan()
	case Sanshokudoukou:
		return s.CheckSanshokudoukou()
	case Toitoi:
		return s.CheckToitoi()
	case Sanankou:
		return s.CheckSanankou()
	case Tsuisou:
		return s.CheckTsuisou()
	case Ryuuisou:
		return s.CheckRyuuisou()
	case Chinroutou:
		return s.CheckChinroutou()
	case Honroutou:
		return s.CheckHonroutou()
	case Chiitoitsu:
		return s.CheckChiitoitsu()
	case Kokushi:
		return s.CheckKokushi()
	case Daisangen:
		return s.CheckDaisangen()
	case Suankou:
		return s.CheckSuankou()
	case Daisuushi:
		return s.CheckDaisuushi()
	case Shosuushi:
		return s.CheckShousuushi()
	case Tenhou:
		return s.CheckTenhou()
	case Chihou:
		return s.CheckChihou()
	case Churen:
		return s.CheckChuren()
	case Dabururiichi:
		return s.CheckDabururiichi()
	case Menzentsumo:
		return s.CheckMenzentsumo()
	default:
		return false
	}
}

// CheckTanyao 检查断幺
func (s *ScoreCounter) CheckTanyao() bool {
	for _, tile := range s.Tiles {
		if Is1hai(tile) || Is9hai(tile) {
			return false
		}
	}
	return true
}

// CheckPinfu 检查平和
func (s *ScoreCounter) CheckPinfu() bool {
	if !s.Player.IsMenzen() {
		return false
	}
	if s.IsSevenPair {
		return false
	}
	// 使用拆分结果判断平和：
	// - 所有面子为顺子
	// - 雀头不是役牌
	// - 胡牌为两面听（在某顺子中移除胡牌后剩下两张非幺九且相差1）
	splitter := GetTileSplitter()
	allCompleted := splitter.GetAllCompletedTiles(s.Tiles)
	for _, ct := range allCompleted {
		if len(ct.Head.Tiles) == 0 {
			continue
		}
		// 雀头不能是役牌
		headTile := ct.Head.Tiles[0]
		if s.Table != nil && IsYakuhai(headTile, s.Table.GameWind, s.Player.Wind) {
			continue
		}
		// 所有身子必须为顺子
		allShuntsu := true
		for _, g := range ct.Body {
			if g.Type != Shuntsu {
				allShuntsu = false
				break
			}
		}
		if !allShuntsu {
			continue
		}
		// 查找包含胡牌的顺子并验证为两面听
		for _, g := range ct.Body {
			if g.Type != Shuntsu {
				continue
			}
			idx := -1
			for i, t := range g.Tiles {
				if t == s.WinTile {
					idx = i
					break
				}
			}
			if idx < 0 {
				continue
			}
			rem := make([]BaseTile, 0, 2)
			for i, t := range g.Tiles {
				if i == idx {
					continue
				}
				rem = append(rem, t)
			}
			if len(rem) != 2 {
				continue
			}
			// 排序
			if rem[0] > rem[1] {
				rem[0], rem[1] = rem[1], rem[0]
			}
			// 剩余两张不能是幺九，并且相差1
			if Is19hai(rem[0]) || Is19hai(rem[1]) {
				continue
			}
			if rem[1]-rem[0] == 1 {
				return true
			}
		}
	}
	return false
}

// CheckIipeikou 检查一对对
func (s *ScoreCounter) CheckIipeikou() bool {
	if !s.Player.IsMenzen() {
		return false
	}
	if s.IsSevenPair {
		return false
	}
	splitter := GetTileSplitter()
	allCompleted := splitter.GetAllCompletedTiles(s.Tiles)
	for _, ct := range allCompleted {
		// 统计相同顺子的出现次数
		seqCount := make(map[string]int)
		for _, g := range ct.Body {
			if g.Type != Shuntsu {
				continue
			}
			key := BaseTileToString(g.Tiles[0]) + "," + BaseTileToString(g.Tiles[1]) + "," + BaseTileToString(g.Tiles[2])
			seqCount[key]++
		}
		for _, v := range seqCount {
			if v >= 2 {
				return true
			}
		}
	}
	return false
}

// CheckRyanpeikou 检查二对对
func (s *ScoreCounter) CheckRyanpeikou() bool {
	if !s.Player.IsMenzen() {
		return false
	}
	if s.IsSevenPair {
		return false
	}
	splitter := GetTileSplitter()
	allCompleted := splitter.GetAllCompletedTiles(s.Tiles)
	for _, ct := range allCompleted {
		seqCount := make(map[string]int)
		for _, g := range ct.Body {
			if g.Type != Shuntsu {
				continue
			}
			key := BaseTileToString(g.Tiles[0]) + "," + BaseTileToString(g.Tiles[1]) + "," + BaseTileToString(g.Tiles[2])
			seqCount[key]++
		}
		dup := 0
		for _, v := range seqCount {
			if v >= 2 {
				dup++
			}
		}
		if dup >= 2 {
			return true
		}
	}
	return false
}

// CheckYakuhai 检查役牌（役牌对子）
func (s *ScoreCounter) CheckYakuhai() bool {
	if s.Table == nil {
		// 无法判断场风，退化为检查手牌中是否有三元牌或自家风/场风对子（保守实现）
	}
	tileCount := make(map[BaseTile]int)
	for _, t := range s.Tiles {
		tileCount[t]++
	}
	for tile, cnt := range tileCount {
		if cnt >= 2 {
			if s.Table != nil {
				if IsYakuhai(tile, s.Table.GameWind, s.Player.Wind) {
					return true
				}
			} else {
				// 只检查三元牌（白、发、中）作为保底
				if tile == _5z || tile == _6z || tile == _7z {
					return true
				}
			}
		}
	}
	return false
}

// CheckHonitsu 检查混一色
func (s *ScoreCounter) CheckHonitsu() bool {
	if s.IsSevenPair {
		return false
	}
	// 检查是否只含有一种花色和字牌
	types := CountTileType(s.Tiles)
	numberTypes := 0
	for i := 0; i < 3; i++ {
		if types[i] > 0 {
			numberTypes++
		}
	}
	return numberTypes == 1 && types[3] > 0
}

// CheckChinitsu 检查清一色
func (s *ScoreCounter) CheckChinitsu() bool {
	if s.IsSevenPair {
		return false
	}
	// 检查是否只含有一种花色
	types := CountTileType(s.Tiles)
	numberTypes := 0
	for i := 0; i < 3; i++ {
		if types[i] > 0 {
			numberTypes++
		}
	}
	return numberTypes == 1 && types[3] == 0
}

// GetBestResult 从多个可能的胡牌结果中选择最佳的
func (s *ScoreCounter) GetBestResult(results []*ScoreCounterResult) *ScoreCounterResult {
	if len(results) == 0 {
		return nil
	}

	best := results[0]
	for _, result := range results[1:] {
		if result.Fan > best.Fan {
			best = result
		} else if result.Fan == best.Fan && result.Fu > best.Fu {
			best = result
		} else if result.Fan == best.Fan && result.Fu == best.Fu && result.RonScore > best.RonScore {
			best = result
		}
	}

	return best
}

// ToString 返回计分结果的字符串表示
func (r *ScoreCounterResult) ToString() string {
	str := fmt.Sprintf("番: %d, 符: %d\n", r.Fan, r.Fu)
	str += fmt.Sprintf("基础分: %d, 荣和: %d\n", r.BaseScore, r.RonScore)
	str += fmt.Sprintf("自摸: %d-%d-%d\n", r.TsumoScore[0], r.TsumoScore[1], r.TsumoScore[2])
	str += "役: "
	for i, yaku := range r.Yakus {
		if i > 0 {
			str += ", "
		}
		str += YakuToString(yaku)
	}
	str += "\n"
	return str
}
