package utils

import (
	"context"
	"os"

	"google.golang.org/api/idtoken"
)

func VerifyGoogleIDToken(ctx context.Context, token string) (*idtoken.Payload, error) {
    clientID := os.Getenv("GOOGLE_CLIENT_ID")
    payload, err := idtoken.Validate(ctx, token, clientID)
    if err != nil {
        return nil, err
    }
    return payload, nil
}