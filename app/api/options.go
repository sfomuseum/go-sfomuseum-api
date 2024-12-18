package api

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"strings"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/runtimevar"
)

type RunOptions struct {
	Verb         string
	APIClientURI string
	Args         *url.Values
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	// START OF move this in to RunWithOptions...?

	if access_token_uri != "" {

		ctx := context.Background()
		token, err := runtimevar.StringVar(ctx, access_token_uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to retrieve access token, %v", err)
		}

		api_client_uri = strings.Replace(api_client_uri, "{ACCESS_TOKEN}", token, 1)
	}

	// END OF move this in to RunWithOptions...?

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
