package bid_controller

import (
	"context"
	"fullcycle-auction_go/configuration/rest_err"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (controlador *ControladorDeLance) BuscarLancesPorLeilao(c *gin.Context) {
	idDoLeilao := c.Param("auctionId")

	if erro := uuid.Validate(idDoLeilao); erro != nil {
		erroRest := rest_err.NovoErroDeRequisicaoInvalida("Campos inválidos", rest_err.Causa{
			Campo:    "auctionId",
			Mensagem: "Valor de UUID inválido",
		})

		c.JSON(erroRest.Codigo, erroRest)
		return
	}

	lances, erro := controlador.casoDeUsoDeLance.BuscarLancesPorLeilao(
		context.Background(), idDoLeilao)
	if erro != nil {
		erroRest := rest_err.ConverterErro(erro)
		c.JSON(erroRest.Codigo, erroRest)
		return
	}

	c.JSON(http.StatusOK, lances)
}
