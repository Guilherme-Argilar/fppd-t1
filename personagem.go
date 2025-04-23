// personagem.go
package main

import (
    "fmt"
    "math"
    "strings"
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
    // move o personagem e captura o elemento que estava na célula‑destino
    jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)
    jogo.PosX, jogo.PosY = nx, ny
    steppedOnMoeda := (jogo.UltimoVisitado.simbolo == Moeda.simbolo)
    // mensagem de status ainda dentro do lock para evitar race‑condition
    jogo.StatusMsg = fmt.Sprintf("Movendo para (%d, %d)", nx, ny)
    jogo.Mu.Unlock()

    // sinaliza à moeda que o jogador acabou de pisar nela (timeout‑select em moeda.go)
    if steppedOnMoeda {
        select {
        case jogo.MoedaEnterChan <- struct{}{}:
        default:
        }
    }

    // notifica inimigos da nova posição
    select {
    case jogo.InimigoPosChan <- [2]int{nx, ny}:
    default:
    }

    // coleta moedas próximas
    coletarMoedasProximas(jogo)

    // força redraw
    select {
    case jogo.RedrawChan <- struct{}{}:
    default:
    }
}

// Função para coletar todas as moedas em um range de 1 passo
func coletarMoedasProximas(jogo *Jogo) {
    dirs := [][2]int{
        {-1, -1}, {0, -1}, {1, -1},
        {-1, 0}, {0, 0}, {1, 0},
        {-1, 1}, {0, 1}, {1, 1},
    }

    moedas := 0
    var sentinelaAlertada bool
    var alertaPos [2]int

    jogo.Mu.RLock()
    px, py := jogo.PosX, jogo.PosY
    jogo.Mu.RUnlock()

    sx, sy := encontrarSentinela(jogo)
    distSentinela := math.Abs(float64(px-sx)) + math.Abs(float64(py-sy))

    for _, dir := range dirs {
        x, y := px+dir[0], py+dir[1]
        if y < 0 || y >= len(jogo.Mapa) || x < 0 || x >= len(jogo.Mapa[0]) {
            continue
        }

        jogo.Mu.Lock()
        elemento := jogo.Mapa[y][x]
        if elemento.simbolo == Moeda.simbolo {
            jogo.Mapa[y][x] = Vazio
            jogo.Score++
            moedas++
            if !sentinelaAlertada && distSentinela <= 10 {
                sentinelaAlertada = true
                alertaPos = [2]int{px, py}
            }
        }
        jogo.Mu.Unlock()
    }

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

    if sentinelaAlertada {
        select {
        case jogo.AlertaChan <- alertaPos:
        default:
        }
    }
}

func personagemInteragir(jogo *Jogo) {
    var ativarSentinela bool
    var proxAlvo [2]int

    jogo.Mu.RLock()
    ultimo := jogo.UltimoVisitado
    jogo.Mu.RUnlock()

    if ultimo.simbolo == Moeda.simbolo {
        sx, sy := encontrarSentinela(jogo)
        dist := math.Abs(float64(jogo.PosX-sx)) + math.Abs(float64(jogo.PosY-sy))
        if dist <= 10 {
            ativarSentinela = true
            proxAlvo = [2]int{jogo.PosX, jogo.PosY}
        }
    }

    jogo.Mu.Lock()
    if ultimo.simbolo == Moeda.simbolo {
        jogo.Score++
        jogo.UltimoVisitado = Vazio
        jogo.StatusMsg = "💰 Moeda coletada!"
    } else {
        jogo.Mu.Unlock()
        coletarMoedasProximas(jogo)
        jogo.Mu.Lock()
        if !strings.HasPrefix(jogo.StatusMsg, "💰 Coletou") {
            jogo.StatusMsg = "Nada aqui para interagir."
        }
    }
    jogo.Mu.Unlock()

    select {
    case jogo.RedrawChan <- struct{}{}:
    default:
    }

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
