// personagem.go
package main

import (
	"fmt"
	"math"
)

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

	jogo.Mu.Lock()
	jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)
	jogo.PosX, jogo.PosY = nx, ny
	jogo.Mu.Unlock()

	select {
	case jogo.InimigoPosChan <- [2]int{nx, ny}:
	default:
	}

	// Coleta todas as moedas em um range de 1 passo automaticamente
	coletarMoedasProximas(jogo)

	jogo.StatusMsg = fmt.Sprintf("Movendo para (%d, %d)", nx, ny)
	select {
	case jogo.RedrawChan <- struct{}{}:
	default:
	}
}

// Função para coletar todas as moedas em um range de 1 passo
func coletarMoedasProximas(jogo *Jogo) {
	// Direções para verificar (incluindo diagonais)
	dirs := [][2]int{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0}, {0, 0}, {1, 0},
		{-1, 1}, {0, 1}, {1, 1},
	}

	moedas := 0
	var sentinelaAlertada bool
	var alertaPos [2]int

	jogo.Mu.Lock()
	px, py := jogo.PosX, jogo.PosY
	jogo.Mu.Unlock()

	// Verifica se há sentinela próxima
	sx, sy := encontrarSentinela(jogo)
	distSentinela := math.Abs(float64(px-sx)) + math.Abs(float64(py-sy))

	for _, dir := range dirs {
		x, y := px+dir[0], py+dir[1]
		
		// Verifica se está dentro dos limites do mapa
		if y < 0 || y >= len(jogo.Mapa) || x < 0 || x >= len(jogo.Mapa[0]) {
			continue
		}
		
		jogo.Mu.Lock()
		elemento := jogo.Mapa[y][x]
		// Se encontrou uma moeda, coleta
		if elemento.simbolo == Moeda.simbolo {
			jogo.Mapa[y][x] = Vazio
			jogo.Score++
			moedas++
			
			// Verifica se esta coleta deve alertar a sentinela
			if !sentinelaAlertada && distSentinela <= 10 {
				sentinelaAlertada = true
				alertaPos = [2]int{px, py}
			}
		}
		jogo.Mu.Unlock()
	}

	// Atualiza mensagem se coletou moedas
	if moedas > 0 {
		jogo.Mu.Lock()
		if moedas == 1 {
			jogo.StatusMsg = "💰 Coletou 1 moeda próxima!"
		} else {
			jogo.StatusMsg = fmt.Sprintf("💰 Coletou %d moedas próximas!", moedas)
		}
		jogo.Mu.Unlock()
		
		select {
		case jogo.RedrawChan <- struct{}{}:
		default:
		}
	}

	// Envia alerta para a sentinela se necessário
	if sentinelaAlertada {
		select {
		case jogo.AlertaChan <- alertaPos:
		default:
		}
	}
}

func personagemInteragir(jogo *Jogo) {
	// 1) Detecta se coletou moeda e se está perto da sentinela
	var ativarSentinela bool
	var proxAlvo [2]int

	// lê fora de qualquer lock para evitar deadlock
	if jogo.UltimoVisitado.simbolo == Moeda.simbolo {
		// calcula distância até a sentinela
		sx, sy := encontrarSentinela(jogo)
		dist := math.Abs(float64(jogo.PosX-sx)) + math.Abs(float64(jogo.PosY-sy))
		if dist <= 10 {
			ativarSentinela = true
			proxAlvo = [2]int{jogo.PosX, jogo.PosY}
		}
	}

	// 2) Atualiza estado do jogo sob lock
	jogo.Mu.Lock()
	if jogo.UltimoVisitado.simbolo == Moeda.simbolo {
		jogo.Score++
		jogo.UltimoVisitado = Vazio
		jogo.StatusMsg = "💰 Moeda coletada!"
	} else {
		// Tenta coletar moedas próximas mesmo se não estiver sobre uma
		jogo.Mu.Unlock()
		coletarMoedasProximas(jogo)
		jogo.Mu.Lock()
		
		if jogo.StatusMsg == "💰 Moeda coletada!" || jogo.StatusMsg == "💰 Coletou 1 moeda próxima!" || jogo.StatusMsg[:5] == "💰 Col" {
			// Status já atualizado pela coleta
		} else {
			jogo.StatusMsg = "Nada aqui para interagir."
		}
	}
	jogo.Mu.Unlock()

	// 3) Redesenha
	select {
	case jogo.RedrawChan <- struct{}{}:
	default:
	}

	// 4) Dispara o alerta (fora do lock!)
	if ativarSentinela {
		jogo.AlertaChan <- proxAlvo
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