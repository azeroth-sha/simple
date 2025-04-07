package grace

import (
	"github.com/azeroth-sha/simple"
	"sort"
	"sync"
	"time"
)

// Grace 优雅关闭服务
type Grace struct {
	mu       *sync.Mutex
	wait     *sync.WaitGroup
	srv      []*service
	logger   simple.Logger
	interval time.Duration
}

// Add 添加服务
func (g *Grace) Add(name string, server Server, priority ...int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.srv = append(g.srv, newServ(name, server, g.logger, g.wait, priority...))
}

// Run 启动服务
func (g *Grace) Run() {
	all := make([]*service, 0, len(g.srv))
	all = append(all, g.srv...)
	sort.SliceStable(all, func(i, j int) bool {
		if all[i].priority != all[j].priority {
			return all[i].priority < all[j].priority
		} else {
			return all[i].name < all[j].name
		}
	})
	g.logger.Info("grace start")
	for _, srv := range all {
		srv.start()
		if g.interval > 0 {
			time.Sleep(g.interval)
		}
	}
	g.logger.Info("grace running")
}

// Stop 停止服务
func (g *Grace) Stop() {
	all := make([]*service, 0, len(g.srv))
	all = append(all, g.srv...)
	sort.SliceStable(all, func(i, j int) bool {
		if all[i].priority != all[j].priority {
			return all[i].priority > all[j].priority
		} else {
			return all[i].name > all[j].name
		}
	})
	g.logger.Info("grace stop")
	defer g.logger.Info("grace stopped")
	for _, srv := range all {
		srv.stop()
		if g.interval > 0 {
			time.Sleep(g.interval)
		}
	}
	g.wait.Wait()
}

// New 创建一个Grace实例
func New(l simple.Logger, opts ...Option) *Grace {
	g := &Grace{
		mu:       new(sync.Mutex),
		wait:     new(sync.WaitGroup),
		srv:      make([]*service, 0),
		logger:   l,
		interval: time.Millisecond * 150,
	}
	for _, option := range opts {
		option(g)
	}
	return g
}
