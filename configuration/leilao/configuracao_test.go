package leilao

import (
	"testing"
	"time"
)

func TestDuracaoDoLeilao(t *testing.T) {
	casos := []struct {
		nome     string
		valor    string
		esperado time.Duration
	}{
		{nome: "valor em segundos", valor: "45s", esperado: 45 * time.Second},
		{nome: "valor em minutos", valor: "10m", esperado: 10 * time.Minute},
		{nome: "valor inválido", valor: "meia hora", esperado: duracaoPadrao},
		{nome: "sem valor", valor: "", esperado: duracaoPadrao},
		{nome: "valor zerado", valor: "0s", esperado: duracaoPadrao},
		{nome: "valor negativo", valor: "-5m", esperado: duracaoPadrao},
	}

	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			t.Setenv(VariavelDuracaoDoLeilao, caso.valor)

			if duracao := DuracaoDoLeilao(); duracao != caso.esperado {
				t.Fatalf("duração obtida = %s, esperada = %s", duracao, caso.esperado)
			}
		})
	}
}
