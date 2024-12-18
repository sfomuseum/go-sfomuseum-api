// api is a command-line tool for invoking a SFO Museum API emitting the results to STDOUT.
package main

import (
	"context"
	"log"

	_ "gocloud.dev/runtimevar/awsparamstore"
	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"

	"github.com/sfomuseum/go-sfomuseum-api/app/api"
)

func main() {

	ctx := context.Background()
	err := api.Run(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
