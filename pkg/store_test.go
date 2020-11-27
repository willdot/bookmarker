package pkg_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/pacedotdev/firesearch-sdk/clients/go/firesearch"
	"github.com/stretchr/testify/assert"
	"github.com/willdot/bookmarker/pkg"
)

type MockIndexer struct {
}

func (m *MockIndexer) PutDoc(ctx context.Context, r firesearch.PutAutocompleteDocRequest) (*firesearch.PutAutocompleteDocResponse, error) {

	return nil, nil
}
func (m *MockIndexer) Complete(ctx context.Context, r firesearch.CompleteRequest) (*firesearch.CompleteResponse, error) {

	res := &firesearch.CompleteResponse{
		Hits: []firesearch.AutocompleteDoc{
			{
				ID: "test",
				Fields: []firesearch.Field{
					{
						Value: "something",
					},
				},
			},
		},
	}
	return res, nil
}

func TestFindBookmark(t *testing.T) {
	var buf bytes.Buffer

	mockIndexer := &MockIndexer{}

	store := pkg.NewStore(mockIndexer, "path", &buf)

	err := store.FindBookmark("test")
	assert.NoError(t, err)

	assert.Contains(t, buf.String(), "test")
	assert.Contains(t, buf.String(), "something")
	assert.NotContains(t, buf.String(), "NOOOOO")
}
