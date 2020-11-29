package main

import (
	"fmt"
	"io"
	"os"

	"github.com/pacedotdev/firesearch-sdk/clients/go/firesearch"
	"github.com/spf13/cobra"
	"github.com/willdot/bookmarker/commands"
	"github.com/willdot/bookmarker/pkg"
)

func main() {
	env := getEnv()

	firesearch, err := createFiresearch(env.indexPath, env.endpoint, env.secret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create firesearch service: %s\n", err.Error())
		os.Exit(1)
	}

	cobraApp := newApp("bookmarker", "A CLI tool for storing bookmarks with tags that can then be searched later", env.indexPath, os.Stdout, firesearch)
	cobraApp.addCommands(
		commands.NewAddCommand(cobraApp.store),
		commands.NewFindCommand(cobraApp.store),
	)

	if err := cobraApp.rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// CobraApp contains all the dependancies to run the bookmarker app
type CobraApp struct {
	rootCmd *cobra.Command
	stdout  io.Writer
	store   *pkg.Store
}

func newApp(name, desc, indexPath string, stdout io.Writer, firesearch *firesearch.AutocompleteService) *CobraApp {
	store := pkg.NewStore(firesearch, indexPath, stdout)

	app := &CobraApp{
		rootCmd: &cobra.Command{
			Use:   name,
			Short: desc,
		},
		stdout: stdout,
		store:  store,
	}
	return app
}

func (a *CobraApp) addCommands(cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		a.rootCmd.AddCommand(cmd)
	}
}

type environment struct {
	indexPath string
	endpoint  string
	secret    string
}

func getEnv() environment {
	indexPath := os.Getenv("INDEXPATH")
	if indexPath == "" {
		indexPath = "bookmarks"
	}

	endpoint := os.Getenv("ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8888/api"
	}

	secret := os.Getenv("SECRET")

	return environment{
		indexPath: indexPath,
		endpoint:  endpoint,
		secret:    secret,
	}
}
