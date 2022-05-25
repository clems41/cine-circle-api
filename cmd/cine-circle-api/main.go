package main

import (
	"cine-circle-api/external/mailService/mailServiceMock"
	"cine-circle-api/external/mediaProvider/theMovieDatabase"
	"cine-circle-api/internal/constant/swaggerConst"
	"cine-circle-api/internal/domain/adminDom/healthCheckDom"
	"cine-circle-api/internal/domain/circleDom"
	"cine-circle-api/internal/domain/mediaDom"
	"cine-circle-api/internal/domain/recommendationDom"
	"cine-circle-api/internal/domain/userDom"
	"cine-circle-api/internal/repository/postgresRepositories"
	"cine-circle-api/pkg/httpServer"
	"cine-circle-api/pkg/logger"
	"cine-circle-api/pkg/sql/sqlConnection"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Handle signal interruption
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT, os.Interrupt, os.Kill)

	// Try to connect to PostgresSQL database using default config (from env variables)
	DB, err := sqlConnection.Open(nil)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	// Since we don't use ISI CI, migration-manager will not be run, so we need to migrate repositories at the start of the application
	tx := DB.Begin()
	err = postgresRepositories.Migrate(tx)
	if err != nil {
		tx.Rollback()
		logger.Fatalf("Cannot migrate repositories : %s", err.Error())
	}
	tx.Commit()

	// Create restful container that will be used to start HTTP server
	restfulContainer := httpServer.NewRestfulContainer()
	restfulContainer.SetLogger(logger.Logger())

	// Create useful services
	serviceMail := mailServiceMock.New()
	mediaProvider := theMovieDatabase.New()
	userRepo := postgresRepositories.NewUser(DB)
	mediaRepo := postgresRepositories.NewMedia(DB)
	circleRepo := postgresRepositories.NewCircle(DB)
	recommendationRepo := postgresRepositories.NewRecommendation(DB)

	// Add all new handlers in restfulContainer : here you can define all project endpoints
	restfulContainer.AddHandlers(
		healthCheckDom.NewHandler(),
		userDom.NewHandler(userDom.NewService(serviceMail, userRepo)),
		mediaDom.NewHandler(mediaDom.NewService(mediaProvider, mediaRepo)),
		circleDom.NewHandler(circleDom.NewService(circleRepo, userRepo)),
		recommendationDom.NewHandler(recommendationDom.NewService(recommendationRepo, userRepo, mediaRepo, circleRepo)),
	)

	// Add endpoint for getting swagger.json file (based on documentation from each webService)
	// You can visualize it here : http://swagger-ui.default.svc.kube.isi/?url=http://localhost:8080/swagger.json
	// Default swagger endpoint : /swagger.json
	// You should maintain tags from swaggerConst.Info updated during project development
	restfulContainer.GenerateSwagger(swaggerConst.Info)

	// Start HTTP server based on handlers defined previously.
	// You can specify some parameters using environment variables :
	// - SWAGGER_URL
	// - HTTP_TRACING_LOG
	// - HTTP_REQUEST_LOG
	// - HTTP_BIND_ADDRESS
	// - HTTP_READ_TIMEOUT_SEC
	// - HTTP_READ_HEADER_TIMEOUT_SEC
	// - HTTP_WRITE_TIMEOUT_SEC
	// - HTTP_IDLE_TIMEOUT_SEC
	// - HTTP_MAX_HEADER_BYTES
	// - TOKEN_EXPIRATION_HOURS
	server, err := restfulContainer.HttpServer()
	if err != nil {
		logger.Fatal(err)
	}
	go func() {
		logger.Infof("Starting HTTP server using configuration : %+v", server)
		err = server.ListenAndServe()
		if err != nil {
			logger.Errorf("Error while trying to serving http server... err: %s", err.Error())
		}
	}()

	// This code below will be executed only when signal is received (shutdown, cancel, kill, etc...).
	// Add all closing methods here.
	<-stopSignal

	logger.Infof("Closing HTTP server...")
	err = server.Close()
	if err != nil {
		logger.Fatalf("Error while trying to close HTTP server... err: %s", err.Error())
	} else {
		logger.Infof("HTTP server has been gracefully closed")
	}

	logger.Infof("Closing database connection...")
	err = sqlConnection.Close(DB)
	if err != nil {
		logger.Fatalf("Error while trying to close database connection... err: %s", err.Error())
	} else {
		logger.Infof("Database connection has been gracefully closed")
	}

	close(stopSignal)
}
