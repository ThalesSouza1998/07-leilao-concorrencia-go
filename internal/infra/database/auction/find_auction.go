package auction

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (repositorio *RepositorioDeLeilao) BuscarLeilaoPorId(
	contexto context.Context, id string) (*auction_entity.Leilao, *internal_error.ErroInterno) {
	filtro := bson.M{"_id": id}

	var leilaoMongo LeilaoMongo
	if erro := repositorio.Colecao.FindOne(contexto, filtro).Decode(&leilaoMongo); erro != nil {
		logger.Erro(fmt.Sprintf("Erro ao buscar leilão por id = %s", id), erro)
		return nil, internal_error.NovoErroInternoDeServidor("Erro ao buscar leilão por id")
	}

	return &auction_entity.Leilao{
		Id:            leilaoMongo.Id,
		NomeDoProduto: leilaoMongo.NomeDoProduto,
		Categoria:     leilaoMongo.Categoria,
		Descricao:     leilaoMongo.Descricao,
		Condicao:      leilaoMongo.Condicao,
		Status:        leilaoMongo.Status,
		Momento:       time.Unix(leilaoMongo.Momento, 0),
	}, nil
}

func (repositorio *RepositorioDeLeilao) BuscarLeiloes(
	contexto context.Context,
	status auction_entity.StatusDoLeilao,
	categoria string,
	nomeDoProduto string) ([]auction_entity.Leilao, *internal_error.ErroInterno) {
	filtro := bson.M{}

	if status != 0 {
		filtro["status"] = status
	}

	if categoria != "" {
		filtro["category"] = categoria
	}

	if nomeDoProduto != "" {
		filtro["product_name"] = primitive.Regex{Pattern: nomeDoProduto, Options: "i"}
	}

	cursor, erro := repositorio.Colecao.Find(contexto, filtro)
	if erro != nil {
		logger.Erro("Erro ao buscar leilões", erro)
		return nil, internal_error.NovoErroInternoDeServidor("Erro ao buscar leilões")
	}
	defer cursor.Close(contexto)

	var leiloesMongo []LeilaoMongo
	if erro := cursor.All(contexto, &leiloesMongo); erro != nil {
		logger.Erro("Erro ao decodificar leilões", erro)
		return nil, internal_error.NovoErroInternoDeServidor("Erro ao decodificar leilões")
	}

	var leiloes []auction_entity.Leilao
	for _, leilaoMongo := range leiloesMongo {
		leiloes = append(leiloes, auction_entity.Leilao{
			Id:            leilaoMongo.Id,
			NomeDoProduto: leilaoMongo.NomeDoProduto,
			Categoria:     leilaoMongo.Categoria,
			Status:        leilaoMongo.Status,
			Descricao:     leilaoMongo.Descricao,
			Condicao:      leilaoMongo.Condicao,
			Momento:       time.Unix(leilaoMongo.Momento, 0),
		})
	}

	return leiloes, nil
}
