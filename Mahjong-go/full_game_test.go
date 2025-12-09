package mahjong

import (
	"fmt"
	"testing"
)

func TestSimulateFullGame(t *testing.T) {
	// 创建桌子和重放器
	table := NewTable()
	replayer := NewPaipuReplayer(table)

	// 初始化并发牌
	replayer.Init()

	// 模拟完整一局
	result := replayer.SimulateToCompletion()

	// 打印结果和分数
	fmt.Printf("游戏结果: %s\n", result.String())
	for i, p := range table.Players {
		fmt.Printf("玩家%d 得分: %d\n", i, p.Score)
	}

	// 简单断言：游戏必须结束
	if result == nil || result.Type == NoResult {
		t.Fatal("游戏未正常结束")
	}
}
