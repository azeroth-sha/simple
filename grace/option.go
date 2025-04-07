package grace

import "time"

type Option func(s *Grace)

func WithInterval(d time.Duration) Option {
	return func(g *Grace) {
		if d < 0 {
			return
		}
		g.interval = d
	}
}
