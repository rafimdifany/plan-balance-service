package security

import (
	"context"

	"google.golang.org/api/idtoken"
)

func VerifyGoogleToken(ctx context.Context, tokenString string, clientID string) (*idtoken.Payload, error) {
	payload, err := idtoken.Validate(ctx, tokenString, clientID)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
