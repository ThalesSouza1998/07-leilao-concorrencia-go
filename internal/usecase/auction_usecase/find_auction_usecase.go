package auction_usecase

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"fullcycle-auction_go/internal/usecase/bid_usecase"
)

func (casoDeUso *casoDeUsoDeLeilao) BuscarLeilaoPorId(
	contexto context.Context, id string) (*SaidaDeLeilaoDTO, *internal_error.ErroInterno) {
	leilao, erro := casoDeUso.repositorioDeLeilao.BuscarLeilaoPorId(contexto, id)
	if erro != nil {
		return nil, erro
	}

	return &SaidaDeLeilaoDTO{
		Id:            leilao.Id,
		NomeDoProduto: leilao.NomeDoProduto,
		Categoria:     leilao.Categoria,
		Descricao:     leilao.Descricao,
		Condicao:      CondicaoDoProduto(leilao.Condicao),
		Status:        StatusDoLeilao(leilao.Status),
		Momento:       leilao.Momento,
	}, nil
}

func (casoDeUso *casoDeUsoDeLeilao) BuscarLeiloes(
	contexto context.Context,
	status StatusDoLeilao,
	categoria, nomeDoProduto string) ([]SaidaDeLeilaoDTO, *internal_error.ErroInterno) {
	leiloes, erro := casoDeUso.repositorioDeLeilao.BuscarLeiloes(
		contexto, auction_entity.StatusDoLeilao(status), categoria, nomeDoProduto)
	if erro != nil {
		return nil, erro
	}

	var saidas []SaidaDeLeilaoDTO
	for _, leilao := range leiloes {
		saidas = append(saidas, SaidaDeLeilaoDTO{
			Id:            leilao.Id,
			NomeDoProduto: leilao.NomeDoProduto,
			Categoria:     leilao.Categoria,
			Descricao:     leilao.Descricao,
			Condicao:      CondicaoDoProduto(leilao.Condicao),
			Status:        StatusDoLeilao(leilao.Status),
			Momento:       leilao.Momento,
		})
	}

	return saidas, nil
}

func (casoDeUso *casoDeUsoDeLeilao) BuscarLanceVencedorPorLeilao(
	contexto context.Context,
	idDoLeilao string) (*InformacoesDoVencedorDTO, *internal_error.ErroInterno) {
	leilao, erro := casoDeUso.repositorioDeLeilao.BuscarLeilaoPorId(contexto, idDoLeilao)
	if erro != nil {
		return nil, erro
	}

	saidaDoLeilao := SaidaDeLeilaoDTO{
		Id:            leilao.Id,
		NomeDoProduto: leilao.NomeDoProduto,
		Categoria:     leilao.Categoria,
		Descricao:     leilao.Descricao,
		Condicao:      CondicaoDoProduto(leilao.Condicao),
		Status:        StatusDoLeilao(leilao.Status),
		Momento:       leilao.Momento,
	}

	lanceVencedor, erro := casoDeUso.repositorioDeLance.BuscarLanceVencedorPorLeilao(
		contexto, leilao.Id)
	if erro != nil {
		logger.Erro("Não foi possível buscar o lance vencedor", erro)
		return &InformacoesDoVencedorDTO{
			Leilao: saidaDoLeilao,
			Lance:  nil,
		}, nil
	}

	saidaDoLance := &bid_usecase.SaidaDeLanceDTO{
		Id:          lanceVencedor.Id,
		IdDoUsuario: lanceVencedor.IdDoUsuario,
		IdDoLeilao:  lanceVencedor.IdDoLeilao,
		Valor:       lanceVencedor.Valor,
		Momento:     lanceVencedor.Momento,
	}

	return &InformacoesDoVencedorDTO{
		Leilao: saidaDoLeilao,
		Lance:  saidaDoLance,
	}, nil
}
