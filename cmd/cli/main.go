package main

import (
	"fmt"
	"os"
	"strings"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configuration config.Config
	environment   string
)

func main() {
	cobra.OnInitialize(initConfig)

	generatorCmdWrapper := newGeneratorCmdWrapper()
	rootCmd := newRootCommand(generatorCmdWrapper.run)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type cobraCmdRunFunc func(cmd *cobra.Command, args []string)

func newRootCommand(runFunc cobraCmdRunFunc) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "fcs",
		Short: "Functional Conformance Suite CLI",
		Long:  `To use with pipelines and reproducible test runs`,
	}

	rootCmd.PersistentFlags().StringVarP(&environment, "environment", "e", "development", "Specify the environment you are running")

	rootCmd.AddCommand(versionCmd())
	rootCmd.AddCommand(generatorCmd(runFunc))

	return rootCmd
}

func initConfig() {
	viper.SetEnvPrefix("FCS")
	viper.SetConfigName(environment)                       // name of config file (without extension), e.g., "development.json"
	viper.SetConfigType("json")                            // or viper.SetConfigType("yaml")
	viper.AddConfigPath("../config/")                      /// path to look for the config file in
	viper.AddConfigPath("./config/")                       // optionally look for config here.
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // so that FCS_SERVER_PORT=9090 becomes FCS_SERVER.PORT=9090
	viper.AutomaticEnv()                                   // read in environment variables that match, e.g., FCS_PORT=80 ./fcs

	setConfigDefaults()

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&configuration); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	printConfig()
}

func printConfig() {
	fmt.Println("Configuration")
	fmt.Println("  Server")
	fmt.Println("    Port:", configuration.Server.Port)

	fmt.Println()
}

func setConfigDefaults() {
	viper.SetDefault("SERVER.PORT", "8080")
}
