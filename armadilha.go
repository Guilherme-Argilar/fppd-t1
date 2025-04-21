// armadilha.go - Armadilha que aparece e desaparece e spawna inimigos
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

var redrawTrap chan struct{}

var Armadilha = Elemento{
	simbolo:   '✖',
	cor:       CorVermelho,
	corFundo:  CorPadrao,
	tangivel:  false,
}

func InitArmadilha(jogo *Jogo, redrawChannel chan struct{}) {
	redrawTrap = redrawChannel
	go armadilhaRoutine(jogo, 50*time.Millisecond)
}

func armadilhaRoutine(jogo *Jogo, tempoCheck time.Duration) {
	rand.Seed(time.Now().UnixNano())

	for {
		x, y := findRandomEmptyCell(jogo)

		jogo.Mutex.Lock()
		jogo.Mapa[y][x] = Armadilha
		jogo.Mutex.Unlock()
		sinalizarRedraw()

		ativada := true

		for ativada {
			triggered := false

			jogo.Mutex.Lock()
			if jogo.PosX == x && jogo.PosY == y {
				jogo.Mutex.Unlock()
				fmt.Println("☠ Jogador pisou na armadilha!")
				os.Exit(0)
			}

			inimigos := findAllInimigos(jogo)
			for _, pos := range inimigos {
				ix, iy := pos[0], pos[1]
				if ix == x && iy == y {
					jogo.Mapa[y][x] = Vazio
					triggered = true
					break
				}
			}
			jogo.Mutex.Unlock()

			if triggered {
				nx1, ny1 := findRandomEmptyCell(jogo)
				nx2, ny2 := findRandomEmptyCell(jogo)

				jogo.Mutex.Lock()
				jogo.Mapa[ny1][nx1] = Inimigo
				jogo.Mapa[ny2][nx2] = Inimigo
				jogo.StatusMsg = "⚠ Dois inimigos surgiram de uma armadilha!"
				jogo.Mutex.Unlock()

				go inimigoRoutine(jogo, inimigoPos, inimigoPause, redrawTrap, nx1, ny1)
				go inimigoRoutine(jogo, inimigoPos, inimigoPause, redrawTrap, nx2, ny2)

				sinalizarRedraw()
				ativada = false
				break
			}

			time.Sleep(tempoCheck)
		}
	}
}

func sinalizarRedraw() {
	select {
	case redrawTrap <- struct{}{}:
	default:
	}
}

func findAllInimigos(jogo *Jogo) [][2]int {
	var posicoes [][2]int
	for y, linha := range jogo.Mapa {
		for x, e := range linha {
			if e.simbolo == Inimigo.simbolo {
				posicoes = append(posicoes, [2]int{x, y})
			}
		}
	}
	return posicoes
}