import pandas as pd
import glob
import re
import os
import matplotlib.pyplot as plt
import seaborn as sns

# --- Configuration ---
RESULTS_DIR = 'results/processed'
PLOTS_DIR = 'results/plots'
TS_SEQUENTIAL = 343.010403

# --- Data Processing (reused from consolidate_results.py) ---

def parse_label(label):
    """Extracts nodes and tasks from the label string."""
    match = re.search(r'nodes=(\d+)-tasks=(\d+)-run=(\d+)', label)
    if match:
        nodes = int(match.group(1))
        tasks = int(match.group(2))
        return pd.Series([nodes, tasks])
    return pd.Series([None, None])

def get_consolidated_results():
    """Loads and processes all CSVs to return a consolidated DataFrame."""
    csv_files = glob.glob(os.path.join(RESULTS_DIR, 'run_*.csv'))
    if not csv_files:
        raise FileNotFoundError(f"No 'run_*.csv' files found in '{RESULTS_DIR}'.")

    df = pd.concat((pd.read_csv(f) for f in csv_files), ignore_index=True)
    df[['nodes', 'tasks']] = df['label'].apply(parse_label)
    df = df.dropna(subset=['nodes', 'tasks']).astype({'nodes': int, 'tasks': int})

    results = df.groupby(['nodes', 'tasks']).agg(
        tp_median=('elapsed_sec', 'median')
    ).reset_index()

    results = results.sort_values(by=['nodes', 'tasks'])
    
    # Add sequential baseline for context
    baseline = pd.DataFrame([{'nodes': 1, 'tasks': 1, 'tp_median': TS_SEQUENTIAL}])
    
    all_results = pd.concat([baseline, results], ignore_index=True).drop_duplicates(subset=['nodes', 'tasks'], keep='last')

    all_results['Speedup'] = TS_SEQUENTIAL / all_results['tp_median']
    all_results['Efficiency'] = all_results['Speedup'] / all_results['tasks']
    
    return all_results

# --- Plotting Functions ---

def plot_speedup(df):
    """Generates and saves the Speedup vs. Processes plot."""
    df_single_node = df[(df['nodes'] == 1) & (df['tasks'] > 1)].sort_values('tasks')
    
    plt.figure(figsize=(10, 6))
    
    max_p = df_single_node['tasks'].max()
    plt.plot([1, max_p], [1, max_p], 'r--', label='Speedup Ideal (Sp = P)')
    
    plt.plot(df_single_node['tasks'], df_single_node['Speedup'], 'o-', label='Speedup Experimental (1 Nó)')
    
    plt.title('Speedup vs. Número de Processos (Intra-Nó)')
    plt.xlabel('Número de Processos (P)')
    plt.ylabel('Speedup (Sp)')
    plt.grid(True, which='both', linestyle='--', linewidth=0.5)
    plt.xticks(df_single_node['tasks'])
    plt.legend()
    
    plt.savefig(os.path.join(PLOTS_DIR, 'speedup_vs_processes.png'))
    print(f"Plot saved to {os.path.join(PLOTS_DIR, 'speedup_vs_processes.png')}")
    plt.close()

def plot_efficiency(df):
    """Generates and saves the Efficiency vs. Processes plot."""
    df_single_node = df[(df['nodes'] == 1) & (df['tasks'] > 1)].sort_values('tasks')
    
    plt.figure(figsize=(10, 6))
    
    plt.axhline(y=1.0, color='r', linestyle='--', label='Eficiência Ideal (E = 1)')
    
    plt.plot(df_single_node['tasks'], df_single_node['Efficiency'], 'o-', label='Eficiência Experimental (1 Nó)')
    
    plt.title('Eficiência vs. Número de Processos (Intra-Nó)')
    plt.xlabel('Número de Processos (P)')
    plt.ylabel('Eficiência (E)')
    plt.ylim(0, 1.1)
    plt.grid(True, which='both', linestyle='--', linewidth=0.5)
    plt.xticks(df_single_node['tasks'])
    plt.legend()
    
    plt.savefig(os.path.join(PLOTS_DIR, 'efficiency_vs_processes.png'))
    print(f"Plot saved to {os.path.join(PLOTS_DIR, 'efficiency_vs_processes.png')}")
    plt.close()

def plot_node_comparison(df):
    """Generates and saves the Intra-node vs. Inter-node time comparison plot."""
    df_compare = df[df['tasks'].isin([8, 16])].copy()
    df_compare['Config'] = df_compare['nodes'].astype(str) + ' Nó(s)'
    
    plt.figure(figsize=(10, 6))
    
    sns.barplot(data=df_compare, x='tasks', y='tp_median', hue='Config')
    
    plt.title('Comparação de Tempo: Intra-Nó vs. Inter-Nós')
    plt.xlabel('Número Total de Processos (P)')
    plt.ylabel('Tempo de Execução (mediana em segundos)')
    plt.grid(axis='y', linestyle='--', linewidth=0.5)
    
    plt.savefig(os.path.join(PLOTS_DIR, 'node_comparison.png'))
    print(f"Plot saved to {os.path.join(PLOTS_DIR, 'node_comparison.png')}")
    plt.close()

def main():
    """Main function to generate all plots."""
    sns.set_theme(style="whitegrid")
    os.makedirs(PLOTS_DIR, exist_ok=True)
    
    try:
        results_df = get_consolidated_results()
        plot_speedup(results_df)
        plot_efficiency(results_df)
        plot_node_comparison(results_df)
    except FileNotFoundError as e:
        print(e)

if __name__ == "__main__":
    main()