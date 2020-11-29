package main

import (
	"context"

	"github.com/pacedotdev/firesearch-sdk/clients/go/firesearch"
	"github.com/pkg/errors"
)

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
