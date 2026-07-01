#!/bin/bash
#
# This is a generic SLURM job script, meant to be called by sbatch
# from an orchestrator script like run_experiments.sh.
#
# Default SBATCH directives (can be overridden by sbatch command line options)
#SBATCH --time=01:00:00
#SBATCH --mem-per-cpu=500M

set -euo pipefail

if [ "$#" -ne 4 ]; then
    echo "Usage: $0 <N> <SEED> <CSV_OUTPUT> <LABEL>"
    exit 1
fi

N_VALUE="$1"
SEED_VALUE="$2"
CSV_OUTPUT="$3"
LABEL="$4"

echo "Running parallel job with N=${N_VALUE}, tasks=${SLURM_NTASKS}, nodes=${SLURM_NNODES}"
mpirun ./bin/paralelo -n "${N_VALUE}" -seed "${SEED_VALUE}" -csv "${CSV_OUTPUT}" -label "${LABEL}"