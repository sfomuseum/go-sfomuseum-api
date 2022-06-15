package client

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-ioutil"
	"io"
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
	http_client        *http.Client
	api_endpoint       *url.URL
	access_token        string
}

func NewOAuth2Client(ctx context.Context, uri string) (Client, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	access_token := q.Get("access_token")
	
	http_client := &http.Client{}

	cl := &OAuth2Client{
		http_client:     http_client,
		access_token: access_token,
	}

	return cl, nil
}

func (cl *OAuth2Client) ExecuteMethod(ctx context.Context, args *url.Values) (io.ReadSeekCloser, error) {

	endpoint, err := url.Parse(API_ENDPOINT)

	if err != nil {
		return nil, err
	}

	http_method := "GET"

	if cl.access_token != "" {
		args.Set("access_token", cl.access_token)
	}

	endpoint.RawQuery = args.Encode()

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
