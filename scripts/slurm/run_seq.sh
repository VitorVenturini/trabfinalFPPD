#!/bin/bash
#SBATCH --job-name=matmul-seq
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --time=01:00:00
#SBATCH --output=results/raw/%x-%j.out
#SBATCH --mem=2G

set -euo pipefail

echo "Compiling sequential binary..."
go build -o bin/sequencial ./sequencial

echo "Running sequential..."
./bin/sequencial -n "${N:-3000}" -seed "${SEED:-42}" -csv "${CSV:-results/processed/seq.csv}"
