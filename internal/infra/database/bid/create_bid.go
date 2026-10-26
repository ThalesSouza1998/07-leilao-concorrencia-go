package bid

import (
	"context"
	"fullcycle-auction_go/configuration/leilao"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/entity/bid_entity"
	"fullcycle-auction_go/internal/infra/database/auction"
	"fullcycle-auction_go/internal/internal_error"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type LanceMongo struct {
	Id          string  `bson:"_id"`
	IdDoUsuario string  `bson:"user_id"`
	IdDoLeilao  string  `bson:"auction_id"`
	Valor       float64 `bson:"amount"`
	Momento     int64   `bson:"timestamp"`
}

type RepositorioDeLance struct {
	Colecao                     *mongo.Collection
	RepositorioDeLeilao         *auction.RepositorioDeLeilao
	duracaoDoLeilao             time.Duration
	statusDosLeiloes            map[string]auction_entity.StatusDoLeilao
	horariosDeEncerramento      map[string]time.Time
	mutexStatusDosLeiloes       *sync.Mutex
	mutexHorariosDeEncerramento *sync.Mutex
}

func NovoRepositorioDeLance(
	banco *mongo.Database,
	repositorioDeLeilao *auction.RepositorioDeLeilao) *RepositorioDeLance {
	return &RepositorioDeLance{
		duracaoDoLeilao:             leilao.DuracaoDoLeilao(),
		statusDosLeiloes:            make(map[string]auction_entity.StatusDoLeilao),
		horariosDeEncerramento:      make(map[string]time.Time),
		mutexStatusDosLeiloes:       &sync.Mutex{},
		mutexHorariosDeEncerramento: &sync.Mutex{},
		Colecao:                     banco.Collection("bids"),
		RepositorioDeLeilao:         repositorioDeLeilao,
	}
}

func (repositorio *RepositorioDeLance) CriarLance(
	contexto context.Context,
	lances []bid_entity.Lance) *internal_error.ErroInterno {
	var grupoDeEspera sync.WaitGroup
	for _, lance := range lances {
		grupoDeEspera.Add(1)
		go func(lanceAtual bid_entity.Lance) {
			defer grupoDeEspera.Done()

			repositorio.mutexStatusDosLeiloes.Lock()
			statusDoLeilao, statusEncontrado := repositorio.statusDosLeiloes[lanceAtual.IdDoLeilao]
			repositorio.mutexStatusDosLeiloes.Unlock()

			repositorio.mutexHorariosDeEncerramento.Lock()
			horarioDeEncerramento, horarioEncontrado :=
				repositorio.horariosDeEncerramento[lanceAtual.IdDoLeilao]
			repositorio.mutexHorariosDeEncerramento.Unlock()

			lanceMongo := &LanceMongo{
				Id:          lanceAtual.Id,
				IdDoUsuario: lanceAtual.IdDoUsuario,
				IdDoLeilao:  lanceAtual.IdDoLeilao,
				Valor:       lanceAtual.Valor,
				Momento:     lanceAtual.Momento.Unix(),
			}

			if horarioEncontrado && statusEncontrado {
				agora := time.Now()
				if statusDoLeilao == auction_entity.Fechado || agora.After(horarioDeEncerramento) {
					return
				}

				if _, erro := repositorio.Colecao.InsertOne(contexto, lanceMongo); erro != nil {
					logger.Erro("Erro ao inserir lance", erro)
					return
				}

				return
			}

			leilao, erro := repositorio.RepositorioDeLeilao.BuscarLeilaoPorId(
				contexto, lanceAtual.IdDoLeilao)
			if erro != nil {
				logger.Erro("Erro ao buscar leilão por id", erro)
				return
			}
			if leilao.Status == auction_entity.Fechado {
				return
			}

			repositorio.mutexStatusDosLeiloes.Lock()
			repositorio.statusDosLeiloes[lanceAtual.IdDoLeilao] = leilao.Status
			repositorio.mutexStatusDosLeiloes.Unlock()

			repositorio.mutexHorariosDeEncerramento.Lock()
			repositorio.horariosDeEncerramento[lanceAtual.IdDoLeilao] =
				leilao.Momento.Add(repositorio.duracaoDoLeilao)
			repositorio.mutexHorariosDeEncerramento.Unlock()

			if _, erro := repositorio.Colecao.InsertOne(contexto, lanceMongo); erro != nil {
				logger.Erro("Erro ao inserir lance", erro)
				return
			}
		}(lance)
	}
	grupoDeEspera.Wait()
	return nil
}
