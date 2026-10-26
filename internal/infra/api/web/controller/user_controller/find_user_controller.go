package user_controller

import (
	"context"
	"fullcycle-auction_go/configuration/rest_err"
	"fullcycle-auction_go/internal/usecase/user_usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type ControladorDeUsuario struct {
	casoDeUsoDeUsuario user_usecase.CasoDeUsoDeUsuario
}

func NovoControladorDeUsuario(
	casoDeUso user_usecase.CasoDeUsoDeUsuario) *ControladorDeUsuario {
	return &ControladorDeUsuario{
		casoDeUsoDeUsuario: casoDeUso,
	}
}

func (controlador *ControladorDeUsuario) BuscarUsuarioPorId(c *gin.Context) {
	idDoUsuario := c.Param("userId")

	if erro := uuid.Validate(idDoUsuario); erro != nil {
		erroRest := rest_err.NovoErroDeRequisicaoInvalida("Campos inválidos", rest_err.Causa{
			Campo:    "userId",
			Mensagem: "Valor de UUID inválido",
		})

		c.JSON(erroRest.Codigo, erroRest)
		return
	}

	usuario, erro := controlador.casoDeUsoDeUsuario.BuscarUsuarioPorId(
		context.Background(), idDoUsuario)
	if erro != nil {
		erroRest := rest_err.ConverterErro(erro)
		c.JSON(erroRest.Codigo, erroRest)
		return
	}

	c.JSON(http.StatusOK, usuario)
}
