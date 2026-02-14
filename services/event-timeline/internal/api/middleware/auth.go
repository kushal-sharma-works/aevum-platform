package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/api/httputil"
)

const ClaimsContextKey = "jwt_claims"

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			httputil.Unauthorized(c, "missing_token", "missing bearer token")
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, errors.New("unexpected signing algorithm")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			httputil.Unauthorized(c, "invalid_token", "invalid token")
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			httputil.Unauthorized(c, "invalid_claims", "invalid token claims")
			c.Abort()
			return
		}
		if _, ok := claims["iss"]; !ok {
			httputil.Unauthorized(c, "missing_claim", "missing iss claim")
			c.Abort()
			return
		}
		if _, ok := claims["sub"]; !ok {
			httputil.Unauthorized(c, "missing_claim", "missing sub claim")
			c.Abort()
			return
		}
		if _, ok := claims["exp"]; !ok {
			httputil.Unauthorized(c, "missing_claim", "missing exp claim")
			c.Abort()
			return
		}
		if _, ok := claims["iat"]; !ok {
			httputil.Unauthorized(c, "missing_claim", "missing iat claim")
			c.Abort()
			return
		}
		c.Set(ClaimsContextKey, claims)
		c.Next()
	}
}
