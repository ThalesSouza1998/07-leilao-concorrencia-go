package leilao

import (
	"fmt"
	"os"
	"time"

	"fullcycle-auction_go/configuration/logger"
)

const (
	VariavelDuracaoDoLeilao = "AUCTION_DURATION"
	duracaoPadrao           = 5 * time.Minute
)

// DuracaoDoLeilao devolve o tempo de vida de um leilão a partir de
// AUCTION_DURATION. Aceita qualquer valor válido para time.ParseDuration,
// como 30s, 10m ou 2h. Sem a variável (ou com valor inválido), usa 5 minutos.
func DuracaoDoLeilao() time.Duration {
	valor := os.Getenv(VariavelDuracaoDoLeilao)
	if valor == "" {
		return duracaoPadrao
	}

	duracao, err := time.ParseDuration(valor)
	if err != nil {
		logger.Erro("AUCTION_DURATION inválida, usando o valor padrão", err)
		return duracaoPadrao
	}

	if duracao <= 0 {
		logger.Erro(
			"AUCTION_DURATION precisa ser maior que zero, usando o valor padrão",
			fmt.Errorf("duração informada: %s", valor))
		return duracaoPadrao
	}

	return duracao
}
