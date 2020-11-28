package commands

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type findRunner struct {
	store Store
}

// NewFindCommand will return a new cobra command used for finding bookmarks
func NewFindCommand(store Store) *cobra.Command {
	r := &findRunner{
		store: store,
	}
	cmd := &cobra.Command{
		Use:   "find",
		Short: "Find a bookmark",
		RunE:  r.run,
	}
	return cmd
}

func (r *findRunner) run(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("need to pass in a search term")
	}

	searchTerm := args[0]

	return r.store.FindBookmark(searchTerm)
}
