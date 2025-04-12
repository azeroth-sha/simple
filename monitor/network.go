package monitor

import (
	"github.com/shirou/gopsutil/v4/net"
	"slices"
)

// NetStat 定义了网络接口信息的结构体，包含接口名称、硬件地址、标志、地址列表、发送字节数和接收字节数等字段
type NetStat struct {
	Name      string   // 网络接口名称（如："eth0", "wlan0"）
	Hardware  string   // 网络接口的硬件地址（MAC地址）
	Flags     []string // 网络接口的标志（如："up", "broadcast"）
	AddrList  []string // 网络接口的IP地址列表
	BytesSent uint64   // 发送的字节数
	BytesRecv uint64   // 接收的字节数
}

// NetStats 获取系统中所有网络接口的信息，包括接口名称、硬件地址、标志、地址列表、发送字节数和接收字节数
// 返回 NetStat 结构体切片和可能的错误信息
func NetStats() ([]*NetStat, error) {
	// 获取网络接口的I/O计数器信息（发送和接收的字节数）
	if counts, err := net.IOCounters(true); err != nil {
		return nil, err
	} else if infos, err := net.Interfaces(); err != nil {
		// 获取网络接口的详细信息
		return nil, err
	} else {
		list := make([]*NetStat, 0, len(infos))
		// 遍历每个网络接口，填充 NetStat 结构体
		for _, info := range infos {
			stat := &NetStat{
				Name:      info.Name,
				Hardware:  info.HardwareAddr,
				Flags:     info.Flags,
				AddrList:  netAddrs(info.Addrs), // 获取网络接口的IP地址列表
				BytesSent: 0,
				BytesRecv: 0,
			}
			// 查找当前网络接口的I/O计数器信息
			if i := slices.IndexFunc(counts, func(stat net.IOCountersStat) bool {
				return stat.Name == info.Name
			}); i >= 0 {
				stat.BytesSent = counts[i].BytesSent
				stat.BytesRecv = counts[i].BytesRecv
			}
			list = append(list, stat)
		}
		return list, nil
	}
}

// MustNetStats 获取系统中所有网络接口的信息，如果获取失败则触发 panic
// 适用于需要确保网络信息必须获取成功的场景
func MustNetStats() []*NetStat {
	if stats, err := NetStats(); err != nil {
		panic(err)
	} else {
		return stats
	}
}

/*
  Package method
*/

// netAddrs 将 net.InterfaceAddrList 转换为字符串切片，提取每个地址的字符串表示
func netAddrs(list net.InterfaceAddrList) []string {
	all := make([]string, 0, len(list))
	for _, addr := range list {
		all = append(all, addr.Addr)
	}
	return all
}
