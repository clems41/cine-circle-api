package main

import (
	"cine-circle/external/omdb"
	"cine-circle/internal/api"
	"cine-circle/internal/db"
	"cine-circle/internal/logger"
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

	// Try to connect to PostgresSQL database
	database, err := db.OpenConnection()
	if err.IsNotNil() {
		logger.Sugar.Fatalf(err.Error())
	}
	err = database.Close()
	if err.IsNotNil() {
		logger.Sugar.Fatalf(err.Error())
	}

	// Test to get some movies from open data movie
	_, movie := omdb.FindMovieByTitle("avengers endgame")
	logger.Sugar.Info("movie : %+v", movie)
	_, movie = omdb.FindMovieByID("tt0848228")
	logger.Sugar.Info("movie : %+v", movie)

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

	for _, ws := range api.DefineRoutes() {
		restful.DefaultContainer.Add(ws)
	}

	config := restfulspec.Config{
		WebServices: restful.DefaultContainer.RegisteredWebServices(),
		APIPath:     "/api/swagger.yaml",
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
