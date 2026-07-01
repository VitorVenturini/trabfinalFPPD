#!/bin/bash

set -euo pipefail

mkdir -p results/raw results/processed

MAX_JOBS=2 # Limite de jobs simultaneos conforme o enunciado

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
    # Espera ate que haja um slot livre na fila do SLURM
    while true; do
      current_jobs=$(squeue -u "$USER" -h | wc -l)
      if [ "$current_jobs" -lt "$MAX_JOBS" ]; then
        break
      fi
      echo "Fila cheia (${current_jobs} jobs). Aguardando 30s por um slot..."
      sleep 30
    done

    label="nodes=${nodes}-tasks=${tasks}-run=${run}"
    output_csv="results/processed/run_${label}.csv"
    output_log="results/raw/run_${label}.out"

    echo "Submitting job: nodes=${nodes}, tasks=${tasks}, run=${run}"
    sbatch --nodes="${nodes}" --ntasks="${tasks}" --job-name="matmul-${nodes}n-${tasks}p" \
      --output="${output_log}" \
      scripts/slurm/run_mpi_job.sh "${N_VALUE}" "${SEED_VALUE}" "${output_csv}" "${label}"
  done
done

echo "Todos os jobs foram submetidos. Monitore a conclusao com 'squeue -u $USER'."
