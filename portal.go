// portal.go
package main

import (
	"math/rand"
	"time"
)

// canal para sinalizar que o jogador entrou no portal
var portalEnter chan struct{}

// InitPortal configura o canal e dispara a goroutine de spawn aleatório
func InitPortal(jogo *Jogo) {
	portalEnter = make(chan struct{}, 1)
	rand.Seed(time.Now().UnixNano())
	// parâmetro: activeDur=5s, minSpawn=5s, spawnRange=10s
	go portalRoutine(jogo, 10*time.Second, 5*time.Second, 10*time.Second)
}

// portalRoutine faz o portal surgir em lugar aleatório, espera uso ou expira
func portalRoutine(jogo *Jogo, activeDur, minSpawn, spawnRange time.Duration) {
	for {
		// espera intervalo aleatório entre minSpawn e minSpawn+spawnRange
		wait := minSpawn + time.Duration(rand.Int63n(int64(spawnRange)))
		time.Sleep(wait)

		// escolhe célula vazia **fora** do lock
		tx, ty := findRandomEmptyCell(jogo)

		// desenha portal
		jogo.Mutex.Lock()
		jogo.Mapa[ty][tx] = Portal
		jogo.Mutex.Unlock()

		// aguarda uso ou expiração
		timer := time.NewTimer(activeDur)
		used := false
		select {
		case <-portalEnter:
			used = true
		case <-timer.C:
		}
		timer.Stop()

		if used {
			// escolhe destino **fora** do lock
			dx, dy := findRandomEmptyCell(jogo)

			jogo.Mutex.Lock()
			jogo.StatusMsg = "✅ Portal usado! Teleportando..."
			// limpa portal
			jogo.Mapa[ty][tx] = Vazio
			// move personagem
			jogo.PosX, jogo.PosY = dx, dy
			jogo.Mutex.Unlock()
		} else {
			jogo.Mutex.Lock()
			jogo.StatusMsg = "⌛ Portal expirou"
			jogo.Mapa[ty][tx] = Vazio
			jogo.Mutex.Unlock()
		}
	}
}

// findRandomEmptyCell encontra uma posição vazia não tangível
func findRandomEmptyCell(jogo *Jogo) (int, int) {
	for {
		x := rand.Intn(len(jogo.Mapa[0]))
		y := rand.Intn(len(jogo.Mapa))

		jogo.Mutex.Lock()
		cell := jogo.Mapa[y][x]
		jogo.Mutex.Unlock()

		if cell.simbolo == Vazio.simbolo {
			return x, y
		}
	}
}
