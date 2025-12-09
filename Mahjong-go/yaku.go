package mahjong

import (
	"fmt"
)

// Yaku 表示一种胡牌役型
type Yaku uint16

const (
	// 一般役（1-2番）
	Tanyao           Yaku = iota // 断幺
	Pinfu                        // 平和
	Iipeikou                     // 一对对
	Ryanpeikou                   // 二对对
	Sanshokusequence             // 三色同顺
	Ittsu                        // 一通贯
	Chanta                       // 全带幺
	Honchanta                    // 混全带幺
	Sanshokudoukou               // 三色同刻
	Toitoi                       // 对对和
	Sanankou                     // 三暗刻
	Sangantu                     // 三杆子
	Sanputsu                     // 三副露

	// 役牌（1番）
	Yakuhai           // 役牌（通用）
	YakuhaiWind       // 役牌（风）
	YakuhaiWhiteBoard // 役牌（白板）
	YakuhaiGreenBoard // 役牌（绿板）
	YakuhaiRedBoard   // 役牌（红板）

	// 混合役（2-3番）
	Honitsu  // 混一色
	Junchan  // 纯全带幺
	Chinitsu // 清一色

	// 役满（13番）
	Kokushi    // 国士无双
	Suankou    // 四暗刻
	Daisangen  // 大三元
	Shosuushi  // 小四喜
	Daisuushi  // 大四喜
	Tsuisou    // 字一色
	Ryuuisou   // 绿一色
	Chinroutou // 清老头
	Honroutou  // 混老头
	Chiitoitsu // 七对子
	Kazoe      // 数役满
	Tenhou     // 天胡
	Chihou     // 地胡
	Churen     // 九莲宝灯

	// 特殊役
	Menzentsumo  // 门清自摸
	Dabururiichi // 双立直

	// 最大值
	MaxYaku
)

// YakuInfo 存储役的信息
type YakuInfo struct {
	Name      string // 役名
	Fan       int    // 番数
	IsOpen    bool   // 是否可以鸣牌
	IsYakuman bool   // 是否为役满
}

// yakuInfoTable 役信息表
var yakuInfoTable = map[Yaku]YakuInfo{
	Tanyao:            {Name: "断幺", Fan: 1, IsOpen: true, IsYakuman: false},
	Pinfu:             {Name: "平和", Fan: 1, IsOpen: false, IsYakuman: false},
	Iipeikou:          {Name: "一对对", Fan: 1, IsOpen: false, IsYakuman: false},
	Ryanpeikou:        {Name: "二对对", Fan: 3, IsOpen: false, IsYakuman: false},
	Sanshokusequence:  {Name: "三色同顺", Fan: 2, IsOpen: true, IsYakuman: false},
	Ittsu:             {Name: "一通贯", Fan: 2, IsOpen: true, IsYakuman: false},
	Chanta:            {Name: "全带幺", Fan: 2, IsOpen: true, IsYakuman: false},
	Honchanta:         {Name: "混全带幺", Fan: 2, IsOpen: true, IsYakuman: false},
	Sanshokudoukou:    {Name: "三色同刻", Fan: 2, IsOpen: true, IsYakuman: false},
	Toitoi:            {Name: "对对和", Fan: 2, IsOpen: true, IsYakuman: false},
	Sanankou:          {Name: "三暗刻", Fan: 2, IsOpen: false, IsYakuman: false},
	Sangantu:          {Name: "三杆子", Fan: 2, IsOpen: true, IsYakuman: false},
	Sanputsu:          {Name: "三副露", Fan: 2, IsOpen: true, IsYakuman: false},
	Yakuhai:           {Name: "役牌", Fan: 1, IsOpen: true, IsYakuman: false},
	YakuhaiWind:       {Name: "役牌（风）", Fan: 1, IsOpen: true, IsYakuman: false},
	YakuhaiWhiteBoard: {Name: "役牌（白板）", Fan: 1, IsOpen: true, IsYakuman: false},
	YakuhaiGreenBoard: {Name: "役牌（绿板）", Fan: 1, IsOpen: true, IsYakuman: false},
	YakuhaiRedBoard:   {Name: "役牌（红板）", Fan: 1, IsOpen: true, IsYakuman: false},
	Honitsu:           {Name: "混一色", Fan: 3, IsOpen: true, IsYakuman: false},
	Junchan:           {Name: "纯全带幺", Fan: 3, IsOpen: true, IsYakuman: false},
	Chinitsu:          {Name: "清一色", Fan: 6, IsOpen: true, IsYakuman: false},
	Kokushi:           {Name: "国士无双", Fan: 13, IsOpen: false, IsYakuman: true},
	Suankou:           {Name: "四暗刻", Fan: 13, IsOpen: false, IsYakuman: true},
	Daisangen:         {Name: "大三元", Fan: 13, IsOpen: true, IsYakuman: true},
	Shosuushi:         {Name: "小四喜", Fan: 13, IsOpen: true, IsYakuman: true},
	Daisuushi:         {Name: "大四喜", Fan: 13, IsOpen: true, IsYakuman: true},
	Tsuisou:           {Name: "字一色", Fan: 13, IsOpen: true, IsYakuman: true},
	Ryuuisou:          {Name: "绿一色", Fan: 13, IsOpen: true, IsYakuman: true},
	Chinroutou:        {Name: "清老头", Fan: 13, IsOpen: true, IsYakuman: true},
	Honroutou:         {Name: "混老头", Fan: 13, IsOpen: true, IsYakuman: true},
	Chiitoitsu:        {Name: "七对子", Fan: 25, IsOpen: false, IsYakuman: true},
	Kazoe:             {Name: "数役满", Fan: 13, IsOpen: true, IsYakuman: false},
	Tenhou:            {Name: "天胡", Fan: 13, IsOpen: false, IsYakuman: true},
	Chihou:            {Name: "地胡", Fan: 13, IsOpen: false, IsYakuman: true},
	Churen:            {Name: "九莲宝灯", Fan: 13, IsOpen: false, IsYakuman: true},
	Menzentsumo:       {Name: "门清自摸", Fan: 1, IsOpen: false, IsYakuman: false},
	Dabururiichi:      {Name: "双立直", Fan: 2, IsOpen: false, IsYakuman: false},
}

// YakuToString 将Yaku转换为字符串
func YakuToString(yaku Yaku) string {
	if info, ok := yakuInfoTable[yaku]; ok {
		return info.Name
	}
	return fmt.Sprintf("未知役%d", yaku)
}

// GetFanCount 获取役的番数
func GetFanCount(yaku Yaku) int {
	if info, ok := yakuInfoTable[yaku]; ok {
		return info.Fan
	}
	return 0
}

// CanAgari 判断是否有役可以胡牌
func CanAgari(yakus []Yaku) bool {
	for _, yaku := range yakus {
		if GetFanCount(yaku) > 0 {
			return true
		}
	}
	return false
}

// IsYakuman 判断是否为役满
func IsYakuman(yaku Yaku) bool {
	if info, ok := yakuInfoTable[yaku]; ok {
		return info.IsYakuman
	}
	return false
}

// CanOpenYaku 判断役是否可以鸣牌成立
func CanOpenYaku(yaku Yaku) bool {
	if info, ok := yakuInfoTable[yaku]; ok {
		return info.IsOpen
	}
	return false
}

// RiichiYaku 是立直状态标记
const RiichiYaku Yaku = 0xFF

// IppatsuYaku 是一发状态标记
const IppatsuYaku Yaku = 0xFE

// GetRiichiInfo 获取立直的役信息
func GetRiichiInfo() YakuInfo {
	return YakuInfo{
		Name:      "立直",
		Fan:       1,
		IsOpen:    false,
		IsYakuman: false,
	}
}

// GetDoraInfo 获取宝牌的役信息
func GetDoraInfo() YakuInfo {
	return YakuInfo{
		Name:      "宝牌",
		Fan:       1,
		IsOpen:    true,
		IsYakuman: false,
	}
}

// GetUradoraInfo 获取里宝牌的役信息
func GetUradoraInfo() YakuInfo {
	return YakuInfo{
		Name:      "里宝牌",
		Fan:       1,
		IsOpen:    false,
		IsYakuman: false,
	}
}

// TotalFan 计算总番数
func TotalFan(yakus []Yaku) int {
	total := 0
	for _, yaku := range yakus {
		total += GetFanCount(yaku)
	}
	return total
}
