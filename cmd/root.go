package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/pacedotdev/firesearch-sdk/clients/go/firesearch"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bookmarker",
	Short: "A CLI tool for storing bookmarks with keywords that can then be searched for and launched",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var endpoint string
var indexPath string
var secret string

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	addCmd.Flags().StringVar(&endpoint, "endpoint", "http://localhost:8888/api", "the Firesearch endpoint address")
	addCmd.Flags().StringVar(&indexPath, "index", "bookmarks", "the index path to use for Firesearch index")
	addCmd.Flags().StringVar(&secret, "secret", "", "the secret API token for Firesearch")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bookmarker.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".bookmarker" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".bookmarker")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func createFiresearch(indexPath, endpoint, secret string) (*firesearch.AutocompleteService, error) {
	client := firesearch.NewClient(endpoint, secret)
	autocompleteService := firesearch.NewAutocompleteService(client)
	err := createIndex(context.TODO(), autocompleteService, indexPath)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating index: %s", indexPath)
	}

	return autocompleteService, nil
}

func createIndex(ctx context.Context, service *firesearch.AutocompleteService, index string) error {
	createAutocompleteIndexReq := firesearch.CreateAutocompleteIndexRequest{
		Index: firesearch.AutocompleteIndex{
			IndexPath:     index,
			Name:          "Bookmark autocomplete index",
			CaseSensitive: false,
		},
	}

	_, err := service.CreateIndex(ctx, createAutocompleteIndexReq)
	if err != nil {
		return errors.Wrap(err, "firesearch: IndexService.CreateIndex")
	}

	return nil
}
