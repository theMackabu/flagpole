package cmd

import (
   "fmt"
   "os"
   "os/signal"

   "flagpole_auth/config"
   "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
   Use:     "flagpole_auth",
   Short:   "",
   Version: config.Version,
}

func Execute() {
   c := make(chan os.Signal, 1)
   signal.Notify(c, os.Interrupt)
   go func() {
      <-c
      fmt.Println("user interrupt")
      os.Exit(0)
   }()

   if err := rootCmd.Execute(); err != nil {
      os.Exit(1)
   }
}