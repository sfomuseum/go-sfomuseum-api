package client

import (
	"context"
	"flag"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"testing"
)

var api_client_uri = flag.String("api-client-uri", "", "A valid sfomuseum/go-sfomuseum-api.Client URI. If non-empty then it will be used to perform additional tests.")

func TestExecuteMethodPaginated(t *testing.T) {

	if *api_client_uri == "" {
		slog.Info("-api-client-uri flag is empty, skipping test")
		t.Skip()
	}

	ctx := context.Background()

	cl, err := NewClient(ctx, *api_client_uri)

	if err != nil {
		t.Fatalf("Failed to create new client for '%s', %v", *api_client_uri, err)
	}

	args := &url.Values{}
	args.Set("method", "sfomuseum.collection.objects.search")
	args.Set("q", "747")

	for r, err := range ExecuteMethodPaginatedWithClient(ctx, cl, http.MethodGet, args) {

		if err != nil {
			t.Fatalf("Failed to execute method paginated, %v", err)
		}

		_, err = io.Copy(io.Discard, r)

		if err != nil {
			t.Fatalf("Failed to read API result, %v", err)
		}

	}

}
