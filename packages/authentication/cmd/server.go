package cmd

import (
	"flagpole_auth/api"
	"flagpole_auth/api/database"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
	"os"
)

var startServer = &cobra.Command{
	Use:   "start",
	Short: "",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cors, _ := cmd.Flags().GetString("cors")
		secret, _ := cmd.Flags().GetString("secret")
		address, _ := cmd.Flags().GetString("address")
		db_name, _ := cmd.Flags().GetString("db_name")
		log := logger.Create("authentication", &logger.StdoutStream{Unbuffered: true})

		os.Setenv("DB_NAME", db_name)
		os.Setenv("SECRET_KEY", secret)

		db := database.Migrations{DB: database.SQLite{}}
		db.MakeMigrations(log)

		server := api.HttpServer{}
		server.CreateRouter()
		server.StartServer(address, cors, log)

	},
}

func init() {
	rootCmd.AddCommand(startServer)

	startServer.Flags().StringP("address", "", "", "")
	startServer.Flags().StringP("db_name", "", "", "")
	startServer.Flags().StringP("secret", "", "", "")
	startServer.Flags().StringP("cors", "", "", "")

	startServer.MarkFlagRequired("address")
	startServer.MarkFlagRequired("db_name")
	startServer.MarkFlagRequired("secret")
	startServer.MarkFlagRequired("cors")
}
