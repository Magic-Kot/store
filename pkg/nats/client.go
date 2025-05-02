package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

type Client struct {
	Url string
}

func NewClient(ctx context.Context, clientSetup *Client) (*nats.Conn, error) {
	natsClient, err := nats.Connect(
		clientSetup.Url,
		nats.ReconnectWait(2*time.Second),
		nats.MaxReconnects(10),
		nats.DisconnectErrHandler(func(c *nats.Conn, err error) {
			zerolog.Ctx(ctx).Warn().Err(err).Msg("Disconnected")
		}),
		nats.ReconnectHandler(func(c *nats.Conn) {
			zerolog.Ctx(ctx).Warn().Str("Reconnected to", "c.ConnectedUrl()").Msg("Reconnected")
		}),
		nats.ClosedHandler(func(c *nats.Conn) {
			zerolog.Ctx(ctx).Warn().Err(c.LastError()).Msg("Connection closed")
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("nats.Connect: %w", err)
	}

	zerolog.Ctx(ctx).Info().Str("url", clientSetup.Url).Msg("connected to nats")

	return natsClient, nil
}
