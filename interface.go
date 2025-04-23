package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

// Cor representa as cores do termbox
type Cor = termbox.Attribute

const (
	CorPadrao      = termbox.ColorDefault
	CorCinzaEscuro = termbox.ColorDarkGray
	CorVermelho    = termbox.ColorRed
	CorVerde       = termbox.ColorGreen
	CorParede      = termbox.ColorBlack | termbox.AttrBold | termbox.AttrDim
	CorFundoParede = termbox.ColorDarkGray
	CorTexto       = termbox.ColorDarkGray
)

// Elemento representa qualquer objeto visível no mapa
type Elemento struct {
	simbolo  rune
	cor      Cor
	corFundo Cor
	tangivel bool
}

var (
	Personagem = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Inimigo    = Elemento{'☠', CorVermelho, CorPadrao, true}
	Parede     = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio      = Elemento{' ', CorPadrao, CorPadrao, false}
	Moeda      = Elemento{'◉', CorVerde, CorPadrao, false}
	Armadilha  = Elemento{'✖', CorVermelho, CorPadrao, false}
)

// EventoTeclado representa uma ação do jogador
type EventoTeclado struct {
	Tipo  string // "sair", "interagir", "mover"
	Tecla rune   // tecla em caso de movimento
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

// Lê um evento de teclado e converte em EventoTeclado
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

// Limpa a tela do terminal
func interfaceLimparTela() {
	termbox.Clear(CorPadrao, CorPadrao)
}

// Desenha um elemento na posição (x, y)
func interfaceDesenharElemento(x, y int, elem Elemento) {
	termbox.SetCell(x, y, elem.simbolo, elem.cor, elem.corFundo)
}

// Atualiza a tela do terminal
func interfaceAtualizarTela() {
	termbox.Flush()
}

// Desenha todo o estado do jogo
func interfaceDesenharJogo(jogo *Jogo) {
	interfaceLimparTela()
	interfaceDesenharBarraDeStatus(jogo)

	for y, linha := range jogo.Mapa {
		for x, elem := range linha {
			interfaceDesenharElemento(x, y+3, elem)
		}
	}
	interfaceDesenharElemento(jogo.PosX, jogo.PosY+3, Personagem)
	interfaceAtualizarTela()
}

// Desenha a barra de status, placar e instruções
func interfaceDesenharBarraDeStatus(jogo *Jogo) {
	for i, c := range jogo.StatusMsg {
		termbox.SetCell(i, 0, c, CorTexto, CorPadrao)
	}

	scoreMsg := fmt.Sprintf("Moedas: %d", jogo.Score)
	for i, c := range scoreMsg {
		termbox.SetCell(i, 1, c, CorTexto, CorPadrao)
	}

	instr := "Use WASD para mover. E para interagir. ESC para sair."
	for i, c := range instr {
		termbox.SetCell(i, 2, c, CorTexto, CorPadrao)
	}
}
