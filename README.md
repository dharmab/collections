# Collections

A module for Go providing collections and operations on those data structures that I find myself reaching for in my projects.

## sets

Inspired by "slices" and "maps" from the Go standard library, "sets" provides helper functions and iterators for working with Go's version of unordered sets (i.e. `map[T]struct{}`). I use this pattern so often that I'm always surprised the standard library doesn't provide it.

## deque

A queue where old elements are automatically removed. I used this data structure in [SkyEye](https://github.com/dharmab/skyeye) for trackfiles (think the little moving dots representing airplanes on an air traffic controller's radar).
