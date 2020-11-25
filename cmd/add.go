package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/willdot/bookmarker/pkg"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new bookmar,",
	RunE: func(cmd *cobra.Command, args []string) error {
		return saveBookmark(args)
	},
}

var bookmark string
var searchTerms []string

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVar(&bookmark, "url", "", "Help message for toggle")
	addCmd.Flags().StringSliceVar(&searchTerms, "terms", []string{}, `pass in any search terms you wish to use, seperated by a comma. Example: --terms="first term" "second"`)
}

func saveBookmark(args []string) error {
	store, err := pkg.NewStore(args)
	if err != nil {
		return errors.Wrap(err, "failed to create store")
	}

	store.SaveBookmark(bookmark, searchTerms)

	return nil
}
