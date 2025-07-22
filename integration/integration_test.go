package integration

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/rinnothing/simple-jwt/internal/api"
	"github.com/rinnothing/simple-jwt/internal/api/schema"
	"github.com/rinnothing/simple-jwt/internal/config"
	"github.com/stretchr/testify/require"

	"go.uber.org/zap"
)

func TestIntegration(t *testing.T) {
	address := "http://localhost:9090"

	cfg, err := config.GetConfig("../config/config.yaml")
	require.NoError(t, err)
	cfg.Webhook.HttpAddress = address
	cfg.Logger.Env = "dev"

	loggerCfg, err := config.ConfigureLogger(cfg.Logger)
	require.NoError(t, err)

	logger, err := loggerCfg.Build()
	require.NoError(t, err)

	server := api.Server{}
	go func() {
		if err := server.Run(cfg, logger); err != nil {
			logger.Fatal("server stopped with error", zap.Error(err))
		}
	}()

	time.Sleep(5 * time.Second)

	// started server
	client, err := schema.NewClientWithResponses(fmt.Sprintf("http://localhost:%s", cfg.Port))
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// put guid
	guid := "111111"
	authResp, err := client.AuthorizeGUIDWithResponse(ctx, guid)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, authResp.StatusCode())

	tokens := *authResp.JSON201

	// get it back
	guidResp, err := client.GetGUIDWithResponse(ctx, &schema.GetGUIDParams{AccessToken: *tokens.AccessToken})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, guidResp.StatusCode())

	require.Equal(t, guid, string(*guidResp.JSON200))

	// do refresh
	refreshResp, err := client.RefreshTokensWithResponse(ctx, tokens)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, refreshResp.StatusCode())

	// not allowed

	guidResp, err = client.GetGUIDWithResponse(ctx, &schema.GetGUIDParams{AccessToken: *tokens.AccessToken})
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, guidResp.StatusCode())

	// allowed

	tokens = *refreshResp.JSON200

	guidResp, err = client.GetGUIDWithResponse(ctx, &schema.GetGUIDParams{AccessToken: *tokens.AccessToken})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, guidResp.StatusCode())

	require.Equal(t, guid, string(*guidResp.JSON200))

	// unauthorize

	unResp, err := client.UnauthorizeWithResponse(ctx, &schema.UnauthorizeParams{AccessToken: *tokens.AccessToken})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, unResp.StatusCode())

	// not allowed

	guidResp, err = client.GetGUIDWithResponse(ctx, &schema.GetGUIDParams{AccessToken: *tokens.AccessToken})
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, guidResp.StatusCode())

	// stopped server

	server.Stop()
}
