package monitor

import (
	"github.com/shirou/gopsutil/v4/disk"
)

// DiskStat 定义了磁盘信息的结构体，包含设备名称、挂载点、文件系统类型、总容量、可用空间、已用空间和使用百分比等字段。
type DiskStat struct {
	Device     string  // 设备名称，表示磁盘设备的标识符
	MountPoint string  // 挂载点，表示磁盘在文件系统中的挂载路径
	FsType     string  // 文件系统类型，表示磁盘使用的文件系统（如ext4、NTFS等）
	Total      uint64  // 总容量，表示磁盘的总空间大小，单位为字节
	Free       uint64  // 可用空间，表示磁盘的剩余可用空间，单位为字节
	Used       uint64  // 已用空间，表示磁盘的已使用空间，单位为字节
	Percent    float64 // 使用百分比，表示磁盘的已用空间占总容量的百分比
}

// DiskStats 获取系统中所有磁盘分区的信息，包括设备名称、挂载点、文件系统类型、总容量、可用空间、已用空间和使用百分比等。
// 如果获取过程中发生错误，则返回错误信息。
func DiskStats() ([]*DiskStat, error) {
	stats := make([]*DiskStat, 0)
	// 获取系统中所有磁盘分区的列表
	list, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}
	// 遍历每个磁盘分区，获取详细信息
	for _, stat := range list {
		item := &DiskStat{
			Device:     stat.Device,
			MountPoint: stat.Mountpoint,
			FsType:     stat.Fstype,
		}
		// 获取磁盘分区的使用情况
		if use, err := disk.Usage(item.MountPoint); err != nil {
			return nil, err
		} else {
			item.Total = use.Total
			item.Free = use.Free
			item.Used = use.Used
			item.Percent = use.UsedPercent
		}
		// 将当前磁盘分区的信息添加到结果列表中
		stats = append(stats, item)
	}
	return stats, nil
}

// MustDiskStats 获取系统中所有磁盘分区的信息，如果获取过程中发生错误，则调用panic函数终止程序。
// 通常用于确保程序在无法获取磁盘信息时立即停止运行。
func MustDiskStats() []*DiskStat {
	stats, err := DiskStats()
	if err != nil {
		panic(err)
	}
	return stats
}
