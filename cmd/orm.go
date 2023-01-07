/*
Copyright © 2023 Killua<captainchengjie@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	host     string
	addr     string
	username string
	password string
	database string
)

// ormCmd represents the orm command
var ormCmd = &cobra.Command{
	Use:   "orm",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("orm called")
	},
}

func init() {
	rootCmd.AddCommand(ormCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ormCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ormCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	ormCmd.Flags().StringVarP(&username, "username", "u", "", "database username")
	ormCmd.Flags().StringVarP(&password, "password", "p", "", "database password")
	ormCmd.Flags().StringVarP(&database, "database", "d", "", "database name")
	ormCmd.Flags().StringVarP(&addr, "addr", "a", "", "database port")
	ormCmd.Flags().StringVarP(&host, "host", "h", "", "database host")
}
