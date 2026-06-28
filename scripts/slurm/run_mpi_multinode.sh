#!/bin/bash
#SBATCH --job-name=matmul-mpi-mn
#SBATCH --nodes=2
#SBATCH --ntasks=8
#SBATCH --time=00:30:00
#SBATCH --output=results/raw/mpi-multinode-%j.out

set -euo pipefail

mkdir -p results/raw results/processed

mpirun -np "${TASKS:-8}" go run ./paralelo -n "${N:-3000}" -seed "${SEED:-42}" -csv "${CSV:-results/processed/parallel-multinode.csv}"
