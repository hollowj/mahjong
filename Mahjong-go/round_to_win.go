package mahjong

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Syanten 计算向听(round to win)
// 向听数表示距离完成形还需要多少步
// 0向听 = 已经是完成形（可以和）
// 1向听 = 还差1步，摸到任何一张可以和的牌就能和
// 2向听 = 还差2步
// ...
// 8向听 = 完全无法和（特殊情况，如国士無双）
type Syanten struct {
	// syanten_map 存储向听数查表
	// key为编码的手牌（uint32），value为长度为4的数组，与 C++ 中 tuple<int,int,int,int> 对应
	// 值的含义与 C++ 一致，用于 normal 向听计算的选择
	syanten_map map[uint32][4]int
	is_loaded   bool
}

var syantenInstance *Syanten

// GetSyanten 获取向听计算单例
func GetSyanten() *Syanten {
	if syantenInstance == nil {
		syantenInstance = &Syanten{
			syanten_map: make(map[uint32][4]int),
			is_loaded:   false,
		}
		// 由于无法加载外部syanten.dat文件，这里提供了一个简化的实现
		// 完整的实现需要依赖syanten.dat文件
		syantenInstance.loadSyantenMap()
	}
	return syantenInstance
}

// loadSyantenMap 加载向听数查表
// 在完整版本中，这会从resource/syanten.dat文件加载
// 这里提供了一个基础实现，可以计算常见的向听情况
func (s *Syanten) loadSyantenMap() {
	// 尝试加载仓库中的 resource/syanten.dat 文件
	// 路径与 C++ 实现一致: "../resource/syanten.dat"
	path := "../resource/syanten.dat"
	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("open syanten.dat error: %v\n请将 'syanten.dat' 放到 %s", err, path))
	}
	defer file.Close()

	// tile_to_bit 与 C++ 中一致 (little endian，每个位置占3位)
	tileToBit := [9]uint32{
		0b000000000000000000000000001,
		0b000000000000000000000001000,
		0b000000000000000000001000000,
		0b000000000000000001000000000,
		0b000000000000001000000000000,
		0b000000000001000000000000000,
		0b000000001000000000000000000,
		0b000001000000000000000000000,
		0b001000000000000000000000000,
	}

	s.syanten_map = make(map[uint32][4]int)
	s.is_loaded = false

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		codeStr := fields[0]
		if len(codeStr) < 9 {
			// 忽略异常行
			continue
		}
		var key uint32 = 0
		for i := 0; i < 9; i++ {
			digit := int(codeStr[i] - '0')
			if digit < 0 {
				digit = 0
			}
			key += tileToBit[i] * uint32(digit)
		}

		// 解析随后的 4 个整数，如果不足则补零
		vals := [4]int{0, 0, 0, 0}
		for i := 1; i <= 4 && i < len(fields); i++ {
			v, err := strconv.Atoi(fields[i])
			if err != nil {
				v = 0
			}
			vals[i-1] = v
		}

		s.syanten_map[key] = vals
		count++
	}
	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("read syanten.dat error: %v", err))
	}

	// C++ 实现校验条目数为 405350
	if count != 405350 {
		panic(fmt.Sprintf("syanten.dat broken or incomplete: expected 405350, got %d", count))
	}
	s.is_loaded = true
}

// HandToCode 将手牌转换为编码表示
// 手牌被分为三类(万、筒、索)，每类用一个uint32表示
// 每个位置占3位(可表示0-7张牌)，共9个位置
func (s *Syanten) HandToCode(hand []*Tile) [4]uint32 {
	var code [4]uint32

	for _, tile := range hand {
		baseTile := tile.Tile

		// 计算分类和位置
		groupIdx := baseTile / 9 // 0: 万, 1: 筒, 2: 索, 3: 字
		tileIdx := baseTile % 9  // 0-8: 代表1-9

		// 每个牌位占3位(可表示0-7张)
		bitPos := tileIdx * 3

		// 将计数加1 (注意最多4张同样的牌)
		code[groupIdx] += (1 << bitPos)
	}

	return code
}

// CodeToString 将编码转换为字符串表示(调试用)
func (s *Syanten) CodeToString(code uint32) string {
	result := ""
	for i := 0; i < 9; i++ {
		result += string(rune('0' + ((code >> (i * 3)) & 0x7)))
	}
	return result
}

// CheckNormal 检查普通牌型的向听数
// handCode: 编码后的手牌
// nCallGroups: 副露的面子数(0-3)
func (s *Syanten) CheckNormal(handCode [4]uint32, nCallGroups int) int {
	// 如果没有加载查表，返回-1表示错误
	if !s.is_loaded || len(s.syanten_map) == 0 {
		return s.checkNormalWithoutMap(handCode, nCallGroups)
	}

	// 从查表中获取向听数
	ptm := 0 // 面子数
	ptt := 0 // 雀头数

	// 对前三类(万、筒、索)分别计算
	for j := 0; j < 3; j++ {
		if values, ok := s.syanten_map[handCode[j]]; ok {
			pt1m, pt1t := values[0], values[1]
			pt2m, pt2t := values[2], values[3]

			// 选择更优的选项（与 C++ 逻辑一致）
			if pt1m*2+pt1t >= pt2m*2+pt2t {
				ptm += pt1m
				ptt += pt1t
			} else {
				ptm += pt2m
				ptt += pt2t
			}
		}
	}

	// 处理字牌(7张)
	for i := 0; i < 7; i++ {
		num := (handCode[3] >> (i * 3)) & 0x7
		if num >= 3 {
			ptm++ // 面子
		} else if num >= 2 {
			ptt++ // 雀头
		}
	}

	// 调整面子和雀头
	// 最多可以有 4-nCallGroups 个面子+雀头
	maxFuShou := 4 - nCallGroups
	for ptm+ptt > maxFuShou && ptt > 0 {
		ptt--
	}
	for ptm+ptt > maxFuShou {
		ptm--
	}

	// 计算向听数: 9 - (面子数*2 + 雀头数) - 副露数*2
	return 9 - ptm*2 - ptt - nCallGroups*2
}

// checkNormalWithoutMap 在没有查表的情况下计算向听数(简化版本)
func (s *Syanten) checkNormalWithoutMap(handCode [4]uint32, nCallGroups int) int {
	// 简化实现: 直接计算当前牌型可能的面子+雀头数
	ptm := 0 // 面子数
	ptt := 0 // 雀头数

	// 对每一类牌(万、筒、索)进行计数
	for j := 0; j < 3; j++ {
		code := handCode[j]

		// 统计各个位置的牌数
		groupPtm := 0 // 这一类的面子数
		groupPtt := 0 // 这一类的雀头数

		for i := 0; i < 9; i++ {
			count := int((code >> (i * 3)) & 0x7)

			// 尽可能地组成面子和雀头
			if count >= 3 {
				groupPtm += count / 3
				count %= 3
			}
			if count >= 2 {
				groupPtt += count / 2
				count %= 2
			}
		}

		ptm += groupPtm
		if groupPtt > 0 {
			ptt = groupPtt // 只能有一个雀头
		}
	}

	// 处理字牌
	for i := 0; i < 7; i++ {
		num := (handCode[3] >> (i * 3)) & 0x7
		if num >= 3 {
			ptm++
		} else if num >= 2 {
			ptt = 1 // 可以作为雀头
		}
	}

	// 调整
	maxFuShou := 4 - nCallGroups
	for ptm+ptt > maxFuShou && ptt > 0 {
		ptt--
	}
	for ptm+ptt > maxFuShou {
		ptm--
	}

	return 9 - ptm*2 - ptt - nCallGroups*2
}

// NormalRoundToWin 计算普通牌型的向听数
func (s *Syanten) NormalRoundToWin(hand []*Tile, nCallGroups int) int {
	code := s.HandToCode(hand)
	return s.CheckNormal(code, nCallGroups)
}

// CalculateRoundToWin 计算手牌距离和牌还有几步(向听数)
func CalculateRoundToWin(hand []*Tile, callGroupCount int) int {
	syanten := GetSyanten()
	return syanten.NormalRoundToWin(hand, callGroupCount)
}
