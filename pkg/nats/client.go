package nats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

type Client struct {
	Url string
}

func NewClient(ctx context.Context, clientSetup *Client) (*nats.Conn, error) {
	natsClient, err := nats.Connect(clientSetup.Url)
	if err != nil {
		return nil, fmt.Errorf("nats.Connect: %w", err)
	}

	defer natsClient.Drain()

	zerolog.Ctx(ctx).Info().Str("url", clientSetup.Url).Msg("connected to nats")

	return natsClient, nil
}
