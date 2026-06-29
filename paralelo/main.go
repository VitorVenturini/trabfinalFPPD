package main

import (
	"fmt"
	"os"

	mpi "github.com/mnlphlp/gompi"

	"trabfinalfppd/internal/config"
	"trabfinalfppd/internal/matrix"
	"trabfinalfppd/internal/metrics"
	"trabfinalfppd/internal/mpiutil"
	"trabfinalfppd/internal/output"
)

func main() {
	mpi.Init()
	defer mpi.Finalize()

	comm := mpi.NewComm(true)
	rank := comm.GetRank()
	processes := comm.GetSize()

	cfg, err := config.Parse()
	if err != nil {
		if rank == 0 {
			fmt.Fprintln(os.Stderr, err)
		}
		return
	}

	// All processes wait here for a clean start before timing.
	comm.Barrier()
	start := mpi.WTime()

	var localWork mpiutil.DistributedWork
	var localC []float64

	if rank == 0 {
		if !cfg.Quiet {
			fmt.Printf("Running parallel matrix multiplication with n=%d and processes=%d\n", cfg.N, processes)
		}
		a, b := matrix.GenerateRandomPair(cfg.N, cfg.Seed)
		partitions := matrix.BuildPartitions(cfg.N, processes)
		localWork = mpiutil.DistributeFromRoot(comm, cfg.N, a, b, partitions)
	} else {
		_, localWork = mpiutil.ReceiveFromRoot(comm)
	}

	// All processes (including rank 0) perform their share of the multiplication.
	localC := make([]float64, localWork.Partition.ElemCount)
	matrix.MultiplyRows(localWork.AChunk, localWork.BFull, localC, cfg.N)

	// Stop the local timer. The total elapsed time will be the maximum of all local timers.
	localElapsed := mpi.WTime() - start

	// The root process collects results and timings, then prints the final output.
	if rank == 0 {
		partitions := matrix.BuildPartitions(cfg.N, processes)
		c := mpiutil.CollectToRoot(comm, cfg.N, partitions, localC)
		elapsed := collectMaxElapsed(comm, localElapsed, processes)

		result := metrics.RunResult{
			Mode:         "parallel",
			N:            cfg.N,
			Seed:         cfg.Seed,
			Nodes:        config.DetectNodeCount(),
			Processes:    processes,
			ElapsedSec:   elapsed,
			Label:        cfg.Label,
			Verification: matrix.ComputeVerification(c, cfg.N),
		}
		output.PrintResult(os.Stdout, result)

		if cfg.CSV != "" {
			if err := metrics.AppendCSV(cfg.CSV, result); err != nil {
				fmt.Fprintf(os.Stderr, "failed to write CSV: %v\n", err)
			}
		}
	} else {
		// Worker processes send their results and elapsed time to the root.
		mpiutil.SendResultToRoot(comm, localC)
		comm.Send([]float64{localElapsed}, 0, mpiutil.TagElapsed)
	}
}

func collectMaxElapsed(comm interface {
	Recv(data interface{}, source int, tag int) int
}, localElapsed float64, processes int) float64 {
	maxElapsed := localElapsed
	for source := 1; source < processes; source++ {
		buffer := make([]float64, 1)
		comm.Recv(buffer, source, mpiutil.TagElapsed)
		if buffer[0] > maxElapsed {
			maxElapsed = buffer[0]
		}
	}

	return maxElapsed
}
