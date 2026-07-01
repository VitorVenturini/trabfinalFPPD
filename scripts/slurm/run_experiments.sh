#!/bin/bash

set -euo pipefail

mkdir -p results/raw results/processed

N_VALUE="${N:-3000}"
SEED_VALUE="${SEED:-42}"

configs=(
  # nodes tasks
  "1 1"   # Baseline P=1
  "1 2"   # Intra-node
  "1 4"   # Intra-node
  "1 8"   # Intra-node
  "1 16"  # Intra-node (full node?)
  "1 32"  # Intra-node (oversubscription/hyperthreading test)
  "2 8"   # Inter-node (4 tasks per node)
  "2 16"  # Inter-node (8 tasks per node)
)

echo "Compiling parallel binary before submitting jobs..."
go build -o bin/paralelo ./paralelo

for config in "${configs[@]}"; do
  read -r nodes tasks <<< "${config}"
  for run in 1 2 3; do
    label="nodes=${nodes}-tasks=${tasks}-run=${run}"
    output_csv="results/processed/run_${label}.csv"
    output_log="results/raw/run_${label}.out"

    echo "Submitting job: nodes=${nodes}, tasks=${tasks}, run=${run}"
    sbatch --nodes="${nodes}" --ntasks="${tasks}" --job-name="matmul-${nodes}n-${tasks}p" \
      --output="${output_log}" \
      scripts/slurm/run_mpi_job.sh "${N_VALUE}" "${SEED_VALUE}" "${output_csv}" "${label}"
  done
done
