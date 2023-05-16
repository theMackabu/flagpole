package cmd

import (
	"log"

	"flagpole_auth/api"
	"flagpole_auth/api/database"
	
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var startServer = &cobra.Command{
	Use:   "start",
	Short: "",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		envError := godotenv.Load()
		address, _ := cmd.Flags().GetString("address")
		cors, _ := cmd.Flags().GetString("cors")
		
		if envError != nil{
			log.Fatal("Failed loading .env file")
		}
		
		db := database.Migrations{DB: database.SQLite{}}
		db.MakeMigrations()
		
		server := api.HttpServer{}
		server.CreateRouter()
		server.StartServer(address, cors)

	},
}

func init() {
	rootCmd.AddCommand(startServer)

	startServer.Flags().StringP("address", "", "", "")
	startServer.Flags().StringP("cors", "", "", "")
   
	startServer.MarkFlagRequired("address")
   startServer.MarkFlagRequired("cors")
}
