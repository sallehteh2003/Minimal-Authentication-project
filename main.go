package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
	"main/Authentication"
	"main/Config"
	"main/Database"
	"main/Handlers"
	"main/Validation"
)

func main() {
	var cfg Config.Config
	logger := logrus.New()
	r := gin.Default()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetReportCaller(true)

	err := cleanenv.ReadConfig("Config/config.json", &cfg)
	if err != nil {
		logger.WithError(err).Panicln("failed to load the configs")
	} else {
		logger.Infof("successful to read the configs: %+v", cfg)
	}

	// Create a new instance of Database sql postgres
	gormDB, err := Database.CreateAndConnectToDb(cfg)
	if err != nil {
		logger.WithError(err).Panicln("can not connect to db")
	}

	//create model of database
	if err := gormDB.CreateModel(); err != nil {
		logger.WithError(err).Fatalln("can not create table in db ")
	}
	logger.Infof("%+v", cfg)

	// Create a new instance of Validation
	valid := Validation.CreateValidation([]string{"gmail.com"})

	//// Create a new instance of Authentication
	auth, err := Authentication.CreateAuthentication(gormDB, 10, logger)
	if err != nil {
		logger.WithError(err).Fatal("can not Create instance od Authentication ")
	}

	//// Create a new instance of server
	server := Handlers.Server{
		Database:       gormDB,
		Logger:         logger,
		Authentication: auth,
		Validation:     valid,
	}

	//// api register
	r.POST("/api/v1/auth/signup", server.HandleSignup)
	r.POST("/api/v1/auth/login", server.HandleLogin)
	r.GET("/api/v1/auth/checkLogin", server.HandleCheckLogin)

	//// RUN SERVER
	if err := r.Run(fmt.Sprintf("localhost:%v", cfg.Server.Port)); err != nil {
		logrus.WithError(err).Fatalln("can not run server ")
	}

}
