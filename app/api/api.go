package api

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/sfomuseum/go-sfomuseum-api/client"
)

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := RunOptionsFromFlagSet(fs)

	if err != nil {
		return err
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	cl, err := client.NewClient(ctx, opts.APIClientURI)

	if err != nil {
		return fmt.Errorf("Failed to create new API client, %w", err)
	}

	rsp, err := cl.ExecuteMethod(ctx, opts.Verb, opts.Args)

	if err != nil {
		return fmt.Errorf("Failed to execute API method, %w", err)
	}

	_, err = io.Copy(os.Stdout, rsp)

	if err != nil {
		return fmt.Errorf("Failed to copy response to STDOUT, %w", err)
	}

	return nil
}
