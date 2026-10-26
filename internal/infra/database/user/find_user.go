package user

import (
	"context"
	"errors"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/user_entity"
	"fullcycle-auction_go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsuarioMongo struct {
	Id   string `bson:"_id"`
	Nome string `bson:"name"`
}

type RepositorioDeUsuario struct {
	Colecao *mongo.Collection
}

func NovoRepositorioDeUsuario(banco *mongo.Database) *RepositorioDeUsuario {
	return &RepositorioDeUsuario{
		Colecao: banco.Collection("users"),
	}
}

func (repositorio *RepositorioDeUsuario) BuscarUsuarioPorId(
	contexto context.Context,
	idDoUsuario string) (*user_entity.Usuario, *internal_error.ErroInterno) {
	filtro := bson.M{"_id": idDoUsuario}

	var usuarioMongo UsuarioMongo
	erro := repositorio.Colecao.FindOne(contexto, filtro).Decode(&usuarioMongo)
	if erro != nil {
		if errors.Is(erro, mongo.ErrNoDocuments) {
			logger.Erro(
				fmt.Sprintf("Usuário não encontrado com o id = %s", idDoUsuario), erro)
			return nil, internal_error.NovoErroNaoEncontrado(
				fmt.Sprintf("Usuário não encontrado com o id = %s", idDoUsuario))
		}

		logger.Erro("Erro ao buscar usuário por id", erro)
		return nil, internal_error.NovoErroInternoDeServidor("Erro ao buscar usuário por id")
	}

	usuario := &user_entity.Usuario{
		Id:   usuarioMongo.Id,
		Nome: usuarioMongo.Nome,
	}

	return usuario, nil
}
