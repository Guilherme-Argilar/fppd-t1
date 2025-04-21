// personagem.go
package main

import "fmt"

// Move o personagem (WASD), notifica inimigo e sinaliza portal
func personagemMover(tecla rune, jogo *Jogo) {
	dx, dy := 0, 0
	switch tecla {
	case 'w':
		dy = -1
	case 'a':
		dx = -1
	case 's':
		dy = 1
	case 'd':
		dx = 1
	}
	nx, ny := jogo.PosX+dx, jogo.PosY+dy
	if !jogoPodeMoverPara(jogo, nx, ny) {
		return
	}

	jogo.Mutex.Lock()
	// Move no mapa
	jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)
	jogo.PosX, jogo.PosY = nx, jogo.PosY+dy
	// Notifica o inimigo da posição atual do jogador
	UpdateInimigoPosition(jogo)
	// Se pisou no portal, sinaliza seu uso
	if jogo.UltimoVisitado.simbolo == Portal.simbolo {
		select {
		case portalEnter <- struct{}{}:
		default:
		}
	}
	jogo.Mutex.Unlock()

	jogo.StatusMsg = fmt.Sprintf("Movendo para (%d, %d)", jogo.PosX, jogo.PosY)
}

// Interação (não afeta portal nem inimigo)
func personagemInteragir(jogo *Jogo) {
	jogo.StatusMsg = fmt.Sprintf("Interagindo em (%d, %d)", jogo.PosX, jogo.PosY)
}

// Processa evento do teclado
func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo) bool {
	switch ev.Tipo {
	case "sair":
		return false
	case "interagir":
		personagemInteragir(jogo)
	case "mover":
		personagemMover(ev.Tecla, jogo)
	}
	return true
}
