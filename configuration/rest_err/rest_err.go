package rest_err

import (
	"fullcycle-auction_go/internal/internal_error"
	"net/http"
)

type ErroRest struct {
	Mensagem string  `json:"message"`
	Erro     string  `json:"err"`
	Codigo   int     `json:"code"`
	Causas   []Causa `json:"causes"`
}

type Causa struct {
	Campo    string `json:"field"`
	Mensagem string `json:"message"`
}

func (erro *ErroRest) Error() string {
	return erro.Mensagem
}

func ConverterErro(erroInterno *internal_error.ErroInterno) *ErroRest {
	switch erroInterno.Codigo {
	case "bad_request":
		return NovoErroDeRequisicaoInvalida(erroInterno.Error())
	case "not_found":
		return NovoErroNaoEncontrado(erroInterno.Error())
	default:
		return NovoErroInternoDeServidor(erroInterno.Error())
	}
}

func NovoErroDeRequisicaoInvalida(mensagem string, causas ...Causa) *ErroRest {
	return &ErroRest{
		Mensagem: mensagem,
		Erro:     "bad_request",
		Codigo:   http.StatusBadRequest,
		Causas:   causas,
	}
}

func NovoErroInternoDeServidor(mensagem string) *ErroRest {
	return &ErroRest{
		Mensagem: mensagem,
		Erro:     "internal_server",
		Codigo:   http.StatusInternalServerError,
		Causas:   nil,
	}
}

func NovoErroNaoEncontrado(mensagem string) *ErroRest {
	return &ErroRest{
		Mensagem: mensagem,
		Erro:     "not_found",
		Codigo:   http.StatusNotFound,
		Causas:   nil,
	}
}
