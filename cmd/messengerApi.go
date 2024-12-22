package cmd

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	echoSwagger "github.com/swaggo/echo-swagger"
	"insider/configs/appConfigs"
	_ "insider/docs/messengerApi"
	"insider/internal/apps/messengerApi"
	"insider/internal/repositories"
	"insider/pkg"
	"time"
)

var messengerCmd = &cobra.Command{
	Use: "messengerApi",
}

// @title Messenger API
// @version 1.0
// @description This is the documentation for the Messenger API.
// @host localhost:3000
// @BasePath /api/messenger
func init() {
	rootCmd.AddCommand(messengerCmd)
	messengerCmd.RunE = func(cmd *cobra.Command, args []string) error {

		cfg := appConfigs.GetConfigs()
		//errCodes := errorConfigs.GetErrorCodes()

		mc, err := pkg.NewMongoClient(time.Second*cfg.Mongo.TimeOutDurationInSeconds, cfg.Mongo.ConnectionURI)
		if err != nil {
			panic(err)
		}

		rc, err := pkg.NewRedisClient(cfg.RedisConfigs.ConnectionURI)
		if err != nil {
			log.Error(err.Error())
		}

		//comment-out if you don't want to use it
		pkg.GenerateMockData(mc, cfg.MessengerApi.MongoDbName, cfg.MessengerApi.MessagesColName, 10)

		logger := pkg.InitLogrusConfig()

		e := echo.New()

		e.Logger = &logger

		pkg.RegisterMiddlewares(e, cfg.MessengerApi.RoutePrefix)

		messengerRepo, err := repositories.NewMessengerRepository(&cfg, mc, rc)
		if err != nil {
			panic(err)
		}

		messengerService := messengerApi.NewMessengerService(messengerRepo)

		messengerHandler := messengerApi.NewHandler(e, messengerService, &cfg)

		messengerHandler.RegisterRoutes(cfg.MessengerApi.ApiPrefix + cfg.MessengerApi.RoutePrefix)

		e.GET(cfg.MessengerApi.RoutePrefix+"/swagger/*", echoSwagger.EchoWrapHandler(echoSwagger.InstanceName("messengerApi")))

		log.Fatal(e.Start(":3000"))

		return nil
	}
}
