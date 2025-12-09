package mahjong

import (
	"fmt"
	"sort"
)

// Wind 表示风向：东、南、西、北
type Wind int

const (
	East Wind = iota
	South
	West
	North
)

// BaseTile 表示基础牌类型
// m表示万(1-9m), p表示筒(1-9p), s表示索(1-9s), z表示字牌(1-7z)
type BaseTile int

const (
	// 万牌 (Man - m)
	_1m BaseTile = iota
	_2m
	_3m
	_4m
	_5m
	_6m
	_7m
	_8m
	_9m
	// 筒牌 (Pin - p)
	_1p
	_2p
	_3p
	_4p
	_5p
	_6p
	_7p
	_8p
	_9p
	// 索牌 (Sou - s)
	_1s
	_2s
	_3s
	_4s
	_5s
	_6s
	_7s
	_8s
	_9s
	// 字牌 (Honors - z)
	_1z // 东
	_2z // 南
	_3z // 西
	_4z // 北
	_5z // 白
	_6z // 发
	_7z // 中
)

// Tile 表示一张实际的牌，包含基本牌类型、赤宝牌标记和ID
type Tile struct {
	Tile    BaseTile // 基础牌类型
	RedDora bool     // 是否为赤宝牌
	ID      int      // 牌的唯一ID
}

// String 返回牌的字符串表示
func (t *Tile) String() string {
	number := t.Tile%9 + 1
	if t.RedDora {
		number = 0
	}

	switch t.Tile / 9 {
	case 0:
		return fmt.Sprintf("%dm", number)
	case 1:
		return fmt.Sprintf("%dp", number)
	case 2:
		return fmt.Sprintf("%ds", number)
	case 3:
		return fmt.Sprintf("%dz", number)
	default:
		panic("Error Tile object")
	}
}

// BaseTileToString 将BaseTile转换为字符串
func BaseTileToString(bt BaseTile) string {
	names := [...]string{
		"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m",
		"1p", "2p", "3p", "4p", "5p", "6p", "7p", "8p", "9p",
		"1s", "2s", "3s", "4s", "5s", "6s", "7s", "8s", "9s",
		"1z", "2z", "3z", "4z", "5z", "6z", "7z",
	}
	if int(bt) < len(names) {
		return names[bt]
	}
	return "unknown"
}

// Char2ToBaseTile 将两个字符转换为BaseTile（例如'1','m' -> _1m）
func Char2ToBaseTile(number byte, color byte) (BaseTile, bool) {
	num := int(number - '0')
	redDora := false

	if num == 0 {
		redDora = true
		num = 5
	}

	switch color {
	case 'm':
		return BaseTile(_1m + BaseTile(num-1)), redDora
	case 'p':
		return BaseTile(_1p + BaseTile(num-1)), redDora
	case 's':
		return BaseTile(_1s + BaseTile(num-1)), redDora
	case 'z':
		return BaseTile(_1z + BaseTile(num-1)), redDora
	default:
		panic("Bad tile string")
	}
}

// GetDoraNext 获取宝牌的下一张牌
func GetDoraNext(tile BaseTile) BaseTile {
	switch tile {
	case _9m:
		return _1m
	case _9s:
		return _1s
	case _9p:
		return _1p
	case _4z:
		return _1z
	case _7z:
		return _5z
	default:
		return tile + 1
	}
}

// Is1hai 判断是否为1或9的边缘牌
func Is1hai(t BaseTile) bool {
	return t == _1m || t == _1p || t == _1s
}

// Is9hai 判断是否为9
func Is9hai(t BaseTile) bool {
	return t == _9m || t == _9p || t == _9s
}

// Is19hai 判断是否为1或9
func Is19hai(t BaseTile) bool {
	return Is1hai(t) || Is9hai(t)
}

// IsTsuhai 判断是否为字牌
func IsTsuhai(t BaseTile) bool {
	return t >= _1z && t <= _7z
}

// IsYaochuhai 判断是否为幺九牌
func IsYaochuhai(t BaseTile) bool {
	return Is19hai(t) || IsTsuhai(t)
}

// IsShuntsu 判断三张牌是否能组成顺子（1-2-3, 2-3-4等）
func IsShuntsu(tiles []BaseTile) bool {
	if len(tiles) != 3 {
		return false
	}

	// 复制并排序
	sorted := make([]BaseTile, 3)
	copy(sorted, tiles)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	// 检查是否连续
	if sorted[1]-sorted[0] != 1 || sorted[2]-sorted[1] != 1 {
		return false
	}

	// 检查顺子的最后一张不能是8m,9m,8p,9p,8s,9s,任何z
	if sorted[2] >= _1z || sorted[2] == _8m || sorted[2] == _9m ||
		sorted[2] == _8p || sorted[2] == _9p || sorted[2] == _8s || sorted[2] == _9s {
		return false
	}

	return true
}

// IsKoutsu 判断三张牌是否相同（刻子）
func IsKoutsu(tiles []BaseTile) bool {
	if len(tiles) != 3 {
		return false
	}
	return tiles[0] == tiles[1] && tiles[1] == tiles[2]
}

// IsKantsu 判断四张牌是否相同（杠子）
func IsKantsu(tiles []BaseTile) bool {
	if len(tiles) != 4 {
		return false
	}
	return tiles[0] == tiles[1] && tiles[1] == tiles[2] && tiles[2] == tiles[3]
}

// IsWindMatch 检查牌是否与特定风向匹配
func IsWindMatch(tile BaseTile, gameWind Wind) bool {
	switch gameWind {
	case East:
		return tile == _1z
	case South:
		return tile == _2z
	case West:
		return tile == _3z
	case North:
		return tile == _4z
	default:
		panic("Unknown wind")
	}
}

// Is567z 检查是否为5、6、7字牌
func Is567z(tile BaseTile) bool {
	return tile == _5z || tile == _6z || tile == _7z
}

// IsYakuhai 检查是否为役牌（风牌或三元牌）
func IsYakuhai(tile BaseTile, gameWind Wind, selfWind Wind) bool {
	if IsWindMatch(tile, gameWind) {
		return true
	}
	if IsWindMatch(tile, selfWind) {
		return true
	}
	if Is567z(tile) {
		return true
	}
	return false
}

// IsIn 检查元素是否在容器中
func IsIn(tiles []BaseTile, t BaseTile) bool {
	for _, tile := range tiles {
		if tile == t {
			return true
		}
	}
	return false
}

// CountTile 计算容器中某个牌出现的次数
func CountTile(tiles []BaseTile, t BaseTile) int {
	count := 0
	for _, tile := range tiles {
		if tile == t {
			count++
		}
	}
	return count
}

// IsSameContainer 检查两个容器是否大小相同且对应值相同
func IsSameContainer(a, b []BaseTile) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
