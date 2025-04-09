package monitor

import "github.com/shirou/gopsutil/v4/mem"

// MemStat 定义了系统内存信息的结构体，包含总内存、已用内存、空闲内存和使用百分比等字段
type MemStat struct {
	Total   uint64  // 总内存，表示系统中物理内存的总量，单位为字节
	Used    uint64  // 已用内存，表示系统中已使用的内存，单位为字节
	Free    uint64  // 空闲内存，表示系统中未使用的内存，单位为字节
	Percent float64 // 使用百分比，表示已用内存占总内存的百分比
}

// MemStats 获取系统的内存信息，包括总内存、已用内存、空闲内存和使用百分比
// 返回 MemStat 结构体指针和可能的错误信息
func MemStats() (*MemStat, error) {
	if info, err := mem.VirtualMemory(); err != nil {
		return nil, err
	} else {
		return &MemStat{
			Total:   info.Total,
			Used:    info.Used,
			Free:    info.Free,
			Percent: info.UsedPercent,
		}, nil
	}
}

// MustMemStats 获取系统的内存信息，如果获取失败则触发 panic
// 适用于需要确保内存信息必须获取成功的场景
func MustMemStats() *MemStat {
	if stat, err := MemStats(); err != nil {
		panic(err)
	} else {
		return stat
	}
}
