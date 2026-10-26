package bid_controller

import (
	"context"
	"fullcycle-auction_go/configuration/rest_err"
	"fullcycle-auction_go/internal/infra/api/web/validation"
	"fullcycle-auction_go/internal/usecase/bid_usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ControladorDeLance struct {
	casoDeUsoDeLance bid_usecase.CasoDeUsoDeLance
}

func NovoControladorDeLance(
	casoDeUso bid_usecase.CasoDeUsoDeLance) *ControladorDeLance {
	return &ControladorDeLance{
		casoDeUsoDeLance: casoDeUso,
	}
}

func (controlador *ControladorDeLance) CriarLance(c *gin.Context) {
	var entrada bid_usecase.EntradaDeLanceDTO

	if erro := c.ShouldBindJSON(&entrada); erro != nil {
		erroRest := validation.ValidarErro(erro)

		c.JSON(erroRest.Codigo, erroRest)
		return
	}

	erro := controlador.casoDeUsoDeLance.CriarLance(context.Background(), entrada)
	if erro != nil {
		erroRest := rest_err.ConverterErro(erro)

		c.JSON(erroRest.Codigo, erroRest)
		return
	}

	c.Status(http.StatusCreated)
}
