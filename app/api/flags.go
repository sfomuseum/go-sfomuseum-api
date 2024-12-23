package api

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var verb string
var api_client_uri string
var access_token_uri string
var params multi.KeyValueString

// DefaultFlagSet returns the default flags (flagset) for running the commandline `api` application.
func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("api")

	fs.StringVar(&api_client_uri, "api-client-uri", "oauth2://?access_token={ACCESS_TOKEN}", "")
	fs.StringVar(&access_token_uri, "access-token-uri", "", "A valid gocloud.dev/runtime variable URI containing a value to replace '{ACCESS_TOKEN}' in the -api-client-uri flag.")

	fs.Var(&params, "param", "One or more KEY=VALUE SFO Museum API parameters")

	fs.StringVar(&verb, "verb", "GET", "The HTTP verb to execute the API method with.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "api is a command-line tool for invoking a SFO Museum API emitting the results to STDOUT.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	return fs
}
