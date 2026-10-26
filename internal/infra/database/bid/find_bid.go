package bid

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/bid_entity"
	"fullcycle-auction_go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (repositorio *RepositorioDeLance) BuscarLancesPorLeilao(
	contexto context.Context,
	idDoLeilao string) ([]bid_entity.Lance, *internal_error.ErroInterno) {
	filtro := bson.M{"auction_id": idDoLeilao}

	cursor, erro := repositorio.Colecao.Find(contexto, filtro)
	if erro != nil {
		logger.Erro(
			fmt.Sprintf("Erro ao buscar lances do leilão %s", idDoLeilao), erro)
		return nil, internal_error.NovoErroInternoDeServidor(
			fmt.Sprintf("Erro ao buscar lances do leilão %s", idDoLeilao))
	}

	var lancesMongo []LanceMongo
	if erro := cursor.All(contexto, &lancesMongo); erro != nil {
		logger.Erro(
			fmt.Sprintf("Erro ao buscar lances do leilão %s", idDoLeilao), erro)
		return nil, internal_error.NovoErroInternoDeServidor(
			fmt.Sprintf("Erro ao buscar lances do leilão %s", idDoLeilao))
	}

	var lances []bid_entity.Lance
	for _, lanceMongo := range lancesMongo {
		lances = append(lances, bid_entity.Lance{
			Id:          lanceMongo.Id,
			IdDoUsuario: lanceMongo.IdDoUsuario,
			IdDoLeilao:  lanceMongo.IdDoLeilao,
			Valor:       lanceMongo.Valor,
			Momento:     time.Unix(lanceMongo.Momento, 0),
		})
	}

	return lances, nil
}

func (repositorio *RepositorioDeLance) BuscarLanceVencedorPorLeilao(
	contexto context.Context,
	idDoLeilao string) (*bid_entity.Lance, *internal_error.ErroInterno) {
	filtro := bson.M{"auction_id": idDoLeilao}

	var lanceMongo LanceMongo
	opcoes := options.FindOne().SetSort(bson.D{{Key: "amount", Value: -1}})
	if erro := repositorio.Colecao.FindOne(
		contexto, filtro, opcoes).Decode(&lanceMongo); erro != nil {
		logger.Erro("Erro ao buscar o lance vencedor do leilão", erro)
		return nil, internal_error.NovoErroInternoDeServidor(
			"Erro ao buscar o lance vencedor do leilão")
	}

	return &bid_entity.Lance{
		Id:          lanceMongo.Id,
		IdDoUsuario: lanceMongo.IdDoUsuario,
		IdDoLeilao:  lanceMongo.IdDoLeilao,
		Valor:       lanceMongo.Valor,
		Momento:     time.Unix(lanceMongo.Momento, 0),
	}, nil
}
