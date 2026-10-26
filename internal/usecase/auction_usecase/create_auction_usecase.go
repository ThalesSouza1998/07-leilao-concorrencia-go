package auction_usecase

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/entity/bid_entity"
	"fullcycle-auction_go/internal/internal_error"
	"fullcycle-auction_go/internal/usecase/bid_usecase"
	"time"
)

type EntradaDeLeilaoDTO struct {
	NomeDoProduto string            `json:"product_name" binding:"required,min=1"`
	Categoria     string            `json:"category" binding:"required,min=2"`
	Descricao     string            `json:"description" binding:"required,min=10,max=200"`
	Condicao      CondicaoDoProduto `json:"condition" binding:"oneof=0 1 2"`
}

type SaidaDeLeilaoDTO struct {
	Id            string            `json:"id"`
	NomeDoProduto string            `json:"product_name"`
	Categoria     string            `json:"category"`
	Descricao     string            `json:"description"`
	Condicao      CondicaoDoProduto `json:"condition"`
	Status        StatusDoLeilao    `json:"status"`
	Momento       time.Time         `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}

type InformacoesDoVencedorDTO struct {
	Leilao SaidaDeLeilaoDTO             `json:"auction"`
	Lance  *bid_usecase.SaidaDeLanceDTO `json:"bid,omitempty"`
}

func NovoCasoDeUsoDeLeilao(
	repositorioDeLeilao auction_entity.RepositorioDeLeilao,
	repositorioDeLance bid_entity.RepositorioDeLance) CasoDeUsoDeLeilao {
	return &casoDeUsoDeLeilao{
		repositorioDeLeilao: repositorioDeLeilao,
		repositorioDeLance:  repositorioDeLance,
	}
}

type CasoDeUsoDeLeilao interface {
	CriarLeilao(
		contexto context.Context,
		entrada EntradaDeLeilaoDTO) *internal_error.ErroInterno

	BuscarLeilaoPorId(
		contexto context.Context, id string) (*SaidaDeLeilaoDTO, *internal_error.ErroInterno)

	BuscarLeiloes(
		contexto context.Context,
		status StatusDoLeilao,
		categoria, nomeDoProduto string) ([]SaidaDeLeilaoDTO, *internal_error.ErroInterno)

	BuscarLanceVencedorPorLeilao(
		contexto context.Context,
		idDoLeilao string) (*InformacoesDoVencedorDTO, *internal_error.ErroInterno)
}

type CondicaoDoProduto int64
type StatusDoLeilao int64

type casoDeUsoDeLeilao struct {
	repositorioDeLeilao auction_entity.RepositorioDeLeilao
	repositorioDeLance  bid_entity.RepositorioDeLance
}

func (casoDeUso *casoDeUsoDeLeilao) CriarLeilao(
	contexto context.Context,
	entrada EntradaDeLeilaoDTO) *internal_error.ErroInterno {
	leilao, erro := auction_entity.CriarLeilao(
		entrada.NomeDoProduto,
		entrada.Categoria,
		entrada.Descricao,
		auction_entity.CondicaoDoProduto(entrada.Condicao))
	if erro != nil {
		return erro
	}

	if erro := casoDeUso.repositorioDeLeilao.CriarLeilao(
		contexto, leilao); erro != nil {
		return erro
	}

	return nil
}
