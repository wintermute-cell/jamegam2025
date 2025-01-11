package lib

// Based on https://stackoverflow.com/a/48330314
//
// A FreeList data structure which is basically an array that lets you remove
// elements from anywhere in constant-time (leaving holes behind which get
// reclaimed upon subsequent insertions in constant-time).
//

type freeListElement[T any] struct {
	element     T
	nextFreeIdx int
	removed     bool
}

// FreeList is a data structure that allows for constant-time removals and
// insertions.
type FreeList[T any] struct {
	data         []freeListElement[T]
	firstFreeIdx int
}

// NewFreeList creates a new FreeList with the given preallocated capacity.
func NewFreeList[T any](prealloc int) *FreeList[T] {
	return &FreeList[T]{
		data:         make([]freeListElement[T], 0, prealloc),
		firstFreeIdx: -1,
	}
}

// Clear removes all elements from the FreeList.
func (fl *FreeList[T]) Clear() {
	fl.data = fl.data[:0]
	fl.firstFreeIdx = -1
}

// Insert adds an element to the FreeList and returns its index.
func (fl *FreeList[T]) Insert(element T) int {
	if fl.firstFreeIdx == -1 {
		// -1 means there are no holes in the array
		elem := freeListElement[T]{element: element, removed: false}
		fl.data = append(fl.data, elem)
		return int(len(fl.data)) - 1
	} else {
		// reuse a hole in this case
		idx := fl.firstFreeIdx
		fl.firstFreeIdx = fl.data[idx].nextFreeIdx
		fl.data[idx] = freeListElement[T]{element: element, removed: false}
		return idx
	}
}

// Remove removes the element at the given index from the FreeList.
func (fl *FreeList[T]) Remove(idx int) {
	fl.data[idx].nextFreeIdx = fl.firstFreeIdx
	fl.data[idx].removed = true
	fl.firstFreeIdx = idx
}

// Get returns the element at the given index.
func (fl *FreeList[T]) Get(idx int) T {
	return fl.data[idx].element
}

// Set sets the element at the given index.
func (fl *FreeList[T]) Set(idx int, element T) {
	fl.data[idx].element = element
}

// GetAll returns all elements in the FreeList.
func (fl *FreeList[T]) GetAll() []T {
	if len(fl.data) == 0 {
		return []T{}
	}

	// build a map of free indices
	freeIndices := make(map[int]bool)
	for idx := fl.firstFreeIdx; idx != -1; idx = fl.data[idx].nextFreeIdx {
		freeIndices[idx] = true
	}

	result := make([]T, 0, len(fl.data)-len(freeIndices))
	for i := range fl.data {
		if !freeIndices[i] {
			result = append(result, fl.data[i].element)
		}
	}

	return result
}

// FuncAll calls the given function for all elements in the FreeList.
func (fl *FreeList[T]) FuncAll(f func(idx int, elem T)) {
	if len(fl.data) == 0 {
		return
	}

	// // build a map of free indices
	// freeIndices := make(map[int]bool)
	// for idx := fl.firstFreeIdx; idx != -1 && idx != fl.data[idx].nextFreeIdx; idx = fl.data[idx].nextFreeIdx {
	// 	freeIndices[idx] = true
	// }

	for i := range fl.data {
		if !fl.data[i].removed {
			f(i, fl.data[i].element)
		}
	}
}
