package cmd

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	echoSwagger "github.com/swaggo/echo-swagger"
	"insider/configs/appConfigs"
	errorConfigs "insider/configs/errorConfigs"
	_ "insider/docs/jobScheduler"
	"insider/internal/apps/jobScheduler/handlers"
	"insider/internal/apps/jobScheduler/jobs"
	"insider/internal/repositories"
	"insider/pkg"
	"time"
)

// @title Job Scheduler
// @version 1.0
// @description This is the API documentation for the Job Scheduler service.
// @host localhost:3001
// @BasePath /job-scheduler
var jobSchedulerCmd = &cobra.Command{
	Use: "jobScheduler",
}

func init() {
	rootCmd.AddCommand(jobSchedulerCmd)
	jobSchedulerCmd.RunE = func(cmd *cobra.Command, args []string) error {

		cfg := appConfigs.GetConfigs()
		errorCodes := errorConfigs.GetErrorCodes()

		mc, err := pkg.NewMongoClient(time.Second*cfg.Mongo.TimeOutDurationInSeconds, cfg.Mongo.ConnectionURI)
		if err != nil {
			panic(err)
		}

		logger := pkg.InitLogrusConfig()

		e := echo.New()

		e.Logger = &logger

		pkg.RegisterMiddlewares(e, cfg.JobScheduler.RoutePrefix)

		jobSchedulerRepo, err := repositories.NewJobSchedulerRepository(&cfg, mc)
		if err != nil {
			panic(err)
		}

		messengerJob := jobs.NewMessengerJob(jobSchedulerRepo, errorCodes, &cfg)
		//xyzJob := jobs.NewXyzJob...

		handler := handlers.NewHandler(e, cfg.JobScheduler.RoutePrefix, messengerJob)

		handler.RegisterRoutesForMessengerJob()
		//handler.RegisterRoutesForXyzJob

		messengerJob.Start()

		e.GET("/swagger/*", echoSwagger.WrapHandler)
		e.GET(cfg.JobScheduler.RoutePrefix+"/swagger/*", echoSwagger.EchoWrapHandler(echoSwagger.InstanceName("jobScheduler")))

		log.Fatal(e.Start(":3001"))

		return nil
	}
}
