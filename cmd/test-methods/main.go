// api is a command-line tool for invoking a SFO Museum API emitting the results to STDOUT.
package main

import (
	_ "gocloud.dev/runtimevar/awsparamstore"
	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"
)

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"strings"

	"github.com/sfomuseum/go-flags/multi"
	"github.com/sfomuseum/go-sfomuseum-api/client"
	"github.com/sfomuseum/go-sfomuseum-api/response"
	"github.com/sfomuseum/runtimevar"
)

func main() {

	api_client_uri := flag.String("api-client-uri", "oauth2://?access_token={ACCESS_TOKEN}", "")
	access_token_uri := flag.String("access-token-uri", "", "A valid gocloud.dev/runtime variable URI containing a value to replace '{ACCESS_TOKEN}' in the -api-client-uri flag.")

	var skip multi.MultiString
	flag.Var(&skip, "skip", "...")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "api is a command-line tool for invoking a SFO Museum API emitting the results to STDOUT.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		flag.PrintDefaults()
	}

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
		log.Fatalf("Failed to create new API client, %v", err)
	}

	q := &url.Values{}
	q.Set("method", "api.spec.methods")

	rsp, err := cl.ExecuteMethod(ctx, q)

	if err != nil {
		log.Fatalf("Failed to execute API method, %v", err)
	}

	type MethodsResponse struct {
		Methods []*response.Method `json:"methods"`
	}

	var methods_rsp *MethodsResponse

	dec := json.NewDecoder(rsp)
	err = dec.Decode(&methods_rsp)

	if err != nil {
		log.Fatalf("Failed to decode methods, %v", err)
	}

	for _, m := range methods_rsp.Methods {

		for _, prefix := range skip {
			if strings.HasPrefix(m.Name, prefix) {
				slog.Info("Method matches skip prefix, skipping", "method", m.Name, "prefix", prefix)
				continue
			}
		}

		if m.RequestMethod != "GET" {
			slog.Info("HTTP method is not GET, skipping", "method", m.Name, "verb", m.RequestMethod)
			continue
		}

		params := &url.Values{}
		params.Set("method", m.Name)

		for _, p := range m.Parameters {
			params.Set(p.Name, fmt.Sprintf("%v", p.Example))
		}

		slog.Info("Execute method", "method", m.RequestMethod, "parameters", params.Encode())

		_, err := cl.ExecuteMethod(ctx, params)

		if err != nil {
			slog.Error("Failed to execute method", "method", m.Name, "error", err)
			continue
		}

		slog.Debug("Method successful", "method", m.Name)
	}
}
