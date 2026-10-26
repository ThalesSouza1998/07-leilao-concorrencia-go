package auction

import (
	"context"
	"fmt"
	"time"

	"fullcycle-auction_go/configuration/leilao"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const tempoLimiteFechamento = 10 * time.Second

type LeilaoMongo struct {
	Id            string                           `bson:"_id"`
	NomeDoProduto string                           `bson:"product_name"`
	Categoria     string                           `bson:"category"`
	Descricao     string                           `bson:"description"`
	Condicao      auction_entity.CondicaoDoProduto `bson:"condition"`
	Status        auction_entity.StatusDoLeilao    `bson:"status"`
	Momento       int64                            `bson:"timestamp"`
}

type RepositorioDeLeilao struct {
	Colecao *mongo.Collection
}

func NovoRepositorioDeLeilao(banco *mongo.Database) *RepositorioDeLeilao {
	return &RepositorioDeLeilao{
		Colecao: banco.Collection("auctions"),
	}
}

func (repositorio *RepositorioDeLeilao) CriarLeilao(
	contexto context.Context,
	leilao *auction_entity.Leilao) *internal_error.ErroInterno {
	leilaoMongo := &LeilaoMongo{
		Id:            leilao.Id,
		NomeDoProduto: leilao.NomeDoProduto,
		Categoria:     leilao.Categoria,
		Descricao:     leilao.Descricao,
		Condicao:      leilao.Condicao,
		Status:        leilao.Status,
		Momento:       leilao.Momento.Unix(),
	}
	_, erro := repositorio.Colecao.InsertOne(contexto, leilaoMongo)
	if erro != nil {
		logger.Erro("Erro ao inserir leilão", erro)
		return internal_error.NovoErroInternoDeServidor("Erro ao inserir leilão")
	}

	repositorio.agendarFechamento(leilao.Id)

	return nil
}

func (repositorio *RepositorioDeLeilao) agendarFechamento(idDoLeilao string) {
	fecharDepoisDoTempo(leilao.DuracaoDoLeilao(), func() {
		// A requisição já terminou quando o prazo acaba, por isso o update
		// não pode reaproveitar o contexto do POST.
		contexto, cancelar := context.WithTimeout(context.Background(), tempoLimiteFechamento)
		defer cancelar()

		repositorio.fecharLeilao(contexto, idDoLeilao)
	})
}

func (repositorio *RepositorioDeLeilao) fecharLeilao(
	contexto context.Context, idDoLeilao string) {
	filtro := bson.M{
		"_id":    idDoLeilao,
		"status": auction_entity.Aberto,
	}
	alteracao := bson.M{
		"$set": bson.M{"status": auction_entity.Fechado},
	}

	resultado, erro := repositorio.Colecao.UpdateOne(contexto, filtro, alteracao)
	if erro != nil {
		logger.Erro(fmt.Sprintf("Erro ao fechar o leilão %s", idDoLeilao), erro)
		return
	}

	if resultado.ModifiedCount == 0 {
		logger.Informacao(fmt.Sprintf("Leilão %s não estava aberto para ser fechado", idDoLeilao))
	}
}

// fecharDepoisDoTempo executa a tarefa em segundo plano, sem segurar a
// requisição que criou o leilão.
func fecharDepoisDoTempo(duracao time.Duration, fechar func()) {
	go func() {
		time.Sleep(duracao)
		fechar()
	}()
}
