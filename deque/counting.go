package deque

import "iter"

// Counting is a FIFO queue that drops the oldest entry whenever a Push
// would cause the queue to exceed its maximum length.
type Counting[T any] struct {
	items  []T
	maxLen int
}

// NewCounting returns a Counting queue with the given maximum length.
// A non-positive maxLen is treated as 0 (the queue accepts no entries).
func NewCounting[T any](maxLen int) *Counting[T] {
	if maxLen < 0 {
		maxLen = 0
	}
	return &Counting[T]{maxLen: maxLen}
}

// Push adds v to the back of the queue, evicting the oldest entry if the
// queue is at capacity.
func (d *Counting[T]) Push(v T) {
	if d.maxLen == 0 {
		return
	}
	d.items = append(d.items, v)
	d.trim(d.maxLen)
}

// Pop removes and returns the oldest entry. The second return value is
// false if the queue is empty.
func (d *Counting[T]) Pop() (T, bool) {
	var zero T
	if len(d.items) == 0 {
		return zero, false
	}
	v := d.items[0]
	n := copy(d.items, d.items[1:])
	clear(d.items[n:])
	d.items = d.items[:n]
	return v, true
}

// Oldest returns the oldest entry without removing it. The second return
// value is false if the queue is empty.
func (d *Counting[T]) Oldest() (T, bool) {
	var zero T
	if len(d.items) == 0 {
		return zero, false
	}
	return d.items[0], true
}

// Newest returns the most recently pushed entry without removing it. The
// second return value is false if the queue is empty.
func (d *Counting[T]) Newest() (T, bool) {
	var zero T
	if len(d.items) == 0 {
		return zero, false
	}
	return d.items[len(d.items)-1], true
}

// At returns the entry at the given index, where 0 is the oldest entry.
// The second return value is false if the index is out of range.
func (d *Counting[T]) At(i int) (T, bool) {
	var zero T
	if i < 0 || i >= len(d.items) {
		return zero, false
	}
	return d.items[i], true
}

// Len returns the current number of entries.
func (d *Counting[T]) Len() int {
	return len(d.items)
}

// Clear removes all entries from the queue.
func (d *Counting[T]) Clear() {
	clear(d.items)
	d.items = d.items[:0]
}

// Cap returns the maximum length configured for the queue.
func (d *Counting[T]) Cap() int {
	return d.maxLen
}

// SetCap updates the maximum length, evicting the oldest entries
// immediately if the new cap is smaller than the current length.
// A non-positive maxLen is treated as 0.
func (d *Counting[T]) SetCap(maxLen int) {
	if maxLen < 0 {
		maxLen = 0
	}
	d.maxLen = maxLen
	d.trim(maxLen)
}

// All returns an iterator over the entries from oldest to newest.
func (d *Counting[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range d.items {
			if !yield(v) {
				return
			}
		}
	}
}

// trim drops the oldest entries until len(items) <= keep.
func (d *Counting[T]) trim(keep int) {
	if len(d.items) <= keep {
		return
	}
	n := copy(d.items, d.items[len(d.items)-keep:])
	clear(d.items[n:])
	d.items = d.items[:n]
}
