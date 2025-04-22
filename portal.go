package main

import (
    "math/rand"
    "time"
)

func InitPortal(jogo *Jogo) {
    go portalRoutine(jogo, 5*time.Second, 5*time.Second, 10*time.Second)
}

func portalRoutine(jogo *Jogo, activeDur, minSpawn, spawnRange time.Duration) {
    for {
        select {
        case <-jogo.Ctx.Done():
            return
        default:
        }

        wait := minSpawn + time.Duration(rand.Int63n(int64(spawnRange)))
        time.Sleep(wait)

        tx, ty := findRandomEmptyCell(jogo)

        jogo.Mu.Lock()
        jogo.Mapa[ty][tx] = Portal
        jogo.Mu.Unlock()
        sinalizarRedraw(jogo)

        timer := time.NewTimer(activeDur)
        used := false
        select {
        case <-jogo.PortalEnterChan:
            used = true
        case <-timer.C:
        case <-jogo.Ctx.Done():
            timer.Stop()
            return
        }
        timer.Stop()

        jogo.Mu.Lock()
        if used {
            dx, dy := findRandomEmptyCell(jogo)
            jogo.StatusMsg = "✅ Portal usado! Teleportando..."
            jogo.Mapa[ty][tx] = Vazio
            jogo.PosX, jogo.PosY = dx, dy
        } else {
            jogo.StatusMsg = "⌛ Portal expirou"
            jogo.Mapa[ty][tx] = Vazio
        }
        jogo.Mu.Unlock()
        sinalizarRedraw(jogo)
    }
}
