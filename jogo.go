// jogo.go
package main

import (
    "bufio"
    "context"
    "os"
    "sync"
)

type Elemento struct {
    simbolo  rune
    cor      Cor
    corFundo Cor
    tangivel bool
}

type Jogo struct {
    Mapa             [][]Elemento
    PosX, PosY       int
    UltimoVisitado   Elemento
    StatusMsg        string
    Mu               sync.RWMutex

    RedrawChan       chan struct{}
    InputChan        chan EventoTeclado
    GameOverChan     chan struct{}
    InimigoPosChan   chan [2]int
    InimigoPauseChan chan struct{}
    PortalEnterChan  chan struct{}

    Ctx        context.Context
    CancelFunc context.CancelFunc
}

func NewJogo(ctx context.Context, cancelFunc context.CancelFunc) *Jogo {
    return &Jogo{
        UltimoVisitado:   Vazio,
        RedrawChan:       make(chan struct{}, 1),
        InputChan:        make(chan EventoTeclado, 1),
        GameOverChan:     make(chan struct{}),
        InimigoPosChan:   make(chan [2]int, 1),
        InimigoPauseChan: make(chan struct{}, 1),
        PortalEnterChan:  make(chan struct{}, 1),
        Ctx:              ctx,
        CancelFunc:       cancelFunc,
    }
}

func jogoCarregarMapa(nome string, jogo *Jogo) error {
    arq, err := os.Open(nome)
    if err != nil {
        return err
    }
    defer arq.Close()

    scanner := bufio.NewScanner(arq)
    y := 0
    for scanner.Scan() {
        linha := scanner.Text()
        var linhaElems []Elemento
        for x, ch := range linha {
            e := Vazio
            switch ch {
            case Parede.simbolo:
                e = Parede
            case Inimigo.simbolo:
                e = Inimigo
            case Vegetacao.simbolo:
                e = Vegetacao
            case Personagem.simbolo:
                jogo.PosX, jogo.PosY = x, y
            }
            linhaElems = append(linhaElems, e)
        }
        jogo.Mu.Lock()
        jogo.Mapa = append(jogo.Mapa, linhaElems)
        jogo.Mu.Unlock()
        y++
    }
    return scanner.Err()
}

// Verifica se o personagem pode se mover para (x, y)
func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
    if y < 0 || y >= len(jogo.Mapa) || x < 0 || x >= len(jogo.Mapa[y]) {
        return false
    }
    jogo.Mu.RLock()
    defer jogo.Mu.RUnlock()
    if jogo.Mapa[y][x].tangivel {
        return false
    }
    return true
}

// Move um elemento no mapa atualizando o ÚltimoVisitado
func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
    nx, ny := x+dx, y+dy
    elemento := jogo.Mapa[y][x]
    jogo.Mapa[y][x] = jogo.UltimoVisitado
    jogo.UltimoVisitado = jogo.Mapa[ny][nx]
    jogo.Mapa[ny][nx] = elemento
}
