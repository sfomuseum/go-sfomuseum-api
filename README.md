# go-sfomuseum-api

Go package providing methods for interacting with the [SFO Museum API](https://api.sfomuseum.org).

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/sfomuseum/go-sfomuseum-api.svg)](https://pkg.go.dev/github.com/sfomuseum/go-sfomuseum-api)

## Example

```
package main

import (
	"context"
	"io"
	"net/url"
	"os"

	"github.com/sfomuseum/go-sfomuseum-api/client"
)

func main() {

	ctx := context.Background()

	client_uri := "oauth2://?access_token={TOKEN}"
	cl, _ := client.NewClient(ctx, client_uri)

	args := &url.Values{}
	args.Set("method", "api.spec.methods")

	fh, _ := cl.ExecuteMethod(ctx, "GET", args)
	defer fh.Close()
	
	io.Copy(os.Stdout, fh)
}
```

## Design

The core of this package's approach to the SFO Museum API is the `ExecuteMethod` method (which is defined in the `client.Client` interface) whose signature looks like this:

```
ExecuteMethod(context.Context, string, *url.Values) (io.ReadSeekCloser, error)
```

This package only defines [a handful of Go types or structs mapping to individual API responses](response). In time there may be others, along with helper methods for unmarshaling API responses in to typed responses but the baseline for all operations will remain: Query paramters (`url.Values`) sent over HTTP returning an `io.ReadSeekCloser` instance that is inspected and validated according to the needs and uses of the tools using the SFO Museum API.

## Clients

The `client.Client` interface provides for common methods for accessing the SFO Museum API. Currently there is only a single client interface that calls the SFO Museum API using the OAuth2 authentication and authorization scheme.

Clients are instantiated using a URI-based syntax where the scheme and query parameters map to specific implementation details.

### OAuth2

The OAuth2 `Client` implementation is instantiated using the `oauth2://` scheme. For example:

```
oauth2://?{QUERY_PARAMETERS}
```

Valid query parameters are:

| Name | Value | Required |
| --- | --- | --- |
| `access_token` | string | yes |

## Tools

```
$> make cli
go build -mod vendor -ldflags="-s -w" -o bin/api cmd/api/main.go
go build -mod vendor -ldflags="-s -w" -o bin/test-methods cmd/test-methods/main.go
```

### api

api is a command-line tool for invoking a SFO Museum API emitting the results to STDOUT.

```
$> ./bin/api -h
api is a command-line tool for invoking a SFO Museum API emitting the results to STDOUT.

Usage:
	 ./bin/api [options]

Valid options are:
  -access-token-uri string
    	A valid gocloud.dev/runtime variable URI containing a value to replace '{ACCESS_TOKEN}' in the -api-client-uri flag.
  -api-client-uri string
    	 (default "oauth2://?access_token={ACCESS_TOKEN}")
  -param value
    	One or more KEY=VALUE SFO Museum API parameters
  -verb string
    	The HTTP verb to execute the API method with. (default "GET")
```

For example:

```
$> ./bin/api \
	-access-token-uri 'constant://?val={TOKEN}' \
	-param method=sfomuseum.collection.objects.getInfo \
	-param object_id=1763138037 \
	| jq -r '.object["wof:name"]'
	
model airplane: Cathay Pacific Airways, Boeing 747-400
```

### Runtime variables

Under the hood the `api` tool uses the [sfomuseum/runtimevar](https://github.com/sfomuseum/runtimevar) package for resolving runtime variables. This package is a thin wrapper around the [gocloud.dev/runtimevar](https://pkg.go.dev/gocloud.dev/runtimevar) package which enables the following `runtimevar` providers by default:

* [AWS ParameterStore values](https://gocloud.dev/howto/runtimevar/#awsps)
* [Constant and file-based values](https://gocloud.dev/howto/runtimevar/#local)

Additionally it provides support for label-based AWS credential string, for example:

```
awsparamstore://hello-world?region=us-west-2&credentials={LABEL}
```

Valid credential labels are:

| Label | Description |
| --- | --- |
| `anon:` | Empty or anonymous credentials. |
| `env:` | Read credentials from AWS defined environment variables. |
| `iam:` | Assume AWS IAM credentials are in effect. |
| `sts:{ARN}` | Assume the role defined by `{ARN}` using STS credentials. |
| `{AWS_PROFILE_NAME}` | This this profile from the default AWS credentials location. |
| `{AWS_CREDENTIALS_PATH}:{AWS_PROFILE_NAME}` | This this profile from a user-defined AWS credentials location. |

## See also

* https://api.sfomuseum.org/
* https://pkg.go.dev/gocloud.dev/runtimevar
* https://github.com/sfomuseum/runtimevar