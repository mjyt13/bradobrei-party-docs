package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bradobrei/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestExtractToken(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		token  string
		ok     bool
	}{
		{name: "empty", input: "", token: "", ok: false},
		{name: "bearer token", input: "Bearer abc.def", token: "abc.def", ok: true},
		{name: "lower bearer token", input: "bearer abc.def", token: "abc.def", ok: true},
		{name: "raw token", input: "abc.def", token: "abc.def", ok: true},
		{name: "bearer without token", input: "Bearer   ", token: "", ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, ok := ExtractToken(tt.input)
			if token != tt.token || ok != tt.ok {
				t.Fatalf("expected (%q, %v), got (%q, %v)", tt.token, tt.ok, token, ok)
			}
		})
	}
}

func TestNormalizeRole(t *testing.T) {
	if got := normalizeRole("ADMIN"); got != models.RoleAdmin {
		t.Fatalf("expected legacy ADMIN to normalize to %s, got %s", models.RoleAdmin, got)
	}

	if got := normalizeRole(models.RoleHR); got != models.RoleHR {
		t.Fatalf("expected role to stay unchanged, got %s", got)
	}
}

func TestAuthRequired(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secret := "unit-test-secret"
	validToken := signTestToken(t, secret, 42, models.RoleAdmin)

	router := gin.New()
	router.GET("/protected", AuthRequired(secret), func(c *gin.Context) {
		claims, ok := GetCurrentClaims(c)
		if !ok {
			t.Fatal("expected claims in context")
		}
		c.JSON(http.StatusOK, gin.H{
			"user_id": claims.UserID,
			"role":    claims.Role,
		})
	})

	t.Run("missing token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
	})

	t.Run("raw token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", validToken)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("bearer token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid.token")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
	})
}

func signTestToken(t *testing.T, secret string, userID uint, role models.UserRole) string {
	t.Helper()

	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return signed
}
