package memory

type Option func(c *Cache)

// DefaultSize is 64M
var DefaultSize = 64 * 1024 * 1024

func Size(size int) Option {
	return func(c *Cache) {
		c.size = size
	}
}
