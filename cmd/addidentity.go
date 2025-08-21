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
		_, err := fmt.Scanln(&token)
		if err != nil {
			fmt.Printf("Could not get user input: %v", err)
			return
		}

		fmt.Print("Insert your zone ID: ")
		_, err = fmt.Scanln(&zoneId)
		if err != nil {
			fmt.Printf("Could not get user input: %v", err)
			return
		}

		fmt.Print("Insert your domain: ")
		_, err = fmt.Scanln(&domain)
		if err != nil {
			fmt.Printf("Could not get user input: %v", err)
			return
		}

		fmt.Print("Where do you want the emails to be redirected: ")
		_, err = fmt.Scanln(&email)
		if err != nil {
			fmt.Printf("Could not get user input: %v", err)
			return
		}

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
