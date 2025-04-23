package main

import (
    "math/rand"
    "time"
)

func InitMoeda(jogo *Jogo) {
    for i := 0; i < 10; i++ {
        go moedaRoutine(jogo, 10*time.Second, 1*time.Second, 5*time.Second)
    }
}

func moedaRoutine(jogo *Jogo, activeDur, minSpawn, spawnRange time.Duration) {
    for {
        select {
        case <-jogo.Ctx.Done():
            return
        default:
        }

        // espera antes de aparecer
        wait := minSpawn + time.Duration(rand.Int63n(int64(spawnRange)))
        time.Sleep(wait)

        x, y := findRandomEmptyCell(jogo)
        jogo.Mu.Lock()
        jogo.Mapa[y][x] = Moeda
        jogo.Mu.Unlock()
        sinalizarRedraw(jogo)

        // aguarda coleta ou expiração
        timer := time.NewTimer(activeDur)
        select {
        case <-timer.C:
            // timer expirou
        case <-jogo.Ctx.Done():
            timer.Stop()
            return
        }
        timer.Stop()

        jogo.Mu.Lock()
        // se quem andou por cima já limpou o mapa, sai sem fazer nada
        if jogo.Mapa[y][x].simbolo != Moeda.simbolo {
            jogo.Mu.Unlock()
            continue
        }
        // caso contrário, expira a moeda
        jogo.Mapa[y][x] = Vazio
        jogo.StatusMsg = "⌛ Moeda expirou"
        jogo.Mu.Unlock()
        sinalizarRedraw(jogo)
    }
}
