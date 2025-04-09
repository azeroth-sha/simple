package monitor

import "github.com/shirou/gopsutil/v4/load"

// LoadStat 定义了系统负载信息的结构体，包含1分钟、5分钟和15分钟的平均负载
type LoadStat struct {
	Load1  float64 // 1分钟平均负载
	Load5  float64 // 5分钟平均负载
	Load15 float64 // 15分钟平均负载
}

// LoadStats 获取系统的负载信息，包括1分钟、5分钟和15分钟的平均负载
// 返回 LoadStat 结构体指针和可能的错误信息
func LoadStats() (*LoadStat, error) {
	if info, err := load.Avg(); err != nil {
		return nil, err
	} else {
		return &LoadStat{
			Load1:  info.Load1,
			Load5:  info.Load5,
			Load15: info.Load15,
		}, nil
	}
}

// MustLoadStats 获取系统的负载信息，如果获取失败则触发 panic
// 适用于需要确保负载信息必须获取成功的场景
func MustLoadStats() *LoadStat {
	if stat, err := LoadStats(); err != nil {
		panic(err)
	} else {
		return stat
	}
}
