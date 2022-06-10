package client

import (
	"context"

	"github.com/featureguards/featureguards-go/v1/internal/meta"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb_auth "github.com/featureguards/featureguards-go/v1/proto/auth"
)

const (
	apiKeyMD = "x-api-key"
)

// authenticate authenticates the client with the given apiKey and returns the accessToken, which is
// short lived and will be used for almost all API calls and refreshToken which is used for getting
// new access and refresh tokens.
func (c *Client) Authenticate(ctx context.Context) (accessToken string, refreshToken string, err error) {
	res, err := c.authClient.Authenticate(withApiKey(ctx, c.apiKey), &pb_auth.AuthenticateRequest{Version: version})
	if err != nil {
		return "", "", err
	}

	return res.AccessToken, res.RefreshToken, nil
}

func (c *Client) Refresh(ctx context.Context, token string) (accesToken string, refreshToken string, err error) {
	res, err := c.authClient.Refresh(ctx, &pb_auth.RefreshRequest{RefreshToken: token})
	if err != nil {
		return "", "", err
	}

	return res.AccessToken, res.RefreshToken, nil
}

// refreshAndAuth uses the given refresh token to get new access and refresh token pairs. Each refresh token
// can be used only once. This is because the server implements refresh token rotation to detect token
// reuse. See https://auth0.com/docs/secure/tokens/refresh-tokens/refresh-token-rotation#automatic-reuse-detection
// for more information.
// Upon failure, the client will need to re-authenticate with the server. The given apiKey will be used
// to re-authenticate. If that also fails, then an error is returned.
func (c *Client) RefreshAndAuth(ctx context.Context, token string) (accessToken string, refreshToken string, err error) {
	accessToken, refreshToken, err = c.Refresh(ctx, token)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return "", "", err
		}
		if st.Code() == codes.PermissionDenied {
			return c.Authenticate(ctx)
		}
		return "", "", err
	}

	return
}

func withApiKey(ctx context.Context, key string) context.Context {
	md := meta.ExtractOutgoing(ctx)
	md = md.Set(apiKeyMD, key)
	return md.ToOutgoing(ctx)

}
