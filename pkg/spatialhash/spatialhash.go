package spatialhash

import (
	"log"
)

type SHBounds struct {
	// Mx and My are the middle point of the bounds.
	Mx, My int32

	// HWidth and HHeight are the half width and half height of the bounds.
	HWidth, HHeight int32
}

type SHElement struct {
	ID             int32
	Bounds         SHBounds
	keysStartIdx   int32
	numKeys        int32
	foundThisQuery bool
}

// hashCoords returns a hash value for the given 2d coordinates.
func hashCoords(x, y int32) int32 {
	// I could use math.Abs here, but let's save the generics overhead.
	// This function will run often.
	k := (x * 92837111) ^ (y * 689287499)
	if k < 0 {
		return -k
	}
	return k
}

// Keys returns the hash keys for the given element. This key depends on the
// elements bounds/position. A large element may span multiple cells and thus
// have multiple keys.
func (b SHBounds) Keys(cellSize int32, ret []int32) []int32 {
	// Snap the bounds to the nearest grid coordinates.
	left := (b.Mx - b.HWidth) / cellSize
	right := (b.Mx + b.HWidth) / cellSize
	top := (b.My - b.HHeight) / cellSize
	bottom := (b.My + b.HHeight) / cellSize

	// Iterate over every 2D cell index that the element spans.
	ret = ret[:0] // Reset the slice, but reuse the underlying array to avoid reallocation
	for x := left; x <= right; x++ {
		for y := top; y <= bottom; y++ {
			ret = append(ret, hashCoords(x, y))
		}
	}

	return ret
}

type SpatialHash struct {
	tableSize   int32
	table       []int32
	cellSize    int32
	maxElements int32
	elements    []*SHElement

	// We re-use these slices to avoid allocations.
	inBoundsRet        []*SHElement
	dedupIndicesHelper []int32
	keysRet            []int32
}

// NewSpatialHash creates a new spatial hash with the given table size and cell
// size.
//
// A good rule of thumb is to define the spacing as two times the radius of the
// average particle. Also, table size should be maxElements*2
func NewSpatialHash(tableSize, cellSize, maxElements int32) *SpatialHash {
	return &SpatialHash{
		tableSize:   tableSize,
		table:       make([]int32, tableSize+1), // we need an extra guard element
		cellSize:    cellSize,
		maxElements: maxElements,
		elements:    make([]*SHElement, maxElements),

		inBoundsRet:        make([]*SHElement, 0, 1024),
		dedupIndicesHelper: make([]int32, 0, 1024),
		keysRet:            make([]int32, 0, 1024),
	}
}

// Construct creates a spatial hash table from the given elements.
// All elements for a collision problem must be passed to this function at once.
func (h *SpatialHash) Construct(elements []*SHElement) {

	// 1. Split each element into keys/cells and insert them into the hash table.
	keys := make([]int32, 0, len(elements)*2) // preallocate for an average of 2 keys per element
	for i, element := range elements {
		elements[i].keysStartIdx = int32(len(keys))
		elemKeys := element.Bounds.Keys(h.cellSize, h.keysRet)
		elements[i].numKeys = int32(len(elemKeys))
		keys = append(keys, elemKeys...)

		// We insert every key of the element into the hash table.
		for _, key := range elemKeys {
			hashIdx := key % h.tableSize
			h.table[hashIdx] += 1
		}
	}

	// 2. Now we do a partial sum pass on the table.
	for i := 1; i < len(h.table); i++ {
		h.table[i] += h.table[i-1]
	}

	// 3. Now insert the elements into the elements array.
	for _, element := range elements {
		for _, key := range keys[element.keysStartIdx : element.keysStartIdx+element.numKeys] {
			hashIdx := key % h.tableSize
			h.table[hashIdx] -= 1
			if h.table[hashIdx] >= h.maxElements {
				log.Fatalf("tried to insert too many elements . Tried %v, max %v", h.table[hashIdx], h.maxElements)
			}
			h.elements[h.table[hashIdx]] = element
		}
	}
}

// Clear removes all elements from the spatial hash.
func (h *SpatialHash) Clear() {
	h.table = make([]int32, h.tableSize+1)
	h.elements = make([]*SHElement, h.maxElements)
}

// InBounds returns all elements that are in the given bounds.
func (h *SpatialHash) InBounds(bounds SHBounds) []*SHElement {
	keys := bounds.Keys(h.cellSize, h.keysRet)

	// Check each hash key for elements.
	//elements := make([]SHElement, 0, len(keys)*2) // preallocate for an average of 2 elements per key.
	h.inBoundsRet = h.inBoundsRet[:0]
	h.dedupIndicesHelper = h.dedupIndicesHelper[:0]
	for _, key := range keys {
		hashIdx := key % h.tableSize
		elemIdx := h.table[hashIdx]
		numElements := h.table[hashIdx+1] - elemIdx // number of elements in this cell

		// Collect all elements in the cell.
		for i := 0; i < int(numElements); i++ {
			idx := elemIdx + int32(i)
			elem := h.elements[idx]
			if !elem.foundThisQuery {
				elem.foundThisQuery = true
				h.inBoundsRet = append(h.inBoundsRet, elem)
				h.dedupIndicesHelper = append(h.dedupIndicesHelper, idx)
			}
		}
	}

	// Reset the deduplication flag for the next query.
	for _, idx := range h.dedupIndicesHelper {
		h.elements[idx].foundThisQuery = false
	}

	return h.inBoundsRet
}

// // Draw draws the grid of the spatial hash using raylib.
// func (h *SpatialHash) Draw(nw, se rl.Vector2, withNumbers bool) {
// 	for x := int32(nw.X); x < int32(se.X); x += h.cellSize {
// 		rl.DrawLine(x, int32(nw.Y), x, int32(se.Y), rl.Gray)
// 	}
// 	for y := int32(nw.Y); y < int32(se.Y); y += h.cellSize {
// 		rl.DrawLine(int32(nw.X), y, int32(se.X), y, rl.Gray)
// 	}
//
// 	if !withNumbers {
// 		return
// 	}
//
// 	// Iterate over the cells in a 1920x1080 window,
// 	// Check the number of elements in each cell and draw the number as text.
// 	for x := int32(nw.X); x < int32(se.X); x += h.cellSize {
// 		for y := int32(nw.Y); y < int32(se.Y); y += h.cellSize {
// 			keys := SHBounds{Mx: x, My: y, HWidth: h.cellSize / 2, HHeight: h.cellSize / 2}.Keys(h.cellSize, h.keysRet)
// 			numElements := int32(0)
// 			for _, key := range keys {
// 				hashIdx := key % h.tableSize
// 				numElements += h.table[hashIdx+1] - h.table[hashIdx]
// 			}
// 			rl.DrawText(fmt.Sprintf("%d", numElements), x, y, 12, rl.Black)
// 		}
// 	}
// }

//type SpatialHash struct {
//	CellSize int32
//	Cells    map[SpatialHashKey][]int32
//}
//
//type Vector struct {
//	X, Y int32
//}
//
//type SpatialHashKey struct {
//	X, Y int32
//}
//
//func NewSpatialHash(cellSize int32) *SpatialHash {
//	cells := make(map[SpatialHashKey][]int32)
//	return &SpatialHash{cellSize, cells}
//}
//
//func (h *SpatialHash) KeyForPoint(point Vector) SpatialHashKey {
//	x := point.X / h.CellSize
//	y := point.Y / h.CellSize
//	return SpatialHashKey{x, y}
//}
//
//func (h *SpatialHash) Add(point Vector, data int32) {
//	key := h.KeyForPoint(point)
//	h.Cells[key] = append(h.Cells[key], data)
//}
//
//func (h *SpatialHash) Nearby(point Vector) []int32 {
//	var result []int32
//	key := h.KeyForPoint(point)
//	for dx := int32(-1); dx <= 1; dx++ {
//		for dy := int32(-1); dy <= 1; dy++ {
//			k := SpatialHashKey{key.X + dx, key.Y + dy}
//			result = append(result, h.Cells[k]...)
//		}
//	}
//	return result
//}
//
//func (h *SpatialHash) Clear() {
//	h.Cells = make(map[SpatialHashKey][]int32)
//}
