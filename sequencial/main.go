package main

import (
	"fmt"
	"os"
	"time"

	"trabfinalfppd/internal/config"
	"trabfinalfppd/internal/matrix"
	"trabfinalfppd/internal/metrics"
	"trabfinalfppd/internal/output"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if !cfg.Quiet {
		fmt.Printf("Running sequential matrix multiplication with n=%d\n", cfg.N)
	}

	a, b := matrix.GenerateRandomPair(cfg.N, cfg.Seed)

	start := time.Now()
	c := matrix.MultiplySequential(a, b, cfg.N)
	elapsed := time.Since(start).Seconds()

	result := metrics.RunResult{
		Mode:        "sequential",
		N:           cfg.N,
		Seed:        cfg.Seed,
		Nodes:       1,
		Processes:   1,
		ElapsedSec:  elapsed,
		Label:       cfg.Label,
		Verification: matrix.ComputeVerification(c, cfg.N),
	}

	output.PrintResult(os.Stdout, result)

	if cfg.CSV != "" {
		if err := metrics.AppendCSV(cfg.CSV, result); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write CSV: %v\n", err)
			os.Exit(1)
		}
	}
}
