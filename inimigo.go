// inimigo.go - Inimigo Perseguidor autônomo
package main

import (
	"fmt"
	"os"
	"time"
)

var (
	inimigoPos   chan [2]int
	inimigoPause chan struct{}
	redrawChan   chan struct{} // canal global para solicitar redesenho
)

func InitInimigo(jogo *Jogo, redrawChannel chan struct{}) {
	inimigoPos = make(chan [2]int, 1)
	inimigoPause = make(chan struct{}, 1)
	redrawChan = redrawChannel

	// Inicializa inimigos já existentes no mapa
	jogo.Mutex.Lock()
	for y, linha := range jogo.Mapa {
		for x, e := range linha {
			if e.simbolo == Inimigo.simbolo {
				go inimigoRoutine(jogo, inimigoPos, inimigoPause, redrawChan, x, y)
			}
		}
	}
	jogo.Mutex.Unlock()
}

func inimigoRoutine(jogo *Jogo, posChan <-chan [2]int, pauseChan <-chan struct{}, redrawChan chan struct{}, ix, iy int) {
	ticker := time.NewTicker(150 * time.Millisecond)
	defer ticker.Stop()

	px, py := jogo.PosX, jogo.PosY
	paused := false

	for {
		select {
		case p := <-posChan:
			px, py = p[0], p[1]

		case <-pauseChan:
			paused = !paused

		case <-ticker.C:
			if paused {
				continue
			}

			jogo.Mutex.Lock()
			cur := jogo.Mapa[iy][ix]
			jogo.Mutex.Unlock()
			if cur.simbolo != Inimigo.simbolo {
				return // inimigo foi removido
			}

			nx, ny, ok := nextStep(jogo, ix, iy, px, py)
			if !ok {
				continue
			}

			if nx == px && ny == py {
				jogo.Mutex.Lock()
				jogo.StatusMsg = "💀 Game Over!"
				jogo.Mutex.Unlock()
				fmt.Println("Game Over!")
				os.Exit(0)
			}

			jogo.Mutex.Lock()
			jogoMoverElemento(jogo, ix, iy, nx-ix, ny-iy)
			jogo.Mapa[ny][nx] = Inimigo
			jogo.Mapa[iy][ix] = Vazio
			ix, iy = nx, ny
			jogo.Mutex.Unlock()

			select {
			case redrawChan <- struct{}{}:
			default:
			}
		}
	}
}

func nextStep(jogo *Jogo, sx, sy, gx, gy int) (int, int, bool) {
	height := len(jogo.Mapa)
	width := len(jogo.Mapa[0])
	type pt struct{ x, y int }

	visited := make([][]bool, height)
	prev := make([][]pt, height)
	for i := range visited {
		visited[i] = make([]bool, width)
		prev[i] = make([]pt, width)
	}

	var queue []pt
	queue = append(queue, pt{sx, sy})
	visited[sy][sx] = true
	dirs := []pt{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	found := false

	for i := 0; i < len(queue); i++ {
		p := queue[i]
		if p.x == gx && p.y == gy {
			found = true
			break
		}
		for _, d := range dirs {
			nx, ny := p.x+d.x, p.y+d.y
			if ny >= 0 && ny < height && nx >= 0 && nx < len(jogo.Mapa[ny]) {
				jogo.Mutex.Lock()
				cell := jogo.Mapa[ny][nx]
				jogo.Mutex.Unlock()
				if !visited[ny][nx] && !cell.tangivel {
					visited[ny][nx] = true
					prev[ny][nx] = p
					queue = append(queue, pt{nx, ny})
				}
			}
		}
	}

	if !found {
		return 0, 0, false
	}

	cx, cy := gx, gy
	for {
		p := prev[cy][cx]
		if p.x == sx && p.y == sy {
			break
		}
		cx, cy = p.x, p.y
	}
	return cx, cy, true
}

func UpdateInimigoPosition(jogo *Jogo) {
	select {
	case inimigoPos <- [2]int{jogo.PosX, jogo.PosY}:
	default:
	}
}

func ToggleInimigoPause() {
	select {
	case inimigoPause <- struct{}{}:
	default:
	}
}
