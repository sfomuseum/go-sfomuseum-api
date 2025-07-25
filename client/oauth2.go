package client

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/whosonfirst/go-ioutil"
)

func init() {

	ctx := context.Background()
	err := RegisterClient(ctx, "oauth2", NewOAuth2Client)

	if err != nil {
		panic(err)
	}
}

// OAuth2Client implements the `Client` interface for API endpoints that require OAuth2 access token authentication.
type OAuth2Client struct {
	Client
	// http_client is the underlying net/http client used to perform API requests.
	http_client *http.Client
	// api_endpoint is the URL of the API endpoint.
	api_endpoint string
	// access_token is the OAuth2 access token to append to API requests.
	access_token string
}

// NewOAuth2Client creates a new `OAuth2Client` instance configured by 'uri' which
// is expected to take the form of:
//
//	oauth2://{HOST}/{PATH}?{PARAMETERS}
//
// Where {PARAMETERS} is:
// - `?access_token={TOKEN}` A valid OAuth2 access token.
// - `?insecure={BOOLEAN}` A boolean flag signaling that TLS verification will be skipped.
func NewOAuth2Client(ctx context.Context, uri string) (Client, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	var api_endpoint string

	switch u.Host {
	case "", "api":
		api_endpoint = API_ENDPOINT
	default:
		api_endpoint = fmt.Sprintf("https://%s/%s", u.Host, u.Path)
	}

	q := u.Query()

	http_client := &http.Client{}

	if q.Has("insecure") {

		v, err := strconv.ParseBool(q.Get("insecure"))

		if err != nil {
			return nil, fmt.Errorf("Failed to parse ?insecure=, %w", err)
		}

		if v {

			config := &tls.Config{
				InsecureSkipVerify: true,
			}

			tr := &http.Transport{
				TLSClientConfig: config,
			}

			http_client = &http.Client{
				Transport: tr,
			}
		}
	}

	access_token := q.Get("access_token")

	cl := &OAuth2Client{
		http_client:  http_client,
		access_token: access_token,
		api_endpoint: api_endpoint,
	}

	return cl, nil
}

// ExecuteMethod will perform an API request derived from 'verb' and 'args'.
func (cl *OAuth2Client) ExecuteMethod(ctx context.Context, verb string, args *url.Values) (io.ReadSeekCloser, error) {

	endpoint, err := url.Parse(cl.api_endpoint)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse endpoint URI, %w", err)
	}

	// Note that we do `req = req.WithContext(ctx)` below so
	// there is no need to do it here too.

	var req *http.Request

	switch verb {
	case http.MethodGet:

		endpoint.RawQuery = args.Encode()

		r, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)

		if err != nil {
			return nil, err
		}

		req = r

	case http.MethodPost:

		args_r := strings.NewReader(args.Encode())

		r, err := http.NewRequest(http.MethodPost, endpoint.String(), args_r)

		if err != nil {
			return nil, err
		}

		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		req = r

	default:
		return nil, fmt.Errorf("Invalid or unsupported verb")
	}

	return cl.executeRequest(ctx, req)
}

// executeRequest will perform an API request derived from 'req'.
func (cl *OAuth2Client) executeRequest(ctx context.Context, req *http.Request) (io.ReadSeekCloser, error) {

	req = req.WithContext(ctx)

	if cl.access_token != "" {
		b64_token := base64.StdEncoding.EncodeToString([]byte(cl.access_token))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", b64_token))
	}

	rsp, err := cl.http_client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Failed to execute API request, %w", err)
	}

	if rsp.StatusCode != http.StatusOK {
		rsp.Body.Close()
		return nil, fmt.Errorf("API call failed with status '%s'", rsp.Status)
	}

	return ioutil.NewReadSeekCloser(rsp.Body)
}
