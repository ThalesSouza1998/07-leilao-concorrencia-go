package auction

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"fullcycle-auction_go/configuration/database/mongodb"
	"fullcycle-auction_go/internal/entity/auction_entity"
)

func TestFechamentoAutomaticoDoLeilao(t *testing.T) {
	const duracaoDoLeilao = 2 * time.Second
	t.Setenv("AUCTION_DURATION", duracaoDoLeilao.String())

	urlDoMongo := os.Getenv("MONGODB_URL")
	if urlDoMongo == "" {
		urlDoMongo = "mongodb://localhost:27017"
	}
	t.Setenv("MONGODB_URL", urlDoMongo)
	t.Setenv("MONGODB_DB", fmt.Sprintf("leilao_teste_%d", time.Now().UnixNano()))

	contextoDaConexao, cancelarConexao := context.WithTimeout(
		context.Background(), 3*time.Second)
	defer cancelarConexao()

	banco, erro := mongodb.NovaConexaoMongoDB(contextoDaConexao)
	if erro != nil {
		t.Skipf("MongoDB indisponível em %s: %v", urlDoMongo, erro)
	}
	defer banco.Client().Disconnect(context.Background())
	defer banco.Drop(context.Background())

	contexto, cancelar := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelar()

	repositorio := NovoRepositorioDeLeilao(banco)
	leilao, erroDaEntidade := auction_entity.CriarLeilao(
		"Notebook Dell",
		"Informática",
		"Notebook usado, com 16GB de memória",
		auction_entity.Usado,
	)
	if erroDaEntidade != nil {
		t.Fatalf("erro ao criar a entidade do leilão: %v", erroDaEntidade)
	}

	if erro := repositorio.CriarLeilao(contexto, leilao); erro != nil {
		t.Fatalf("erro ao criar o leilão: %v", erro)
	}

	leilaoSalvo, erroDaConsulta := repositorio.BuscarLeilaoPorId(contexto, leilao.Id)
	if erroDaConsulta != nil {
		t.Fatalf("erro ao consultar o leilão recém-criado: %v", erroDaConsulta)
	}
	if leilaoSalvo.Status != auction_entity.Aberto {
		t.Fatalf("o leilão deveria ser criado aberto, status atual = %d", leilaoSalvo.Status)
	}

	aguardarStatus(t, contexto, repositorio, leilao.Id, auction_entity.Fechado)
}

func TestAgendamentoDoFechamentoNaoBloqueia(t *testing.T) {
	const espera = 200 * time.Millisecond
	fechou := make(chan time.Time, 1)
	inicio := time.Now()

	fecharDepoisDoTempo(espera, func() {
		fechou <- time.Now()
	})

	if decorrido := time.Since(inicio); decorrido >= espera {
		t.Fatalf("a chamada bloqueou a execução por %s", decorrido)
	}

	select {
	case momento := <-fechou:
		if decorrido := momento.Sub(inicio); decorrido < espera {
			t.Fatalf("o fechamento foi executado antes do prazo: %s", decorrido)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("o fechamento não foi executado")
	}
}

func aguardarStatus(
	t *testing.T,
	contexto context.Context,
	repositorio *RepositorioDeLeilao,
	idDoLeilao string,
	statusEsperado auction_entity.StatusDoLeilao) {
	t.Helper()

	var statusAtual auction_entity.StatusDoLeilao
	for limite := time.Now().Add(10 * time.Second); time.Now().Before(limite); {
		leilao, erro := repositorio.BuscarLeilaoPorId(contexto, idDoLeilao)
		if erro != nil {
			t.Fatalf("erro ao consultar o status do leilão: %v", erro)
		}

		statusAtual = leilao.Status
		if statusAtual == statusEsperado {
			return
		}

		time.Sleep(100 * time.Millisecond)
	}

	t.Fatalf("o leilão não foi fechado automaticamente, status atual = %d", statusAtual)
}
