package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/pacedotdev/firesearch-sdk/clients/go/firesearch"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/willdot/bookmarker/commands"
	"github.com/willdot/bookmarker/pkg"
)

// CobraApp contains all the dependancies to run the bookmarker app
type CobraApp struct {
	rootCmd *cobra.Command
	stdout  io.Writer
	store   *pkg.Store
}

func newApp(name, desc, indexPath string, stdout io.Writer, firesearch *firesearch.AutocompleteService) *CobraApp {
	store := pkg.NewStore(firesearch, indexPath, stdout)

	a := &CobraApp{
		rootCmd: &cobra.Command{
			Use:   name,
			Short: desc,
		},
		stdout: stdout,
		store:  store,
	}
	return a
}

func (a *CobraApp) addCommands(cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		a.rootCmd.AddCommand(cmd)
	}
}

func main() {
	// TODO: get these from env
	endpoint := "http://localhost:8888/api"
	indexPath := "bookmarks"
	secret := ""

	firesearch, err := createFiresearch(endpoint, indexPath, secret)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to crete firesearch service")
		os.Exit(1)
	}

	cobraApp := newApp("bookmarker", "A CLI tool for storing bookmarks with tags that can then be searched later", indexPath, os.Stdout, firesearch)
	cobraApp.addCommands(
		commands.NewAddCommand(cobraApp.store),
		commands.NewFindCommand(cobraApp.store),
	)

	if err := cobraApp.rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
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
