package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// removeidentityCmd represents the removeidentity command
var removeidentityCmd = &cobra.Command{
	Use:   "removeidentity <identity>",
	Short: "Remove a previously added identity",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		identities := GetAppContext().Identities
		err := identities.Remove(name)
		if err == nil {
			identities.Save()
			fmt.Printf("Removed the %s identity", name)
		} else {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeidentityCmd)
}
