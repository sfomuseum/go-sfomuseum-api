package client

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/url"
	"testing"
)

var api_client_uri = flag.String("api-client-uri", "", "A valid sfomuseum/go-sfomuseum-api.Client URI. If non-empty then it will be used to perform additional tests.")

func TestExecuteMethodPaginated(t *testing.T) {

	if *api_client_uri == "" {
		t.Skip()
	}

	ctx := context.Background()

	cl, err := NewClient(ctx, *api_client_uri)

	if err != nil {
		t.Fatalf("Failed to create new client for '%s', %v", *api_client_uri, err)
	}

	args := &url.Values{}
	args.Set("method", "sfomuseum.collection.search")
	args.Set("q", "747")

	cb := func(ctx context.Context, r io.ReadSeekCloser, err error) error {

		if err != nil {
			fmt.Println(err)
		}

		return nil
	}

	err = ExecuteMethodPaginatedWithClient(ctx, cl, args, cb)

	if err != nil {
		t.Fatalf("Failed to execute method paginated, %v", err)
	}

}
