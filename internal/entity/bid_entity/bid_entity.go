package bid_entity

import (
	"context"
	"fullcycle-auction_go/internal/internal_error"
	"github.com/google/uuid"
	"time"
)

type Lance struct {
	Id          string
	IdDoUsuario string
	IdDoLeilao  string
	Valor       float64
	Momento     time.Time
}

func CriarLance(
	idDoUsuario, idDoLeilao string,
	valor float64) (*Lance, *internal_error.ErroInterno) {
	lance := &Lance{
		Id:          uuid.New().String(),
		IdDoUsuario: idDoUsuario,
		IdDoLeilao:  idDoLeilao,
		Valor:       valor,
		Momento:     time.Now(),
	}

	if erro := lance.Validar(); erro != nil {
		return nil, erro
	}

	return lance, nil
}

func (lance *Lance) Validar() *internal_error.ErroInterno {
	if erro := uuid.Validate(lance.IdDoUsuario); erro != nil {
		return internal_error.NovoErroDeRequisicaoInvalida("IdDoUsuario não é um id válido")
	} else if erro := uuid.Validate(lance.IdDoLeilao); erro != nil {
		return internal_error.NovoErroDeRequisicaoInvalida("IdDoLeilao não é um id válido")
	} else if lance.Valor <= 0 {
		return internal_error.NovoErroDeRequisicaoInvalida("Valor do lance inválido")
	}

	return nil
}

type RepositorioDeLance interface {
	CriarLance(
		contexto context.Context,
		lances []Lance) *internal_error.ErroInterno

	BuscarLancesPorLeilao(
		contexto context.Context, idDoLeilao string) ([]Lance, *internal_error.ErroInterno)

	BuscarLanceVencedorPorLeilao(
		contexto context.Context, idDoLeilao string) (*Lance, *internal_error.ErroInterno)
}
