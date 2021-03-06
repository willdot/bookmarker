package pkg

import (
	"context"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/pkg/errors"

	"github.com/pacedotdev/firesearch-sdk/clients/go/firesearch"
)

// Store handles interactions of putting and getting bookmarks from storage
type Store struct {
	indexPath string
	indexer   Indexer
	writer    io.Writer
}

// Indexer defines the functions required to interact with a store of documents
type Indexer interface {
	PutDoc(ctx context.Context, r firesearch.PutAutocompleteDocRequest) (*firesearch.PutAutocompleteDocResponse, error)
	Complete(ctx context.Context, r firesearch.CompleteRequest) (*firesearch.CompleteResponse, error)
}

// NewStore returns a new store
func NewStore(indexer Indexer, indexPath string, writer io.Writer) *Store {
	return &Store{
		indexer:   indexer,
		indexPath: indexPath,
		writer:    writer,
	}
}

// SaveBookmark will save a new bookmark, or update an existing one with the given details
func (s *Store) SaveBookmark(bookmarkName, url string, tags []string) error {
	// if the bookmark name isn't in the tags, add it in so it can be searched for
	if checkIfBookmarkNameIsInTags(tags, bookmarkName) == false {
		tags = append(tags, bookmarkName)
	}

	err := s.addDocument(context.TODO(), bookmarkName, url, tags)
	if err != nil {
		return errors.Wrap(err, "error storing bookmark document")
	}

	return nil
}

func checkIfBookmarkNameIsInTags(tags []string, bookmarkName string) bool {
	for _, item := range tags {
		if item == bookmarkName {
			return true
		}
	}
	return false
}

func (s *Store) addDocument(ctx context.Context, bookmarkName, url string, tags []string) error {
	_, err := s.indexer.PutDoc(ctx, firesearch.PutAutocompleteDocRequest{
		IndexPath: s.indexPath,
		Doc: firesearch.AutocompleteDoc{
			ID:   bookmarkName,
			Text: strings.Join(tags, " "),
			Fields: []firesearch.Field{
				{
					Key:   "url",
					Value: url,
				},
			},
		},
	})
	if err != nil {
		return errors.Wrap(err, "error putting document in the indexer")
	}

	return nil
}

// FindBookmark will attempt to search for a bookmark using the given search term and will print the results
func (s *Store) FindBookmark(searchTerm string) error {
	return s.getDocuments(context.TODO(), searchTerm)
}

func (s *Store) getDocuments(ctx context.Context, searchTerm string) error {
	res, err := s.indexer.Complete(ctx, firesearch.CompleteRequest{
		Query: firesearch.CompleteQuery{
			IndexPath: s.indexPath,
			Text:      searchTerm,
		},
	})

	if err != nil {
		return errors.Wrap(err, "error calling complate on the indexer")
	}

	if len(res.Hits) == 0 {
		fmt.Fprintf(s.writer, "no search result found for: %s\n", searchTerm)
		return nil
	}

	w := new(tabwriter.Writer)
	w.Init(s.writer, 0, 8, 0, '\t', 0)
	defer w.Flush()

	fmt.Fprintln(w, "Name\tURL\t")
	fmt.Fprintln(w, "\t")
	for _, hit := range res.Hits {
		fmt.Fprintf(w, "%s\t%s\n", hit.ID, hit.Fields[0].Value)
	}

	fmt.Fprintln(w)
	return nil
}
