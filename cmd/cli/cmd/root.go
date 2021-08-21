package cmd

import (
	"fmt"
	"os"

	"github.com/justcompile/tnl/pkg/socketclient"
	"github.com/justcompile/tnl/pkg/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use: "run",
		Run: func(cmd *cobra.Command, args []string) {
			opts := &socketclient.Options{}

			opts.WebsocketServerBindAddress, _ = cmd.Flags().GetString("websocket-address")
			opts.LocalBindAddress, _ = cmd.Flags().GetString("port")

			opts.Protocol = "https"

			if useSSL, _ := cmd.Flags().GetBool("use-ssl"); !useSSL {
				opts.Protocol = "http"
			}

			client := socketclient.New(opts)

			defer client.Close()

			window, err := ui.ConstructUI("0:" + opts.LocalBindAddress)
			if err != nil {
				log.Fatal(err)
			}

			log.Infoln("connecting...")

			client.Connect(
				window.GetComponent(ui.ComponentIDInfo).GetUpdateChannel(),
				window.GetComponent(ui.ComponentIDRequests).GetUpdateChannel(),
			)

			log.Infoln("running...")

			if runErr := window.Run(); runErr != nil {
				log.Fatal(runErr)
			}
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringP("port", "p", "3333", "local port to forward traffic onto")
	rootCmd.Flags().StringP("websocket-address", "w", "tnl.justcompile.io:8081", "Address of Websocket Server")
	rootCmd.Flags().BoolP("use-ssl", "s", true, "remote domain operates over ssl")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
