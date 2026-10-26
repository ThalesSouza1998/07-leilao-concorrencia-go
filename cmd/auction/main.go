package main

import (
	"context"
	"fullcycle-auction_go/configuration/database/mongodb"
	"fullcycle-auction_go/internal/infra/api/web/controller/auction_controller"
	"fullcycle-auction_go/internal/infra/api/web/controller/bid_controller"
	"fullcycle-auction_go/internal/infra/api/web/controller/user_controller"
	"fullcycle-auction_go/internal/infra/database/auction"
	"fullcycle-auction_go/internal/infra/database/bid"
	"fullcycle-auction_go/internal/infra/database/user"
	"fullcycle-auction_go/internal/usecase/auction_usecase"
	"fullcycle-auction_go/internal/usecase/bid_usecase"
	"fullcycle-auction_go/internal/usecase/user_usecase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func main() {
	contexto := context.Background()

	if erro := godotenv.Load("cmd/auction/.env"); erro != nil {
		log.Println("Arquivo .env não encontrado, usando as variáveis de ambiente do sistema")
	}

	banco, erro := mongodb.NovaConexaoMongoDB(contexto)
	if erro != nil {
		log.Fatal(erro.Error())
		return
	}

	roteador := gin.Default()

	controladorDeUsuario, controladorDeLance, controladorDeLeilao := iniciarDependencias(banco)

	roteador.GET("/auction", controladorDeLeilao.BuscarLeiloes)
	roteador.GET("/auction/:auctionId", controladorDeLeilao.BuscarLeilaoPorId)
	roteador.POST("/auction", controladorDeLeilao.CriarLeilao)
	roteador.GET("/auction/winner/:auctionId", controladorDeLeilao.BuscarLanceVencedorPorLeilao)
	roteador.POST("/bid", controladorDeLance.CriarLance)
	roteador.GET("/bid/:auctionId", controladorDeLance.BuscarLancesPorLeilao)
	roteador.GET("/user/:userId", controladorDeUsuario.BuscarUsuarioPorId)

	if erro := roteador.Run(":8080"); erro != nil {
		log.Fatal(erro)
	}
}

func iniciarDependencias(banco *mongo.Database) (
	controladorDeUsuario *user_controller.ControladorDeUsuario,
	controladorDeLance *bid_controller.ControladorDeLance,
	controladorDeLeilao *auction_controller.ControladorDeLeilao) {

	repositorioDeLeilao := auction.NovoRepositorioDeLeilao(banco)
	repositorioDeLance := bid.NovoRepositorioDeLance(banco, repositorioDeLeilao)
	repositorioDeUsuario := user.NovoRepositorioDeUsuario(banco)

	controladorDeUsuario = user_controller.NovoControladorDeUsuario(
		user_usecase.NovoCasoDeUsoDeUsuario(repositorioDeUsuario))
	controladorDeLeilao = auction_controller.NovoControladorDeLeilao(
		auction_usecase.NovoCasoDeUsoDeLeilao(repositorioDeLeilao, repositorioDeLance))
	controladorDeLance = bid_controller.NovoControladorDeLance(
		bid_usecase.NovoCasoDeUsoDeLance(repositorioDeLance))

	return
}
