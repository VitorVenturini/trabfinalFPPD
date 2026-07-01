package mpiutil

import (
	gompicomm "github.com/mvneves/gompi/comm"

	"trabfinalfppd/internal/matrix"
)

type DistributedWork struct {
	Partition matrix.WorkPartition
	AChunk    []float64
	BFull     []float64
}

func DistributeFromRoot(comm gompicomm.Communicator, n int, a, b []float64, partitions []matrix.WorkPartition) DistributedWork {
	for _, partition := range partitions {
		if partition.Rank == 0 {
			continue
		}

		meta := []int{n, partition.StartRow, partition.RowCount}
		comm.Send(meta, partition.Rank, TagMeta)
		if partition.RowCount == 0 {
			continue
		}

		comm.Send(b, partition.Rank, TagMatrixB)
		comm.Send(matrix.SliceRows(a, n, partition), partition.Rank, TagMatrixA)
	}

	local := partitions[0]
	return DistributedWork{
		Partition: local,
		AChunk:    matrix.SliceRows(a, n, local),
		BFull:     b,
	}
}

func ReceiveFromRoot(comm gompicomm.Communicator) (int, DistributedWork) {
	meta := make([]int, 3)
	comm.Recv(meta, 0, TagMeta)

	work := DistributedWork{
		Partition: matrix.WorkPartition{
			Rank:     comm.GetRank(),
			StartRow: meta[1],
			RowCount: meta[2],
		},
	}
	work.Partition.ElemCount = work.Partition.RowCount * meta[0]

	if work.Partition.RowCount == 0 {
		return meta[0], work
	}

	work.BFull = make([]float64, meta[0]*meta[0])
	work.AChunk = make([]float64, work.Partition.ElemCount)
	comm.Recv(work.BFull, 0, TagMatrixB)
	comm.Recv(work.AChunk, 0, TagMatrixA)

	return meta[0], work
}
