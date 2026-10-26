package mongodb

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

const (
	variavelURLDoMongo   = "MONGODB_URL"
	variavelBancoDoMongo = "MONGODB_DB"
)

func NovaConexaoMongoDB(contexto context.Context) (*mongo.Database, error) {
	urlDoMongo := os.Getenv(variavelURLDoMongo)
	nomeDoBanco := os.Getenv(variavelBancoDoMongo)

	cliente, erro := mongo.Connect(
		contexto, options.Client().ApplyURI(urlDoMongo))
	if erro != nil {
		logger.Erro("Erro ao conectar ao banco MongoDB", erro)
		return nil, erro
	}

	if erro := cliente.Ping(contexto, nil); erro != nil {
		logger.Erro("Erro ao verificar a conexão com o MongoDB", erro)
		return nil, erro
	}

	return cliente.Database(nomeDoBanco), nil
}
