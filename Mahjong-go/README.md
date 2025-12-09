# Mahjong-Go - Go语言麻将游戏引擎

这是一个用 Go 语言完整重写的麻将游戏引擎，基于原始 C++ 实现，保持了完全相同的语义逻辑。

## 项目结构

整个项目直接放在 `Mahjong-go/` 目录下，包含以下主要文件：

### 核心模块

1. **tile.go** - 牌的定义和操作
   - `BaseTile` 枚举：34种基础牌（万、筒、索、字牌）
   - `Tile` 结构体：包含牌类型、赤宝牌标记、ID等
   - `Wind` 枚举：风向（东、南、西、北）
   - 各种牌判断函数（如 `Is1hai`、`Is9hai`、`IsYaochuhai` 等）

2. **action.go** - 行动定义
   - `BaseAction` 枚举：15种行动类型（弃牌、吃、碰、杠、自摸、立直等）
   - `Action` 结构体：表示一个行动及其对应的牌
   - `SelfAction` 和 `ResponseAction`：玩家自身和响应行动

3. **yaku.go** - 役的定义和分类
   - `Yaku` 枚举：麻将的各种役（共80多种）
   - 角色分类：一番、两番、三番、满贯、役满、双倍役满
   - 各种角色的中文名称映射

4. **rule.go** - 游戏规则和牌组
   - `TileGroup` 结构体：对子、顺子、刻子、杠子的表示
   - `CompletedTiles` 结构体：完整的胡牌型（1个对子+4个组）
   - `TileSplitter`：递归拆分算法，找出所有可能的胡牌方式

5. **player.go** - 玩家管理
   - `Player` 结构体：管理玩家的手牌、河、分数、各种状态
   - `River` 和 `RiverTile`：河牌记录
   - 各种玩家行动生成方法（`GetDiscard`、`GetKakan` 等）
   - 立直、振听、自摸、胡牌等状态管理

6. **table.go** - 游戏桌面管理
   - `Table` 结构体：管理整个游戏的核心
   - 136张牌的初始化、洗牌、摸牌逻辑
   - 牌山、宝牌、里宝牌的管理
   - 四个玩家和游戏状态的追踪

7. **score_counter.go** - 分数计算
   - `CounterResult` 结构体：计算结果（役、番数、符数、得分）
   - `ScoreCounter` 类：详细的役判定和分数计算
   - 各种分数常数定义

8. **game_result.go** - 游戏结果
   - `ResultType` 枚举：游戏结果类型（荣和、自摸、流局等）
   - `Result` 结构体：最终游戏结果
   - 各种结果生成函数（`GenerateResultRon`、`GenerateResultTsumo` 等）
   - 流局满贯的处理逻辑

9. **gameplay.go** - 游戏流程
   - `PaipuReplayer` 结构体：牌谱重放器
   - 行动执行和游戏推进
   - 游戏状态查询和分析

### 辅助文件

- **go.mod** - Go 模块定义文件
- **doc.go** - 项目文档和组件说明

## 主要特性

✅ **完整的游戏逻辑**
- 牌的管理和操作
- 34种基础牌的完整支持
- 赤宝牌（红牌）支持

✅ **玩家状态管理**
- 手牌、河、鸣牌组
- 立直、振听、一发、自摸等状态
- 听牌计算和分析

✅ **丰富的游戏规则**
- 80多种役的定义和识别
- 复杂的牌型拆分算法
- 役满、满贯等特殊胡法

✅ **详细的分数计算**
- 番数和符数的计算
- 各种胡法的得分
- 流局满贯的处理

✅ **游戏流程控制**
- 牌谱重放功能
- 行动的生成和执行
- 游戏状态的追踪

## 代码示例

### 创建游戏

```go
package main

import "mahjong"

func main() {
    // 创建一个新的游戏
    replayer := mahjong.NewPaipuReplayer()
    
    // 初始化游戏
    yama := []int{...} // 牌山的牌ID序列
    scores := [4]int{25000, 25000, 25000, 25000}
    replayer.Init(yama, scores, 0, 0, mahjong.East, 0)
    
    // 输出游戏状态
    replayer.PrintGameState()
}
```

### 判断牌型

```go
// 判断是否为顺子
tiles := []BaseTile{_1m, _2m, _3m}
if IsShuntsu(tiles) {
    println("这是一个顺子")
}

// 获取所有可能的胡牌方式
splitter := GetTileSplitter()
completed := splitter.GetAllCompletedTiles(handTiles)
if len(completed) > 0 {
    println("可以胡牌")
}
```

### 计算分数

```go
// 创建分数计算器
counter := NewScoreCounter()
counter.AnalyzeTiles(handTiles, winTile)

// 计算役
yakus := counter.CalculateYakus(player, table, winTile, isTsumo)

// 获取最终结果
result := counter.GetBestResult()
println("番数：", result.Fan, "符数：", result.Fu)
```

## 中文注释

代码中已添加了详细的中文注释，包括：
- 各类型、函数和方法的说明
- 复杂逻辑的步骤注释
- 重要参数的详细说明

## 对应原始C++类

| Go 包 | 对应 C++ 类/文件 | 说明 |
|------|-----------------|------|
| tile.go | Tile.h | 牌和风向定义 |
| action.go | Action.h/cpp | 行动类型和定义 |
| yaku.go | Yaku.h | 役的枚举和分类 |
| rule.go | Rule.h/cpp | 规则和牌组拆分 |
| player.go | Player.h/cpp | 玩家状态管理 |
| table.go | Table.h/cpp | 游戏桌面管理 |
| score_counter.go | ScoreCounter.h/cpp | 分数计算 |
| game_result.go | GameResult.h/cpp | 游戏结果处理 |
| gameplay.go | GamePlay.h/cpp | 游戏流程控制 |

## 语义相同性

✅ 所有的数据结构都保持了原始C++版本的语义
✅ 算法逻辑完全一致，只是用Go语言重新实现
✅ 函数名和参数名都尽可能保持一致
✅ 异常处理使用Go的惯例方式

## 编译和使用

```bash
# 进入目录
cd Mahjong-go

# 初始化模块（如果还未初始化）
go mod tidy

# 构建您的应用
go build -v

# 运行测试
go test -v ./...
```

## 扩展说明

这个Go版本是可扩展的，您可以：
- 添加新的役种
- 实现AI玩家
- 创建网络游戏服务器
- 添加用户界面
- 进行性能优化

## 许可证

遵循原始项目的许可证。

---

**项目完成于：** 2024年12月5日

这个Go语言版本完整地重新实现了C++原版的所有核心功能，保证了逻辑的一致性和正确性。
