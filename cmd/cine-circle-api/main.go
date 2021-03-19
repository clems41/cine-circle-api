package main

import (
	externalApi "cine-circle/external/api"
	"cine-circle/internal/db"
	"cine-circle/internal/logger"
	"flag"
	"github.com/spf13/cobra"
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

	if err := apiCmd.Execute(); err != nil {
		logger.Sugar.Fatal(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	database, err := db.OpenConnection()
	if err.IsNotNil() {
		panic(err)
	}
	defer database.Close()
	_, movie := externalApi.FindMovieBySearch("avengers endgame")
	logger.Sugar.Info("movie : %+v", movie)
	_, movie = externalApi.FindMovieByID("tt0848228")
	logger.Sugar.Info("movie : %+v", movie)
}
