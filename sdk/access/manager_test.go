package access

import (
	"context"
	"net/http"
	"testing"
)

type errorProvider struct {
	err *AuthError
}

func (p errorProvider) Identifier() string {
	return "error-provider"
}

func (p errorProvider) Authenticate(context.Context, *http.Request) (*Result, *AuthError) {
	return nil, p.err
}

func TestManagerAuthenticateFallbackErrors(t *testing.T) {
	tests := []struct {
		name      string
		providers []Provider
		wantCode  AuthErrorCode
	}{
		{
			name:      "not handled only",
			providers: []Provider{errorProvider{err: NewNotHandledError()}},
			wantCode:  AuthErrorCodeNoCredentials,
		},
		{
			name:      "missing credentials",
			providers: []Provider{errorProvider{err: NewNotHandledError()}, errorProvider{err: NewNoCredentialsError()}},
			wantCode:  AuthErrorCodeNoCredentials,
		},
		{
			name:      "invalid wins over missing",
			providers: []Provider{errorProvider{err: NewNoCredentialsError()}, errorProvider{err: NewInvalidCredentialError()}},
			wantCode:  AuthErrorCodeInvalidCredential,
		},
		{
			name:      "unexpected error returns immediately",
			providers: []Provider{errorProvider{err: NewInternalAuthError("boom", nil)}, testProvider{id: "never-reached"}},
			wantCode:  AuthErrorCodeInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewManager()
			manager.SetProviders(tt.providers)

			res, authErr := manager.Authenticate(context.Background(), nil)
			if res != nil {
				t.Fatalf("Authenticate() result = %+v, want nil", res)
			}
			if !IsAuthErrorCode(authErr, tt.wantCode) {
				t.Fatalf("Authenticate() error = %v, want code %s", authErr, tt.wantCode)
			}
		})
	}
}

func TestManagerAuthenticateReturnsFirstSuccess(t *testing.T) {
	manager := NewManager()
	manager.SetProviders([]Provider{errorProvider{err: NewNoCredentialsError()}, testProvider{id: "ok"}})

	res, authErr := manager.Authenticate(context.Background(), nil)
	if authErr != nil {
		t.Fatalf("Authenticate() error = %v", authErr)
	}
	if res == nil || res.Provider != "ok" {
		t.Fatalf("Authenticate() result = %+v, want provider ok", res)
	}
}
