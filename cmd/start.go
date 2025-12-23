package cmd

import (
	"context"
	"errors"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pharos-exporter/internal"
)

func runStart(args []string) error {
	fs := flag.NewFlagSet("start", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)

	rpcURL := fs.String("rpc", "https://atlantic.dplabs-internal.com/", "JSON-RPC endpoint")
	myBlsKey := fs.String("my-bls-key", "", "my BLS pubkey (0x...)")
	checkBlockProof := fs.Bool("check-block-proof", true, "check signedBlsKeys trigger")
	checkValidatorSet := fs.Bool("check-validator-set", true, "check validator set trigger")
	pollInterval := fs.Duration("poll-interval", 1*time.Second, "poll interval for latest block")
	if err := fs.Parse(args); err != nil {
		return err
	}

	tracker, err := internal.NewBlockTracker(internal.BlockTrackerConfig{
		RPCURL:            *rpcURL,
		MyBlsKey:          *myBlsKey,
		CheckBlockProof:   *checkBlockProof,
		CheckValidatorSet: *checkValidatorSet,
		PollInterval:      *pollInterval,
	})
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := tracker.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}
