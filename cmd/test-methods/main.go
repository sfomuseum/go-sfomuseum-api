// test-methods is a command line tool for iterating through the list of available SFO Museum API
// methods and executing all methods which use the HTTP "GET" verb testing for a successful response.
package main

/*

go run cmd/test-methods/main.go \
	-api-client-uri 'oauth2://?access_token={TOKEN}' \
	-skip sfomuseum.you.shoebox \
	-break-on-error

*/

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	_ "gocloud.dev/runtimevar/awsparamstore"
	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"

	"github.com/sfomuseum/go-flags/multi"
	"github.com/sfomuseum/go-sfomuseum-api/v2/client"
	"github.com/sfomuseum/go-sfomuseum-api/v2/response"
	"github.com/sfomuseum/runtimevar"
	"github.com/tidwall/gjson"
)

func main() {

	api_client_uri := flag.String("api-client-uri", "oauth2://?access_token={ACCESS_TOKEN}", "")
	access_token_uri := flag.String("access-token-uri", "", "A valid gocloud.dev/runtime variable URI containing a value to replace '{ACCESS_TOKEN}' in the -api-client-uri flag.")

	var skip multi.MultiString
	flag.Var(&skip, "skip", "Zero or more matching prefixes for method names to skip.")

	var break_on_error bool
	flag.BoolVar(&break_on_error, "break-on-error", false, "...")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "test-methods is a command-line tool...\n\n")
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

	rsp, err := cl.ExecuteMethod(ctx, http.MethodGet, q)

	if err != nil {
		log.Fatalf("Failed to execute API method, %v", err)
	}

	type MethodsResponse struct {
		Methods []*response.Method `json:"methods"`
	}

	var methods_rsp *MethodsResponse

	dec := json.NewDecoder(rsp)
	err = dec.Decode(&methods_rsp)

	// If the methods can be resolved as a set loop through them individually
	// and report the definitions that can't be parsed.

	if err != nil {

		slog.Error("Failed to decode methods", "error", err)

		_, err = rsp.Seek(0, 0)

		if err != nil {
			log.Fatalf("Failed to rewind response reader, %v", err)
		}

		rsp_body, err := io.ReadAll(rsp)

		if err != nil {
			log.Fatalf("Failed to read response body, %v", err)
		}

		method_rsp := gjson.GetBytes(rsp_body, "methods")

		for _, m_rsp := range method_rsp.Array() {

			var m *response.Method

			err := json.Unmarshal([]byte(m_rsp.String()), &m)

			if err != nil {
				slog.Error("Invalid method definition", "def", m_rsp.String())
				log.Fatalf("Failed to unmarshal method, %v", err)
			}

		}
	}

	// Execute each method in turn

	for _, m := range methods_rsp.Methods {

		skip_method := false

		for _, prefix := range skip {
			if strings.HasPrefix(m.Name, prefix) {
				slog.Info("Method matches skip prefix, skipping", "method", m.Name, "prefix", prefix)
				skip_method = true
				break
			}
		}

		if skip_method {
			continue
		}

		if m.RequestMethod != "GET" {
			slog.Info("HTTP method is not GET, skipping", "method", m.Name, "verb", m.RequestMethod)
			continue
		}

		params := &url.Values{}
		params.Set("method", m.Name)

		for _, p := range m.Parameters {

			switch m.Name {
			case "sfomuseum.collection.objects.search":

				if p.Required {
					params.Set(p.Name, fmt.Sprintf("%s", p.Example))
				}

			default:
				params.Set(p.Name, fmt.Sprintf("%s", p.Example))
			}

		}

		// slog.Debug("Execute method", "method", m.RequestMethod, "parameters", params.Encode())

		_, err := cl.ExecuteMethod(ctx, m.RequestMethod, params)

		if err != nil {

			if m.Name != "api.test.error" {
				slog.Error("Failed to execute method", "method", m.Name, "parameters", params.Encode(), "error", err)

				if break_on_error {
					break
				}

				continue
			}
		}

		slog.Debug("Method successful", "method", m.Name, "parameters", params.Encode())
	}
}
