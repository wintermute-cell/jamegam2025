package lib

import (
	"testing"
)

// TestFreeList_Insert verifies that elements are inserted correctly.
func TestFreeList_Insert(t *testing.T) {
	fl := NewFreeList[int](10)

	// Test inserting elements
	idx1 := fl.Insert(10)
	if fl.Get(idx1) != 10 {
		t.Fatalf("expected 10, got %d", fl.Get(idx1))
	}

	idx2 := fl.Insert(20)
	if fl.Get(idx2) != 20 {
		t.Fatalf("expected 20, got %d", fl.Get(idx2))
	}

	// Test that indices are sequential
	if idx2 != idx1+1 {
		t.Fatalf("expected idx2 to be %d, got %d", idx1+1, idx2)
	}
}

// TestFreeList_Remove verifies that elements are removed correctly.
func TestFreeList_Remove(t *testing.T) {
	fl := NewFreeList[int](10)

	// Insert elements
	idx1 := fl.Insert(10)
	idx2 := fl.Insert(20)
	idx3 := fl.Insert(30)

	// Remove the middle element
	fl.Remove(idx2)

	// Reuse the removed index
	idx4 := fl.Insert(40)
	if idx4 != idx2 {
		t.Fatalf("expected idx4 to be %d, got %d", idx2, idx4)
	}

	// Verify that the reused index contains the new value
	if fl.Get(idx4) != 40 {
		t.Fatalf("expected 40, got %d", fl.Get(idx4))
	}

	// Verify that other elements are unaffected
	if fl.Get(idx1) != 10 {
		t.Fatalf("expected 10, got %d", fl.Get(idx1))
	}
	if fl.Get(idx3) != 30 {
		t.Fatalf("expected 30, got %d", fl.Get(idx3))
	}
}

// TestFreeList_Clear verifies that the freelist is cleared correctly.
func TestFreeList_Clear(t *testing.T) {
	fl := NewFreeList[int](10)

	// Insert
	fl.Insert(10)
	fl.Insert(20)
	fl.Insert(30)

	// Clear the freelist
	fl.Clear()

	// Test that the freelist is empty
	if len(fl.data) != 0 {
		t.Fatalf("expected len(fl.data) to be 0, got %d", len(fl.data))
	}

	// Insert a new element after clearing
	idx := fl.Insert(40)
	if fl.Get(idx) != 40 {
		t.Fatalf("expected 40, got %d", fl.Get(idx))
	}
}

// TestFreeList_GetSet verifies the Get and Set methods.
func TestFreeList_GetSet(t *testing.T) {
	fl := NewFreeList[int](10)

	// Insert an elem
	idx := fl.Insert(10)

	// Set a new value
	fl.Set(idx, 20)

	// Verify the value was set correctly
	if fl.Get(idx) != 20 {
		t.Fatalf("expected 20, got %d", fl.Get(idx))
	}
}

// TestFreeList_ReuseAfterMultipleRemovals verifies that the freelist correctly reuses indices after multiple removals.
func TestFreeList_ReuseAfterMultipleRemovals(t *testing.T) {
	fl := NewFreeList[int](10)

	// Insert elements
	idx1 := fl.Insert(10)
	idx2 := fl.Insert(20)
	idx3 := fl.Insert(30)

	// Remove two elements
	fl.Remove(idx2)
	fl.Remove(idx1)

	// Insert new elements and check if indices are reused correctly
	idx4 := fl.Insert(40)
	idx5 := fl.Insert(50)

	if idx4 != idx1 {
		t.Fatalf("expected idx4 to be %d, got %d", idx1, idx4)
	}
	if idx5 != idx2 {
		t.Fatalf("expected idx5 to be %d, got %d", idx2, idx5)
	}

	// Verify the values at the reused indices
	if fl.Get(idx4) != 40 {
		t.Fatalf("expected 40, got %d", fl.Get(idx4))
	}
	if fl.Get(idx5) != 50 {
		t.Fatalf("expected 50, got %d", fl.Get(idx5))
	}

	// Check that the remaining element is still correct
	if fl.Get(idx3) != 30 {
		t.Fatalf("expected 30, got %d", fl.Get(idx3))
	}
}

// FuzzFreeList is a fuzz test for the FreeList data structure.
func FuzzFreeList(f *testing.F) {
	// Seed the fuzzer with initial values.
	f.Add(int(10))
	f.Add(int(20))
	f.Add(int(30))

	f.Fuzz(func(t *testing.T, value int) {
		// Create a new FreeList
		fl := NewFreeList[int](10)

		// insert the fuzzed value
		idx := fl.Insert(value)

		// Check that the value was inserted correctly
		if fl.Get(idx) != value {
			t.Errorf("expected value %d, got %d", value, fl.Get(idx))
		}

		// Randomly decide to remove the element
		if value%2 == 0 { // for some variability, remove if value is even
			fl.Remove(idx)

			// After removal, reinsert and check if the index is reused properly
			newIdx := fl.Insert(value)
			if newIdx != idx {
				t.Errorf("expected index reuse, but got new index %d instead of %d", newIdx, idx)
			}
		}

		// Optionally clear the FreeList to test behavior after clearing
		if value%3 == 0 { // every third value, clear the list
			fl.Clear()

			// After clearing, check if a new insert works properly
			newIdx := fl.Insert(value)
			if fl.Get(newIdx) != value {
				t.Errorf("expected value %d after clearing, got %d", value, fl.Get(newIdx))
			}
		}
	})
}
