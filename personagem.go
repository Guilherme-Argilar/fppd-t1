// personagem.g

package main

import "fmt"

func personagemMover(tecla rune, jogo *Jogo) {
	dx, dy := 0, 0
	switch tecla {
	case 'w': dy = -1
	case 'a': dx = -1
	case 's': dy = 1
	case 'd': dx = 1
	}
	nx, ny := jogo.PosX+dx, jogo.PosY+dy

	if !jogoPodeMoverPara(jogo, nx, ny) {
		return
	}

	jogo.Mu.Lock()
	jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)
	jogo.PosX, jogo.PosY = nx, ny
	jogo.Mu.Unlock()

	select {
	case jogo.InimigoPosChan <- [2]int{nx, ny}:
	default:
	}

	jogo.StatusMsg = fmt.Sprintf("Movendo para (%d, %d)", nx, ny)
	select {
	case jogo.RedrawChan <- struct{}{}:
	default:
	}
}

func personagemInteragir(jogo *Jogo) {
	jogo.Mu.Lock()
	defer jogo.Mu.Unlock()

	if jogo.UltimoVisitado.simbolo == Moeda.simbolo {
		jogo.Score++
		jogo.UltimoVisitado = Vazio
		jogo.StatusMsg = "💰 Moeda coletada!"
	} else {
		jogo.StatusMsg = "Nada aqui para interagir."
	}

	select {
	case jogo.RedrawChan <- struct{}{}:
	default:
	}
}

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