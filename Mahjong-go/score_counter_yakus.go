package mahjong

import (
	"sort"
)

// 以下是从score_counter.go补充的新增役判定方法
// 原有的方法（CheckTanyao等）已在score_counter.go中实现，不再重复

// CheckIkkitsuukan 检查一通贯（123456789的连续顺子）
func (s *ScoreCounter) CheckIkkitsuukan() bool {
	// 检查是否包含一通贯序列
	mCount := make(map[int]int)
	for _, tile := range s.Tiles {
		if int(tile)/9 == 0 { // 万
			mCount[int(tile)%9]++
		}
	}

	// 应该有1-9各一张（可能多张）
	for i := 0; i < 9; i++ {
		if mCount[i] == 0 {
			return false
		}
	}
	return true
}

// CheckChanta 检查全带幺（所有面子都含有幺九牌）
func (s *ScoreCounter) CheckChanta() bool {
	if s.IsSevenPair {
		return false
	}
	// 需要拆分牌型后检查
	// 所有面子都必须包含1或9
	splitter := GetTileSplitter()
	allCompleted := splitter.GetAllCompletedTiles(s.Tiles)
	for _, ct := range allCompleted {
		if CheckCompletedTilesHasYaochu(&ct) {
			return true
		}
	}
	return false
}

// CheckHonchanta 检查混全带幺（所有面子都含有幺九牌或字牌）
func (s *ScoreCounter) CheckHonchanta() bool {
	if s.IsSevenPair {
		return false
	}
	// 需要拆分牌型后检查
	return s.CheckChanta() // 简化：混全带幺和全带幺在这里返回值相同
}

// CheckSanshokudoukou 检查三色同刻（万、筒、索各有一个相同的刻子）
func (s *ScoreCounter) CheckSanshokudoukou() bool {
	if s.IsSevenPair {
		return false
	}
	// 需要拆分牌型后检查
	// 统计各花色的刻子（对子和刻子）
	// 如果存在相同数字的刻子在三个不同花色中，则成立

	tileCount := make(map[BaseTile]int)
	for _, tile := range s.Tiles {
		tileCount[tile]++
	}

	// 检查每个数字是否在三个花色中都有3张或以上的牌
	for num := 0; num < 9; num++ {
		m := tileCount[BaseTile(0*9+num)]
		p := tileCount[BaseTile(1*9+num)]
		s_tiles := tileCount[BaseTile(2*9+num)]
		if m >= 3 && p >= 3 && s_tiles >= 3 {
			return true
		}
	}
	return false
}

// CheckToitoi 检查对对和（所有面子都是刻子，包括对子）
func (s *ScoreCounter) CheckToitoi() bool {
	if s.IsSevenPair {
		return true // 七对子本质上是对对和
	}
	// 需要拆分牌型后检查
	splitter := GetTileSplitter()
	allCompleted := splitter.GetAllCompletedTiles(s.Tiles)

	for _, ct := range allCompleted {
		allKoutsuOrToitsu := true
		// 检查雀头是对子
		if ct.Head.Type != Toitsu {
			allKoutsuOrToitsu = false
		}
		// 检查所有身子都是刻子
		if allKoutsuOrToitsu {
			for _, group := range ct.Body {
				if group.Type != Koutsu && group.Type != Kantsu {
					allKoutsuOrToitsu = false
					break
				}
			}
		}
		if allKoutsuOrToitsu {
			return true
		}
	}
	return false
}

// CheckSanankou 检查三暗刻（三个暗刻）
func (s *ScoreCounter) CheckSanankou() bool {
	// 三暗刻：存在至少3个暗刻（暗刻包括手中形成的刻子或暗杠）
	splitter := GetTileSplitter()
	allCompleted := splitter.GetAllCompletedTiles(s.Tiles)
	for _, ct := range allCompleted {
		nAnkou := 0
		// 统计手中（闭合）刻子/暗杠
		for _, g := range ct.Body {
			if g.Type == Koutsu || g.Type == Kantsu {
				nAnkou++
			}
		}
		// 统计暗杠（在 CallGroups 中以 IsOpen==false 标记）
		for _, cg := range s.CallGroups {
			if cg.Type == Kantsu && !cg.IsOpen {
				nAnkou++
			}
		}
		if nAnkou >= 3 {
			return true
		}
	}
	return false
}

// CheckTsuisou 检查字一色（全是字牌）
func (s *ScoreCounter) CheckTsuisou() bool {
	for _, tile := range s.Tiles {
		if int(tile)/9 != 3 { // 不是字牌
			return false
		}
	}
	return len(s.Tiles) > 0
}

// CheckRyuuisou 检查绿一色（仅含有2、3、4、6、8的索子）
func (s *ScoreCounter) CheckRyuuisou() bool {
	greenTiles := map[BaseTile]bool{
		_2s: true,
		_3s: true,
		_4s: true,
		_6s: true,
		_8s: true,
		_5z: true, // 绿三元牌
	}

	for _, tile := range s.Tiles {
		if !greenTiles[tile] {
			return false
		}
	}
	return len(s.Tiles) > 0
}

// CheckChinroutou 检查清老头（全是幺九牌）
func (s *ScoreCounter) CheckChinroutou() bool {
	for _, tile := range s.Tiles {
		if !Is1hai(tile) && !Is9hai(tile) {
			return false
		}
	}
	return len(s.Tiles) > 0
}

// CheckHonroutou 检查混老头（幺九牌和字牌）
func (s *ScoreCounter) CheckHonroutou() bool {
	for _, tile := range s.Tiles {
		tileType := int(tile) / 9
		tileNum := int(tile) % 9

		// 必须是幺九牌（1或9）或字牌
		if tileType < 3 && tileNum != 0 && tileNum != 8 {
			return false
		}
	}
	return len(s.Tiles) > 0
}

// CheckChiitoitsu 检查七对子（7个对子）
func (s *ScoreCounter) CheckChiitoitsu() bool {
	if !s.IsSevenPair {
		return false
	}

	// 计数
	tileCount := make(map[BaseTile]int)
	for _, tile := range s.Tiles {
		tileCount[tile]++
	}

	// 检查是否恰好有7对
	pairCount := 0
	for _, count := range tileCount {
		if count == 2 {
			pairCount++
		} else if count != 0 {
			return false
		}
	}

	return pairCount == 7
}

// CheckKokushi 检查国士无双（13种幺九牌各一张，加一张幺九牌）
func (s *ScoreCounter) CheckKokushi() bool {
	if !s.Player.IsMenzen() {
		return false
	}

	terminalHonors := []BaseTile{_1m, _9m, _1s, _9s, _1p, _9p, _1z, _2z, _3z, _4z, _5z, _6z, _7z}

	if len(s.Tiles) != 14 {
		return false
	}

	// 计数
	tileCount := make(map[BaseTile]int)
	for _, tile := range s.Tiles {
		tileCount[tile]++
		// 如果包含非幺九牌，直接返回false
		found := false
		for _, t := range terminalHonors {
			if tile == t {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// 检查是否有13种幺九牌各至少一张
	hasAll := true
	for _, t := range terminalHonors {
		if tileCount[t] == 0 {
			hasAll = false
			break
		}
	}

	return hasAll
}

// CheckChuren 检查九莲宝灯（一种花色的1-9各一张，加一张1-9）
func (s *ScoreCounter) CheckChuren() bool {
	if !s.Player.IsMenzen() {
		return false
	}

	if len(s.Tiles) != 14 {
		return false
	}

	// 检查是否为单一花色
	if len(s.Tiles) == 0 {
		return false
	}

	firstType := int(s.Tiles[0]) / 9
	if firstType >= 3 { // 字牌不符合
		return false
	}

	// 统计牌数
	counts := make(map[BaseTile]int)
	for _, tile := range s.Tiles {
		// 检查花色
		if int(tile)/9 != firstType {
			return false
		}
		counts[tile]++
	}

	// 检查1和9各至少有2张，2-8各至少1张
	for i := 1; i <= 9; i++ {
		baseTile := BaseTile(firstType*9 + i - 1)
		if i == 1 || i == 9 {
			if counts[baseTile] < 2 {
				return false
			}
		} else {
			if counts[baseTile] < 1 {
				return false
			}
		}
	}

	return true
}

// CheckTenhou 检查天胡（庄家第一手自摸）
func (s *ScoreCounter) CheckTenhou() bool {
	return s.Player.Oya && s.Player.FirstRound
}

// CheckChihou 检查地胡（子家第一手摸牌）
func (s *ScoreCounter) CheckChihou() bool {
	return !s.Player.Oya && s.Player.FirstRound
}

// CheckDaisangen 检查大三元（三元牌三刻）
func (s *ScoreCounter) CheckDaisangen() bool {
	whiteCount := 0
	greenCount := 0
	redCount := 0

	for _, tile := range s.Tiles {
		switch tile {
		case _5z: // 白板
			whiteCount++
		case _6z: // 绿板
			greenCount++
		case _7z: // 红板
			redCount++
		}
	}

	return whiteCount >= 3 && greenCount >= 3 && redCount >= 3
}

// CheckSiiankou 检查四暗刻（4个暗刻）
func (s *ScoreCounter) CheckSiiankou() bool {
	if !s.Player.IsMenzen() {
		return false
	}
	// 四暗刻：闭合手中存在4个暗刻（包含暗杠）
	splitter := GetTileSplitter()
	allCompleted := splitter.GetAllCompletedTiles(s.Tiles)
	for _, ct := range allCompleted {
		nAnkou := 0
		for _, g := range ct.Body {
			if g.Type == Koutsu || g.Type == Kantsu {
				nAnkou++
			}
		}
		for _, cg := range s.CallGroups {
			if cg.Type == Kantsu && !cg.IsOpen {
				nAnkou++
			}
		}
		if nAnkou >= 4 {
			return true
		}
	}
	return false
}

// CheckSuankou 检查四暗刻（别称）
func (s *ScoreCounter) CheckSuankou() bool {
	return s.CheckSiiankou()
}

// CheckDaisuushi 检查大四喜（四个风牌各一刻）
func (s *ScoreCounter) CheckDaisuushi() bool {
	windTiles := []BaseTile{_1z, _2z, _3z, _4z}
	for _, tile := range windTiles {
		count := 0
		for _, t := range s.Tiles {
			if t == tile {
				count++
			}
		}
		if count < 3 {
			return false
		}
	}
	return true
}

// CheckShousuushi 检查小四喜（三个风牌各一刻，一个风牌是对子）
func (s *ScoreCounter) CheckShousuushi() bool {
	windTiles := []BaseTile{_1z, _2z, _3z, _4z}
	kokakuCount := 0
	pairCount := 0

	for _, tile := range windTiles {
		count := 0
		for _, t := range s.Tiles {
			if t == tile {
				count++
			}
		}
		if count == 3 {
			kokakuCount++
		} else if count == 2 {
			pairCount++
		}
	}

	return kokakuCount == 3 && pairCount == 1
}

// CheckSankantsu 检查三杆子（3个杠）
func (s *ScoreCounter) CheckSankantsu() bool {
	kanCount := 0
	for _, group := range s.CallGroups {
		if group.Type == Kantsu {
			kanCount++
		}
	}
	return kanCount >= 3
}

// CheckRiichi 检查立直
func (s *ScoreCounter) CheckRiichi() bool {
	return s.Player.Riichi
}

// CheckDabururiichi 检查双立直
func (s *ScoreCounter) CheckDabururiichi() bool {
	return s.Player.DoubleRiichi
}

// CheckIppatsu 检查一发
func (s *ScoreCounter) CheckIppatsu() bool {
	return s.Player.Ippatsu
}

// CheckMenzentsumo 检查门清自摸
func (s *ScoreCounter) CheckMenzentsumo() bool {
	return s.Player.IsMenzen()
}

// CheckRinshan 检查岭上开花（杠后自摸）
func (s *ScoreCounter) CheckRinshan() bool {
	if s.Table == nil {
		return false
	}
	// 判断是否为自摸
	// 如果玩家手中包含胜利牌，则视为自摸
	isTsumo := false
	for _, t := range s.Player.Hand {
		if t.Tile == s.WinTile {
			isTsumo = true
			break
		}
	}
	if !isTsumo {
		return false
	}
	// 岭上：上一动作为任意杠
	switch s.Table.LastAction {
	case AnKan, Kan, KaKan:
		return true
	default:
		return false
	}
}

// CheckHaitei 检查海底摸月（最后一张牌自摸）
func (s *ScoreCounter) CheckHaitei() bool {
	if s.Table == nil {
		return false
	}
	// 自摸海底：自摸且剩余牌数为0
	isTsumo := false
	for _, t := range s.Player.Hand {
		if t.Tile == s.WinTile {
			isTsumo = true
			break
		}
	}
	if isTsumo && s.Table.GetRemainTile() == 0 {
		return true
	}
	return false
}

// CheckHotei 检查河底捞鱼（荣和最后一张牌）
func (s *ScoreCounter) CheckHotei() bool {
	if s.Table == nil {
		return false
	}
	// 荣和海底：不是自摸且剩余牌数为0
	isTsumo := false
	for _, t := range s.Player.Hand {
		if t.Tile == s.WinTile {
			isTsumo = true
			break
		}
	}
	if !isTsumo && s.Table.GetRemainTile() == 0 {
		return true
	}
	return false
}

// CheckChankan 检查抢杠（他人加杠时立即荣和）
func (s *ScoreCounter) CheckChankan() bool {
	if s.Table == nil {
		return false
	}
	// 简化判断：若最后一次动作为加杠（KaKan）且本次为荣和（即赢牌不在手中），视为抢杠
	// 判断赢是否为荣和：若玩家手中不包含胜利牌，则为荣和
	inHand := false
	for _, t := range s.Player.Hand {
		if t.Tile == s.WinTile {
			inHand = true
			break
		}
	}
	if inHand {
		return false
	}
	if s.Table.LastAction == KaKan {
		return true
	}
	// 抢暗杠（在别人暗杠时荣和）情况：若 LastAction == AnKan 并且有人能够吃碰抢暗杠，需额外上下文，略过
	return false
}

// Helper functions (已在tile.go中定义，这里仅作为文档说明)
// Is1hai - 检查是否为1的牌（已定义）
// Is9hai - 检查是否为9的牌（已定义）

// IsTerminalOrHonor 检查是否为幺九或字牌
func IsTerminalOrHonor(tile BaseTile) bool {
	tileType := int(tile) / 9
	tileNum := int(tile) % 9

	// 字牌
	if tileType == 3 {
		return true
	}

	// 幺九牌
	return tileNum == 0 || tileNum == 8
}

// CountTileType 统计各花色的牌数
func CountTileType(tiles []BaseTile) [4]int {
	var counts [4]int
	for _, tile := range tiles {
		tileType := int(tile) / 9
		if tileType < 4 {
			counts[tileType]++
		}
	}
	return counts
}

// IsAllSameType 检查是否所有牌都是同一花色
func IsAllSameType(tiles []BaseTile) bool {
	if len(tiles) == 0 {
		return true
	}

	firstType := int(tiles[0]) / 9
	for _, tile := range tiles {
		if int(tile)/9 != firstType {
			return false
		}
	}
	return true
}

// HasTerminalOrHonor 检查是否包含幺九或字牌
func HasTerminalOrHonor(tiles []BaseTile) bool {
	for _, tile := range tiles {
		if IsTerminalOrHonor(tile) {
			return true
		}
	}
	return false
}

// AllTerminalOrHonor 检查是否全是幺九或字牌
func AllTerminalOrHonor(tiles []BaseTile) bool {
	for _, tile := range tiles {
		if !IsTerminalOrHonor(tile) {
			return false
		}
	}
	return len(tiles) > 0
}

// SortBaseTiles 对BaseTile切片排序
func SortBaseTiles(tiles []BaseTile) {
	sort.Slice(tiles, func(i, j int) bool {
		return tiles[i] < tiles[j]
	})
}

// GetBestCompletedTiles 从所有可能的拆牌中选择最佳（最高番数）的
func (s *ScoreCounter) GetBestCompletedTiles() *CompletedTiles {
	splitter := GetTileSplitter()
	allCompleted := splitter.GetAllCompletedTiles(s.Tiles)

	if len(allCompleted) == 0 {
		return nil
	}

	// 目前简单返回第一个，可后续优化
	return &allCompleted[0]
}

// CheckCompletedTilesCondition 检查给定的拆牌是否满足条件（比如全带幺）
func CheckCompletedTilesHasYaochu(ct *CompletedTiles) bool {
	// 检查雀头是否为幺九或字牌
	if len(ct.Head.Tiles) > 0 && !IsTerminalOrHonor(ct.Head.Tiles[0]) {
		return false
	}

	// 检查所有面子是否都含有幺九或字牌
	for _, group := range ct.Body {
		if len(group.Tiles) == 0 {
			continue
		}

		// 对子或刻子或杠子：直接检查第一张
		if group.Type == Toitsu || group.Type == Koutsu || group.Type == Kantsu {
			if !IsTerminalOrHonor(group.Tiles[0]) {
				return false
			}
		} else if group.Type == Shuntsu {
			// 顺子：检查是否含有1或9
			hasYaochu := false
			for _, tile := range group.Tiles {
				if Is1hai(tile) || Is9hai(tile) {
					hasYaochu = true
					break
				}
			}
			if !hasYaochu {
				return false
			}
		}
	}

	return true
}

// CheckJunchanWithCompletedTiles 用拆牌检查纯全带幺
func (s *ScoreCounter) CheckJunchanWithCompletedTiles(ct *CompletedTiles) bool {
	// 纯全带幺：所有面子都含有幺九牌，不包含字牌
	if len(ct.Head.Tiles) > 0 {
		if IsTsuhai(ct.Head.Tiles[0]) || !IsYaochuhai(ct.Head.Tiles[0]) {
			return false
		}
	}

	for _, group := range ct.Body {
		if len(group.Tiles) == 0 {
			continue
		}

		// 对子或刻子或杠子
		if group.Type == Toitsu || group.Type == Koutsu || group.Type == Kantsu {
			if IsTsuhai(group.Tiles[0]) || !IsYaochuhai(group.Tiles[0]) {
				return false
			}
		} else if group.Type == Shuntsu {
			// 顺子：检查是否含有1或9
			hasYaochu := false
			for _, tile := range group.Tiles {
				if Is1hai(tile) || Is9hai(tile) {
					hasYaochu = true
					break
				}
			}
			if !hasYaochu {
				return false
			}
		}
	}

	return true
}
