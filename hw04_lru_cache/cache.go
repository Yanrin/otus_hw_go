package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

// NewCache initializes the cache.
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

// Set sets an element to the cache by the key.
func (c *lruCache) Set(key Key, value interface{}) bool {
	if item, ok := c.items[key]; ok { // element already exists
		c.queue.MoveToFront(item)
		item.Value = value
		c.queue.Front().Value = value

		return true
	}
	if c.queue.Len() >= c.capacity {
		lastItem := c.queue.Back()
		var removingKey Key
		for k, item := range c.items {
			if item == lastItem {
				removingKey = k
				break
			}
		}
		delete(c.items, removingKey)
		c.queue.Remove(lastItem)
	}
	c.items[key] = c.queue.PushFront(value)
	return false
}

// Get gets an element from the cache by the key.
func (c *lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := c.items[key]; ok { // element already exists
		c.queue.MoveToFront(item)
		return item.Value, true
	}
	return nil, false
}

// Clear clears the cache.
func (c *lruCache) Clear() {
	for k, item := range c.items {
		c.queue.Remove(item)
		delete(c.items, k)
	}
}
