#!/bin/bash
#SBATCH --job-name=matmul-seq
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --time=00:30:00
#SBATCH --output=results/raw/seq-%j.out

set -euo pipefail

mkdir -p results/raw results/processed

echo "Compiling sequential binary..."
go build -o bin/sequencial ./sequencial

echo "Running sequential..."
./bin/sequencial -n "${N:-3000}" -seed "${SEED:-42}" -csv "${CSV:-results/processed/seq.csv}"
