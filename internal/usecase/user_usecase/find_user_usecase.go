package user_usecase

import (
	"context"
	"fullcycle-auction_go/internal/entity/user_entity"
	"fullcycle-auction_go/internal/internal_error"
)

func NovoCasoDeUsoDeUsuario(
	repositorioDeUsuario user_entity.RepositorioDeUsuario) CasoDeUsoDeUsuario {
	return &casoDeUsoDeUsuario{
		repositorioDeUsuario,
	}
}

type casoDeUsoDeUsuario struct {
	repositorioDeUsuario user_entity.RepositorioDeUsuario
}

type SaidaDeUsuarioDTO struct {
	Id   string `json:"id"`
	Nome string `json:"name"`
}

type CasoDeUsoDeUsuario interface {
	BuscarUsuarioPorId(
		contexto context.Context,
		id string) (*SaidaDeUsuarioDTO, *internal_error.ErroInterno)
}

func (casoDeUso *casoDeUsoDeUsuario) BuscarUsuarioPorId(
	contexto context.Context, id string) (*SaidaDeUsuarioDTO, *internal_error.ErroInterno) {
	usuario, erro := casoDeUso.repositorioDeUsuario.BuscarUsuarioPorId(contexto, id)
	if erro != nil {
		return nil, erro
	}

	return &SaidaDeUsuarioDTO{
		Id:   usuario.Id,
		Nome: usuario.Nome,
	}, nil
}
