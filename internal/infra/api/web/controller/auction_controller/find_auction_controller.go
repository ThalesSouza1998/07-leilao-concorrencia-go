package auction_controller

import (
	"context"
	"fullcycle-auction_go/configuration/rest_err"
	"fullcycle-auction_go/internal/usecase/auction_usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

func (controlador *ControladorDeLeilao) BuscarLeilaoPorId(c *gin.Context) {
	idDoLeilao := c.Param("auctionId")

	if erro := uuid.Validate(idDoLeilao); erro != nil {
		erroRest := rest_err.NovoErroDeRequisicaoInvalida("Campos inválidos", rest_err.Causa{
			Campo:    "auctionId",
			Mensagem: "Valor de UUID inválido",
		})

		c.JSON(erroRest.Codigo, erroRest)
		return
	}

	leilao, erro := controlador.casoDeUsoDeLeilao.BuscarLeilaoPorId(
		context.Background(), idDoLeilao)
	if erro != nil {
		erroRest := rest_err.ConverterErro(erro)
		c.JSON(erroRest.Codigo, erroRest)
		return
	}

	c.JSON(http.StatusOK, leilao)
}

func (controlador *ControladorDeLeilao) BuscarLeiloes(c *gin.Context) {
	status := c.Query("status")
	categoria := c.Query("category")
	nomeDoProduto := c.Query("productName")

	statusNumerico, erroConversao := strconv.Atoi(status)
	if erroConversao != nil {
		erroRest := rest_err.NovoErroDeRequisicaoInvalida(
			"Erro ao validar o parâmetro de status do leilão")
		c.JSON(erroRest.Codigo, erroRest)
		return
	}

	leiloes, erro := controlador.casoDeUsoDeLeilao.BuscarLeiloes(
		context.Background(),
		auction_usecase.StatusDoLeilao(statusNumerico), categoria, nomeDoProduto)
	if erro != nil {
		erroRest := rest_err.ConverterErro(erro)
		c.JSON(erroRest.Codigo, erroRest)
		return
	}

	c.JSON(http.StatusOK, leiloes)
}

func (controlador *ControladorDeLeilao) BuscarLanceVencedorPorLeilao(c *gin.Context) {
	idDoLeilao := c.Param("auctionId")

	if erro := uuid.Validate(idDoLeilao); erro != nil {
		erroRest := rest_err.NovoErroDeRequisicaoInvalida("Campos inválidos", rest_err.Causa{
			Campo:    "auctionId",
			Mensagem: "Valor de UUID inválido",
		})

		c.JSON(erroRest.Codigo, erroRest)
		return
	}

	informacoesDoVencedor, erro := controlador.casoDeUsoDeLeilao.BuscarLanceVencedorPorLeilao(
		context.Background(), idDoLeilao)
	if erro != nil {
		erroRest := rest_err.ConverterErro(erro)
		c.JSON(erroRest.Codigo, erroRest)
		return
	}

	c.JSON(http.StatusOK, informacoesDoVencedor)
}
