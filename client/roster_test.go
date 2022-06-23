package client

import (
	"context"
	"testing"
)

func TestClientRoster(t *testing.T) {

	ctx := context.Background()

	err := RegisterClient(ctx, "oauth2", NewOAuth2Client)

	if err == nil {
		t.Fatalf("Expected NewOAuth2Client to be registered already")
	}

}
