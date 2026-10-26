package user_entity

import (
	"context"
	"fullcycle-auction_go/internal/internal_error"
)

type Usuario struct {
	Id   string
	Nome string
}

type RepositorioDeUsuario interface {
	BuscarUsuarioPorId(
		contexto context.Context, idDoUsuario string) (*Usuario, *internal_error.ErroInterno)
}
