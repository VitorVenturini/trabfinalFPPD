#!/bin/bash
#SBATCH --job-name=matmul-seq
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --time=00:30:00
#SBATCH --output=%x-%j.out

set -euo pipefail

echo "Compiling sequential binary..."
go build -o bin/sequencial ./sequencial

echo "Running sequential..."
./bin/sequencial -n "${N:-3000}" -seed "${SEED:-42}" -csv "${CSV:-results/processed/seq.csv}"
