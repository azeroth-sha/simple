package monitor

import (
	"github.com/shirou/gopsutil/v4/cpu"
)

// CPUStat 定义了CPU信息的结构体，包含CPU编号、核心ID、核心数、模型、模型名称、主频和使用百分比等字段。
type CPUStat struct {
	CPU       int32   // CPU编号，标识当前CPU的序号
	CoreID    string  // 核心ID，标识当前核心的唯一ID
	Cores     int32   // CPU核心数，表示当前CPU的核心数量
	Model     string  // CPU型号，表示CPU的具体型号
	ModelName string  // CPU模型名称，表示CPU的完整型号名称
	Mhz       float64 // CPU主频，表示CPU的运行频率，单位为MHz
	Percent   float64 // CPU使用百分比，表示当前CPU的使用率
}

// CPUStats 获取系统的CPU信息，包括CPU编号、核心ID、核心数、模型、模型名称、主频和使用百分比等。
// 如果获取过程中发生错误，则返回错误信息。
func CPUStats() (*CPUStat, error) {
	// 获取CPU的详细信息
	if r1, e1 := cpu.Info(); e1 != nil {
		return nil, e1
	} else if r2, e2 := cpu.Percent(0, false); e2 != nil {
		// 获取CPU的使用百分比
		return nil, e2
	} else {
		stats := new(CPUStat)
		if len(r1) > 0 {
			// 填充CPUStat结构体的详细信息
			stats.CPU = r1[0].CPU
			stats.CoreID = r1[0].CoreID
			stats.Cores = r1[0].Cores
			stats.Model = r1[0].Model
			stats.ModelName = r1[0].ModelName
			stats.Mhz = r1[0].Mhz
		}
		if len(r2) > 0 {
			// 填充CPUStat结构体的使用百分比
			stats.Percent = r2[0]
		}
		return stats, nil
	}
}

// MustCPUStats 获取系统的CPU信息，如果获取过程中发生错误，则调用panic函数终止程序。
// 通常用于确保程序在无法获取CPU信息时立即停止运行。
func MustCPUStats() *CPUStat {
	if stats, err := CPUStats(); err != nil {
		panic(err)
	} else {
		return stats
	}
}
