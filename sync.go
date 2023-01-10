package featureguards

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (cw *clientWrapper) listenLoop(bgCtx context.Context) {
	for {
		err := cw.listen(bgCtx)
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				// Something odd happened
				log.Errorf("Unexpected error (%s). Sleeping a bit.\n", err)
				// Sleep a bit
				time.Sleep(3 * time.Second)
			}
			if st.Code() == codes.PermissionDenied {
				// Need to re-authenticate
				if err := cw.refreshTokens(bgCtx); err != nil {
					// Something odd happened
					log.Errorf("Unexpected error (%s). Sleeping a bit.\n", err)
					// Sleep a bit
					time.Sleep(3 * time.Second)
				}
			}
		}
		select {
		case <-bgCtx.Done():
			return
		default:
			// Continue
		}
	}
}

func (cw *clientWrapper) listen(bgCtx context.Context) error {
	ctx, cancel, err := cw.listenCtx(bgCtx)
	if err != nil {
		return err
	}
	defer cancel()
	cw.mu.RLock()
	accessToken := cw.accessToken
	cw.mu.RUnlock()
	ch, err := cw.client.Listen(ctx, accessToken, atomic.LoadInt64(&cw.ftVersion), atomic.LoadInt64(&cw.dsVersion))
	if err != nil {
		return err
	}
	for payload := range ch {
		cw.ft.process(payload.FeatureToggles, payload.Version)
		atomic.StoreInt64(&cw.ftVersion, payload.Version)

		// TODO: process dynamic settings
	}
	return nil
}

func (ft *clientWrapper) listenCtx(ctx context.Context) (context.Context, context.CancelFunc, error) {
	ft.mu.RLock()
	accessToken := ft.accessToken
	ft.mu.RUnlock()
	token, err := parseToken(accessToken)
	if err != nil {
		return nil, nil, err
	}
	newCtx, cancel := context.WithDeadline(ctx, token.Expiration())
	return newCtx, cancel, nil
}

func (ft *clientWrapper) refreshTokens(bgCtx context.Context) error {
	ctx, cancel := context.WithTimeout(bgCtx, time.Second)
	defer cancel()

	ft.mu.RLock()
	refreshToken := ft.refreshToken
	ft.mu.RUnlock()
	accessToken, refreshToken, err := ft.client.RefreshAndAuth(ctx, refreshToken)
	if err != nil {
		return err
	}

	ft.mu.Lock()
	defer ft.mu.Unlock()
	ft.accessToken = accessToken
	ft.refreshToken = refreshToken
	return nil
}

func parseToken(token string) (jwt.Token, error) {
	return jwt.Parse([]byte(token), jwt.WithVerify(false))
}
