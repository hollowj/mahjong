package mahjong

import (
	"fmt"
	"sort"
)

// BaseAction 表示基础行动类型
type BaseAction uint8

const (
	// 响应行动（回应他人）
	Pass      BaseAction = iota
	Chi                  // 吃
	Pon                  // 碰
	Kan                  // 大明杠
	Ron                  // 荣和
	ChanAnKan            // 抢暗杠
	ChanKan              // 抢杠

	// 自身行动
	AnKan        // 暗杠
	KaKan        // 加杠
	Discard      // 弃牌
	Ippatsu      // 一发
	Tsumo        // 自摸
	Kyushukyuhai // 九种九牌流局
	Riichi       // 立直
)

// Action 表示一个行动及其对应的牌
type Action struct {
	Action          BaseAction // 行动类型
	CorrespondTiles []*Tile    // 对应的牌（已排序）
}

// GetAction 获取行动类型
func (a *Action) GetAction() BaseAction {
	return a.Action
}

// GetCorrespondTiles 获取对应的牌
func (a *Action) GetCorrespondTiles() []*Tile {
	return a.CorrespondTiles
}

// String 返回行动的字符串表示
func (a *Action) String() string {
	switch a.Action {
	case Pass:
		return "Pass"
	case Chi:
		return fmt.Sprintf("Chi %s %s", a.CorrespondTiles[0].String(), a.CorrespondTiles[1].String())
	case Pon:
		return fmt.Sprintf("Pon %s %s", a.CorrespondTiles[0].String(), a.CorrespondTiles[1].String())
	case Kan:
		if len(a.CorrespondTiles) >= 3 {
			return fmt.Sprintf("Kan %s %s %s", a.CorrespondTiles[0].String(),
				a.CorrespondTiles[1].String(), a.CorrespondTiles[2].String())
		}
		return "Kan"
	case Ron:
		return fmt.Sprintf("Ron %s", a.CorrespondTiles[0].String())
	case Discard:
		return fmt.Sprintf("Discard %s", a.CorrespondTiles[0].String())
	case Tsumo:
		return "Tsumo"
	case Riichi:
		return "Riichi"
	case AnKan:
		return fmt.Sprintf("AnKan %s", a.CorrespondTiles[0].String())
	case KaKan:
		return fmt.Sprintf("KaKan %s", a.CorrespondTiles[0].String())
	case Ippatsu:
		return "Ippatsu"
	case Kyushukyuhai:
		return "Kyushukyuhai"
	case ChanAnKan:
		return fmt.Sprintf("ChanAnKan %s", a.CorrespondTiles[0].String())
	case ChanKan:
		return fmt.Sprintf("ChanKan %s", a.CorrespondTiles[0].String())
	default:
		return "Unknown"
	}
}

// SelfAction 表示玩家的自主行动
type SelfAction struct {
	Action
}

// ResponseAction 表示对他人行动的响应
type ResponseAction struct {
	Action
}

// Compare 比较两个行动的优先级
func (a *Action) Compare(b *Action) int {
	if a.Action < b.Action {
		return -1
	} else if a.Action > b.Action {
		return 1
	}

	// 如果行动类型相同，比较对应的牌
	if len(a.CorrespondTiles) != len(b.CorrespondTiles) {
		return len(a.CorrespondTiles) - len(b.CorrespondTiles)
	}

	for i := 0; i < len(a.CorrespondTiles); i++ {
		if a.CorrespondTiles[i].ID < b.CorrespondTiles[i].ID {
			return -1
		} else if a.CorrespondTiles[i].ID > b.CorrespondTiles[i].ID {
			return 1
		}
	}

	return 0
}

// ActionSortByTile 按牌排序行动
type ActionSortByTile []*Action

func (a ActionSortByTile) Len() int {
	return len(a)
}

func (a ActionSortByTile) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ActionSortByTile) Less(i, j int) bool {
	return a[i].Compare(a[j]) < 0
}

// SortActions 对行动进行排序
func SortActions(actions []*Action) {
	sort.Sort(ActionSortByTile(actions))
}

// SelfActionSortByTile 按牌排序自主行动
type SelfActionSortByTile []*SelfAction

func (a SelfActionSortByTile) Len() int {
	return len(a)
}

func (a SelfActionSortByTile) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a SelfActionSortByTile) Less(i, j int) bool {
	return a[i].Compare(&a[j].Action) < 0
}

// SortSelfActions 对自主行动进行排序
func SortSelfActions(actions []*SelfAction) {
	sort.Sort(SelfActionSortByTile(actions))
}

// ResponseActionSortByTile 按牌排序响应行动
type ResponseActionSortByTile []*ResponseAction

func (a ResponseActionSortByTile) Len() int {
	return len(a)
}

func (a ResponseActionSortByTile) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ResponseActionSortByTile) Less(i, j int) bool {
	return a[i].Compare(&a[j].Action) < 0
}

// SortResponseActions 对响应行动进行排序
func SortResponseActions(actions []*ResponseAction) {
	sort.Sort(ResponseActionSortByTile(actions))
}

// GetActionIndex 获取行动的索引值
// 这个方法主要用于处理Red Dora的特殊情况
func (a *Action) GetActionIndex() int {
	// 如果没有对应的牌，返回行动类型的索引
	if len(a.CorrespondTiles) == 0 {
		return int(a.Action)
	}

	// 获取第一张牌的ID
	baseIndex := int(a.Action)*100 + int(a.CorrespondTiles[0].ID)

	// 处理赤宝牌的情况
	for _, tile := range a.CorrespondTiles {
		if tile.RedDora {
			baseIndex += 1000 // 赤宝牌的特殊标记
			break
		}
	}

	return baseIndex
}

// BaseActionToString 将BaseAction转换为字符串
func BaseActionToString(action BaseAction) string {
	switch action {
	case Pass:
		return "Pass"
	case Discard:
		return "弃牌"
	case Pon:
		return "碰"
	case Chi:
		return "吃"
	case Kan:
		return "杠"
	case Tsumo:
		return "自摸"
	case Ron:
		return "荣和"
	case Riichi:
		return "立直"
	case Ippatsu:
		return "一发"
	case AnKan:
		return "暗杠"
	case KaKan:
		return "加杠"
	case Kyushukyuhai:
		return "九种九牌"
	case ChanKan:
		return "抢杠"
	case ChanAnKan:
		return "抢暗杠"
	default:
		return "未知"
	}
}
