package commands

import (
	"errors"

	"github.com/spf13/cobra"
)

// Store defines the functions required for adding and searching bookmarks
type Store interface {
	SaveBookmark(bookmarkName, url string, tags []string) error
	FindBookmark(searchTerm string) error
}

type addRunner struct {
	bookmark string
	url      string
	tags     []string

	store Store
}

// NewAddCommand will return a cobra command used for adding bookmarks
func NewAddCommand(store Store) *cobra.Command {
	r := &addRunner{
		store: store,
	}
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new bookmark",
		RunE:  r.run,
	}

	cmd.Flags().StringVar(&r.url, "url", "", "the url of the bookmark")
	cmd.Flags().StringVarP(&r.bookmark, "name", "n", "", "the name of the bookmark")
	cmd.Flags().StringSliceVarP(&r.tags, "tags", "t", []string{}, `pass in any tags you wish to use, seperated by a comma. Example: --tags="github, source control"`)
	return cmd
}

func (r *addRunner) run(_ *cobra.Command, _ []string) error {
	if r.url == "" {
		return errors.New("url flag is required")
	}

	if r.bookmark == "" {
		return errors.New("bookmark name is required")
	}

	return r.store.SaveBookmark(r.bookmark, r.url, r.tags)
}
