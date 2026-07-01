# trabfinalFPPD

Trabalho 2 de FPPD: multiplicacao de matrizes em Go com baseline sequencial e versao paralela com MPI.

## Objetivo

Atender ao enunciado em `enunciado.md`:

- versao sequencial (baseline);
- versao paralela com MPI;
- validacao de corretude por cantos e checksum;
- experimentos no cluster com repeticoes;
- analise com mediana, speedup e eficiencia;
- geracao dos 3 graficos obrigatorios.

## Estrutura do repositorio

- `sequencial/`: executavel baseline.
- `paralelo/`: executavel MPI.
- `internal/`: logica compartilhada (matriz, particionamento, MPI util, metricas e saida).
- `scripts/slurm/`: scripts de execucao no cluster.
- `scripts/analysis/`: consolidacao e geracao de graficos.
- `scripts/local/`: validacao local.
- `results/raw/`: logs brutos dos jobs.
- `results/processed/`: CSV por execucao e baseline sequencial.
- `results/plots/`: graficos finais.

## Requisitos

- Go 1.22+
- MPI (`mpirun`)
- dependencia Go `github.com/mvneves/gompi`

## Build e validacao local

```bash
go mod tidy
go test ./...
go build -o bin/sequencial ./sequencial
go build -o bin/paralelo ./paralelo
```

PowerShell:

```powershell
.\scripts\local\validate.ps1
```

## Execucao

Sequencial:

```bash
go run ./sequencial -n 3000 -seed 42 -csv results/processed/seq.csv
```

Paralelo:

```bash
mpirun -np 8 ./bin/paralelo -n 3000 -seed 42 -csv results/processed/run_nodes=1-tasks=8-run=1.csv -label nodes=1-tasks=8-run=1
```

No cluster (SLURM):

```bash
sbatch scripts/slurm/run_seq.sh
nohup ./scripts/slurm/run_experiments.sh &
```

## O que foi rodado no cluster (SLURM)

Os scripts em `scripts/slurm/` foram usados da seguinte forma:

- `run_seq.sh`: baseline sequencial (`seq.csv`).
- `run_experiments.sh`: orquestrador com limite de 2 jobs simultaneos (como no enunciado), compilando `bin/paralelo` e submetendo 3 repeticoes por configuracao.
- `run_mpi_job.sh`: job generico chamado pelo orquestrador, executando `mpirun ./bin/paralelo` com `-n`, `-seed`, `-csv` e `-label`.
- `run_mpi_1node.sh` e `run_mpi_multinode.sh`: scripts de execucao direta (um unico job), mantidos como alternativas.

Importante: o `run_experiments.sh` atual esta configurado para submeter apenas as configuracoes pendentes em 2 nos (`2x8` e `2x16`), pois as configuracoes de 1 no ja tinham sido executadas anteriormente e salvas em `results/processed/run_nodes=1-*.csv`.

Evidencias nos logs brutos:

- `results/raw/run_nodes=2-tasks=8-run=1.out`
- `results/raw/run_nodes=2-tasks=8-run=2.out`
- `results/raw/run_nodes=2-tasks=8-run=3.out`
- `results/raw/run_nodes=2-tasks=16-run=1.out`
- `results/raw/run_nodes=2-tasks=16-run=2.out`
- `results/raw/run_nodes=2-tasks=16-run=3.out`

## Processamento dos resultados

Consolidar tabela:

```bash
python scripts/analysis/consolidate_results.py
```

Gerar graficos:

```bash
python scripts/analysis/generate_plots.py
```

Arquivos gerados:

- `results/plots/speedup_vs_processes.png`
- `results/plots/efficiency_vs_processes.png`
- `results/plots/node_comparison.png`

## Status atual dos resultados (N=3000, seed=42)

Tempo sequencial baseline (`results/processed/seq.csv`):

- `Ts = 343.010403 s`

Configuracoes paralelas encontradas em `results/processed/run_*.csv`:

- 7 configuracoes paralelas
- 3 repeticoes por configuracao
- 21 execucoes paralelas no total

Tabela consolidada (mediana por configuracao):

|   Nós |   Processos (P) |   Tp (mediana) |   Speedup |   Efficiency |   Repetições | Obs                 |
|------:|----------------:|---------------:|----------:|-------------:|-------------:|:--------------------|
|     1 |               1 |       343.0104 |    1.0000 |       1.0000 |            1 | Baseline Sequencial |
|     1 |               1 |       702.6829 |    0.4881 |       0.4881 |            3 | Paralelo (P=1)      |
|     1 |               2 |       342.5101 |    1.0015 |       0.5007 |            3 |                     |
|     1 |               4 |       164.2645 |    2.0882 |       0.5220 |            3 |                     |
|     1 |               8 |        93.3386 |    3.6749 |       0.4594 |            3 |                     |
|     1 |              16 |        47.6061 |    7.2052 |       0.4503 |            3 |                     |
|     2 |               8 |        94.8261 |    3.6173 |       0.4522 |            3 |                     |
|     2 |              16 |        40.5280 |    8.4635 |       0.5290 |            3 |                     |

Observacoes rapidas:

- A configuração `1 nó / 32 processos` não foi executada, pois o cluster não disponibiliza essa quantidade de tarefas para um único job, o que constitui um resultado experimental sobre os limites do ambiente.
- Corretude: todos os CSVs paralelos possuem o mesmo `checksum` e os mesmos cantos (`c00`, `c0n`, `cn0`, `cnn`) da versao sequencial.
- Cobertura minima do enunciado: 8 configuracoes totais considerando baseline sequencial + 7 paralelas.
- Graficos obrigatorios do enunciado estao presentes em `results/plots/`.

## Operacao no cluster (resumo pratico)

Para manter o orquestrador rodando mesmo apos desconexao:

```bash
nohup ./scripts/slurm/run_experiments.sh &
tail -f nohup.out
squeue -u $USER
```

Transferencia de resultados para a maquina local:

```bash
scp usuario@atlantica:~/trabfinalFPPD/results/plots/*.png ./
```

## Notas para o relatorio

- No texto final, explicitar que o baseline usado para speedup e o sequencial (`Ts = 343.010403 s`), nao o caso MPI com `P=1`.
- Documentar no relatorio a escolha de `N=3000` e justificar se houve ou nao necessidade de subir para `N=4000`.
- Incluir discussao dos 3 fatores pedidos no enunciado: escalabilidade, impacto de rede (intra-no vs inter-nos) e hyperthreading.
