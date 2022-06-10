package featureguards

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb_ft "github.com/featureguards/featureguards-go/v2/proto/feature_toggle"
)

func (ft *featureToggles) listenLoop(bgCtx context.Context) {
	for {
		err := ft.listen(bgCtx)
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
				if err := ft.refreshTokens(bgCtx); err != nil {
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

func (ft *featureToggles) listen(bgCtx context.Context) error {
	ctx, cancel, err := ft.listenCtx(bgCtx)
	if err != nil {
		return err
	}
	defer cancel()
	ft.mu.RLock()
	accessToken := ft.accessToken
	ft.mu.RUnlock()
	ch, err := ft.client.Listen(ctx, accessToken, atomic.LoadInt64(&ft.clientVersion))
	if err != nil {
		return err
	}
	for payload := range ch {
		ft.process(payload.FeatureToggles, payload.Version)
	}
	return nil
}

func (ft *featureToggles) listenCtx(ctx context.Context) (context.Context, context.CancelFunc, error) {
	ft.mu.RLock()
	accessToken := ft.accessToken
	ft.mu.RUnlock()
	token, err := ft.parse(accessToken)
	if err != nil {
		return nil, nil, err
	}
	newCtx, cancel := context.WithDeadline(ctx, token.Expiration())
	return newCtx, cancel, nil
}

func (ft *featureToggles) refreshTokens(bgCtx context.Context) error {
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

func (ft *featureToggles) parse(token string) (jwt.Token, error) {
	return jwt.Parse([]byte(token))
}

func (ft *featureToggles) process(fts []*pb_ft.FeatureToggle, version int64) {
	ft.mu.Lock()
	defer ft.mu.Unlock()
	for _, toggle := range fts {
		if toggle.DeletedAt.IsValid() {
			delete(ft.ftByName, toggle.Name)
		} else {
			ft.ftByName[toggle.Name] = toggle
		}
	}
	atomic.StoreInt64(&ft.clientVersion, version)
}
