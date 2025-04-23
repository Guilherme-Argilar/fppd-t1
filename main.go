// main.go
package main

import (
	"context"
	"math/rand"
	"os"
	"time"
)

func main() {
	interfaceIniciar()
	defer interfaceFinalizar()

	rand.Seed(time.Now().UnixNano())
	ctx, cancel := context.WithCancel(context.Background())
	jogo := NewJogo(ctx, cancel)

	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}
	if err := jogoCarregarMapa(mapaFile, jogo); err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				ev := interfaceLerEventoTeclado()
				jogo.InputChan <- ev
			}
		}
	}()

	InitMoeda(jogo)
	InitInimigo(jogo)
	InitArmadilha(jogo)

	interfaceDesenharJogo(jogo)

	for {
		select {
		case ev := <-jogo.InputChan:
			if !personagemExecutarAcao(ev, jogo) {
				cancel()
				return
			}
			interfaceDesenharJogo(jogo)

		case <-jogo.RedrawChan:
			interfaceDesenharJogo(jogo)

		case <-jogo.GameOverChan:
			jogo.StatusMsg = "💀 Game Over!"
			interfaceDesenharJogo(jogo)
			time.Sleep(2 * time.Second)
			cancel()
			return

		case <-ctx.Done():
			return
		}
	}
}