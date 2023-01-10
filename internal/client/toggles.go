package client

import (
	"context"

	"github.com/featureguards/featureguards-go/v2/internal/meta"
	pb_toggles "github.com/featureguards/featureguards-go/v2/proto/toggles"
)

// Fetch fetches new feature toggles since ftVersion and settingsVersion
func (c *Client) Fetch(ctx context.Context, accessToken string, ftVersion, settingsVersion int64) (*pb_toggles.FetchResponse, error) {
	return c.togglesClient.Fetch(withJwtToken(ctx, accessToken), &pb_toggles.FetchRequest{Version: ftVersion, SettingsVersion: settingsVersion})
}

func (c *Client) Listen(ctx context.Context, accessToken string, ftVersion, settingsVersion int64) (<-chan *pb_toggles.ListenPayload, error) {
	stream, err := c.togglesClient.Listen(withJwtToken(ctx, accessToken), &pb_toggles.ListenRequest{Version: ftVersion, SettingsVersion: settingsVersion})
	if err != nil {
		return nil, err
	}
	ch := make(chan *pb_toggles.ListenPayload, 10)
	go func() {
		defer close(ch)
		for {
			res, err := stream.Recv()
			if err != nil {
				return
			}
			ch <- res

		}
	}()
	return ch, nil
}

func withJwtToken(ctx context.Context, token string) context.Context {
	md := meta.ExtractOutgoing(ctx)
	md = md.Set("Authorization", "Bearer "+token)
	return md.ToOutgoing(ctx)
}
