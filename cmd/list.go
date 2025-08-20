package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"starless.dev/smokescreen/cloudflare"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list <identity>",
	Short: "List the emails you created",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		identity, err := GetAppContext().Identities.Get(name)
		if err != nil {
			fmt.Println(err)
			return
		}

		list, err := cloudflare.ListEmails(identity)
		if err != nil {
			fmt.Printf("An error occurred while listing emails: %v", err)
		} else if list.Success {
			emails := list.Result
			fmt.Printf("Found %d emails:", len(emails))
			for _, email := range emails {
				emailName := email.Name
				if !strings.HasPrefix(emailName, cloudflare.EmailNamePrefix) {
					continue
				}
				matchers := email.Matchers
				if len(matchers) == 0 {
					continue
				}

				realName, _ := strings.CutPrefix(emailName, cloudflare.EmailNamePrefix)
				fmt.Printf("\n- %s: %s (enabled: %v)", realName, matchers[0].Value, email.Enabled)
			}
		} else {
			fmt.Printf("Cloudflare returned an error: %v", list.Errors)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
