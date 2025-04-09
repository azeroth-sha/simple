package monitor

// SystemStat 定义了系统整体监控信息的结构体，包含CPU、磁盘、主机、负载、内存和网络等信息
type SystemStat struct {
	CPU  *CPUStat    // CPU信息，包含CPU编号、核心数、使用百分比等
	Disk []*DiskStat // 磁盘信息列表，包含每个磁盘的设备名称、挂载点、使用情况等
	Host *HostStat   // 主机信息，包含主机名、运行时间、操作系统信息等
	Load *LoadStat   // 系统负载信息，包含1分钟、5分钟和15分钟的平均负载
	Mem  *MemStat    // 内存信息，包含总内存、已用内存、空闲内存和使用百分比
	Net  []*NetStat  // 网络信息列表，包含每个网络接口的名称、硬件地址、发送和接收字节数等
}

// SysStats 获取系统的整体监控信息，包括CPU、磁盘、主机、负载、内存和网络等信息
// 返回 SystemStat 结构体指针和可能的错误信息
func SysStats() (*SystemStat, error) {
	stats := new(SystemStat)
	// 获取CPU信息
	if cpu, err := CPUStats(); err != nil {
		return nil, err
	} else {
		stats.CPU = cpu
	}
	// 获取磁盘信息
	if disk, err := DiskStats(); err != nil {
		return nil, err
	} else {
		stats.Disk = disk
	}
	// 获取主机信息
	if host, err := HostStats(); err != nil {
		return nil, err
	} else {
		stats.Host = host
	}
	// 获取系统负载信息
	if load, err := LoadStats(); err != nil {
		return nil, err
	} else {
		stats.Load = load
	}
	// 获取内存信息
	if mem, err := MemStats(); err != nil {
		return nil, err
	} else {
		stats.Mem = mem
	}
	// 获取网络信息
	if net, err := NetStats(); err != nil {
		return nil, err
	} else {
		stats.Net = net
	}
	return stats, nil
}

// MustSysStats 获取系统的整体监控信息，如果获取失败则触发 panic
// 适用于需要确保系统监控信息必须获取成功的场景
func MustSysStats() *SystemStat {
	stats, err := SysStats()
	if err != nil {
		panic(err)
	}
	return stats
}
