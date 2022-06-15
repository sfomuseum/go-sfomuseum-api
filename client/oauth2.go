package client

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-ioutil"
	"io"
	"log"
	"net/http"
	"net/url"
)

func init() {

	ctx := context.Background()
	err := RegisterClient(ctx, "oauth2", NewOAuth2Client)

	if err != nil {
		panic(err)
	}
}

type OAuth2Client struct {
	http_client  *http.Client
	api_endpoint string
	access_token string
}

func NewOAuth2Client(ctx context.Context, uri string) (Client, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	var api_endpoint string

	switch u.Host {
	case "collection":
		api_endpoint = COLLECTION_ENDPOINT
	case "millsfield":
		api_endpoint = MILLSFIELD_ENDPOINT
	default:
		api_endpoint = fmt.Sprintf("https://%s/%s", u.Host, u.Path)
	}

	http_client := &http.Client{}

	q := u.Query()

	access_token := q.Get("access_token")

	cl := &OAuth2Client{
		http_client:  http_client,
		access_token: access_token,
		api_endpoint: api_endpoint,
	}

	return cl, nil
}

func (cl *OAuth2Client) ExecuteMethod(ctx context.Context, args *url.Values) (io.ReadSeekCloser, error) {

	endpoint, err := url.Parse(cl.api_endpoint)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse endpoint URI, %w", err)
	}

	http_method := "GET"

	if cl.access_token != "" {
		args.Set("access_token", cl.access_token)
	}

	endpoint.RawQuery = args.Encode()

	log.Println(endpoint.String())

	req, err := http.NewRequest(http_method, endpoint.String(), nil)

	if err != nil {
		return nil, err
	}

	return cl.call(ctx, req)
}

func (cl *OAuth2Client) call(ctx context.Context, req *http.Request) (io.ReadSeekCloser, error) {

	req = req.WithContext(ctx)

	rsp, err := cl.http_client.Do(req)

	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != http.StatusOK {
		rsp.Body.Close()
		return nil, fmt.Errorf("API call failed with status '%s'", rsp.Status)
	}

	return ioutil.NewReadSeekCloser(rsp.Body)
}
