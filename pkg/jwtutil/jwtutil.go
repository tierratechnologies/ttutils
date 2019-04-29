package jwtutil

import (
	"context"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"
)

var (
	TokenCtxKey = &contextKey{"Token"}
	ErrorCtxKey = &contextKey{"Error"}
)

func FromContext(ctx context.Context) (*jwt.Token, jwt.MapClaims, error) {
	token, _ := ctx.Value(TokenCtxKey).(*jwt.Token)

	var claims jwt.MapClaims
	if token != nil {
		if tokenClaims, ok := token.Claims.(jwt.MapClaims); ok {
			claims = tokenClaims
		} else {
			panic(fmt.Sprintf("jwtauth: unknown type of Claims: %T", token.Claims))
		}
	} else {
		claims = jwt.MapClaims{}
	}

	err, _ := ctx.Value(ErrorCtxKey).(error)

	return token, claims, err
}

// UnixTime returns the given time in UTC milliseconds
func UnixTime(tm time.Time) int64 {
	return tm.UTC().Unix()
}

// EpochNow is a helper function that returns the NumericDate time value used by the spec
func EpochNow() int64 {
	return time.Now().UTC().Unix()
}

// ExpireIn is a helper function to return calculated time in the future for "exp" claim
func ExpireIn(tm time.Duration) int64 {
	return EpochNow() + int64(tm.Seconds())
}

// Set issued at ("iat") to specified time in the claims
func SetIssuedAt(claims jwt.MapClaims, tm time.Time) {
	claims["iat"] = tm.UTC().Unix()
}

// Set issued at ("iat") to present time in the claims
func SetIssuedNow(claims jwt.MapClaims) {
	claims["iat"] = EpochNow()
}

// Set expiry ("exp") in the claims
func SetExpiry(claims jwt.MapClaims, tm time.Time) {
	claims["exp"] = tm.UTC().Unix()
}

// Set expiry ("exp") in the claims to some duration from the present time
func SetExpiryIn(claims jwt.MapClaims, tm time.Duration) {
	claims["exp"] = ExpireIn(tm)
}

// TokenFromCookie tries to retreive the token string from a cookie named
// "jwt".
func TokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return ""
	}
	return cookie.Value
}

// TokenFromHeader tries to retreive the token string from the
// "Authorization" reqeust header: "Authorization: BEARER T".
func TokenFromHeader(r *http.Request) string {
	// Get token from authorization header.
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

// TokenFromQuery tries to retreive the token string from the "jwt" URI
// query parameter.
func TokenFromQuery(r *http.Request) string {
	// Get token from query param named "jwt".
	return r.URL.Query().Get("jwt")
}

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "jwtauth context value " + k.name
}
