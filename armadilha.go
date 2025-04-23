
package main

import (
    "time"
)

func InitArmadilha(jogo *Jogo) {
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

            jogo.Mu.RLock()
            if jogo.PosX == x && jogo.PosY == y {
                jogo.Mu.RUnlock()
                jogo.GameOverChan <- struct{}{}
                return
            }
            jogo.Mu.RUnlock()

            triggered := false
            inimigos := findAllInimigos(jogo)
            for _, pos := range inimigos {
                if pos[0] == x && pos[1] == y {
                    triggered = true
                    break
                }
            }

            if triggered {
                jogo.Mu.Lock()
                jogo.Mapa[y][x] = Vazio
                jogo.Mu.Unlock()

                tx, ty := findRandomEmptyCell(jogo)
                jogo.Mu.Lock()
                jogo.Mapa[ty][tx] = Inimigo
                jogo.StatusMsg = "⚠ Inimigo teletransportado!"
                jogo.Mu.Unlock()

                go inimigoRoutine(jogo, jogo.InimigoPosChan, jogo.RedrawChan, tx, ty)
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
