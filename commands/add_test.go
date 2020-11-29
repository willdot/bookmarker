package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStore struct {
}

func (m *mockStore) SaveBookmark(bookmarkName, url string, tags []string) error {
	return nil
}
func (m *mockStore) FindBookmark(searchTerm string) error {
	return nil
}

func TestAddCommandRunInvalidRunner(t *testing.T) {
	mockStore := &mockStore{}

	addRunner := addRunner{
		store: mockStore,
	}

	tt := map[string]struct {
		urlFlag              string
		nameFlag             string
		expectedErrorMessage string
	}{
		"missing url":           {urlFlag: "", nameFlag: "", expectedErrorMessage: "url flag is required"},
		"missing bookmark name": {urlFlag: "http://interwebz", nameFlag: "", expectedErrorMessage: "bookmark name is required"},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			addRunner.url = tc.urlFlag
			addRunner.bookmark = tc.nameFlag

			err := addRunner.run(nil, nil)

			require.Error(t, err)
			assert.Equal(t, tc.expectedErrorMessage, err.Error())
		})
	}
}
