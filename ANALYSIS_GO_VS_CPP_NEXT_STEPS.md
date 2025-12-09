# Mahjong Go vs C++ — 验证与重放（操作步骤）

下面给出一套可执行的验证步骤，用于在本地对比 C++ 与 Go 的运行行为（确定性重放与日志对比）。如果你愿意，我也可以为你在仓库中添加两个小工具：一个 C++ 导出牌山/执行一次对局并输出 `export_yama` 的可执行程序；一个 Go 的重放脚本，读取该 `export_yama` 并在 Go 实现中重放并输出日志，便于 1:1 对比。

1) 在本机运行 Go 单元测试（验证 Go 实现自测）：

```powershell
cd D:\download\mahjong\tmp\mahjong
go test ./Mahjong-go -v
```

如果 `go` 命令不可用，请先安装 Go（https://go.dev/dl/），或在 PowerShell 中确认 `go env` 能工作。

2) 在 C++ 中导出一个确定性牌山（示例思路）：
- 使用 C++ 的 `Table::set_seed(seed)` 或 `game_init_with_config(...)` 初始化，并调用 `export_yama()` 将牌山序列保存到文本文件 `yama.txt`。如果仓库中没有现成可执行文件，需要通过下面命令构建一个小工具（示例）：

```powershell
md build; cd build
cmake ..
cmake --build . --config Release
# 生成后运行（示例，可执行文件名按实际生成修改）
.\tools\export_yama.exe > ..\yama.txt
```

3) 在 Go 中重放该牌山并比较日志：
- 在 `Mahjong-go` 中，可以编写一个小脚本/测试读取 `yama.txt`（逗号或空格分隔的 tile id 列表），调用 `NewTable()`, `ImportYama(yamaSlice)`，然后按同样的决策策略驱动游戏或使用相同的选择序列重放，并将 `GameLogRecord` 输出到 `go_log.txt`。

示例（PowerShell）:

```powershell
cd D:\download\mahjong\tmp\mahjong\Mahjong-go
# 运行自定义的重放程序（需实现），假设命名为 replay_with_yama.go 的 main
go run replay_with_yama.go ..\yama.txt > ..\go_log.txt
```

4) 对比 C++ 与 Go 的日志：

```powershell
cd D:\download\mahjong\tmp\mahjong
fc.exe cpp_log.txt go_log.txt
```

或者使用 `diff`、`git diff --no-index` 等工具做文本比较，关注 `GameResult`、每条 `Action` 的时间序列以及最终分数/输赢信息是否一致。

5) 可选：基于相同随机种子比较整体统计
- 在不能完全保证每一步行为一致的情况下，可以设置相同的随机种子（C++: `set_seed(int)`, Go: `SetSeed(int64)`），然后运行大量对局（N=1000+），比较统计指标（胡牌率、流局率、平均分、yakus 分布）。统计级别一致性可以检测逻辑偏差而非逐步差异。

---

如果你希望，我可以：
- (A) 在仓库中增加一个 C++ 小工具 `tools/export_yama.cpp`，编译并生成 `export_yama.exe`，用于导出 `export_yama()` 的输出；
- (B) 增加一个 Go 重放工具 `Mahjong-go/cmd/replay_with_yama/main.go`，用于读取 `yama.txt` 并在 Go 实现中重放；
- (C) 运行 `go test` 并把测试输出（失败的测试名称和堆栈）贴回给你以便进一步诊断（我可以在仓库内添加并运行这些工具，但需要你确认是否要我修改/添加文件）。

请选择下一步（例如："请生成 C++/Go 重放工具" 或 "我先在本地运行 go test 并把输出来给你"），我会继续实现或指导你运行。
