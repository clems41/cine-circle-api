package main

import (
	"cine-circle/internal/domain/userDom"
	"cine-circle/internal/repository"
	"cine-circle/internal/repository/postgres"
	"cine-circle/internal/webService"
	logger "cine-circle/pkg/logger"
	"context"
	"flag"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/signal"
)

var (
	Signals chan os.Signal
	cancel context.CancelFunc
	CtxGlobal context.Context
	bind string
)

func main() {
	logger.InitLogger()

	apiCmd := &cobra.Command{
		Use:  "cine-circle",
		Long: `Cine-circle API`,
		Run:  run,
	}

	// include standard flags
	apiCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	flags := apiCmd.Flags()
	flags.StringVar(&bind, "bind", ":8080", "HTTP bind specification")

	if err := apiCmd.Execute(); err != nil {
		logger.Sugar.Fatal(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	// Create context and handle signal interruption
	CtxGlobal, cancel = context.WithCancel(context.Background())
	Signals = make(chan os.Signal, 1)
	signal.Notify(Signals, os.Interrupt)

	// Try to connect to PostgresSQL repository
	DB, err := postgres.OpenConnection()
	if err != nil {
		logger.Sugar.Fatalf(err.Error())
	}
	defer func() {
		err = postgres.CloseConnection(DB)
		if err != nil {
			logger.Sugar.Fatalf("Error while trying to close database connection... err: %s", err.Error())
		}
	}()

	repositories := repository.NewAllRepositories(DB)
	repositories.Migrate()

	// Add container filter to enable CORS
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST"},
		CookiesAllowed: false,
		Container:      restful.DefaultContainer}
	restful.DefaultContainer.Filter(cors.Filter)

	// Add container filter to respond to OPTIONS
	restful.DefaultContainer.Filter(restful.DefaultContainer.OPTIONSFilter)
	// accept and respond in JSON unless told otherwise
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)
	// gzip if accepted
	restful.DefaultContainer.EnableContentEncoding(true)

	// Adding all new handlers here
/*	webService.AddWebService(restful_utils.DefaultContainer, webService.NewRootHandler())

	webService.AddWebService(restful_utils.DefaultContainer,
		webService.NewCircleHandler(circleDom.NewService(repositories.Circle)))

	webService.AddWebService(restful_utils.DefaultContainer,
		webService.NewMovieHandler(movieDom.NewService(repositories.Movie)))

	webService.AddWebService(restful_utils.DefaultContainer,
		webService.NewRecommendationHandler(recommendationDom.NewService(repositories.Recommendation)))

	webService.AddWebService(restful_utils.DefaultContainer, userDom.NewUserHandler(userService))

	webService.AddWebService(restful_utils.DefaultContainer,
		webService.NewWatchlistHandler(watchlistDom.NewService(repositories.Watchlist)))*/

	webService.AddHandler(restful.DefaultContainer,
		userDom.NewHandler(userDom.NewService(userDom.NewRepository(DB))))

	config := restfulspec.Config{
		WebServices: restful.DefaultContainer.RegisteredWebServices(),
		APIPath:     "/webService/swagger.yaml",
	}
	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))

	srv := http.Server{
		Addr:    bind,
		Handler: restful.DefaultContainer,
	}

	logger.Sugar.Infof("Routes are now accessibles...")

	go func() {
		<-Signals
		cancel()
		logger.Sugar.Infof("Closing http server...")
		err = srv.Close()
		if err != nil {
			logger.Sugar.Fatalf("Error while trying to close http server... err: %s", err.Error())
		}
		err = postgres.CloseConnection(DB)
		if err != nil {
			logger.Sugar.Fatalf("Error while trying to close database connection... err: %s", err.Error())
		}
	}()

	err = srv.ListenAndServe()
	if err != nil {
		logger.Sugar.Fatalf("Error while trying to serving http server... err: %s", err.Error())
	}
}
