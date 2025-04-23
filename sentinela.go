// sentinela.go
package main

import (
	"math/rand"
	"time"
)

var Sentinela = Elemento{
	simbolo:  'Σ',
	cor:      CorVerde,
	corFundo: CorPadrao,
	tangivel: true,
}

// InitSentinela dispara a goroutine da sentinela, que escuta alertaChan
func InitSentinela(jogo *Jogo, alertaChan <-chan [2]int) {
	go sentinelaRoutine(jogo, alertaChan)
}

func sentinelaRoutine(jogo *Jogo, alertaChan <-chan [2]int) {
	// posiciona a sentinela em célula aleatória
	x, y := findRandomEmptyCell(jogo)
	jogo.Mu.Lock()
	jogo.Mapa[y][x] = Sentinela
	jogo.Mu.Unlock()
	sinalizarRedraw(jogo)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var perseguindo bool
	var perseguindoTimer *time.Timer

	for {
		select {
		case <-jogo.Ctx.Done():
			// encerra quando o contexto é cancelado
			if perseguindoTimer != nil {
				perseguindoTimer.Stop()
			}
			return

		case <-alertaChan:
			// recebe o trigger e entra em modo perseguição
			perseguindo = true
			jogo.Mu.Lock()
			jogo.StatusMsg = "Σ Sentinela está perseguindo você!"
			jogo.Mu.Unlock()
			sinalizarRedraw(jogo)

			// Cancela timer anterior se existir
			if perseguindoTimer != nil {
				perseguindoTimer.Stop()
			}

			// Cria novo timer para desativar perseguição após 10 segundos
			perseguindoTimer = time.AfterFunc(10*time.Second, func() {
				perseguindo = false
				jogo.Mu.Lock()
				jogo.StatusMsg = "Sentinela voltou a patrulhar"
				jogo.Mu.Unlock()
				sinalizarRedraw(jogo)
			})
		case <-ticker.C:
			if perseguindo {
				// lê posição atual do jogador
				jogo.Mu.RLock()
				playerX, playerY := jogo.PosX, jogo.PosY
				jogo.Mu.RUnlock()

				// calcula próximo passo em direção ao jogador
				nx, ny, ok := nextStep(jogo, x, y, playerX, playerY)
				if !ok {
					// sem caminho válido, volta a patrulhar
					perseguindo = false
					continue
				}

				// se encostar no jogador, game over
				if nx == playerX && ny == playerY {
					jogo.Mu.Lock()
					jogo.StatusMsg = "Σ Sentinela te pegou! GAME OVER!"
					jogo.Mu.Unlock()
					sinalizarRedraw(jogo)  // Forçar um redesenho imediato
					time.Sleep(500 * time.Millisecond)  // Pequena pausa para mostrar a mensagem
					jogo.GameOverChan <- struct{}{}
					return
				}

				// move a sentinela
				jogo.Mu.Lock()
				jogo.Mapa[y][x] = Vazio
				jogo.Mapa[ny][nx] = Sentinela
				x, y = nx, ny
				jogo.Mu.Unlock()
				sinalizarRedraw(jogo)
			} else {
				// patrulhamento aleatório
				dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
				rand.Shuffle(len(dirs), func(i, j int) { dirs[i], dirs[j] = dirs[j], dirs[i] })
				for _, d := range dirs {
					nx, ny := x+d[0], y+d[1]
					if jogoPodeMoverPara(jogo, nx, ny) {
						jogo.Mu.Lock()
						jogo.Mapa[y][x] = Vazio
						jogo.Mapa[ny][nx] = Sentinela
						x, y = nx, ny
						jogo.Mu.Unlock()
						sinalizarRedraw(jogo)
						break
					}
				}
			}
		}
	}
}
