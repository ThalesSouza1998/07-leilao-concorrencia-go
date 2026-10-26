package validation

import (
	"encoding/json"
	"errors"
	"fullcycle-auction_go/configuration/rest_err"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	traducao_ptbr "github.com/go-playground/validator/v10/translations/pt_BR"
)

var (
	Validador = validator.New()
	tradutor  ut.Translator
)

func init() {
	if valor, ok := binding.Validator.Engine().(*validator.Validate); ok {
		locale := pt_BR.New()
		universo := ut.New(locale, locale)
		tradutor, _ = universo.GetTranslator("pt_BR")
		traducao_ptbr.RegisterDefaultTranslations(valor, tradutor)
	}
}

func ValidarErro(erroDeValidacao error) *rest_err.ErroRest {
	var erroDeJSON *json.UnmarshalTypeError
	var erroDeCampos validator.ValidationErrors

	if errors.As(erroDeValidacao, &erroDeJSON) {
		return rest_err.NovoErroNaoEncontrado("Tipo de dado inválido")
	}

	if errors.As(erroDeValidacao, &erroDeCampos) {
		causas := []rest_err.Causa{}

		for _, erro := range erroDeValidacao.(validator.ValidationErrors) {
			causas = append(causas, rest_err.Causa{
				Campo:    erro.Field(),
				Mensagem: erro.Translate(tradutor),
			})
		}

		return rest_err.NovoErroDeRequisicaoInvalida("Valores de campos inválidos", causas...)
	}

	return rest_err.NovoErroDeRequisicaoInvalida("Erro ao converter os campos")
}
