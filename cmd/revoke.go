package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"starless.dev/smokescreen/cloudflare"
)

// revokeCmd represents the revoke command
var revokeCmd = &cobra.Command{
	Use:   "revoke <identity> <email>",
	Short: "Revoke a generated email. This will completely delete the address!",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		email := cloudflare.EmailNamePrefix + args[1]
		identity, err := GetAppContext().Identities.Get(name)
		if err != nil {
			fmt.Println(err)
			return
		}

		list, err := cloudflare.ListEmails(identity)
		if err != nil {
			fmt.Printf("An error occurred while revoking the email: %v", err)
			return
		} else if !list.Success {
			fmt.Printf("Cloudflare returned an error: %v", list.Errors)
			return
		}

		emailId, err := getAddressId(list, email)
		if err != nil {
			fmt.Println(err)
			return
		}

		res, err := cloudflare.RevokeEmail(identity, emailId)
		if err != nil {
			fmt.Printf("An error occurred while revoking the email: %v", err)
		} else if res.Success {
			fmt.Printf("The email was revoked successfully")
		} else {
			fmt.Printf("Cloudflare returned an error: %v", res.Errors)
		}
	},
}

func getAddressId(list *cloudflare.ListEmailResponse, name string) (string, error) {
	for _, email := range list.Result {
		if strings.EqualFold(email.Name, name) {
			return email.Id, nil
		}
	}
	return "", errors.New("email address not found")
}

func init() {
	rootCmd.AddCommand(revokeCmd)
}
