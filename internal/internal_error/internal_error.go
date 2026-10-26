package internal_error

type ErroInterno struct {
	Mensagem string
	Codigo   string
}

func (erro *ErroInterno) Error() string {
	return erro.Mensagem
}

func NovoErroNaoEncontrado(mensagem string) *ErroInterno {
	return &ErroInterno{
		Mensagem: mensagem,
		Codigo:   "not_found",
	}
}

func NovoErroInternoDeServidor(mensagem string) *ErroInterno {
	return &ErroInterno{
		Mensagem: mensagem,
		Codigo:   "internal_server_error",
	}
}

func NovoErroDeRequisicaoInvalida(mensagem string) *ErroInterno {
	return &ErroInterno{
		Mensagem: mensagem,
		Codigo:   "bad_request",
	}
}
