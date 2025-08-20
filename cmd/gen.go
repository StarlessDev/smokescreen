package cmd

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/spf13/cobra"
	"starless.dev/smokescreen/cloudflare"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen <identity> <name>",
	Short: "Generate an email address using a certain identity.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		email := args[1]

		identity, err := GetAppContext().Identities.Get(name)
		if err != nil {
			fmt.Println(err)
			return
		}

		fullEmail, err := generateEmail(identity, email)
		if err != nil {
			fmt.Println(err)
			return
		}

		res, err := cloudflare.GenerateEmail(identity, cloudflare.EmailNamePrefix+email, fullEmail)
		if err != nil {
			fmt.Printf("An error occurred while generating the email address: %v", err)
		} else if res.Success {
			fmt.Printf("Generated email: %s", fullEmail)
		} else {
			fmt.Printf("Cloudflare returned an error: %v", res.Errors)
		}
	},
}

func generateEmail(id *cloudflare.Identity, email string) (string, error) {
	i, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err == nil {
		return fmt.Sprintf("%s-%s@%s", email, i.String(), id.Domain), nil
	} else {
		return "", err
	}
}

func init() {
	rootCmd.AddCommand(genCmd)
}
