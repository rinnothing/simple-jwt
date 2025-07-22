package jwt_test

import (
	"math/rand"
	"testing"

	"github.com/rinnothing/simple-jwt/utils/jwt"
	"github.com/stretchr/testify/require"
)

var (
	accessKey      = "a-string-secret-at-least-512-bits-long-a-string-secret-at-least-"
	refreshKey     = []byte{104, 140, 116, 191, 173, 106, 188, 170, 227, 80, 174, 1, 170, 28, 54, 58, 57, 229, 182, 253, 164, 252, 32, 121, 213, 250, 154, 134, 106, 96, 43, 111, 112, 222, 71, 171, 184, 114, 146, 210, 229, 149, 115, 203, 10, 52, 98, 198, 251, 191, 112, 172, 18, 141, 238, 5, 173, 1, 63, 169, 28, 226, 232, 34}
	refreshHashKey = []byte{155, 195, 6, 73, 11, 176, 25, 209, 85, 153, 65, 246, 109, 4, 138, 98, 141, 149, 89, 58, 11, 253, 244, 251, 145, 241, 144, 215, 166, 236, 168, 26, 177, 145, 146, 96, 196, 107, 54, 225, 83, 72, 16, 49, 93, 193, 143, 176, 189, 84, 45, 193, 157, 154, 208, 199, 78, 177, 41, 74, 160, 207, 0, 19}

	rightToken jwt.AccessToken = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJyYW5kb21fdmFsdWUiOiI1NDMyMSIsInV1aWQiOiIxMjM0NSJ9.nUE1A10XUytenD9d5q_ZPDJwe5n1d-S1qnU6MlZnxmZRL5OpMP7bPDSXYKW5Q80stKmtNP0YNkpMLNZLjVSHAg"
)

func TestJWT(t *testing.T) {
	tool := jwt.NewJWTTool(accessKey, string(refreshKey), string(refreshHashKey))
	tool.RandomString = "54321"

	uuid := "12345"
	access, refresh := tool.IssueTokens(uuid)
	require.Equal(t, rightToken, access)

	payload, err := access.GetPayload()
	require.NoError(t, err)
	require.Equal(t, uuid, payload.UUID)

	require.True(t, tool.CheckAccess(access))
	require.True(t, tool.CheckRefresh(access, refresh))
	require.Equal(t, tool.AccessToRefresh(access), refresh)

	refresh = jwt.RefreshToken(brakeOneChar(string(refresh)))
	require.False(t, refresh.Validate(access, string(refreshKey), string(refreshHashKey)))

	access = jwt.AccessToken(brakeOneChar(string(access)))
	require.False(t, access.Validate(accessKey))
}

func brakeOneChar(s string) string {
	n := rand.Int() % len(s)
	accessBytes := []byte(s)
	accessBytes[n] = 255 - accessBytes[n]

	return string(accessBytes)
}
