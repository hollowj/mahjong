package mahjong

import (
	"fmt"
	"sort"
	"strings"
)

// RiverTile 表示河里的一张牌及其信息
type RiverTile struct {
	Tile     *Tile // 牌
	Number   int   // 第几张牌丢下去的
	Riichi   bool  // 是否为立直后弃牌
	Remain   bool  // 这张牌明面上还在不在河里
	FromHand bool  // true为手切，false为摸切
}

// River 表示河（已弃牌的记录）
type River struct {
	River []RiverTile
}

// ToBaseTile 将河中的所有牌转换为BaseTile数组
func (r *River) ToBaseTile() []BaseTile {
	result := make([]BaseTile, 0, len(r.River))
	for _, tile := range r.River {
		result = append(result, tile.Tile.Tile)
	}
	return result
}

// String 返回河的字符串表示
func (r *River) String() string {
	sb := strings.Builder{}
	for _, tile := range r.River {
		sb.WriteString(tile.Tile.String())
		sb.WriteString(fmt.Sprintf("%d", tile.Number))
		if tile.FromHand {
			sb.WriteString("h")
		}
		if tile.Riichi {
			sb.WriteString("r")
		}
		if !tile.Remain {
			sb.WriteString("-")
		}
		sb.WriteString(" ")
	}
	return sb.String()
}

// At 获取河中指定位置的牌
func (r *River) At(index int) *RiverTile {
	if index >= 0 && index < len(r.River) {
		return &r.River[index]
	}
	return nil
}

// Size 获取河中的牌数
func (r *River) Size() int {
	return len(r.River)
}

// PushBack 将牌添加到河中
func (r *River) PushBack(rt RiverTile) {
	r.River = append(r.River, rt)
}

// SetNotRemain 设置最后一张牌为不在河里（被取走）
func (r *River) SetNotRemain() {
	if len(r.River) > 0 {
		r.River[len(r.River)-1].Remain = false
	}
}

// Player 表示一个玩家
type Player struct {
	// 基本状态
	DoubleRiichi bool // 是否为两立直
	Riichi       bool // 是否已立直
	Menzen       bool // 是否为门前清（未鸣牌）
	Wind         Wind // 玩家风向
	Oya          bool // 是否为庄家
	Score        int  // 点数

	// 陷阱（复合状态）
	FuritenRound  bool // 回合振听（本回合弃过同牌）
	FuritenRiver  bool // 河振听（历史弃过同牌）
	FuritenRiichi bool // 立直振听（立直后弃过同牌）

	// 其他标记
	Ippatsu    bool // 是否有一发权
	FirstRound bool // 是否为第一回合

	// 牌
	Hand       []*Tile     // 手中的牌
	River      River       // 河
	CallGroups []CallGroup // 鸣牌组
	AtariTiles []BaseTile  // 听牌的牌

	// 分析工具
	counter *ScoreCounter // 分数计算器
}

// NewPlayer 创建一个新的玩家
func NewPlayer(wind Wind, oya bool) *Player {
	return &Player{
		Wind:       wind,
		Oya:        oya,
		Menzen:     true,
		Score:      25000,
		FirstRound: true,
		Ippatsu:    false,
		Hand:       make([]*Tile, 0, 14),
		CallGroups: make([]CallGroup, 0),
		AtariTiles: make([]BaseTile, 0),
		counter:    &ScoreCounter{},
	}
}

// IsRiichi 判断是否处于立直状态
func (p *Player) IsRiichi() bool {
	return p.Riichi || p.DoubleRiichi
}

// IsFuriten 判断是否处于振听状态
func (p *Player) IsFuriten() bool {
	return p.FuritenRound || p.FuritenRiver || p.FuritenRiichi
}

// GetFuuros 获取所有的鸣牌组
func (p *Player) GetFuuros() []CallGroup {
	return p.CallGroups
}

// IsMenzen 判断是否为门前清
func (p *Player) IsMenzen() bool {
	return p.Menzen
}

// IsTenpai 判断是否处于听牌状态
func (p *Player) IsTenpai() bool {
	return len(p.AtariTiles) > 0
}

// GetRiver 获取河
func (p *Player) GetRiver() *River {
	return &p.River
}

// HandToString 将手中的牌转换为字符串
func (p *Player) HandToString() string {
	sb := strings.Builder{}
	tiles := make([]string, 0, len(p.Hand))
	for _, tile := range p.Hand {
		tiles = append(tiles, tile.String())
	}
	sort.Strings(tiles)
	for i, tile := range tiles {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(tile)
	}
	return sb.String()
}

// RiverToString 将河转换为字符串
func (p *Player) RiverToString() string {
	return p.River.String()
}

// String 返回玩家的字符串表示
func (p *Player) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("风向: %d\n", p.Wind))
	sb.WriteString(fmt.Sprintf("点数: %d\n", p.Score))
	sb.WriteString(fmt.Sprintf("手牌: %s\n", p.HandToString()))
	sb.WriteString(fmt.Sprintf("河: %s\n", p.RiverToString()))
	sb.WriteString(fmt.Sprintf("门前清: %v\n", p.Menzen))
	sb.WriteString(fmt.Sprintf("立直: %v\n", p.IsRiichi()))
	sb.WriteString(fmt.Sprintf("听牌: %v\n", p.IsTenpai()))
	return sb.String()
}

// TenpaiToString 返回听牌的字符串表示
func (p *Player) TenpaiToString() string {
	sb := strings.Builder{}
	for i, tile := range p.AtariTiles {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(BaseTileToString(tile))
	}
	return sb.String()
}

// UpdateAtariTiles 更新听牌列表
func (p *Player) UpdateAtariTiles() {
	p.AtariTiles = p.AtariTiles[:0]

	// 获取手中牌的BaseTile列表
	baseTiles := ConvertTilesToBaseTiles(p.Hand)

	// 对于每一张可能的牌，检查是否能胡
	for tile := BaseTile(0); tile < 34; tile++ {
		testTiles := make([]BaseTile, len(baseTiles)+1)
		copy(testTiles, baseTiles)
		testTiles[len(testTiles)-1] = tile
		sort.Slice(testTiles, func(i, j int) bool { return testTiles[i] < testTiles[j] })

		// 如果能成形胡牌，则加入听牌列表
		if CanWinWithTiles(testTiles) {
			p.AtariTiles = append(p.AtariTiles, tile)
		}
	}
}

// UpdateFuritenRiver 更新河振听状态
func (p *Player) UpdateFuritenRiver() {
	// 检查是否弃过听牌的牌
	for _, atariTile := range p.AtariTiles {
		for _, riverTile := range p.River.River {
			if riverTile.Tile.Tile == atariTile && riverTile.Remain {
				p.FuritenRiver = true
				return
			}
		}
	}
}

// RemoveAtariTiles 移除某个特定的听牌
func (p *Player) RemoveAtariTiles(tile BaseTile) {
	for i := 0; i < len(p.AtariTiles); i++ {
		if p.AtariTiles[i] == tile {
			p.AtariTiles = append(p.AtariTiles[:i], p.AtariTiles[i+1:]...)
			i--
		}
	}
}

// Minogashi 当立直状态下过掉自己的听牌时调用
func (p *Player) Minogashi() {
	if p.IsRiichi() {
		p.FuritenRiichi = true
	} else {
		p.FuritenRound = true
	}
}

// GetKakan 获取可能的加杠列表
func (p *Player) GetKakan() []*SelfAction {
	actions := make([]*SelfAction, 0)
	for _, group := range p.CallGroups {
		if group.Type == Koutsu && len(group.Tiles) == 3 {
			if idx := FindInTiles(p.Hand, &Tile{Tile: group.Tiles[0]}); idx >= 0 {
				action := &SelfAction{Action{Action: KaKan, CorrespondTiles: []*Tile{p.Hand[idx]}}}
				actions = append(actions, action)
			}
		}
	}
	return actions
}

// GetAnkan 获取可能的暗杠列表
func (p *Player) GetAnkan() []*SelfAction {
	actions := make([]*SelfAction, 0)
	tileCount := make(map[BaseTile]int)
	for _, tile := range p.Hand {
		tileCount[tile.Tile]++
	}

	for baseTile, count := range tileCount {
		if count >= 4 {
			tiles := GetNCopies(p.Hand, baseTile, 4)
			if len(tiles) == 4 {
				action := &SelfAction{Action{Action: AnKan, CorrespondTiles: tiles}}
				actions = append(actions, action)
			}
		}
	}
	return actions
}

// GetDiscard 获取可能的弃牌列表
// 去重弃牌，避免重复弃牌选项
func (p *Player) GetDiscard(afterChipon bool) []*SelfAction {
	actions := make([]*SelfAction, 0)
	seen := make(map[BaseTile]bool)

	// 创建一个临时副本用于排序
	handTiles := make([]*Tile, len(p.Hand))
	copy(handTiles, p.Hand)

	// 排序手牌以保证一致的弃牌顺序
	sort.Slice(handTiles, func(i, j int) bool {
		if handTiles[i].Tile != handTiles[j].Tile {
			return handTiles[i].Tile < handTiles[j].Tile
		}
		return handTiles[i].RedDora && !handTiles[j].RedDora
	})

	// 遍历排序后的手牌，为每种牌添加一个弃牌选项
	for _, tile := range handTiles {
		if !seen[tile.Tile] {
			action := &SelfAction{Action{Action: Discard, CorrespondTiles: []*Tile{tile}}}
			actions = append(actions, action)
			seen[tile.Tile] = true
		}
	}

	return actions
}

// GetTsumo 获取自摸胡牌的选项
// 只有听牌时才能自摸，并且需要有役
func (p *Player) GetTsumo(table *Table) []*SelfAction {
	if p.IsTenpai() {
		counter := &ScoreCounter{}
		baseTiles := ConvertTilesToBaseTiles(p.Hand)
		isSevenPair := IsSevenPairPattern(baseTiles)
		result := counter.CalculateScore(table, p, baseTiles, p.CallGroups, baseTiles[len(baseTiles)-1], isSevenPair)
		if result != nil {
			action := &SelfAction{Action{Action: Tsumo, CorrespondTiles: []*Tile{}}}
			return []*SelfAction{action}
		}
	}
	return nil
}

// GetRiichi 获取立直的选项
// 立直后自动生成可弃牌的列表
func (p *Player) GetRiichi() []*SelfAction {
	// 与 C++ 实现一致：立直必须是门清并且尚未立直，返回针对每个可弃牌的立直动作
	if p.IsRiichi() || !p.IsMenzen() {
		return nil
	}

	// 只有处于听牌状态才可立直
	if !p.IsTenpai() {
		return nil
	}

	actions := make([]*SelfAction, 0)
	// 利用已有的弃牌选项生成立直对应的弃牌（去重）
	discards := p.GetDiscard(false)
	seen := make(map[BaseTile]bool)
	for _, d := range discards {
		if len(d.CorrespondTiles) == 0 {
			continue
		}
		tile := d.CorrespondTiles[0]
		if seen[tile.Tile] {
			continue
		}
		seen[tile.Tile] = true
		action := &SelfAction{Action{Action: Riichi, CorrespondTiles: []*Tile{tile}}}
		actions = append(actions, action)
	}
	if len(actions) == 0 {
		return nil
	}
	return actions
}

// GetKyushukyuhai 获取九种九牌流局的选项
func (p *Player) GetKyushukyuhai() []*SelfAction {
	yaochuTiles := make(map[BaseTile]bool)
	for _, tile := range p.Hand {
		if IsYaochuhai(tile.Tile) {
			yaochuTiles[tile.Tile] = true
		}
	}

	if len(yaochuTiles) >= 9 {
		action := &SelfAction{Action{Action: Kyushukyuhai, CorrespondTiles: p.Hand}}
		return []*SelfAction{action}
	}
	return nil
}

// GetRon 生成荣和行动
func (p *Player) GetRon(table *Table, tile *Tile) []*ResponseAction {
	if !p.IsFuriten() {
		action := &ResponseAction{Action{Action: Ron, CorrespondTiles: []*Tile{tile}}}
		return []*ResponseAction{action}
	}
	return nil
}

// GetChi 生成吃的行动
func (p *Player) GetChi(tile *Tile) []*ResponseAction {
	actions := make([]*ResponseAction, 0)
	baseTile := tile.Tile
	if baseTile >= _1m && baseTile <= _7m {
		action := &ResponseAction{Action{Action: Chi, CorrespondTiles: []*Tile{tile}}}
		actions = append(actions, action)
	}
	return actions
}

// GetPon 生成碰的行动
func (p *Player) GetPon(tile *Tile) []*ResponseAction {
	count := 0
	for _, handTile := range p.Hand {
		if handTile.Tile == tile.Tile {
			count++
		}
	}
	if count >= 2 {
		action := &ResponseAction{Action{Action: Pon, CorrespondTiles: []*Tile{tile}}}
		return []*ResponseAction{action}
	}
	return nil
}

// GetKan 生成大明杠的行动
func (p *Player) GetKan(tile *Tile) []*ResponseAction {
	count := 0
	for _, handTile := range p.Hand {
		if handTile.Tile == tile.Tile {
			count++
		}
	}
	if count >= 3 {
		action := &ResponseAction{Action{Action: Kan, CorrespondTiles: []*Tile{tile}}}
		return []*ResponseAction{action}
	}
	return nil
}

// GetChanAnkan 生成抢暗杠的行动
func (p *Player) GetChanAnkan(tile *Tile) []*ResponseAction {
	if !p.IsFuriten() {
		action := &ResponseAction{Action{Action: ChanAnKan, CorrespondTiles: []*Tile{tile}}}
		return []*ResponseAction{action}
	}
	return nil
}

// GetChankan 生成抢杠的行动
func (p *Player) GetChankan(tile *Tile) []*ResponseAction {
	if !p.IsFuriten() {
		action := &ResponseAction{Action{Action: ChanKan, CorrespondTiles: []*Tile{tile}}}
		return []*ResponseAction{action}
	}
	return nil
}

// RemoveFromHand 从手中移除一张牌
func (p *Player) RemoveFromHand(tile *Tile) {
	for i, t := range p.Hand {
		if t.ID == tile.ID {
			p.Hand = append(p.Hand[:i], p.Hand[i+1:]...)
			return
		}
	}
}

// ExecuteAnkan 执行暗杠
func (p *Player) ExecuteAnkan(tile BaseTile) {
	tiles := GetNCopies(p.Hand, tile, 4)
	for _, t := range tiles {
		p.RemoveFromHand(t)
	}
	p.CallGroups = append(p.CallGroups, CallGroup{Type: Kantsu, Tiles: []BaseTile{tile, tile, tile, tile}, IsOpen: false})
}

// ExecuteKakan 执行加杠
func (p *Player) ExecuteKakan(tile *Tile) {
	p.RemoveFromHand(tile)
	for i := range p.CallGroups {
		if p.CallGroups[i].Type == Koutsu && len(p.CallGroups[i].Tiles) == 3 && p.CallGroups[i].Tiles[0] == tile.Tile {
			p.CallGroups[i] = CallGroup{Type: Kantsu, Tiles: []BaseTile{tile.Tile, tile.Tile, tile.Tile, tile.Tile}, IsOpen: true}
			break
		}
	}
}

// CanWinWithTiles 判断给定的牌是否能胡牌
func CanWinWithTiles(tiles []BaseTile) bool {
	if len(tiles) != 14 {
		return false
	}
	splitter := GetTileSplitter()
	completed := splitter.GetAllCompletedTiles(tiles)
	return len(completed) > 0
}

// SortHand 对手牌进行排序
// 按照牌的类型排序，相同牌号则赤宝牌排在后面
func (p *Player) SortHand() {
	sort.Slice(p.Hand, func(i, j int) bool {
		if p.Hand[i].Tile != p.Hand[j].Tile {
			return p.Hand[i].Tile < p.Hand[j].Tile
		}
		// 相同牌号时，红宝排在后面
		return !p.Hand[i].RedDora && p.Hand[j].RedDora
	})
}

// RiichiGetAnkan 立直后的暗杠选项
// 立直后只能进行特定的暗杠（新摸到的牌同花色）
func (p *Player) RiichiGetAnkan() []*SelfAction {
	actions := make([]*SelfAction, 0)
	if !p.IsRiichi() {
		return actions
	}

	// 立直后只能杠新摸到的牌
	if len(p.Hand) > 0 {
		lastTile := p.Hand[len(p.Hand)-1]
		count := 0
		for _, tile := range p.Hand {
			if tile.Tile == lastTile.Tile {
				count++
			}
		}
		if count >= 4 {
			tiles := GetNCopies(p.Hand, lastTile.Tile, 4)
			if len(tiles) == 4 {
				action := &SelfAction{Action{Action: AnKan, CorrespondTiles: tiles}}
				actions = append(actions, action)
			}
		}
	}
	return actions
}

// RiichiGetDiscard 立直后的弃牌选项
// 立直后只能从新摸到的牌中选择弃牌
func (p *Player) RiichiGetDiscard() []*SelfAction {
	actions := make([]*SelfAction, 0)
	if !p.IsRiichi() {
		return actions
	}

	// 立直后只能弃新摸到的牌
	if len(p.Hand) > 0 {
		lastTile := p.Hand[len(p.Hand)-1]
		action := &SelfAction{Action{Action: Discard, CorrespondTiles: []*Tile{lastTile}}}
		actions = append(actions, action)
	}
	return actions
}
