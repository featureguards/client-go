package featureguards

import (
	"context"
	"os"
	"testing"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func Test_refreshTokens(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	addr := os.Getenv("GRPC_ADDR")
	apiKey := os.Getenv("API_KEY")
	fg := New(ctx, withDomain(addr), withoutListen(), withTestCerts(),
		WithApiKey(apiKey), WithDefaults(map[string]bool{"BAR": true}))
	fg.cw.mu.Lock()
	accessToken := fg.cw.accessToken
	refreshToken := fg.cw.refreshToken
	fg.cw.mu.Unlock()
	err := fg.cw.refreshTokens(ctx)
	if err != nil {
		t.Errorf("refreshTokens() error = %v, wantErr nil", err)
	}
	fg.cw.mu.Lock()
	defer fg.cw.mu.Unlock()
	if accessToken == fg.cw.accessToken {
		t.Errorf("accessToken = %v, want != %v", fg.cw.accessToken, accessToken)
	}
	if refreshToken == fg.cw.refreshToken {
		t.Errorf("refreshToken = %v, want != %v", fg.cw.refreshToken, refreshToken)
	}
}

func Test_IsOn(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	addr := os.Getenv("GRPC_ADDR")
	apiKey := os.Getenv("API_KEY")
	fg := New(ctx, withDomain(addr), withoutListen(), withTestCerts(),
		WithApiKey(apiKey), WithDefaults(map[string]bool{"BAR": true}))
	if fg.cw == nil {
		t.Error("feature toggles should not be nil")
	}

	// test listen
	deadlineCtx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()
	if err := fg.cw.listen(deadlineCtx); err != nil {
		t.Errorf("listen() error = %v, wantErr not found", err)
	}

	on, err := fg.IsOn("TEST")
	if err != nil {
		t.Errorf("IsOn() error = %v, wantErr nil", err)
	}
	if !on {
		t.Errorf("on() on = %v, wantOn %v", on, true)
	}

	on, err = fg.IsOn("BAR")
	if err == nil {
		t.Errorf("IsOn() error = %v, wantErr not found", err)
	}
	if !on {
		t.Errorf("on() on = %v, wantOn %v", on, true)
	}
}
