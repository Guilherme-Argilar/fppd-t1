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

    // Seed global random generator
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

    // Start keyboard input loop
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

    // Initialize autonomous elements
    InitPortal(jogo)
    InitInimigo(jogo)
    InitArmadilha(jogo)

    interfaceDesenharJogo(jogo)

    for {
        select {
        case ev := <-jogo.InputChan:
            if continuar := personagemExecutarAcao(ev, jogo); !continuar {
                cancel()
                return
            }
            interfaceDesenharJogo(jogo)

        case <-jogo.RedrawChan:
            interfaceDesenharJogo(jogo)

        case <-jogo.GameOverChan:
            cancel()
            return

        case <-ctx.Done():
            return

        default:
            time.Sleep(10 * time.Millisecond)
        }
    }
}
