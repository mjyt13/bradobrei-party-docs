package services

import "testing"

func TestUniqueSalonIDs(t *testing.T) {
	input := []uint{0, 3, 5, 3, 7, 5, 0, 9}
	got := uniqueSalonIDs(input)
	want := []uint{3, 5, 7, 9}

	if len(got) != len(want) {
		t.Fatalf("expected %d items, got %d: %#v", len(want), len(got), got)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected %v, got %v", want, got)
		}
	}
}
