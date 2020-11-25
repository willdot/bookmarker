package pkg

import (
	"errors"
	"flag"
	"fmt"

	"github.com/pacedotdev/firesearch-sdk/clients/go/firesearch"
)

type Store struct {
	autocompleteService *firesearch.AutocompleteService
}

func NewStore(args []string) (*Store, error) {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	var (
		endpoint  = flags.String("endpoint", "http://localhost:8888/api", "Firesearch API endpoint")
		secret    = flags.String("secret", "firesearch-playground", "Secret API key")
		indexPath = flags.String("index", "firesearch-tutorial/indexes/movies-index", "Firesearch index path")
	)
	if err := flags.Parse(args[1:]); err != nil {
		return nil, err
	}
	if *endpoint == "" {
		return nil, errors.New("missing endpoint")
	}
	if *indexPath == "" {
		return nil, errors.New("missing indexPath")
	}

	client := firesearch.NewClient(*endpoint, *secret)
	autocompleteService := firesearch.NewAutocompleteService(client)

	return &Store{
		autocompleteService: autocompleteService,
	}, nil
}

func (s *Store) SaveBookmark(url string, searchTerms []string) {
	for _, st := range searchTerms {
		fmt.Println(st)
	}

	fmt.Printf("URL: %s\n", url)
}
