package bid_usecase

import (
	"context"
	"fullcycle-auction_go/internal/internal_error"
)

func (casoDeUso *casoDeUsoDeLance) BuscarLancesPorLeilao(
	contexto context.Context,
	idDoLeilao string) ([]SaidaDeLanceDTO, *internal_error.ErroInterno) {
	lances, erro := casoDeUso.repositorioDeLance.BuscarLancesPorLeilao(contexto, idDoLeilao)
	if erro != nil {
		return nil, erro
	}

	var saidas []SaidaDeLanceDTO
	for _, lance := range lances {
		saidas = append(saidas, SaidaDeLanceDTO{
			Id:          lance.Id,
			IdDoUsuario: lance.IdDoUsuario,
			IdDoLeilao:  lance.IdDoLeilao,
			Valor:       lance.Valor,
			Momento:     lance.Momento,
		})
	}

	return saidas, nil
}

func (casoDeUso *casoDeUsoDeLance) BuscarLanceVencedorPorLeilao(
	contexto context.Context,
	idDoLeilao string) (*SaidaDeLanceDTO, *internal_error.ErroInterno) {
	lance, erro := casoDeUso.repositorioDeLance.BuscarLanceVencedorPorLeilao(
		contexto, idDoLeilao)
	if erro != nil {
		return nil, erro
	}

	saida := &SaidaDeLanceDTO{
		Id:          lance.Id,
		IdDoUsuario: lance.IdDoUsuario,
		IdDoLeilao:  lance.IdDoLeilao,
		Valor:       lance.Valor,
		Momento:     lance.Momento,
	}

	return saida, nil
}
