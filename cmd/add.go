package cmd

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/willdot/bookmarker/pkg"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new bookmark",
	RunE: func(cmd *cobra.Command, args []string) error {
		return saveBookmark(args)
	},
}

var bookmark string
var url string
var tags []string

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVar(&url, "url", "", "the url of the bookmark")
	addCmd.Flags().StringVarP(&bookmark, "name", "n", "", "the name of the bookmark")
	addCmd.Flags().StringSliceVarP(&tags, "tags", "t", []string{}, `pass in any tags you wish to use, seperated by a comma. Example: --tags="github, source control"`)
}

func saveBookmark(args []string) error {
	//TODO: verify flags are provided
	firesearchIndexer, err := createFiresearch(indexPath, endpoint, secret)
	if err != nil {
		return errors.Wrap(err, "failed to create firesearch")
	}
	store := pkg.NewStore(firesearchIndexer, indexPath, os.Stdout)

	return store.SaveBookmark(bookmark, url, tags)
}
