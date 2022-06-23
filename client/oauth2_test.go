package client

import (
	"context"
	"encoding/json"
	_ "fmt"
	"net/url"
	"testing"
)

func TestOAuth2ClientRoster(t *testing.T) {

	ctx := context.Background()

	uris := []string{
		"oauth2://collection?access_token={TOKEN}",
		"oauth2://millsfield?access_token={TOKEN}",
		"oauth2://localhost:8080/api/rest?access_token={TOKEN}",
	}

	for _, uri := range uris {

		_, err := NewClient(ctx, uri)

		if err != nil {
			t.Fatalf("Failed to create new client for '%s', %v", uri, err)
		}
	}

	if *api_client_uri != "" {

		cl, err := NewClient(ctx, *api_client_uri)

		if err != nil {
			t.Fatalf("Failed to create new client for '%s', %v", *api_client_uri, err)
		}

		args := &url.Values{}
		args.Set("method", "api.test.echo")
		args.Set("hello", "world")

		rsp, err := cl.ExecuteMethod(ctx, args)

		if err != nil {
			t.Fatalf("Failed to execute method, %v", err)
		}

		var tmp map[string]string

		dec := json.NewDecoder(rsp)
		err = dec.Decode(&tmp)

		if err != nil {
			t.Fatalf("Failed to decode API results, %v", err)
		}

		v, ok := tmp["hello"]

		if !ok {
			t.Fatalf("API result missing 'hello' key")
		}

		if v != "world" {
			t.Fatalf("Unexpected value for 'hello' key: %s", v)
		}

	}
}
