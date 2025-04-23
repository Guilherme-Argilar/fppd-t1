# Jogo de Terminal em Go

Este projeto é um pequeno jogo desenvolvido em Go que roda no terminal usando a biblioteca [termbox-go](https://github.com/nsf/termbox-go). O jogador controla um personagem que pode se mover por um mapa carregado de um arquivo de texto.

## Como funciona

- O mapa é carregado de um arquivo `.txt` contendo caracteres que representam diferentes elementos do jogo.
- O personagem se move com as teclas **W**, **A**, **S**, **D**.
- Pressione **E** para interagir com o ambiente.
- Pressione **ESC** para sair do jogo.

### Controles

| Tecla | Ação |
|-------|---------------------------|
| W     | Mover para cima           |
| A     | Mover para esquerda       |
| S     | Mover para baixo          |
| D     | Mover para direita        |
| E     | Interagir                 |
| ESC   | Sair do jogo              |

## Como compilar

1. Instale o Go e clone este repositório.
2. Inicialize um novo módulo "jogo":

```bash
go mod init jogo
go get -u github.com/nsf/termbox-go
```

3. Compile o programa:

Linux:

```bash
go build -o jogo
```

Windows:

```bash
go build -o jogo.exe
```

Também é possível compilar o projeto usando `make` (Linux) ou `build.bat` (Windows).

## Como executar

1. Certifique-se de ter o arquivo `mapa.txt` com um mapa válido.
2. Execute no terminal:

```bash
./jogo
```

## Estrutura do projeto

- **main.go** — ponto de entrada e loop principal  
- **interface.go** — I/O e renderização com termbox  
- **jogo.go** — estruturas e lógica de estado  
- **personagem.go** — ações do jogador  
- **moeda.go / armadilha.go / inimigo.go / sentinela.go** — elementos concorrentes  
- **util.go** — utilidades diversas  
- **mapa.txt** — exemplo de mapa

---

## Relatório — Elementos Concorrentes & Mecânicas Implementadas

### Visão geral
O jogo foi expandido com **quatro** elementos autônomos executados em *goroutines* independentes: **Moeda**, **Armadilha**, **Inimigo** e **Sentinela**.  Cada elemento se comunica por **canais** e todas as regiões críticas são protegidas por `sync.RWMutex` (`jogo.Mu`).

### Elementos autônomos

| Elemento | Símbolo | Comportamento | Canais |
|----------|---------|---------------|--------|
| **Moeda** | `◉` | Surge em posições aleatórias. Espera até **10 s** pela mensagem em `MoedaEnterChan`; se o jogador pisar antes, é coletada; caso contrário, expira. | `MoedaEnterChan`, `RedrawChan` |
| **Armadilha** | `✖` | 20 instâncias simultâneas. Teleporta um inimigo que pisa nela; se o jogador pisar, envia `GameOver`. | `RedrawChan`, `GameOverChan` |
| **Inimigo** | `☠` | Recebe a posição do jogador via `InimigoPosChan` e persegue usando BFS a cada 1 s; finaliza o jogo se tocar o jogador. | `InimigoPosChan`, `RedrawChan`, `GameOverChan` |
| **Sentinela** | `Σ` | Patrulha; ao receber alerta (`AlertaChan`) persegue o jogador por 10 s. | `AlertaChan`, `RedrawChan`, `GameOverChan` |

### Comunicação e sincronização

* **Mutex (`jogo.Mu`)** — protege `Mapa`, posições, pontuação e mensagens.
* **Canais principais**
  * `RedrawChan` — força redesenho não bloqueante.
  * `InimigoPosChan` — broadcast da posição do jogador para inimigos.
  * `AlertaChan` — ativa perseguição da sentinela.
  * `MoedaEnterChan` — confirma coleta de moeda.
  * `GameOverChan` — sinaliza término da partida.
* **Escuta múltipla** — inimigo e sentinela fazem `select` em **3** canais (`Ctx.Done`, canal de evento, `ticker.C`).
* **Timeout real** — a moeda usa `select { case <-timer.C … case <-MoedaEnterChan … }` atendendo ao requisito de tempo‑limite.
