package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listidentitiesCmd = &cobra.Command{
	Use:   "listidentities",
	Short: "List your identities",
	Run: func(cmd *cobra.Command, args []string) {
		identities := *GetAppContext().Identities
		fmt.Printf("You added %d identities:", len(identities))
		for _, identity := range identities {
			fmt.Printf("\n- %s (%s)", identity.Name, identity.Domain)
		}
	},
}

func init() {
	rootCmd.AddCommand(listidentitiesCmd)
}
