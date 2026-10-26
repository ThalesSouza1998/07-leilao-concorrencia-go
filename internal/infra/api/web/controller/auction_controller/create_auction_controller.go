package auction_controller

import (
	"context"
	"fullcycle-auction_go/configuration/rest_err"
	"fullcycle-auction_go/internal/infra/api/web/validation"
	"fullcycle-auction_go/internal/usecase/auction_usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ControladorDeLeilao struct {
	casoDeUsoDeLeilao auction_usecase.CasoDeUsoDeLeilao
}

func NovoControladorDeLeilao(
	casoDeUso auction_usecase.CasoDeUsoDeLeilao) *ControladorDeLeilao {
	return &ControladorDeLeilao{
		casoDeUsoDeLeilao: casoDeUso,
	}
}

func (controlador *ControladorDeLeilao) CriarLeilao(c *gin.Context) {
	var entrada auction_usecase.EntradaDeLeilaoDTO

	if erro := c.ShouldBindJSON(&entrada); erro != nil {
		erroRest := validation.ValidarErro(erro)

		c.JSON(erroRest.Codigo, erroRest)
		return
	}

	erro := controlador.casoDeUsoDeLeilao.CriarLeilao(context.Background(), entrada)
	if erro != nil {
		erroRest := rest_err.ConverterErro(erro)

		c.JSON(erroRest.Codigo, erroRest)
		return
	}

	c.Status(http.StatusCreated)
}
