package mpiutil

import (
	gompicomm "github.com/mvneves/gompi/comm"

	"trabfinalfppd/internal/matrix"
)

func CollectToRoot(comm gompicomm.Communicator, n int, partitions []matrix.WorkPartition, localC []float64) []float64 {
	full := make([]float64, n*n)
	if partitions[0].RowCount > 0 {
		copy(full[:partitions[0].ElemCount], localC)
	}

	for _, partition := range partitions[1:] {
		if partition.RowCount == 0 {
			continue
		}

		buffer := make([]float64, partition.ElemCount)
		comm.Recv(buffer, partition.Rank, TagMatrixC)

		start := partition.StartRow * n
		end := start + partition.ElemCount
		copy(full[start:end], buffer)
	}

	return full
}

func SendResultToRoot(comm gompicomm.Communicator, localC []float64) {
	if len(localC) == 0 {
		return
	}

	comm.Send(localC, 0, TagMatrixC)
}
