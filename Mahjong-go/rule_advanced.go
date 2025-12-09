package mahjong

import (
	"sort"
)

// TileShape 表示牌的形状/特征
type TileShape int

const (
	ShuntsuShape TileShape = iota // 顺子形
	KoutsuShape                   // 刻子形
	ToitsuShape                   // 对子形
	KantsuShape                   // 杠子形
)

// AnalyzeHandPattern 分析手牌的模式
func AnalyzeHandPattern(tiles []BaseTile) string {
	// 返回手牌模式的字符串描述
	// 例如: "顺子 顺子 刻子 对子" 等
	if len(tiles) != 14 {
		return ""
	}

	pattern := ""
	// 需要拆分成完成形和对子
	// 这是简化版本

	return pattern
}

// HasCompletedTiles 检查是否存在完成的牌型
func HasCompletedTiles(tiles []BaseTile) bool {
	if len(tiles) != 14 {
		return false
	}

	splitter := &TileSplitter{}
	completed := splitter.GetAllCompletedTiles(tiles)
	return len(completed) > 0
}

// FindAllWinTiles 找到所有可以和牌的牌
func FindAllWinTiles(tiles []BaseTile) []BaseTile {
	winTiles := make([]BaseTile, 0)

	for i := 0; i < 34; i++ {
		testTile := BaseTile(i)

		// 尝试加入这张牌
		testTiles := append(tiles, testTile)
		if len(testTiles) == 14 {
			splitter := &TileSplitter{}
			completed := splitter.GetAllCompletedTiles(testTiles)
			if len(completed) > 0 {
				winTiles = append(winTiles, testTile)
			}
		}
	}

	return winTiles
}

// GetTenpaiPattern 获取听牌的模式
func GetTenpaiPattern(tiles []BaseTile) []BaseTile {
	tenpaiTiles := make([]BaseTile, 0)

	if len(tiles) != 13 {
		return tenpaiTiles
	}

	// 对每张可能的牌，检查是否能和
	for i := 0; i < 34; i++ {
		testTile := BaseTile(i)
		testTiles := append(tiles, testTile)

		splitter := &TileSplitter{}
		completed := splitter.GetAllCompletedTiles(testTiles)
		if len(completed) > 0 {
			tenpaiTiles = append(tenpaiTiles, testTile)
		}
	}

	return tenpaiTiles
}

// IsTenpaiPattern 检查是否是听牌格子
func IsTenpaiPattern(tiles []BaseTile) bool {
	return len(GetTenpaiPattern(tiles)) > 0
}

// CountTiles 统计某种牌的数量
func CountTiles(tiles []BaseTile, target BaseTile) int {
	count := 0
	for _, tile := range tiles {
		if tile == target {
			count++
		}
	}
	return count
}

// HasYaochuTile 检查是否包含幺九牌
func HasYaochuTile(tiles []BaseTile) bool {
	for _, tile := range tiles {
		if Is1hai(tile) || Is9hai(tile) {
			return true
		}
	}
	return false
}

// AllYaochuTiles 检查是否全是幺九牌
func AllYaochuTiles(tiles []BaseTile) bool {
	for _, tile := range tiles {
		if !Is1hai(tile) && !Is9hai(tile) {
			return false
		}
	}
	return len(tiles) > 0
}

// GetSuitTiles 获取特定花色的牌
func GetSuitTiles(tiles []BaseTile, suit int) []BaseTile {
	suitTiles := make([]BaseTile, 0)
	for _, tile := range tiles {
		if int(tile)/9 == suit {
			suitTiles = append(suitTiles, tile)
		}
	}
	return suitTiles
}

// GetZTiles 获取字牌
func GetZTiles(tiles []BaseTile) []BaseTile {
	return GetSuitTiles(tiles, 3)
}

// GetMTiles 获取万牌
func GetMTiles(tiles []BaseTile) []BaseTile {
	return GetSuitTiles(tiles, 0)
}

// GetPTiles 获取饼牌
func GetPTiles(tiles []BaseTile) []BaseTile {
	return GetSuitTiles(tiles, 1)
}

// GetSTiles 获取索牌
func GetSTiles(tiles []BaseTile) []BaseTile {
	return GetSuitTiles(tiles, 2)
}

// SortTiles 对牌进行排序
func SortTiles(tiles []BaseTile) []BaseTile {
	sorted := make([]BaseTile, len(tiles))
	copy(sorted, tiles)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})
	return sorted
}

// RemoveTile 从牌列表中移除一张牌
func RemoveTile(tiles []BaseTile, target BaseTile) []BaseTile {
	result := make([]BaseTile, 0, len(tiles))
	found := false

	for _, tile := range tiles {
		if tile == target && !found {
			found = true
			continue
		}
		result = append(result, tile)
	}

	return result
}

// RemoveTiles 从牌列表中移除多张牌
func RemoveTiles(tiles []BaseTile, toRemove []BaseTile) []BaseTile {
	result := tiles
	for _, tile := range toRemove {
		result = RemoveTile(result, tile)
	}
	return result
}

// GetSequenceTiles 获取连续的牌序列
func GetSequenceTiles(start BaseTile, length int) []BaseTile {
	sequence := make([]BaseTile, length)
	for i := 0; i < length; i++ {
		sequence[i] = BaseTile(int(start) + i)
	}
	return sequence
}

// CheckSequenceExists 检查是否存在特定的顺子
func CheckSequenceExists(tiles []BaseTile, sequence []BaseTile) bool {
	for _, requiredTile := range sequence {
		found := false
		for _, tile := range tiles {
			if tile == requiredTile {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// GetMissingTilesForSequence 获取缺少的牌以组成顺子
func GetMissingTilesForSequence(tiles []BaseTile, sequence []BaseTile) []BaseTile {
	missing := make([]BaseTile, 0)

	for _, requiredTile := range sequence {
		found := false
		for _, tile := range tiles {
			if tile == requiredTile {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, requiredTile)
		}
	}

	return missing
}

// CanFormShuntsu 检查是否能组成顺子
func CanFormShuntsu(tiles []BaseTile, start BaseTile) bool {
	sequence := GetSequenceTiles(start, 3)
	return CheckSequenceExists(tiles, sequence)
}

// CanFormKoutsu 检查是否能组成刻子
func CanFormKoutsu(tiles []BaseTile, tile BaseTile) bool {
	count := CountTiles(tiles, tile)
	return count >= 3
}

// CanFormToitsu 检查是否能组成对子
func CanFormToitsu(tiles []BaseTile, tile BaseTile) bool {
	count := CountTiles(tiles, tile)
	return count >= 2
}

// GetPossibleShuntsuStarts 获取所有可能的顺子起始位置
func GetPossibleShuntsuStarts(tiles []BaseTile) []BaseTile {
	starts := make([]BaseTile, 0)

	// 对每个花色，检查可能的顺子起点（1-7，因为8-9也能开始但会越界）
	for suit := 0; suit < 3; suit++ { // 字牌不能组成顺子
		for num := 0; num < 7; num++ { // 0代表1，6代表7
			start := BaseTile(suit*9 + num)
			if CanFormShuntsu(tiles, start) {
				starts = append(starts, start)
			}
		}
	}

	return starts
}

// GetPossibleKoutsuTiles 获取所有可能的刻子牌
func GetPossibleKoutsuTiles(tiles []BaseTile) []BaseTile {
	koutsuTiles := make([]BaseTile, 0)
	seen := make(map[BaseTile]bool)

	for _, tile := range tiles {
		if !seen[tile] {
			if CanFormKoutsu(tiles, tile) {
				koutsuTiles = append(koutsuTiles, tile)
				seen[tile] = true
			}
		}
	}

	return koutsuTiles
}

// GetPossibleToitsuTiles 获取所有可能的对子牌
func GetPossibleToitsuTiles(tiles []BaseTile) []BaseTile {
	toitsuTiles := make([]BaseTile, 0)
	seen := make(map[BaseTile]bool)

	for _, tile := range tiles {
		if !seen[tile] {
			if CanFormToitsu(tiles, tile) {
				toitsuTiles = append(toitsuTiles, tile)
				seen[tile] = true
			}
		}
	}

	return toitsuTiles
}

// AnalyzeRoughPattern 分析粗略的手牌模式
func AnalyzeRoughPattern(tiles []BaseTile) map[string]int {
	pattern := make(map[string]int)

	// 统计各花色的牌数
	pattern["m"] = len(GetMTiles(tiles))
	pattern["p"] = len(GetPTiles(tiles))
	pattern["s"] = len(GetSTiles(tiles))
	pattern["z"] = len(GetZTiles(tiles))

	// 统计可能的面子
	pattern["shuntsu"] = len(GetPossibleShuntsuStarts(tiles))
	pattern["koutsu"] = len(GetPossibleKoutsuTiles(tiles))

	return pattern
}

// CheckRichiiCondition 检查立直条件
func CheckRichiiCondition(tiles []BaseTile) bool {
	// 需要13张牌，能听牌
	if len(tiles) != 13 {
		return false
	}

	return IsTenpaiPattern(tiles)
}

// GenerateTileCombination 生成牌的组合
func GenerateTileCombination(tiles []BaseTile, n int) [][]BaseTile {
	var result [][]BaseTile

	if n == 0 {
		return [][]BaseTile{{}}
	}

	if len(tiles) == 0 {
		return result
	}

	// 递归生成组合
	first := tiles[0]
	rest := tiles[1:]

	// 包含第一个元素的组合
	for _, comb := range GenerateTileCombination(rest, n-1) {
		result = append(result, append([]BaseTile{first}, comb...))
	}

	// 不包含第一个元素的组合
	result = append(result, GenerateTileCombination(rest, n)...)

	return result
}

// GetComplimentaryTiles 获取补充牌（使13张牌变成14张能和的牌）
func GetComplimentaryTiles(tiles []BaseTile) []BaseTile {
	if len(tiles) != 13 {
		return []BaseTile{}
	}

	return FindAllWinTiles(tiles)
}

// EstimateHandSafety 估计手牌安全性（简化版）
func EstimateHandSafety(tiles []BaseTile, riverTiles []BaseTile) int {
	// 0 = 危险, 1 = 中等, 2 = 安全

	// 计数被弃出的幺九牌
	discardsYaochu := 0
	for _, tile := range riverTiles {
		if Is1hai(tile) || Is9hai(tile) {
			discardsYaochu++
		}
	}

	// 如果弃出了很多幺九牌，幺九牌相对较安全
	if discardsYaochu > 6 {
		return 2
	}

	return 1
}

// GetAdvancedTiles 获取先进牌（可能让对手听牌的牌）
func GetAdvancedTiles(tiles []BaseTile, opponentTenpaiTiles []BaseTile) []BaseTile {
	advancedTiles := make([]BaseTile, 0)

	for _, tile := range tiles {
		// 如果弃这张牌会被对手和牌
		for _, opponentTile := range opponentTenpaiTiles {
			if tile == opponentTile {
				advancedTiles = append(advancedTiles, tile)
				break
			}
		}
	}

	return advancedTiles
}

// CalculateRoughWaitDistance 计算粗略的等待距离
func CalculateRoughWaitDistance(tiles []BaseTile) int {
	// 返回需要多少张牌才能听牌
	// 0 = 已听牌
	// 1 = 一张牌的距离
	// 2 = 两张牌的距离
	// 等等

	if IsTenpaiPattern(tiles) {
		return 0
	}

	// 简化估算：基于手牌的组织程度
	// 这是非常粗略的估算

	possibleShuntsu := len(GetPossibleShuntsuStarts(tiles))
	possibleKoutsu := len(GetPossibleKoutsuTiles(tiles))

	if possibleShuntsu > 1 && possibleKoutsu > 0 {
		return 1
	}

	return 2
}

// SimulateDiscard 模拟弃牌
func SimulateDiscard(tiles []BaseTile, discardTile BaseTile) []BaseTile {
	return RemoveTile(tiles, discardTile)
}

// SimulateDraw 模拟摸牌
func SimulateDraw(tiles []BaseTile, drawTile BaseTile) []BaseTile {
	return append(tiles, drawTile)
}

// EvaluateHandValue 评估手牌价值（简化版）
func EvaluateHandValue(tiles []BaseTile) int {
	// 返回手牌的评估值
	// 基于各种因素：听牌距离、和牌概率等

	value := 0

	// 加分：可能的面子
	value += len(GetPossibleShuntsuStarts(tiles))
	value += len(GetPossibleKoutsuTiles(tiles))

	// 如果已听牌，大幅加分
	if IsTenpaiPattern(tiles) {
		value += 50
	}

	return value
}
