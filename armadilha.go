// armadilha.go
package main

import "time"

// Armadilha: aparece e desaparece, spawna inimigos
var Armadilha = Elemento{
    simbolo:  '✖',
    cor:      CorVermelho,
    corFundo: CorPadrao,
    tangivel: false,
}

func InitArmadilha(jogo *Jogo) {
    // Cria 20 armadilhas concorrentes ao iniciar
    for i := 0; i < 20; i++ {
        go armadilhaRoutine(jogo, 50*time.Millisecond)
    }
}

func armadilhaRoutine(jogo *Jogo, tempoCheck time.Duration) {
    for {
        select {
        case <-jogo.Ctx.Done():
            return
        default:
        }

        x, y := findRandomEmptyCell(jogo)

        jogo.Mu.Lock()
        jogo.Mapa[y][x] = Armadilha
        jogo.Mu.Unlock()
        sinalizarRedraw(jogo)

        ativada := true
        for ativada {
            select {
            case <-jogo.Ctx.Done():
                return
            default:
            }

            // jogador pisa na armadilha?
            jogo.Mu.RLock()
            if jogo.PosX == x && jogo.PosY == y {
                jogo.Mu.RUnlock()
                jogo.GameOverChan <- struct{}{}
                return
            }
            jogo.Mu.RUnlock()

            // inimigo ativa armadilha?
            triggered := false
            inimigos := findAllInimigos(jogo)
            for _, pos := range inimigos {
                if pos[0] == x && pos[1] == y {
                    triggered = true
                    break
                }
            }

            if triggered {
                nx1, ny1 := findRandomEmptyCell(jogo)
                nx2, ny2 := findRandomEmptyCell(jogo)

                jogo.Mu.Lock()
                jogo.Mapa[ny1][nx1] = Inimigo
                jogo.Mapa[ny2][nx2] = Inimigo
                jogo.StatusMsg = "⚠ Dois inimigos surgiram de uma armadilha!"
                jogo.Mu.Unlock()

                go inimigoRoutine(jogo, jogo.InimigoPosChan, jogo.InimigoPauseChan, jogo.RedrawChan, nx1, ny1)
                go inimigoRoutine(jogo, jogo.InimigoPosChan, jogo.InimigoPauseChan, jogo.RedrawChan, nx2, ny2)

                sinalizarRedraw(jogo)
                ativada = false
            } else {
                time.Sleep(tempoCheck)
            }
        }
    }
}

func sinalizarRedraw(jogo *Jogo) {
    select {
    case jogo.RedrawChan <- struct{}{}:
    default:
    }
}

func findAllInimigos(jogo *Jogo) [][2]int {
    var posicoes [][2]int
    jogo.Mu.RLock()
    defer jogo.Mu.RUnlock()
    for y, linha := range jogo.Mapa {
        for x, e := range linha {
            if e.simbolo == Inimigo.simbolo {
                posicoes = append(posicoes, [2]int{x, y})
            }
        }
    }
    return posicoes
}
