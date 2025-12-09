package mahjong

import (
	"fmt"
	"sort"
	"time"
)

// TimeUnit 表示时间单位
type TimeUnit float64

const (
	Nanosecond  TimeUnit = 1
	Microsecond TimeUnit = 1e-3
	Millisecond TimeUnit = 1e-6
	Second      TimeUnit = 1e-9
)

// ProfileEvent 性能分析事件
type ProfileEvent struct {
	Name      string        // 事件名称
	StartTime time.Time     // 开始时间
	EndTime   time.Time     // 结束时间
	Duration  time.Duration // 持续时间
	Count     int           // 调用次数
}

// Profiler 性能分析器
// 用于测量和分析代码执行时间
type Profiler struct {
	events map[string]*ProfileEvent // 事件映射
	active map[string]time.Time     // 当前活跃的计时
}

var profilerInstance *Profiler

// GetProfiler 获取全局分析器实例
func GetProfiler() *Profiler {
	if profilerInstance == nil {
		profilerInstance = &Profiler{
			events: make(map[string]*ProfileEvent),
			active: make(map[string]time.Time),
		}
	}
	return profilerInstance
}

// Begin 开始对一个事件进行计时
func (p *Profiler) Begin(name string) {
	p.active[name] = time.Now()
}

// End 结束计时并记录事件
func (p *Profiler) End(name string) time.Duration {
	if startTime, exists := p.active[name]; exists {
		duration := time.Since(startTime)

		if event, ok := p.events[name]; ok {
			event.EndTime = time.Now()
			event.Duration += duration
			event.Count++
		} else {
			p.events[name] = &ProfileEvent{
				Name:      name,
				StartTime: startTime,
				EndTime:   time.Now(),
				Duration:  duration,
				Count:     1,
			}
		}

		delete(p.active, name)
		return duration
	}
	return 0
}

// GetEvent 获取指定名称的事件
func (p *Profiler) GetEvent(name string) *ProfileEvent {
	if event, exists := p.events[name]; exists {
		return event
	}
	return nil
}

// GetAverageDuration 获取平均执行时间(纳秒)
func (p *Profiler) GetAverageDuration(name string) float64 {
	if event, exists := p.events[name]; exists && event.Count > 0 {
		return float64(event.Duration.Nanoseconds()) / float64(event.Count)
	}
	return 0
}

// GetTotalDuration 获取总执行时间(纳秒)
func (p *Profiler) GetTotalDuration(name string) int64 {
	if event, exists := p.events[name]; exists {
		return event.Duration.Nanoseconds()
	}
	return 0
}

// GetCallCount 获取调用次数
func (p *Profiler) GetCallCount(name string) int {
	if event, exists := p.events[name]; exists {
		return event.Count
	}
	return 0
}

// Reset 重置所有事件
func (p *Profiler) Reset() {
	p.events = make(map[string]*ProfileEvent)
	p.active = make(map[string]time.Time)
}

// ResetEvent 重置指定的事件
func (p *Profiler) ResetEvent(name string) {
	delete(p.events, name)
	delete(p.active, name)
}

// PrintReport 打印性能报告
func (p *Profiler) PrintReport() {
	if len(p.events) == 0 {
		fmt.Println("No profile events recorded")
		return
	}

	// 排序事件
	names := make([]string, 0, len(p.events))
	for name := range p.events {
		names = append(names, name)
	}
	sort.Strings(names)

	// 打印表头
	fmt.Println("=== Profiler Report ===")
	fmt.Printf("%-30s %15s %15s %15s\n", "Event Name", "Total(ns)", "Average(ns)", "Count")
	fmt.Println(string(make([]byte, 75)))

	// 打印每个事件
	for _, name := range names {
		event := p.events[name]
		totalNs := event.Duration.Nanoseconds()
		avgNs := float64(totalNs) / float64(event.Count)

		fmt.Printf("%-30s %15d %15.2f %15d\n",
			name, totalNs, avgNs, event.Count)
	}
	fmt.Println(string(make([]byte, 75)))
}

// PrintReportMilliseconds 打印性能报告(毫秒)
func (p *Profiler) PrintReportMilliseconds() {
	if len(p.events) == 0 {
		fmt.Println("No profile events recorded")
		return
	}

	// 排序事件
	names := make([]string, 0, len(p.events))
	for name := range p.events {
		names = append(names, name)
	}
	sort.Strings(names)

	// 打印表头
	fmt.Println("=== Profiler Report (Milliseconds) ===")
	fmt.Printf("%-30s %15s %15s %15s\n", "Event Name", "Total(ms)", "Average(ms)", "Count")
	fmt.Println(string(make([]byte, 75)))

	// 打印每个事件
	for _, name := range names {
		event := p.events[name]
		totalMs := float64(event.Duration.Nanoseconds()) / 1e6
		avgMs := totalMs / float64(event.Count)

		fmt.Printf("%-30s %15.4f %15.6f %15d\n",
			name, totalMs, avgMs, event.Count)
	}
	fmt.Println(string(make([]byte, 75)))
}

// GetAllEvents 获取所有事件
func (p *Profiler) GetAllEvents() map[string]*ProfileEvent {
	return p.events
}

// ScopedTimer 作用域计时器，用于自动计时
type ScopedTimer struct {
	name string
}

// NewScopedTimer 创建一个作用域计时器
func NewScopedTimer(name string) *ScopedTimer {
	GetProfiler().Begin(name)
	return &ScopedTimer{name: name}
}

// Stop 停止计时
func (st *ScopedTimer) Stop() {
	GetProfiler().End(st.name)
}

// ResetProfiler 重置全局分析器
func ResetProfiler() {
	profilerInstance = nil
}
