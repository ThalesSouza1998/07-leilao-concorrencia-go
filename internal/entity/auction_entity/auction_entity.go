package auction_entity

import (
	"context"
	"fullcycle-auction_go/internal/internal_error"
	"github.com/google/uuid"
	"time"
)

func CriarLeilao(
	nomeDoProduto, categoria, descricao string,
	condicao CondicaoDoProduto) (*Leilao, *internal_error.ErroInterno) {
	leilao := &Leilao{
		Id:            uuid.New().String(),
		NomeDoProduto: nomeDoProduto,
		Categoria:     categoria,
		Descricao:     descricao,
		Condicao:      condicao,
		Status:        Aberto,
		Momento:       time.Now(),
	}

	if erro := leilao.Validar(); erro != nil {
		return nil, erro
	}

	return leilao, nil
}

func (leilao *Leilao) Validar() *internal_error.ErroInterno {
	if len(leilao.NomeDoProduto) <= 1 ||
		len(leilao.Categoria) <= 2 ||
		len(leilao.Descricao) <= 10 && (leilao.Condicao != Novo &&
			leilao.Condicao != Reformado &&
			leilao.Condicao != Usado) {
		return internal_error.NovoErroDeRequisicaoInvalida("leilão inválido")
	}

	return nil
}

type Leilao struct {
	Id            string
	NomeDoProduto string
	Categoria     string
	Descricao     string
	Condicao      CondicaoDoProduto
	Status        StatusDoLeilao
	Momento       time.Time
}

type CondicaoDoProduto int
type StatusDoLeilao int

const (
	Aberto StatusDoLeilao = iota
	Fechado
)

const (
	Novo CondicaoDoProduto = iota + 1
	Usado
	Reformado
)

type RepositorioDeLeilao interface {
	CriarLeilao(
		contexto context.Context,
		leilao *Leilao) *internal_error.ErroInterno

	BuscarLeiloes(
		contexto context.Context,
		status StatusDoLeilao,
		categoria, nomeDoProduto string) ([]Leilao, *internal_error.ErroInterno)

	BuscarLeilaoPorId(
		contexto context.Context, id string) (*Leilao, *internal_error.ErroInterno)
}
