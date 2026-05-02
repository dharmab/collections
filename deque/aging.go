package deque

import (
	"iter"
	"time"
)

type agingEntry[T any] struct {
	value   T
	addedAt time.Time
}

// Aging is a FIFO queue that drops entries older than a configured
// maximum age. Eviction happens lazily on every operation.
type Aging[T any] struct {
	items  []agingEntry[T]
	maxAge time.Duration
}

// NewAging returns an Aging queue with the given maximum age. A
// non-positive maxAge is treated as 0 (every entry expires immediately).
func NewAging[T any](maxAge time.Duration) *Aging[T] {
	if maxAge < 0 {
		maxAge = 0
	}
	return &Aging[T]{maxAge: maxAge}
}

// Push adds v to the back of the queue with the current timestamp.
func (d *Aging[T]) Push(v T) {
	d.items = append(d.items, agingEntry[T]{value: v, addedAt: time.Now()})
	d.evict()
}

// Pop removes and returns the oldest non-expired entry. The second
// return value is false if no such entry exists.
func (d *Aging[T]) Pop() (T, bool) {
	d.evict()
	var zero T
	if len(d.items) == 0 {
		return zero, false
	}
	v := d.items[0].value
	n := copy(d.items, d.items[1:])
	clear(d.items[n:])
	d.items = d.items[:n]
	return v, true
}

// Oldest returns the oldest non-expired entry without removing it. The
// second return value is false if no such entry exists.
func (d *Aging[T]) Oldest() (T, bool) {
	d.evict()
	var zero T
	if len(d.items) == 0 {
		return zero, false
	}
	return d.items[0].value, true
}

// Newest returns the most recently pushed non-expired entry without
// removing it. The second return value is false if no such entry exists.
func (d *Aging[T]) Newest() (T, bool) {
	d.evict()
	var zero T
	if len(d.items) == 0 {
		return zero, false
	}
	return d.items[len(d.items)-1].value, true
}

// At returns the non-expired entry at the given index, where 0 is the
// oldest non-expired entry. The second return value is false if the
// index is out of range.
func (d *Aging[T]) At(i int) (T, bool) {
	d.evict()
	var zero T
	if i < 0 || i >= len(d.items) {
		return zero, false
	}
	return d.items[i].value, true
}

// Len returns the number of non-expired entries.
func (d *Aging[T]) Len() int {
	d.evict()
	return len(d.items)
}

// Clear removes all entries from the queue.
func (d *Aging[T]) Clear() {
	clear(d.items)
	d.items = d.items[:0]
}

// MaxAge returns the configured maximum entry age.
func (d *Aging[T]) MaxAge() time.Duration {
	return d.maxAge
}

// SetMaxAge updates the maximum entry age. Entries that exceed the new
// age are evicted lazily on the next operation. A non-positive maxAge is
// treated as 0 (every entry expires immediately).
func (d *Aging[T]) SetMaxAge(maxAge time.Duration) {
	if maxAge < 0 {
		maxAge = 0
	}
	d.maxAge = maxAge
	d.evict()
}

// All returns an iterator over the non-expired entries from oldest to newest.
func (d *Aging[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		d.evict()
		for _, e := range d.items {
			if !yield(e.value) {
				return
			}
		}
	}
}

// evict drops every entry whose age exceeds maxAge.
func (d *Aging[T]) evict() {
	if len(d.items) == 0 {
		return
	}
	if d.maxAge == 0 {
		clear(d.items)
		d.items = d.items[:0]
		return
	}
	now := time.Now()
	i := 0
	for i < len(d.items) && now.Sub(d.items[i].addedAt) > d.maxAge {
		i++
	}
	if i == 0 {
		return
	}
	n := copy(d.items, d.items[i:])
	clear(d.items[n:])
	d.items = d.items[:n]
}
