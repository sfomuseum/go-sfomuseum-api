package main

import (
	_ "gocloud.dev/runtimevar/awsparamstore"
	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"
)

import (
	"context"
	"flag"
	_ "fmt"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/sfomuseum/go-sfomuseum-api/client"
	"github.com/sfomuseum/runtimevar"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
)

func main() {

	api_client_uri := flag.String("api-client-uri", "oauth2://collection?access_token={ACCESS_TOKEN}", "")
	access_token_uri := flag.String("access-token-uri", "", "A valid gocloud.dev/runtime variable URI containing a value to replace '{ACCESS_TOKEN}' in the -api-client-uri flag.")

	var params multi.KeyValueString
	flag.Var(&params, "param", "")

	flag.Parse()

	ctx := context.Background()

	if *access_token_uri != "" {

		token, err := runtimevar.StringVar(ctx, *access_token_uri)

		if err != nil {
			log.Fatalf("Failed to retrieve access token, %v", err)
		}

		*api_client_uri = strings.Replace(*api_client_uri, "{ACCESS_TOKEN}", token, 1)
	}

	cl, err := client.NewClient(ctx, *api_client_uri)

	if err != nil {
		log.Fatalf("Failed to create new API client, %w", err)
	}

	q := &url.Values{}

	for _, kv := range params {
		q.Set(kv.Key(), kv.Value().(string))
	}

	rsp, err := cl.ExecuteMethod(ctx, q)

	if err != nil {
		log.Fatalf("Failed to execute API method, %v", err)
	}

	_, err = io.Copy(os.Stdout, rsp)

	if err != nil {
		log.Fatalf("Failed to copy response to STDOUT, %v", err)
	}
}
