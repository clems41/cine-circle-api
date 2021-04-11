package main

import (
	"cine-circle/internal/domain/authenticationDom"
	"cine-circle/internal/domain/circleDom"
	"cine-circle/internal/domain/movieDom"
	"cine-circle/internal/domain/recommendationDom"
	"cine-circle/internal/domain/userDom"
	"cine-circle/internal/domain/watchlistDom"
	"cine-circle/internal/handler"
	"cine-circle/internal/logger"
	"cine-circle/internal/repository"
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

	flag := apiCmd.Flags()
	flag.StringVar(&bind, "bind", ":8080", "HTTP bind specification")

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
	database, err := repository.OpenConnection()
	if err.IsNotNil() {
		logger.Sugar.Fatalf(err.Error())
	}
	defer database.Close()

	repositories := repository.NewAllRepositories(database.DB())
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

	handler.AddWebService(restful.DefaultContainer,
		handler.NewAuthenticationHandler(authenticationDom.NewService(repositories.Authentication)))
	handler.AddWebService(restful.DefaultContainer,
		handler.NewCircleHandler(circleDom.NewService(repositories.Circle)))
	handler.AddWebService(restful.DefaultContainer,
		handler.NewRootHandler())
	handler.AddWebService(restful.DefaultContainer,
		handler.NewMovieHandler(movieDom.NewService(repositories.Movie)))
	handler.AddWebService(restful.DefaultContainer,
		handler.NewRecommendationHandler(recommendationDom.NewService(repositories.Recommendation)))
	handler.AddWebService(restful.DefaultContainer,
		handler.NewUserHandler(userDom.NewService(repositories.User)))
	handler.AddWebService(restful.DefaultContainer,
		handler.NewWatchlistHandler(watchlistDom.NewService(repositories.Watchlist)))

	config := restfulspec.Config{
		WebServices: restful.DefaultContainer.RegisteredWebServices(),
		APIPath:     "/handler/swagger.yaml",
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
		err2 := srv.Close()
		if err2 != nil {
			logger.Sugar.Fatalf("Error while trying to close http server... err: %v", err2)
		}
	}()

	err3 := srv.ListenAndServe()
	if err3 != nil {
		logger.Sugar.Fatalf("Error while trying to serving http server... err: %v", err3)
	}
}
