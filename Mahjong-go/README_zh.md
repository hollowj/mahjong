# 麻将游戏引擎 - Go 版本

一个完整的日本麻将（麻雀）游戏引擎的 Go 实现，从 C++ 版本移植而来。

## 快速开始

### 编译
```bash
cd Mahjong-go
go build
```

### 运行测试
```bash
go test -v
```

### 使用示例

```go
package main

import (
    "fmt"
    "github.com/Agony5757/mahjong-go"
)

func main() {
    // 创建一个游戏桌子
    table := mahjong.NewTable()
    
    // 初始化牌牌
    table.InitTiles()
    table.InitRedDora3()
    table.ShuffleTiles()
    table.InitYama()
    table.InitDora()
    
    // 创建玩家
    for i := 0; i < mahjong.NPlayers; i++ {
        player := &mahjong.Player{
            Wind:  mahjong.Wind(i),
            Oya:   (i == 0),
            Score: 30000,
        }
        table.Players[i] = player
    }
    
    // 抽取初始牌
    for i := 0; i < 3; i++ {
        for j := 0; j < mahjong.NPlayers; j++ {
            tile := table.DrawTenhouStyle()
            table.Players[j].Hand = append(table.Players[j].Hand, tile)
        }
    }
    
    // 额外的一张牌
    for j := 0; j < mahjong.NPlayers; j++ {
        if j == 0 {
            tile := table.DrawTenhouStyle()
            table.Players[j].Hand = append(table.Players[j].Hand, tile)
        }
    }
    
    fmt.Println("游戏初始化完成！")
}
```

## 项目结构

| 文件 | 描述 |
|------|------|
| `tile.go` | 瓷牌类型、枚举和基础工具函数 |
| `action.go` | 玩家行动系统 |
| `yaku.go` | 役型（胜利手牌类型）定义 |
| `rule.go` | 游戏规则和验证 |
| `player.go` | 玩家状态管理 |
| `table.go` | 游戏桌子管理 |
| `score_counter.go` | 积分计算 |
| `game_result.go` | 游戏结果处理 |
| `gameplay.go` | 游戏流程管理 |
| `mahjong_test.go` | 单元测试 |

## 核心概念

### 瓷牌 (Tiles)
- **万牌** (Man): 1m - 9m
- **筒牌** (Pin): 1p - 9p  
- **索牌** (Sou): 1s - 9s
- **字牌** (Honors): 1z - 7z (东南西北白发中)

### 风向 (Winds)
- 东 (East)
- 南 (South)
- 西 (West)
- 北 (North)

### 役型 (Yaku)
31 种不同的完成形，从基础的断幺到高级的国士無双。

### 动作 (Actions)
- 打牌 (Discard)
- 吃牌 (Chi)
- 碰牌 (Pon)
- 杠牌 (Kan)
- 和牌 (Ron)
- 自摸 (Tsumo)
- 立直 (Riichi)

## 主要特性

✅ 完整的日本麻将规则实现
✅ 34 种瓷牌类型支持
✅ 31 种役型的自动检测
✅ 番数和符数的自动计算
✅ 河牌管理和振听检测
✅ 宝牌机制实现
✅ 纯 Go 标准库实现，无外部依赖
✅ 清晰的中文注释和文档
✅ 完整的单元测试覆盖

## API 概览

### 主要类型

```go
// 瓷牌
type Tile struct {
    Tile    BaseTile // 基础牌类型
    RedDora bool     // 是否为赤宝牌
    ID      int      // 唯一标识
}

// 玩家
type Player struct {
    Wind      Wind      // 方位
    Oya       bool      // 是否为庄家
    Score     int       // 分数
    Hand      []*Tile   // 手牌
    River     River     // 弃牌河
    // ... 其他字段
}

// 游戏桌子
type Table struct {
    Tiles            [NTiles]*Tile
    Players          [NPlayers]*Player
    DoraIndicator    []*Tile
    UraDoraIndicator []*Tile
    Yama             []*Tile
    // ... 其他字段
}

// 游戏结果
type GameResult struct {
    Type      ResultType
    WinnerIdx int
    LoserIdx  int
    // ... 其他字段
}
```

### 主要方法

```go
// 创建新的游戏桌子
table := NewTable()

// 初始化牌
table.InitTiles()
table.InitRedDora3()
table.ShuffleTiles()
table.InitYama()
table.InitDora()

// 抽取牌
tile := table.DrawTenhouStyle()  // 初始抽取
tile := table.DrawNormal()       // 普通抽取
tile := table.DrawRinshan()      // 死牌堆抽取

// 获取玩家行动
actions := player.GetDiscard()
actions := player.GetRon()
// ... 等等

// 检查完成形
splitter := GetTileSplitter()
completed := splitter.GetAllCompletedTiles(tiles)

// 计算积分
counter := &ScoreCounter{}
result := counter.CalculateScore(player, tiles, callGroups, winTile, false)
```

## 测试

运行所有测试：
```bash
go test -v
```

运行特定测试：
```bash
go test -v -run TestTileConstants
```

生成测试覆盖率报告：
```bash
go test -cover
```

## 文档

详细的项目文档请查看 `GO_PROJECT_SUMMARY.md`。

## 许可证

遵循原 C++ 版本的许可证。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 更新日志

### v1.0 (2024-12-05)
- 完整的从 C++ 到 Go 的移植
- 10 个单元测试全部通过
- 完整的 API 文档

## 联系方式

如有问题或建议，请提交 Issue。

