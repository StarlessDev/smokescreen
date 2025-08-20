package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"starless.dev/smokescreen/cloudflare"
)

var addidentityCmd = &cobra.Command{
	Use:   "addidentity <identity>",
	Short: "Add your Cloudflare token and zone id to start managing your emails",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		var token, zoneId, domain, email string

		fmt.Print("Insert your API Token: ")
		fmt.Scanln(&token)

		fmt.Print("Insert your zone ID: ")
		fmt.Scanln(&zoneId)

		fmt.Print("Insert your domain: ")
		fmt.Scanln(&domain)

		fmt.Print("Where do you want the emails to be redirected: ")
		fmt.Scanln(&email)

		identities := GetAppContext().Identities
		identities.Add(&cloudflare.Identity{
			Name:   name,
			Token:  token,
			ZoneId: zoneId,
			Domain: domain,
			Email:  email,
		})
		identities.Save()
		fmt.Printf("Added %s as an identity", name)
	},
}

func init() {
	rootCmd.AddCommand(addidentityCmd)
}
