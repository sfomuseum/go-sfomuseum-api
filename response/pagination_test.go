package response

import (
	"bytes"
	"context"
	_ "embed"
	"github.com/whosonfirst/go-ioutil"
	"testing"
)

//go:embed response.json
var rsp []byte

func TestDervivePagination(t *testing.T) {

	ctx := context.Background()

	r := bytes.NewReader(rsp)

	fh, err := ioutil.NewReadSeekCloser(r)

	if err != nil {
		t.Fatalf("Failed to create NewReadSeekCloser, %v", err)
	}

	pg, err := DerivePagination(ctx, fh)

	if err != nil {
		t.Fatalf("Failed to derive pagination, %v", err)
	}

	if pg.Pages != 688 {
		t.Fatalf("Unexpected pages count: %d", pg.Pages)
	}

}
