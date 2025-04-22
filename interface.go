// interface.go
package main

import "github.com/nsf/termbox-go"

// Cor representa as cores do termbox
type Cor = termbox.Attribute

const (
    CorPadrao      Cor = termbox.ColorDefault
    CorCinzaEscuro     = termbox.ColorDarkGray
    CorVermelho        = termbox.ColorRed
    CorVerde           = termbox.ColorGreen
    CorParede          = termbox.ColorBlack | termbox.AttrBold | termbox.AttrDim
    CorFundoParede     = termbox.ColorDarkGray
    CorTexto           = termbox.ColorDarkGray
    CorRoxo            = termbox.ColorMagenta
)

// Elementos visuais do jogo
var (
    Personagem = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
    Inimigo    = Elemento{'☠', CorVermelho, CorPadrao, true}
    Parede     = Elemento{'▤', CorParede, CorFundoParede, true}
    Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false}
    Vazio      = Elemento{' ', CorPadrao, CorPadrao, false}
    Portal     = Elemento{'O', CorRoxo, CorPadrao, false}
)

// EventoTeclado representa uma ação do jogador lida do teclado
type EventoTeclado struct {
    Tipo  string // "sair", "interagir", "mover"
    Tecla rune   // tecla pressionada, usada em caso de movimento
}

// Inicializa o termbox
func interfaceIniciar() {
    if err := termbox.Init(); err != nil {
        panic(err)
    }
}

// Fecha o termbox
func interfaceFinalizar() {
    termbox.Close()
}

// Lê um evento de teclado e converte para EventoTeclado
func interfaceLerEventoTeclado() EventoTeclado {
    ev := termbox.PollEvent()
    if ev.Type != termbox.EventKey {
        return EventoTeclado{}
    }
    if ev.Key == termbox.KeyEsc {
        return EventoTeclado{Tipo: "sair"}
    }
    if ev.Ch == 'e' {
        return EventoTeclado{Tipo: "interagir"}
    }
    return EventoTeclado{Tipo: "mover", Tecla: ev.Ch}
}

// Desenha o estado atual do jogo na tela
func interfaceDesenharJogo(jogo *Jogo) {
    interfaceLimparTela()
    for y, linha := range jogo.Mapa {
        for x, elem := range linha {
            interfaceDesenharElemento(x, y, elem)
        }
    }
    interfaceDesenharElemento(jogo.PosX, jogo.PosY, Personagem)
    interfaceDesenharBarraDeStatus(jogo)
    interfaceAtualizarTela()
}




// Limpa a tela do terminal
func interfaceLimparTela() {
	termbox.Clear(CorPadrao, CorPadrao)
}

// Força a atualização da tela do terminal
func interfaceAtualizarTela() {
	termbox.Flush()
}

// Desenha um elemento na posição (x, y)
func interfaceDesenharElemento(x, y int, elem Elemento) {
	termbox.SetCell(x, y, elem.simbolo, elem.cor, elem.corFundo)
}

// Exibe uma barra de status com informações úteis ao jogador
func interfaceDesenharBarraDeStatus(jogo *Jogo) {
	for i, c := range jogo.StatusMsg {
		termbox.SetCell(i, len(jogo.Mapa)+1, c, CorTexto, CorPadrao)
	}
	msg := "Use WASD para mover, E para portal. ESC para sair."
	for i, c := range msg {
		termbox.SetCell(i, len(jogo.Mapa)+3, c, CorTexto, CorPadrao)
	}
}
