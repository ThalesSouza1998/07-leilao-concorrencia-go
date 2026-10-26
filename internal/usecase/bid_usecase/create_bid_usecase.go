package bid_usecase

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/bid_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"strconv"
	"time"
)

type EntradaDeLanceDTO struct {
	IdDoUsuario string  `json:"user_id"`
	IdDoLeilao  string  `json:"auction_id"`
	Valor       float64 `json:"amount"`
}

type SaidaDeLanceDTO struct {
	Id          string    `json:"id"`
	IdDoUsuario string    `json:"user_id"`
	IdDoLeilao  string    `json:"auction_id"`
	Valor       float64   `json:"amount"`
	Momento     time.Time `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}

type casoDeUsoDeLance struct {
	repositorioDeLance bid_entity.RepositorioDeLance

	temporizador              *time.Timer
	tamanhoMaximoDoLote       int
	intervaloDeInsercaoDoLote time.Duration
	canalDeLances             chan bid_entity.Lance
}

func NovoCasoDeUsoDeLance(
	repositorioDeLance bid_entity.RepositorioDeLance) CasoDeUsoDeLance {
	intervalo := obterIntervaloDeInsercaoDoLote()
	tamanhoMaximo := obterTamanhoMaximoDoLote()

	casoDeUso := &casoDeUsoDeLance{
		repositorioDeLance:        repositorioDeLance,
		tamanhoMaximoDoLote:       tamanhoMaximo,
		intervaloDeInsercaoDoLote: intervalo,
		temporizador:              time.NewTimer(intervalo),
		canalDeLances:             make(chan bid_entity.Lance, tamanhoMaximo),
	}

	casoDeUso.iniciarRotinaDeCriacao(context.Background())

	return casoDeUso
}

var loteDeLances []bid_entity.Lance

type CasoDeUsoDeLance interface {
	CriarLance(
		contexto context.Context,
		entrada EntradaDeLanceDTO) *internal_error.ErroInterno

	BuscarLanceVencedorPorLeilao(
		contexto context.Context, idDoLeilao string) (*SaidaDeLanceDTO, *internal_error.ErroInterno)

	BuscarLancesPorLeilao(
		contexto context.Context, idDoLeilao string) ([]SaidaDeLanceDTO, *internal_error.ErroInterno)
}

func (casoDeUso *casoDeUsoDeLance) iniciarRotinaDeCriacao(contexto context.Context) {
	go func() {
		defer close(casoDeUso.canalDeLances)

		for {
			select {
			case lance, aberto := <-casoDeUso.canalDeLances:
				if !aberto {
					if len(loteDeLances) > 0 {
						if erro := casoDeUso.repositorioDeLance.CriarLance(
							contexto, loteDeLances); erro != nil {
							logger.Erro("Erro ao processar o lote de lances", erro)
						}
					}
					return
				}

				loteDeLances = append(loteDeLances, lance)

				if len(loteDeLances) >= casoDeUso.tamanhoMaximoDoLote {
					if erro := casoDeUso.repositorioDeLance.CriarLance(
						contexto, loteDeLances); erro != nil {
						logger.Erro("Erro ao processar o lote de lances", erro)
					}

					loteDeLances = nil
					casoDeUso.temporizador.Reset(casoDeUso.intervaloDeInsercaoDoLote)
				}
			case <-casoDeUso.temporizador.C:
				if erro := casoDeUso.repositorioDeLance.CriarLance(
					contexto, loteDeLances); erro != nil {
					logger.Erro("Erro ao processar o lote de lances", erro)
				}
				loteDeLances = nil
				casoDeUso.temporizador.Reset(casoDeUso.intervaloDeInsercaoDoLote)
			}
		}
	}()
}

func (casoDeUso *casoDeUsoDeLance) CriarLance(
	contexto context.Context,
	entrada EntradaDeLanceDTO) *internal_error.ErroInterno {

	lance, erro := bid_entity.CriarLance(
		entrada.IdDoUsuario, entrada.IdDoLeilao, entrada.Valor)
	if erro != nil {
		return erro
	}

	casoDeUso.canalDeLances <- *lance

	return nil
}

func obterIntervaloDeInsercaoDoLote() time.Duration {
	intervalo := os.Getenv("BATCH_INSERT_INTERVAL")
	duracao, erro := time.ParseDuration(intervalo)
	if erro != nil {
		return 3 * time.Minute
	}

	return duracao
}

func obterTamanhoMaximoDoLote() int {
	valor, erro := strconv.Atoi(os.Getenv("MAX_BATCH_SIZE"))
	if erro != nil {
		return 5
	}

	return valor
}
