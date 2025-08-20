package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"starless.dev/smokescreen/cloudflare"
)

// AppContext holds shared application state and resources
type AppContext struct {
	Identities *cloudflare.Identities
}

// NewAppContext creates a new application context with loaded identities
func NewAppContext() (*AppContext, error) {
	identities, err := cloudflare.ReadIdentities()
	if err != nil {
		return nil, err
	}
	return &AppContext{
		Identities: identities,
	}, nil
}

// Global app context instance
var appCtx *AppContext

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "smokescreen",
	Short: "CLI utility to manage email aliases using Cloudflare's email routing feature.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Initialize application context
	var err error
	appCtx, err = NewAppContext()
	if err != nil {
		panic("Could not read your identities")
	}

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// GetAppContext returns the global application context
func GetAppContext() *AppContext {
	return appCtx
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.smokescreen.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
