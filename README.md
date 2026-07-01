# trabfinalFPPD

Projeto de multiplicacao de matrizes em Go com versoes sequencial e paralela usando MPI, conforme o enunciado de FPPD.

## Objetivo

Chegar aos resultados esperados do trabalho:

- implementar a versao sequencial baseline;
- implementar a versao paralela com MPI;
- validar que ambas produzem a mesma saida numerica;
- executar os experimentos no cluster;
- calcular speedup e eficiencia;
- montar os graficos e o relatorio final.

## Estrutura do repositorio

- `PLANO_PROJETO.md`: arquitetura e plano tecnico.
- `sequencial/`: executavel baseline.
- `paralelo/`: executavel MPI.
- `internal/`: logica compartilhada.
- `scripts/slurm/`: scripts de execucao no cluster.
- `scripts/local/`: validacao local.
- `results/raw/`: logs brutos.
- `results/processed/`: CSVs consolidados.
- `report/`: local para o relatorio final em PDF e arquivos relacionados.
- `results/plots/`: graficos finais.

## Passo a passo

## 1. Ler o enunciado e alinhar o escopo

Antes de rodar qualquer experimento, confirmar os requisitos em:

- `enunciado.md`
- `PLANO_PROJETO.md`

O que precisa obrigatoriamente existir no final:

- versao sequencial;
- versao paralela com MPI;
- pelo menos 8 configuracoes de experimento contando o baseline;
- 3 repeticoes por configuracao;
- mediana dos tempos;
- tabela com tempo, speedup e eficiencia;
- graficos pedidos no enunciado;
- relatorio final.

## 2. Preparar o ambiente local

Instalar e validar:

- Go 1.22 ou superior;
- MPI com `mpirun`;
- acesso ao pacote `github.com/mvneves/gompi`.

Conferencias minimas:

```bash
go version
mpirun --version
```

Se o projeto ainda nao tiver dependencias baixadas:

```bash
go mod tidy
```

## 3. Validar a base localmente

Rodar os testes das partes puras:

```bash
go test ./...
```

Ou no PowerShell:

```powershell
.\scripts\local\validate.ps1
```

Objetivo desta etapa:

- confirmar que a multiplicacao sequencial esta correta;
- confirmar que o particionamento por linhas esta correto;
- confirmar que a verificacao de checksum e cantos funciona.

## 4. Validar a versao sequencial

Executar primeiro com matriz pequena:

```bash
go run ./sequencial -n 4 -seed 42
go run ./sequencial -n 8 -seed 42
go run ./sequencial -n 16 -seed 42
```

Depois executar com um tamanho maior de teste:

```bash
go run ./sequencial -n 256 -seed 42
```

Conferir:

- o programa executa sem erro;
- a saida mostra tempo;
- a saida mostra `c00`, `c0n`, `cn0`, `cnn` e `checksum`.

## 5. Validar a versao paralela localmente

Primeiro, compile o executavel:

```bash
go build -o bin/paralelo ./paralelo
```

Depois, execute com tamanhos pequenos e poucos processos:

```bash
mpirun -np 2 ./bin/paralelo -n 4 -seed 42
mpirun --oversubscribe -np 4 ./bin/paralelo -n 16 -seed 42
```

> **Observacao:** Se ao rodar localmente com `-np` maior que o numero de nucleos do seu processador voce receber um erro sobre "not enough slots", adicione a flag `--oversubscribe` ao comando `mpirun`.

Nesta etapa, comparar manualmente a saida da versao paralela com a sequencial para o mesmo `n` e a mesma `seed`.

Os valores abaixo devem ser identicos entre as duas execucoes:

- `c00`
- `c0n`
- `cn0`
- `cnn`
- `checksum`

Se os valores divergirem, nao avancar para o cluster antes de corrigir.

## 6. Testar casos importantes antes do cluster

Rodar casos que cobrem divisao desigual:
```bash
mpirun --oversubscribe -np 3 ./bin/paralelo -n 10 -seed 42
mpirun --oversubscribe -np 4 ./bin/paralelo -n 10 -seed 42
mpirun --oversubscribe -np 6 ./bin/paralelo -n 10 -seed 42
```

Objetivo:

- validar `N % P != 0`;
- validar ranks com menos linhas;
- reduzir risco de erro de particionamento.

## 7. Definir o tamanho final do problema

O enunciado pede:

- usar `N = 3000` como padrao;
- se o sequencial no cluster ficar abaixo de 3 minutos, subir para `N = 4000`.

Procedimento:

1. Rodar o baseline sequencial no cluster com `N=3000`.
2. Medir o tempo.
3. Se o tempo for menor que 3 minutos, repetir com `N=4000`.
4. Fixar o valor final de `N` para todos os experimentos.

Esse valor deve ser documentado no relatorio.

## 8. Preparar o cluster

Antes de submeter jobs, ajustar no cluster:

1.  **Dar permissao de execucao aos scripts:**
    ```bash
    chmod +x scripts/slurm/*.sh
    ```
2.  **Carregar modulos de ambiente (se necessario):**
    Verificar se `go` e `mpirun` estao disponiveis. Se nao, usar `module load` ou `export PATH=...`.

3.  **Revisar os scripts SLURM:**
    Conferir os scripts em `scripts/slurm/` para ajustar parametros como tempo de walltime, numero de nos/tarefas, etc., conforme a necessidade do cluster.

## 9. Rodar o baseline sequencial no cluster

Submeter primeiro o baseline:

```bash
sbatch scripts/slurm/run_seq.sh
```

Ou executar manualmente:

```bash
go run ./sequencial -n 3000 -seed 42 -csv results/processed/seq.csv
```

Guardar:

- tempo do baseline `Ts`;
- saida de verificacao;
- valor final de `N`.

## 10. Rodar a versao paralela no cluster

Executar configuracoes que cubram os 3 fatores do enunciado:

- escalabilidade por numero de processos;
- comparacao intra-no vs inter-nos;
- impacto de hyperthreading.

Conjunto inicial sugerido:
O script `scripts/slurm/run_experiments.sh` ja esta configurado com um conjunto de experimentos que cobre esses fatores.

### Executando o Orquestrador de Forma Segura

Como o script orquestrador precisa rodar por um longo tempo para submeter todos os jobs (respeitando o limite de 2 jobs na fila do cluster), ele deve ser executado de forma que não seja interrompido se sua conexão cair. Para isso, usamos o comando `nohup`.

A partir do diretorio raiz do projeto (`~/trabfinalFPPD`), execute:

```bash
nohup ./scripts/slurm/run_experiments.sh &
```

Isso iniciara o script em segundo plano. Voce pode acompanhar o progresso vendo o arquivo `nohup.out` (`tail -f nohup.out`) e monitorar a fila com `squeue -u $USER`. Apos iniciar o script, voce pode se desconectar do cluster com seguranca.

## 11. Repetir cada configuracao 3 vezes

O enunciado exige pelo menos 3 execucoes por configuracao.

Para cada combinacao de nos e processos:

1. rodar 3 vezes;
2. registrar os 3 tempos;
3. calcular a mediana;
4. usar a mediana nas tabelas e graficos.

O script `scripts/slurm/run_experiments.sh` ja esta configurado para automatizar este processo, submetendo 3 jobs para cada uma das configuracoes definidas.

## 12. Salvar todos os resultados

Durante os experimentos, organizar os arquivos assim:

- logs brutos em `results/raw/`;
- CSVs em `results/processed/`;
- graficos em `results/plots/`.

Cada execucao deve preservar:

- modo;
- `N`;
- numero de nos;
- numero de processos;
- tempo;
- checksum;
- valores de verificacao.

Nao misturar resultados de testes pequenos com resultados finais do relatorio.

## 13. Conferir corretude antes da analise final

Antes de calcular speedup, verificar que as execucoes finais mantiveram a mesma saida numerica.

Conferir:

- baseline sequencial;
- paralelo com `P=1`;
- paralelo com `P>1`.

Se checksum ou cantos diferirem, os tempos nao servem para o relatorio.

## 14. Calcular speedup e eficiencia

Com os tempos medianos:

- `Speedup = Ts / Tp`
- `Eficiencia = Speedup / P`

Onde:

- `Ts` = tempo sequencial baseline;
- `Tp` = tempo paralelo da configuracao;
- `P` = numero de processos.

Esses dados devem ser colocados em uma tabela final.

## 15. Montar os graficos obrigatorios

Produzir pelo menos estes 3 graficos:

1. `Speedup vs numero de processos`
2. `Eficiencia vs numero de processos`
3. `Tempo intra-no vs inter-nos`

Adicionar referencias ideais:

- reta `speedup ideal = P`
- reta `eficiencia ideal = 1`

## 16. Responder as perguntas do relatorio

Com os dados em maos, responder:

1. o speedup foi sub-linear, linear ou super-linear;
2. a partir de que ponto a eficiencia caiu;
3. qual foi o impacto da comunicacao em rede;
4. se hyperthreading ajudou ou nao;
5. qual a fracao paralelizavel estimada pela Lei de Amdahl.

Nao responder de forma generica. As respostas devem citar os resultados obtidos.

## 17. Fechar o relatorio final

O relatorio deve incluir:

- modelo de paralelismo escolhido;
- justificativa da escolha;
- como os dados sao distribuidos;
- como os resultados sao coletados;
- tabela com tempos, speedup e eficiencia;
- graficos obrigatorios;
- discussao das perguntas do enunciado;
- valor final de `N`;
- instrucoes de compilacao e execucao.

## 18. Checklist final de entrega

Antes de entregar, confirmar:

- `go.mod` presente;
- versao sequencial funcionando;
- versao paralela funcionando;
- instrucoes no `README.md`;
- scripts SLURM revisados;
- resultados brutos guardados;
- tabela final pronta;
- graficos prontos;
- relatorio pronto.

## 19. Baixando os Resultados do Cluster
 
 Após gerar os resultados (CSVs, gráficos, etc.) no cluster, você precisará transferi-los para sua máquina local para incluí-los no relatório.
 
 ### Método 1: Usando Git (Recomendado)
 
 Se você já está usando Git, a forma mais simples e organizada é comitar os resultados.
 
 1.  **No cluster:** Adicione, comite e envie os resultados para o seu repositório.
     ```bash
     git add results/ scripts/analysis/
     git commit -m "Adiciona resultados e gráficos dos experimentos"
     git push
     ```
 2.  **No seu PC local:** Abra o projeto no VS Code e puxe as atualizações.
     ```bash
     git pull
     ```
 
 ### Método 2: Usando `scp` ou Cliente SFTP
 
 Se preferir não comitar os arquivos de resultado, você pode usar `scp` (Secure Copy) ou um cliente SFTP como WinSCP ou FileZilla para transferir os arquivos manualmente.
 
 ```bash
 # Exemplo de SCP para baixar os gráficos
 scp fppd3217@atlantica:~/trabfinalFPPD/results/plots/*.png ./
 ```

## Comandos principais

Sequencial:

```bash
go run ./sequencial -n 3000 -seed 42
```

Paralelo:

```bash
mpirun -np 4 go run ./paralelo -n 3000 -seed 42
```

Testes:

```bash
go test ./...
```

Script local no PowerShell:

```powershell
.\scripts\local\validate.ps1
```

## Observacao

O ambiente desta sessao nao possui `go` nem `mpirun` no `PATH`, entao o codigo foi implementado, mas nao validado aqui por compilacao nem execucao real. A validacao precisa ser feita no ambiente local de voces ou diretamente no cluster.
