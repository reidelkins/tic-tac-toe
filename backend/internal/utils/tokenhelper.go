package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// GetToken returns the token from the Authorization header
func GetToken(w http.ResponseWriter, r *http.Request, production bool) (string, error) {
    // Declare tokenCookie variable
    var tokenCookie *http.Cookie
    var token string
    if production {
        var err error
        tokenCookie, err = r.Cookie("token")    
        if err != nil {            
            // Handle the case when the cookie is not found            
            http.Error(w, "Token cookie not found", http.StatusUnauthorized)
            return "", err
        }
        token = tokenCookie.Value
    } else {
        fmt.Println("DEV MODE")
        body, err := io.ReadAll(r.Body)
        if err != nil {
            fmt.Println("Failed to read request body", err)
            http.Error(w, "Failed to read request body", http.StatusBadRequest)
            return "", err
        }
        defer r.Body.Close()
        fmt.Println("Request body:", string(body))
        // Parse the request body
        var requestBody struct {
            Token string `json:"token"`
        }
        err = json.Unmarshal(body, &requestBody)
        if err != nil {
            fmt.Println("Failed to unmarshal request body", err)
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return "", err
        }
        token = requestBody.Token
        
    }
    return token, nil
}