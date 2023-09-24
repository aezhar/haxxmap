package haxxmap

type Options[K Hashable, V any] struct {
	initialSize uintptr
	hasher      HashFn[K]
	comparator  EqualFn[K]
}

func (c *Options[K, V]) init() {
	c.initialSize = defaultSize
}

func (c *Options[K, V]) setDefaults() {
	if c.hasher == nil {
		c.hasher = defaultHasher[K]()
	}

	if c.comparator == nil {
		c.comparator = defaultComparator[K]()
	}
}

// Option changes the behavior of the hashmap
type Option[K Hashable, V any] func(c *Options[K, V])

// WithInitialSize specifies the initial size of the hashmap
func WithInitialSize[K Hashable, V any](s uintptr) Option[K, V] {
	return Option[K, V](func(c *Options[K, V]) {
		c.initialSize = s
	})
}

// WithHasher specifies the hash function to use
func WithHasher[K Hashable, V any](fn HashFn[K]) Option[K, V] {
	return Option[K, V](func(c *Options[K, V]) {
		c.hasher = fn
	})
}

// WithComparator specifies the compare function to test keys for equality
func WithComparator[K Hashable, V any](fn EqualFn[K]) Option[K, V] {
	return Option[K, V](func(c *Options[K, V]) {
		c.comparator = fn
	})
}

func newOptions[K Hashable, V any](opts ...Option[K, V]) Options[K, V] {
	var c Options[K, V]
	c.init()

	for _, opt := range opts {
		opt(&c)
	}

	c.setDefaults()

	return c
}
