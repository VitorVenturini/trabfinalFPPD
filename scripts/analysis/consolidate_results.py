import pandas as pd
import glob
import re
import os

# --- Configuration ---
# Path to the directory containing the CSV results
RESULTS_DIR = 'results/processed'

# Sequential baseline time (Ts) in seconds.
# From log: matmul-seq-35426.out
TS_SEQUENTIAL = 343.010403

# --- Script ---

def parse_label(label):
    """Extracts nodes and tasks from the label string."""
    match = re.search(r'nodes=(\d+)-tasks=(\d+)-run=(\d+)', label)
    if match:
        nodes = int(match.group(1))
        tasks = int(match.group(2))
        return pd.Series([nodes, tasks])
    return pd.Series([None, None])

def main():
    """
    Consolidates parallel run results, calculates median times, speedup,
    and efficiency, and prints a summary table.
    """
    # Find all parallel run CSV files
    csv_files = glob.glob(os.path.join(RESULTS_DIR, 'run_*.csv'))
    if not csv_files:
        print(f"Error: No 'run_*.csv' files found in '{RESULTS_DIR}'.")
        print("Please make sure the experiments have run and generated results.")
        return

    # Load all CSVs into a single DataFrame
    df = pd.concat((pd.read_csv(f) for f in csv_files), ignore_index=True)

    # Parse the label to get nodes and tasks
    df[['nodes', 'tasks']] = df['label'].apply(parse_label)
    df = df.dropna(subset=['nodes', 'tasks']).astype({'nodes': int, 'tasks': int})

    # Group by configuration and calculate the median elapsed time
    # also count the number of runs to ensure we have enough data
    results = df.groupby(['nodes', 'tasks']).agg(
        tp_median=('elapsed_sec', 'median'),
        runs=('elapsed_sec', 'count')
    ).reset_index()

    # Sort the results for a clean presentation
    results = results.sort_values(by=['nodes', 'tasks'])

    # Add the sequential baseline result as the first row
    baseline = pd.DataFrame([{'nodes': 1, 'tasks': 1, 'tp_median': TS_SEQUENTIAL, 'runs': 1, 'Obs': 'Baseline Sequencial'}])
    
    # Separate the parallel run with P=1 to highlight its overhead
    parallel_p1 = results[results['tasks'] == 1].copy()
    if not parallel_p1.empty:
        parallel_p1['Obs'] = 'Paralelo (P=1)'
        results = results[results['tasks'] != 1] # Remove it from the main results
        results = pd.concat([baseline, parallel_p1, results], ignore_index=True)
    else:
        results = pd.concat([baseline, results], ignore_index=True)

    # Calculate Speedup and Efficiency
    results['Speedup'] = TS_SEQUENTIAL / results['tp_median']
    results['Efficiency'] = results['Speedup'] / results['tasks']

    # Format for final presentation
    results_final = results[['nodes', 'tasks', 'tp_median', 'Speedup', 'Efficiency', 'runs', 'Obs']].copy()
    results_final.rename(columns={
        'nodes': 'Nós', 'tasks': 'Processos (P)', 'tp_median': 'Tp (mediana)', 'runs': 'Repetições',
    }, inplace=True)
    
    results_final['Obs'].fillna('', inplace=True)

    print("--- Tabela de Resultados Consolidados ---")
    print(f"Tempo Sequencial (Ts) usado como baseline: {TS_SEQUENTIAL:.4f} segundos")
    print("\n")
    # Use tabulate format for better alignment in most terminals
    print(results_final.to_markdown(index=False, floatfmt=",.4f"))

if __name__ == "__main__":
    main()