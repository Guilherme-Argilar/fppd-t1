package main

import (
	"os"
	"time"
)

func main() {
	interfaceIniciar()
	defer interfaceFinalizar()

	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// Canal de redesenho do inimigo
	redrawChan := make(chan struct{}, 1)

	// Canal de input do jogador
	inputChan := make(chan EventoTeclado, 1)

	// Inicia leitor de teclado numa goroutine
	go func() {
		for {
			ev := interfaceLerEventoTeclado()
			inputChan <- ev
		}
	}()

	// Inicia elementos autônomos
	InitPortal(&jogo)
	InitInimigo(&jogo, redrawChan)
	InitArmadilha(&jogo, redrawChan)
	// Desenha primeira vez
	interfaceDesenharJogo(&jogo)

	for {
		select {
		case ev := <-inputChan:
			if continuar := personagemExecutarAcao(ev, &jogo); !continuar {
				return
			}
			interfaceDesenharJogo(&jogo)

		case <-redrawChan:
			interfaceDesenharJogo(&jogo)

		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
