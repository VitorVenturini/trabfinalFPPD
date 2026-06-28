#!/bin/bash
#SBATCH --job-name=matmul-mpi-1n
#SBATCH --nodes=1
#SBATCH --ntasks=4
#SBATCH --time=00:30:00
#SBATCH --output=results/raw/mpi-1node-%j.out

set -euo pipefail

mkdir -p results/raw results/processed

mpirun -np "${TASKS:-4}" go run ./paralelo -n "${N:-3000}" -seed "${SEED:-42}" -csv "${CSV:-results/processed/parallel-1node.csv}"
