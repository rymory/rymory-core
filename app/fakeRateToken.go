package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	u "github.com/lemoras/goutils/api"
)

func RateTokenCheck(rateToken string) (bool, map[string]interface{}) {
	now := time.Now().Unix()

	var rl RateLimitToken
	var valid bool

	if rateToken == "" {
		rl = RateLimitToken{
			Remaining: maxRequests,
			ResetAt:   now + windowSeconds,
		}
	} else {
		rl, valid = verifyToken(rateToken)
		if !valid {
			return false, u.Message(false, "Invalid rate token")
		}

		if now > rl.ResetAt {
			rl.Remaining = maxRequests
			rl.ResetAt = now + windowSeconds
		}

		if rl.Remaining <= 0 {
			return false, u.Message(false, fmt.Sprintf("Rate limit exceeded. Retry-After ", strconv.FormatInt(rl.ResetAt-now, 10)))
		}

		rl.Remaining--
	}

	newToken, _ := generateToken(rl.Remaining, rl.ResetAt)

	resp := u.Message(true, "Generated Rate Token by Successful")
	resp["rateToken"] = newToken

	return true, resp
}

func RateTokenhandler(w http.ResponseWriter, r *http.Request) bool {
	token := r.Header.Get("X-RateToken")
	now := time.Now().Unix()

	var rl RateLimitToken
	var valid bool

	if token == "" {
		rl = RateLimitToken{
			Remaining: maxRequests,
			ResetAt:   now + windowSeconds,
		}
	} else {
		rl, valid = verifyToken(token)
		if !valid {
			http.Error(w, "Invalid rate token", http.StatusUnauthorized)
			return false
		}

		if now > rl.ResetAt {
			rl.Remaining = maxRequests
			rl.ResetAt = now + windowSeconds
		}

		if rl.Remaining <= 0 {
			w.Header().Set("Retry-After", strconv.FormatInt(rl.ResetAt-now, 10))
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return false
		}

		rl.Remaining--
	}

	newToken, _ := generateToken(rl.Remaining, rl.ResetAt)
	w.Header().Set("X-RateToken", newToken)
	w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(rl.Remaining))
	w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte(`{"message": "OK"}`))

	return true
}

func computeHMAC(data []byte) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(data)
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func verifyToken(token string) (RateLimitToken, bool) {
	decoded, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return RateLimitToken{}, false
	}

	parts := strings.SplitN(string(decoded), ".", 2)
	if len(parts) != 2 {
		return RateLimitToken{}, false
	}

	payloadStr := parts[0]
	sig := parts[1]

	if computeHMAC([]byte(payloadStr)) != sig {
		return RateLimitToken{}, false
	}

	var payload RateLimitToken
	err = json.Unmarshal([]byte(payloadStr), &payload)
	if err != nil {
		return RateLimitToken{}, false
	}

	return payload, true
}

func generateToken(remaining int, resetAt int64) (string, error) {
	payload := RateLimitToken{
		Remaining: remaining,
		ResetAt:   resetAt,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	sig := computeHMAC(jsonBytes)
	full := append(jsonBytes, []byte("."+sig)...)

	return base64.URLEncoding.EncodeToString(full), nil
}

const (
	secretKey     = "rate_secret_key"
	maxRequests   = 5
	windowSeconds = 60
)

type RateLimitToken struct {
	Remaining int   `json:"remaining"`
	ResetAt   int64 `json:"resetAt"`
}
