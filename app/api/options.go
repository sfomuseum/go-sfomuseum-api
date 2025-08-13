package api

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"strings"

	"github.com/aaronland/gocloud/runtimevar"
	"github.com/sfomuseum/go-flags/flagset"
)

// RunOptions defines options and configurations to execute the commandline `api` application with.
type RunOptions struct {
	// The HTTP verb to execute an API method with.
	Verb string
	// A valid `sfomuseum/go-sfomuseum-api/client.Client` URI.
	APIClientURI string
	// One or more query parameters to execute an API method with (at a minimum the "method" parameter is required).
	Args *url.Values
}

// RunOptionsFromFlagSet will return a new `RunOptions` for use with the commandline `api` application derived from 'fs'.
func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	if access_token_uri != "" {

		ctx := context.Background()
		token, err := runtimevar.StringVar(ctx, access_token_uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to retrieve access token, %v", err)
		}

		api_client_uri = strings.Replace(api_client_uri, "{ACCESS_TOKEN}", token, 1)
	}

	args := &url.Values{}

	for _, kv := range params {
		args.Set(kv.Key(), kv.Value().(string))
	}

	opts := &RunOptions{
		APIClientURI: api_client_uri,
		Args:         args,
		Verb:         verb,
	}

	return opts, nil
}
