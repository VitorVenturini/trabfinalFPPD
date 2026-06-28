#!/bin/bash

set -euo pipefail

mkdir -p results/raw results/processed

configs=(
  "1 1"
  "1 2"
  "1 4"
  "1 8"
  "2 8"
  "2 16"
)

for config in "${configs[@]}"; do
  read -r nodes tasks <<< "${config}"
  for run in 1 2 3; do
    output="results/processed/parallel-n${nodes}-p${tasks}.csv"
    echo "nodes=${nodes} tasks=${tasks} run=${run}"
    mpirun -np "${tasks}" go run ./paralelo -n "${N:-3000}" -seed "${SEED:-42}" -label "nodes=${nodes}-tasks=${tasks}-run=${run}" -csv "${output}"
  done
done
