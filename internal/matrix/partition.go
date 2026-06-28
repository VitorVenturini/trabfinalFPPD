package matrix

type WorkPartition struct {
	Rank      int
	StartRow  int
	RowCount  int
	ElemCount int
}

func BuildPartitions(n, processes int) []WorkPartition {
	partitions := make([]WorkPartition, processes)
	baseRows := n / processes
	remainder := n % processes
	startRow := 0

	for rank := 0; rank < processes; rank++ {
		rowCount := baseRows
		if rank < remainder {
			rowCount++
		}

		partitions[rank] = WorkPartition{
			Rank:      rank,
			StartRow:  startRow,
			RowCount:  rowCount,
			ElemCount: rowCount * n,
		}

		startRow += rowCount
	}

	return partitions
}

func SliceRows(data []float64, n int, partition WorkPartition) []float64 {
	if partition.RowCount == 0 {
		return nil
	}

	start := partition.StartRow * n
	end := start + partition.ElemCount
	return data[start:end]
}
