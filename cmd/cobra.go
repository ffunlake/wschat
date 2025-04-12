package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)
var (
	configFile string
	rootCmd = &cobra.Command{Use: "wschat"}
)

func init() {
	
	rootCmd.AddCommand(clientCommand)
	rootCmd.AddCommand(serverCommand)
}

//Execute : Cobra entrance
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}