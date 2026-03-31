package handlers

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestParsePeriod(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("explicit period", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		req := httptest.NewRequest("GET", "/reports/salon-activity?from=2026-03-01&to=2026-03-31", nil)
		c.Request = req

		from, to, err := parsePeriod(c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if from.Format("2006-01-02") != "2026-03-01" || to.Format("2006-01-02") != "2026-03-31" {
			t.Fatalf("unexpected range: %s - %s", from, to)
		}
	})

	t.Run("bad from date", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		req := httptest.NewRequest("GET", "/reports/salon-activity?from=bad-date&to=2026-03-31", nil)
		c.Request = req

		_, _, err := parsePeriod(c)
		if err == nil {
			t.Fatal("expected error for invalid date")
		}
	})

	t.Run("default period", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		req := httptest.NewRequest("GET", "/reports/salon-activity", nil)
		c.Request = req

		from, to, err := parsePeriod(c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if to.Before(from) {
			t.Fatalf("expected to >= from, got from=%s to=%s", from, to)
		}

		diff := to.Sub(from)
		if diff < 27*24*time.Hour || diff > 32*24*time.Hour {
			t.Fatalf("expected about one month default range, got %s", diff)
		}
	})
}
