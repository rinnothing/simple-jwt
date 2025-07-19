package jwt

type Tool struct {
	accessKey      string
	refreshKey     string
	refreshHashKey string
}

func NewJWTTool(accessKey, refreshKey, refreshHashKey string) *Tool {
	return &Tool{
		accessKey:      accessKey,
		refreshKey:     refreshKey,
		refreshHashKey: refreshHashKey,
	}
}

func (t *Tool) IssueTokens(uuid string) (AccessToken, RefreshToken) {
	preAccess := PreAccessToken{
		Header: Header{
			Algorithm: "HS512",
			Type:      "JWT",
		},
		Payload: Payload{
			UUID: uuid,
		},
	}
	access := preAccess.Encode(t.accessKey)

	preRefresh := PreRefreshToken{
		Access: access,
	}
	refresh := preRefresh.Encode(t.refreshKey, t.refreshHashKey)

	return access, refresh
}

func (t *Tool) AccessToRefresh(access AccessToken) RefreshToken {
	preRefresh := PreRefreshToken{
		Access: access,
	}
	return preRefresh.Encode(t.refreshKey, t.refreshHashKey)
}

func (t *Tool) CheckAccess(access AccessToken) bool {
	return access.Validate(t.accessKey)
}

// remember: the refresh key could be already used, the method only checks for access and refresh tokens compatibility
func (t *Tool) ChechRefresh(access AccessToken, refresh RefreshToken) bool {
	return refresh.Validate(access, t.refreshKey, t.refreshHashKey)
}
