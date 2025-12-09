package mahjong

import (
	"sort"
	"strings"
)

// TileGroupType 表示牌组的类型
type TileGroupType int

const (
	Toitsu  TileGroupType = iota // 对子
	Shuntsu                      // 顺子
	Koutsu                       // 刻子
	Kantsu                       // 杠子
)

// TileGroup 表示一组牌（对子、顺子、刻子、杠子）
type TileGroup struct {
	Type  TileGroupType // 组类型
	Tiles []BaseTile    // 牌组内的牌
}

// String 返回TileGroup的字符串表示
func (tg *TileGroup) String() string {
	sb := strings.Builder{}
	sb.WriteString("[")
	for _, tile := range tg.Tiles {
		sb.WriteString(BaseTileToString(tile))
	}
	sb.WriteString("]")
	return sb.String()
}

// Find 在TileGroup中查找某个牌
func (tg *TileGroup) Find(tile BaseTile) int {
	for i, t := range tg.Tiles {
		if t == tile {
			return i
		}
	}
	return -1
}

// SortTiles 对组内的牌进行排序
func (tg *TileGroup) SortTiles() {
	sort.Slice(tg.Tiles, func(i, j int) bool {
		return tg.Tiles[i] < tg.Tiles[j]
	})
}

// SetTiles 设置组内的牌并进行排序
func (tg *TileGroup) SetTiles(tiles []BaseTile) {
	tg.Tiles = make([]BaseTile, len(tiles))
	copy(tg.Tiles, tiles)
	tg.SortTiles()
}

// CompletedTiles 表示一个完成的牌型（包含雀头和4个组）
type CompletedTiles struct {
	Head TileGroup   // 雀头
	Body []TileGroup // 4个组
}

// String 返回CompletedTiles的字符串表示
func (ct *CompletedTiles) String() string {
	sb := strings.Builder{}
	sb.WriteString(ct.Head.String())
	sb.WriteString(" ")
	for _, group := range ct.Body {
		sb.WriteString(group.String())
		sb.WriteString(" ")
	}
	return sb.String()
}

// SortBody 对body中的组进行排序
func (ct *CompletedTiles) SortBody() {
	for i := range ct.Body {
		ct.Body[i].SortTiles()
	}
	sort.Slice(ct.Body, func(i, j int) bool {
		if ct.Body[i].Tiles[0] != ct.Body[j].Tiles[0] {
			return ct.Body[i].Tiles[0] < ct.Body[j].Tiles[0]
		}
		return ct.Body[i].Type < ct.Body[j].Type
	})
}

// TileSplitter 用于拆分牌组，找出所有可能的胡牌方式
type TileSplitter struct {
	completedTiles CompletedTiles   // 当前完成的牌型（雀头 + body）
	hasHead        bool             // 是否已经有雀头
	allCompleted   []CompletedTiles // 所有完成的牌型
}

// GetInstance 获取TileSplitter的单例
var tileSplitterInstance *TileSplitter

func init() {
	tileSplitterInstance = &TileSplitter{}
}

// GetTileSplitter 获取TileSplitter实例
func GetTileSplitter() *TileSplitter {
	return tileSplitterInstance
}

// Reset 重置拆分器状态
func (ts *TileSplitter) Reset() {
	ts.completedTiles = CompletedTiles{}
	ts.hasHead = false
	ts.allCompleted = nil
}

// GetAllCompletedTiles 获取所有可能的完成牌型
func (ts *TileSplitter) GetAllCompletedTiles(tiles []BaseTile) []CompletedTiles {
	ts.Reset()

	if len(tiles) == 0 {
		return []CompletedTiles{ts.completedTiles}
	}

	// 排序输入的牌
	sortedTiles := make([]BaseTile, len(tiles))
	copy(sortedTiles, tiles)
	sort.Slice(sortedTiles, func(i, j int) bool { return sortedTiles[i] < sortedTiles[j] })

	result := ts.getAllCompletedTilesRecursive(sortedTiles)
	// 对每个结果进行排序和去重将在调用方处理；此处直接返回
	return result
}

// getAllCompletedTilesRecursive 递归地获取所有完成的牌型
func (ts *TileSplitter) getAllCompletedTilesRecursive(tiles []BaseTile) []CompletedTiles {
	if len(tiles) == 0 {
		// 返回当前记录的 completedTiles
		return []CompletedTiles{ts.completedTiles}
	}

	result := []CompletedTiles{}
	processed := make(map[BaseTile]bool)

	for _, tile := range tiles {
		if processed[tile] {
			continue
		}
		processed[tile] = true

		// 1. 尝试作为对子（雀头）
		if !ts.hasHead && countInSlice(tiles, tile) >= 2 {
			tmpTiles := removeFromSlice(tiles, tile, 2)
			// 设置雀头
			ts.completedTiles.Head = TileGroup{Type: Toitsu, Tiles: []BaseTile{tile, tile}}
			ts.hasHead = true

			subResults := ts.getAllCompletedTilesRecursive(tmpTiles)
			result = append(result, subResults...)

			// 恢复状态
			ts.hasHead = false
			ts.completedTiles.Head = TileGroup{}
		}

		// 2. 尝试作为刻子
		if countInSlice(tiles, tile) >= 3 {
			tmpTiles := removeFromSlice(tiles, tile, 3)
			grp := TileGroup{Type: Koutsu, Tiles: []BaseTile{tile, tile, tile}}

			// 添加到 body
			ts.completedTiles.Body = append(ts.completedTiles.Body, grp)
			subResults := ts.getAllCompletedTilesRecursive(tmpTiles)
			result = append(result, subResults...)
			// 恢复 body
			ts.completedTiles.Body = ts.completedTiles.Body[:len(ts.completedTiles.Body)-1]
		}

		// 3. 尝试作为顺子
		if !isShuntsuBadHead(tile) && IsIn(tiles, tile+1) && IsIn(tiles, tile+2) {
			tmpTiles := removeFromSlice(tiles, tile, 1)
			tmpTiles = removeFromSlice(tmpTiles, tile+1, 1)
			tmpTiles = removeFromSlice(tmpTiles, tile+2, 1)

			grp := TileGroup{Type: Shuntsu, Tiles: []BaseTile{tile, tile + 1, tile + 2}}
			ts.completedTiles.Body = append(ts.completedTiles.Body, grp)
			subResults := ts.getAllCompletedTilesRecursive(tmpTiles)
			result = append(result, subResults...)
			ts.completedTiles.Body = ts.completedTiles.Body[:len(ts.completedTiles.Body)-1]
		}
	}

	return result
}

// isShuntsuBadHead 检查是否为不能作为顺子开头的牌
func isShuntsuBadHead(tile BaseTile) bool {
	badHeads := []BaseTile{_8m, _9m, _8p, _9p, _8s, _9s, _1z, _2z, _3z, _4z, _5z, _6z, _7z}
	for _, h := range badHeads {
		if tile == h {
			return true
		}
	}
	return false
}

// countInSlice 计算某个元素在切片中出现的次数
func countInSlice(tiles []BaseTile, tile BaseTile) int {
	count := 0
	for _, t := range tiles {
		if t == tile {
			count++
		}
	}
	return count
}

// removeFromSlice 从切片中移除指定个数的元素
func removeFromSlice(tiles []BaseTile, tile BaseTile, count int) []BaseTile {
	result := make([]BaseTile, 0, len(tiles))
	removed := 0
	for _, t := range tiles {
		if t == tile && removed < count {
			removed++
			continue
		}
		result = append(result, t)
	}
	return result
}

// CheckColorCount 检查颜色数量的有效性（用于初步筛选胡牌型）
func CheckColorCount(tiles []BaseTile) bool {
	// 统计万、筒、索三色的数量以及字牌（1z..7z）各自的数量
	var suitCounts [3]int
	var honorCounts [7]int
	for _, t := range tiles {
		switch {
		case t <= _9m:
			suitCounts[0]++
		case t <= _9p:
			suitCounts[1]++
		case t <= _9s:
			suitCounts[2]++
		default:
			idx := int(t - _1z)
			if idx >= 0 && idx < 7 {
				honorCounts[idx]++
			}
		}
	}

	hasHead := false
	// 检查三色总数和字牌的模3关系
	for i := 0; i < 3; i++ {
		if suitCounts[i]%3 == 1 {
			return false
		}
		if suitCounts[i]%3 == 2 {
			if hasHead {
				return false
			}
			hasHead = true
		}
	}
	for i := 0; i < 7; i++ {
		if honorCounts[i]%3 == 1 {
			return false
		}
		if honorCounts[i]%3 == 2 {
			if hasHead {
				return false
			}
			hasHead = true
		}
	}
	return true
}

// NewTileGroup 创建一个新的TileGroup
func NewTileGroup(groupType TileGroupType, tiles []BaseTile) *TileGroup {
	tg := &TileGroup{Type: groupType}
	tg.SetTiles(tiles)
	return tg
}

// NewCompletedTiles 创建一个新的CompletedTiles
func NewCompletedTiles(head *TileGroup, body []*TileGroup) *CompletedTiles {
	ct := &CompletedTiles{
		Head: *head,
		Body: make([]TileGroup, len(body)),
	}
	for i, group := range body {
		ct.Body[i] = *group
	}
	ct.SortBody()
	return ct
}

// CallGroup 表示鸣牌组（吃、碰、杠）
type CallGroup struct {
	Type   TileGroupType // 组类型
	Tiles  []BaseTile    // 鸣牌的牌
	IsOpen bool          // 是否为明牌
}

// String 返回CallGroup的字符串表示
func (cg *CallGroup) String() string {
	sb := strings.Builder{}
	typeStr := ""
	switch cg.Type {
	case Shuntsu:
		typeStr = "吃"
	case Koutsu:
		typeStr = "碰"
	case Kantsu:
		typeStr = "杠"
	}
	sb.WriteString(typeStr)
	sb.WriteString("[")
	for _, tile := range cg.Tiles {
		sb.WriteString(BaseTileToString(tile))
	}
	sb.WriteString("]")
	return sb.String()
}
