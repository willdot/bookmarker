package pkg

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/pkg/errors"

	"github.com/pacedotdev/firesearch-sdk/clients/go/firesearch"
)

const (
	endpoint  = "http://localhost:8888/api"
	secret    = "firesearch-playground"
	indexPath = "bookmarks"
)

type Store struct {
	autocompleteService *firesearch.AutocompleteService
	indexPath           string
}

func NewStore() (*Store, error) {

	client := firesearch.NewClient(endpoint, secret)
	autocompleteService := firesearch.NewAutocompleteService(client)
	createIndex(context.TODO(), autocompleteService, indexPath)

	return &Store{
		autocompleteService: autocompleteService,
		indexPath:           indexPath,
	}, nil
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

func (s *Store) SaveBookmark(bookmarkName, url string, searchTerms []string) error {

	// if the bookmark name isn't in the search terms, add it in so it can be searched for
	if checkIfBookmarkNameIsInSearchTerms(searchTerms, bookmarkName) == false {
		searchTerms = append(searchTerms, bookmarkName)
	}

	err := s.addDocument(context.TODO(), bookmarkName, url, searchTerms)
	if err != nil {
		return errors.Wrap(err, "error storing bookmark document")
	}

	return nil
}

func checkIfBookmarkNameIsInSearchTerms(searchTerms []string, bookmarkName string) bool {
	for _, item := range searchTerms {
		if item == bookmarkName {
			return true
		}
	}
	return false
}

func (s *Store) addDocument(ctx context.Context, bookmarkName, url string, searchTerms []string) error {
	_, err := s.autocompleteService.PutDoc(ctx, firesearch.PutAutocompleteDocRequest{
		IndexPath: s.indexPath,
		Doc: firesearch.AutocompleteDoc{
			ID:   bookmarkName,
			Text: strings.Join(searchTerms, " "),
			Fields: []firesearch.Field{
				{
					Key:   "url",
					Value: url,
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) FindBookmark(searchTerm string) error {
	return s.getDocuments(context.TODO(), searchTerm)
}

func (s *Store) getDocuments(ctx context.Context, searchTerm string) error {
	res, err := s.autocompleteService.Complete(ctx, firesearch.CompleteRequest{
		Query: firesearch.CompleteQuery{
			IndexPath: s.indexPath,
			Text:      searchTerm,
		},
	})

	if err != nil {
		return err
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	defer w.Flush()

	fmt.Fprintln(w, "Name\tURL\t")
	fmt.Fprintln(w, "\t")
	for _, hit := range res.Hits {
		fmt.Fprintf(w, "%s\t%s\n", hit.ID, hit.Fields[0].Value)
	}

	fmt.Fprintln(w)
	return nil
}
