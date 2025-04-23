// sentinela.go
package main

import (
	"math/rand"
	"time"
)



// InitSentinela dispara a goroutine da sentinela, que escuta alertaChan
func InitSentinela(jogo *Jogo, alertaChan <-chan [2]int) {
	go sentinelaRoutine(jogo, alertaChan)
}

func encontrarSentinela(jogo *Jogo) (int, int) {
	jogo.Mu.RLock()
	defer jogo.Mu.RUnlock()
	for y, linha := range jogo.Mapa {
		for x, e := range linha {
			if e.simbolo == Sentinela.simbolo {
				return x, y
			}
		}
	}
	return -1, -1 // se não encontrar
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
	var perseguirAte time.Time // deadline de perseguição

	for {
		select {
		case <-jogo.Ctx.Done():
			return

		case <-alertaChan:
			// entra em modo perseguição por 10 s
			perseguindo = true
			perseguirAte = time.Now().Add(10 * time.Second)
			jogo.Mu.Lock()
			jogo.StatusMsg = "Σ Sentinela está perseguindo você!"
			jogo.Mu.Unlock()
			sinalizarRedraw(jogo)

		case <-ticker.C:
			// expira perseguição?
			if perseguindo && time.Now().After(perseguirAte) {
				perseguindo = false
				jogo.Mu.Lock()
				jogo.StatusMsg = "Sentinela voltou a patrulhar"
				jogo.Mu.Unlock()
				sinalizarRedraw(jogo)
			}

			if perseguindo {
				// posição do jogador
				jogo.Mu.RLock()
				playerX, playerY := jogo.PosX, jogo.PosY
				jogo.Mu.RUnlock()

				nx, ny, ok := nextStep(jogo, x, y, playerX, playerY)
				if !ok {
					// sem caminho válido
					perseguindo = false
					continue
				}

				// pegou o jogador?
				if nx == playerX && ny == playerY {
					jogo.Mu.Lock()
					jogo.StatusMsg = "Σ Sentinela te pegou! GAME OVER!"
					jogo.Mu.Unlock()
					sinalizarRedraw(jogo)
					time.Sleep(500 * time.Millisecond)
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
