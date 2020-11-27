package cmd

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/willdot/bookmarker/pkg"
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Find a bookmark",
	RunE: func(cmd *cobra.Command, args []string) error {
		return findBookmark(args)
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
}

func findBookmark(args []string) error {
	if len(args) < 1 {
		return errors.New("need to pass in a search term")
	}

	searchTerm := args[0]

	firesearchIndexer, err := createFiresearch(indexPath, endpoint, secret)
	if err != nil {
		return errors.Wrap(err, "failed to create firesearch")
	}
	store := pkg.NewStore(firesearchIndexer, indexPath, os.Stdout)

	return store.FindBookmark(searchTerm)
}
