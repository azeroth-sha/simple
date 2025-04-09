package monitor

import (
	"github.com/shirou/gopsutil/v4/host"
)

// HostStat 定义了主机信息的结构体，包含主机名、运行时间、启动时间、操作系统信息等字段
type HostStat struct {
	Hostname        string // 主机名，表示当前设备在网络中的名称
	Uptime          uint64 // 系统运行时间，单位为秒
	BootTime        uint64 // 系统启动时间（Unix时间戳），单位为秒
	OS              string // 操作系统类型（如："linux", "windows"）
	Platform        string // 操作系统平台名称（如："ubuntu", "centos"）
	PlatformVersion string // 操作系统平台版本号
	KernelVersion   string // 操作系统内核版本
	KernelArch      string // 操作系统内核架构（如："x86_64", "arm64"）
}

// HostStats 获取系统主机信息，包括主机名、运行时间、启动时间、操作系统信息等
// 返回 HostStat 结构体指针和可能的错误信息
func HostStats() (*HostStat, error) {
	if info, err := host.Info(); err != nil {
		return nil, err
	} else {
		return &HostStat{
			Hostname:        info.Hostname,
			Uptime:          info.Uptime,
			BootTime:        info.BootTime,
			OS:              info.OS,
			Platform:        info.Platform,
			PlatformVersion: info.PlatformVersion,
			KernelVersion:   info.KernelVersion,
			KernelArch:      info.KernelArch,
		}, nil
	}
}

// MustHostStats 获取系统主机信息，如果获取失败则触发 panic
// 适用于需要确保主机信息必须获取成功的场景
func MustHostStats() *HostStat {
	if stat, err := HostStats(); err != nil {
		panic(err)
	} else {
		return stat
	}
}
