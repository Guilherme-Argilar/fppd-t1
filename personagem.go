package main

import "fmt"

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

    jogo.Mu.RLock()
    canMove := jogoPodeMoverPara(jogo, nx, ny)
    jogo.Mu.RUnlock()
    if !canMove {
        return
    }

    jogo.Mu.Lock()
    jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)
    jogo.PosX, jogo.PosY = nx, ny
    select {
    case jogo.InimigoPosChan <- [2]int{jogo.PosX, jogo.PosY}:
    default:
    }
    if jogo.UltimoVisitado.simbolo == Portal.simbolo {
        select {
        case jogo.PortalEnterChan <- struct{}{}:
        default:
        }
    }
    jogo.Mu.Unlock()

    jogo.StatusMsg = fmt.Sprintf("Movendo para (%d, %d)", jogo.PosX, jogo.PosY)
}

func personagemInteragir(jogo *Jogo) {
    jogo.StatusMsg = fmt.Sprintf("Interagindo em (%d, %d)", jogo.PosX, jogo.PosY)
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
