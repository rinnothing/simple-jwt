package jwt

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Header struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type Payload struct {
	UUID string `json:"uuid"`
}

type Signature string

type PreAccessToken struct {
	Header  Header
	Payload Payload
}

func (p PreAccessToken) Encode(key string) AccessToken {
	headerJSON, err := json.Marshal(p.Header)
	if err != nil {
		panic(err)
	}

	payloadJSON, err := json.Marshal(p.Payload)
	if err != nil {
		panic(err)
	}

	unsigned := fmt.Sprintf("%s.%s",
		base64.URLEncoding.EncodeToString(headerJSON),
		base64.URLEncoding.EncodeToString(payloadJSON),
	)

	signature := hmac.New(sha512.New, []byte(key)).Sum([]byte(unsigned))

	return AccessToken(fmt.Sprintf("%s.%s", unsigned, base64.URLEncoding.EncodeToString(signature)))
}

type AccessToken string

func (a AccessToken) Validate(key string) bool {
	if strings.Count(string(a), ".") != 2 {
		return false
	}

	signature := strings.Split(string(a), ".")[2]
	body := a[:len(a)-len(signature)-1]

	expectedSignature := base64.URLEncoding.EncodeToString(hmac.New(sha512.New, []byte(key)).Sum([]byte(body)))
	return expectedSignature == signature
}

func (a AccessToken) GetPayload() (*Payload, error) {
	parts := strings.Split(string(a), ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token")
	}
	middle := parts[1]

	decoded, err := base64.URLEncoding.DecodeString(middle)
	if err != nil {
		return nil, fmt.Errorf("middle isn't in base64url: %w", err)
	}

	var payload Payload
	err = json.Unmarshal(decoded, &payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload json: %w", err)
	}

	return &payload, nil
}

// TODO: dont't forget to explain reasons on your hash token

type PreRefreshToken struct {
	Access AccessToken
}

func (p PreRefreshToken) Encode(refreshKey, refreshHashKey string) RefreshToken {
	unsigned := base64.URLEncoding.EncodeToString(hmac.New(sha512.New, []byte(refreshKey)).Sum([]byte(p.Access)))

	signature := base64.URLEncoding.EncodeToString(hmac.New(sha512.New, []byte(refreshHashKey)).Sum([]byte(unsigned)))

	return RefreshToken(base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s.%s", unsigned, signature))))
}

type RefreshToken string

func (r RefreshToken) Validate(access AccessToken, refreshKey, refreshHashKey string) bool {
	repeatRefresh := PreRefreshToken{Access: access}.Encode(refreshKey, refreshHashKey)

	return r == repeatRefresh
}
