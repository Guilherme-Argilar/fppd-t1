// inimigo.go

package main

import "time"

func InitInimigo(jogo *Jogo) {
	jogo.Mu.RLock()
	for y, linha := range jogo.Mapa {
		for x, e := range linha {
			if e.simbolo == Inimigo.simbolo {
				go inimigoRoutine(jogo, jogo.InimigoPosChan, jogo.RedrawChan, x, y)
			}
		}
	}
	jogo.Mu.RUnlock()
}

func inimigoRoutine(jogo *Jogo, posChan <-chan [2]int, redrawChan chan struct{}, ix, iy int) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	px, py := jogo.PosX, jogo.PosY

	for {
		select {
		case <-jogo.Ctx.Done():
			return
		case p := <-posChan:
			px, py = p[0], p[1]
		case <-ticker.C:
			jogo.Mu.RLock()
			cur := jogo.Mapa[iy][ix]
			jogo.Mu.RUnlock()
			if cur.simbolo != Inimigo.simbolo {
				return
			}

			nx, ny, ok := nextStep(jogo, ix, iy, px, py)
			if !ok {
				continue
			}

			if nx == px && ny == py {
				jogo.Mu.Lock()
				jogo.StatusMsg = "💀 Game Over!"
				jogo.Mu.Unlock()
				jogo.GameOverChan <- struct{}{}
				return
			}

			jogo.Mu.Lock()
			jogoMoverElemento(jogo, ix, iy, nx-ix, ny-iy)
			jogo.Mapa[ny][nx] = Inimigo
			jogo.Mapa[iy][ix] = Vazio
			ix, iy = nx, ny
			jogo.Mu.Unlock()

			select {
			case redrawChan <- struct{}{}:
			default:
			}
		}
	}
}