package matrix

import "testing"

func TestMultiplySequential2x2(t *testing.T) {
	a := []float64{
		1, 2,
		3, 4,
	}
	b := []float64{
		5, 6,
		7, 8,
	}

	got := MultiplySequential(a, b, 2)
	want := []float64{
		19, 22,
		43, 50,
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected result at index %d: got %.2f want %.2f", i, got[i], want[i])
		}
	}
}

func TestBuildPartitionsUneven(t *testing.T) {
	partitions := BuildPartitions(10, 3)

	wantRows := []int{4, 3, 3}
	wantStarts := []int{0, 4, 7}

	if len(partitions) != 3 {
		t.Fatalf("unexpected partition count: got %d want 3", len(partitions))
	}

	for i := range partitions {
		if partitions[i].RowCount != wantRows[i] {
			t.Fatalf("unexpected row count for rank %d: got %d want %d", i, partitions[i].RowCount, wantRows[i])
		}
		if partitions[i].StartRow != wantStarts[i] {
			t.Fatalf("unexpected start row for rank %d: got %d want %d", i, partitions[i].StartRow, wantStarts[i])
		}
	}
}

func TestComputeVerification(t *testing.T) {
	c := []float64{
		1, 2,
		3, 4,
	}

	got := ComputeVerification(c, 2)

	if got.TopLeft != 1 {
		t.Fatalf("unexpected top left: got %.2f want 1", got.TopLeft)
	}
	if got.TopRight != 2 {
		t.Fatalf("unexpected top right: got %.2f want 2", got.TopRight)
	}
	if got.BottomLeft != 3 {
		t.Fatalf("unexpected bottom left: got %.2f want 3", got.BottomLeft)
	}
	if got.BottomRight != 4 {
		t.Fatalf("unexpected bottom right: got %.2f want 4", got.BottomRight)
	}
	if got.Checksum != 10 {
		t.Fatalf("unexpected checksum: got %.2f want 10", got.Checksum)
	}
}
