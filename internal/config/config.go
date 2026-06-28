package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type RunConfig struct {
	N     int
	Seed  int64
	Quiet bool
	CSV   string
	Label string
}

func Parse() (RunConfig, error) {
	cfg := RunConfig{}

	flag.IntVar(&cfg.N, "n", 3000, "matrix dimension")
	flag.Int64Var(&cfg.Seed, "seed", 42, "random seed")
	flag.BoolVar(&cfg.Quiet, "quiet", false, "reduce terminal output")
	flag.StringVar(&cfg.CSV, "csv", "", "optional CSV output file")
	flag.StringVar(&cfg.Label, "label", "", "optional run label")
	flag.Parse()

	if cfg.N <= 0 {
		return RunConfig{}, fmt.Errorf("invalid matrix dimension %d", cfg.N)
	}

	return cfg, nil
}

func DetectNodeCount() int {
	candidates := []string{
		"SLURM_JOB_NUM_NODES",
		"SLURM_NNODES",
	}

	for _, key := range candidates {
		if value := os.Getenv(key); value != "" {
			parsed, err := strconv.Atoi(value)
			if err == nil && parsed > 0 {
				return parsed
			}
		}
	}

	return 1
}
