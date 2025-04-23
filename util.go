// util.go
package main

import (
    "math/rand"
)

// findRandomEmptyCell retorna coordenadas de uma célula vazia no mapa
func findRandomEmptyCell(jogo *Jogo) (int, int) {
    height := len(jogo.Mapa)
    width := len(jogo.Mapa[0])
    for {
        x := rand.Intn(width)
        y := rand.Intn(height)

        jogo.Mu.RLock()
        cell := jogo.Mapa[y][x]
        jogo.Mu.RUnlock()

        if cell.simbolo == Vazio.simbolo {
            return x, y
        }
    }
}

// nextStep calcula o próximo passo na direção do jogador (gx, gy) a partir de (sx, sy)
// usando BFS para garantir um caminho válido
func nextStep(jogo *Jogo, sx, sy, gx, gy int) (int, int, bool) {
    height := len(jogo.Mapa)
    width := len(jogo.Mapa[0])

    type pt struct{ x, y int }
    visited := make([][]bool, height)
    prev := make([][]pt, height)
    for i := range visited {
        visited[i] = make([]bool, width)
        prev[i] = make([]pt, width)
    }

    var queue []pt
    queue = append(queue, pt{sx, sy})
    visited[sy][sx] = true
    dirs := []pt{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
    found := false

    for i := 0; i < len(queue); i++ {
        p := queue[i]
        if p.x == gx && p.y == gy {
            found = true
            break
        }
        for _, d := range dirs {
            nx, ny := p.x+d.x, p.y+d.y
            if ny >= 0 && ny < height && nx >= 0 && nx < width {
                jogo.Mu.RLock()
                cell := jogo.Mapa[ny][nx]
                jogo.Mu.RUnlock()
                if !visited[ny][nx] && !cell.tangivel {
                    visited[ny][nx] = true
                    prev[ny][nx] = p
                    queue = append(queue, pt{nx, ny})
                }
            }
        }
    }

    if !found {
        return 0, 0, false
    }

    cx, cy := gx, gy
    for {
        p := prev[cy][cx]
        if p.x == sx && p.y == sy {
            break
        }
        cx, cy = p.x, p.y
    }
    return cx, cy, true
}

func encontrarSentinela(jogo *Jogo) (int, int) {
	jogo.Mu.RLock()
	defer jogo.Mu.RUnlock()
	for y, linha := range jogo.Mapa {
		for x, e := range linha {
			if e.simbolo == Sentinela.simbolo {
				return x, y
			}
		}
	}
	return -1, -1 // se não encontrar
}
