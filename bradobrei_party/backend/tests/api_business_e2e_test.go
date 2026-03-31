package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"bradobrei/backend/internal/models"

	"testing"
)

func TestBookingAccessAndFlow(t *testing.T) {
	app := newTestApp(t)
	defer app.close(t)

	adminToken := loginAndGetToken(t, app, "admin", "password")
	_, workerToken := hireEmployeeAndLogin(t, app, adminToken)
	clientToken := registerClientAndLogin(t, app, "client_petr")
	serviceID := seedBookableFixture(t, app)

	createPayload := bookingPayload(app.baseSalonID, serviceID, time.Now().Add(24*time.Hour).UTC())

	status, body := app.request(t, http.MethodPost, "/api/v1/bookings", "", createPayload)
	if status != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token for booking create, got %d: %s", status, string(body))
	}

	status, body = app.request(t, http.MethodPost, "/api/v1/bookings", workerToken, createPayload)
	if status != http.StatusForbidden {
		t.Fatalf("expected 403 for worker on booking create, got %d: %s", status, string(body))
	}

	badPayload := bookingPayload(app.baseSalonID, serviceID, time.Now().Add(25*time.Hour).UTC())
	badPayload["start_time"] = "bad-time"
	status, body = app.request(t, http.MethodPost, "/api/v1/bookings", clientToken, badPayload)
	if status != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid booking payload, got %d: %s", status, string(body))
	}

	bookingID, createBody := createBookingAsClient(t, app, clientToken, createPayload)
	app.saveArtifact(t, "booking_create", createBody)

	status, body = app.request(t, http.MethodGet, "/api/v1/bookings/my", clientToken, nil)
	if status != http.StatusOK {
		t.Fatalf("expected 200 from /bookings/my, got %d: %s", status, string(body))
	}
	app.saveArtifact(t, "booking_get_my", body)

	status, body = app.request(t, http.MethodPost, fmt.Sprintf("/api/v1/bookings/%d/confirm", bookingID), workerToken, nil)
	if status != http.StatusOK {
		t.Fatalf("expected 200 for booking confirm, got %d: %s", status, string(body))
	}
	app.saveArtifact(t, "booking_confirm", body)

	status, body = app.request(t, http.MethodGet, "/api/v1/bookings/master", workerToken, nil)
	if status != http.StatusOK {
		t.Fatalf("expected 200 from /bookings/master, got %d: %s", status, string(body))
	}
	app.saveArtifact(t, "booking_master_list", body)
}

func TestPaymentAccessAndFlow(t *testing.T) {
	app := newTestApp(t)
	defer app.close(t)

	adminToken := loginAndGetToken(t, app, "admin", "password")
	clientToken := registerClientAndLogin(t, app, "client_olga")
	serviceID := seedBookableFixture(t, app)
	bookingID, _ := createBookingAsClient(t, app, clientToken, bookingPayload(app.baseSalonID, serviceID, time.Now().Add(24*time.Hour).UTC()))

	status, body := app.request(t, http.MethodGet, "/api/v1/payments", "", nil)
	if status != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token for payments list, got %d: %s", status, string(body))
	}

	status, body = app.request(t, http.MethodGet, "/api/v1/payments", clientToken, nil)
	if status != http.StatusForbidden {
		t.Fatalf("expected 403 for client on payments list, got %d: %s", status, string(body))
	}

	badPayload := map[string]any{
		"booking_id": 999999,
		"amount":     1200,
		"status":     "SUCCESS",
	}
	status, body = app.request(t, http.MethodPost, "/api/v1/payments", adminToken, badPayload)
	if status != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid payment payload, got %d: %s", status, string(body))
	}

	createPayload := map[string]any{
		"booking_id":               bookingID,
		"amount":                   0,
		"status":                   "SUCCESS",
		"external_transaction_id":  "txn_test_success_001",
	}
	status, body = app.request(t, http.MethodPost, "/api/v1/payments", adminToken, createPayload)
	if status != http.StatusCreated {
		t.Fatalf("expected 201 for payment create, got %d: %s", status, string(body))
	}
	app.saveArtifact(t, "payment_create", body)

	var createdPayment struct {
		ID uint `json:"id"`
	}
	if err := json.Unmarshal(body, &createdPayment); err != nil {
		t.Fatalf("failed to decode payment create response: %v", err)
	}

	status, body = app.request(t, http.MethodGet, fmt.Sprintf("/api/v1/payments/%d", createdPayment.ID), adminToken, nil)
	if status != http.StatusOK {
		t.Fatalf("expected 200 for payment get by id, got %d: %s", status, string(body))
	}
	app.saveArtifact(t, "payment_get", body)
}

func TestReportsAccessAndFlow(t *testing.T) {
	app := newTestApp(t)
	defer app.close(t)

	adminToken := loginAndGetToken(t, app, "admin", "password")
	_, workerToken := hireEmployeeAndLogin(t, app, adminToken)
	clientToken := registerClientAndLogin(t, app, "client_maria")
	serviceID := seedBookableFixture(t, app)

	reportTime := time.Now().Add(48 * time.Hour).UTC()
	bookingID, _ := createBookingAsClient(t, app, clientToken, bookingPayload(app.baseSalonID, serviceID, reportTime))

	status, body := app.request(t, http.MethodGet, "/api/v1/reports/employees", "", nil)
	if status != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token for reports, got %d: %s", status, string(body))
	}

	status, body = app.request(t, http.MethodGet, "/api/v1/reports/employees", clientToken, nil)
	if status != http.StatusForbidden {
		t.Fatalf("expected 403 for client on reports, got %d: %s", status, string(body))
	}

	status, body = app.request(t, http.MethodGet, "/api/v1/reports/salon-activity?from=bad-date&to=2026-12-31", adminToken, nil)
	if status != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid report period, got %d: %s", status, string(body))
	}

	status, body = app.request(t, http.MethodPost, fmt.Sprintf("/api/v1/bookings/%d/confirm", bookingID), workerToken, nil)
	if status != http.StatusOK {
		t.Fatalf("expected 200 for booking confirm before reports, got %d: %s", status, string(body))
	}

	markBookingCompleted(t, app, bookingID)
	reviewID := createReviewAsClient(t, app, clientToken)

	from := reportTime.Add(-24 * time.Hour).Format("2006-01-02")
	to := reportTime.Add(24 * time.Hour).Format("2006-01-02")

	status, body = app.request(t, http.MethodGet, "/api/v1/reports/employees", adminToken, nil)
	if status != http.StatusOK {
		t.Fatalf("expected 200 for employee report, got %d: %s", status, string(body))
	}
	app.saveArtifact(t, "report_employees", body)

	status, body = app.request(t, http.MethodGet, fmt.Sprintf("/api/v1/reports/salon-activity?from=%s&to=%s", from, to), adminToken, nil)
	if status != http.StatusOK {
		t.Fatalf("expected 200 for salon activity report, got %d: %s", status, string(body))
	}
	app.saveArtifact(t, "report_salon_activity", body)

	status, body = app.request(t, http.MethodGet, fmt.Sprintf("/api/v1/reports/master-activity?from=%s&to=%s", from, to), adminToken, nil)
	if status != http.StatusOK {
		t.Fatalf("expected 200 for master activity report, got %d: %s", status, string(body))
	}
	app.saveArtifact(t, "report_master_activity", body)

	status, body = app.request(t, http.MethodGet, "/api/v1/reports/reviews", adminToken, nil)
	if status != http.StatusOK {
		t.Fatalf("expected 200 for reviews report, got %d: %s", status, string(body))
	}
	app.saveArtifact(t, "report_reviews", body)

	status, body = app.request(t, http.MethodGet, fmt.Sprintf("/api/v1/reviews/%d", reviewID), adminToken, nil)
	if status != http.StatusOK {
		t.Fatalf("expected 200 for review get by id, got %d: %s", status, string(body))
	}
	app.saveArtifact(t, "review_get", body)
}

func TestReviewAccessAndFlow(t *testing.T) {
	app := newTestApp(t)
	defer app.close(t)

	adminToken := loginAndGetToken(t, app, "admin", "password")
	clientToken := registerClientAndLogin(t, app, "client_review")
	_, workerToken := hireEmployeeAndLogin(t, app, adminToken)

	status, body := app.request(t, http.MethodPost, "/api/v1/reviews", "", map[string]any{
		"text":   "Отзыв без токена",
		"rating": 5,
	})
	if status != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token for review create, got %d: %s", status, string(body))
	}

	status, body = app.request(t, http.MethodGet, "/api/v1/reviews", clientToken, nil)
	if status != http.StatusForbidden {
		t.Fatalf("expected 403 for client on reviews list, got %d: %s", status, string(body))
	}

	status, body = app.request(t, http.MethodGet, "/api/v1/reviews", workerToken, nil)
	if status != http.StatusForbidden {
		t.Fatalf("expected 403 for worker on reviews list, got %d: %s", status, string(body))
	}

	status, body = app.request(t, http.MethodPost, "/api/v1/reviews", clientToken, map[string]any{
		"text":   "",
		"rating": 10,
	})
	if status != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid review payload, got %d: %s", status, string(body))
	}

	reviewID := createReviewAsClient(t, app, clientToken)

	status, body = app.request(t, http.MethodGet, "/api/v1/reviews", adminToken, nil)
	if status != http.StatusOK {
		t.Fatalf("expected 200 for admin reviews list, got %d: %s", status, string(body))
	}
	app.saveArtifact(t, "review_list", body)

	status, body = app.request(t, http.MethodGet, fmt.Sprintf("/api/v1/reviews/%d", reviewID), adminToken, nil)
	if status != http.StatusOK {
		t.Fatalf("expected 200 for admin review get by id, got %d: %s", status, string(body))
	}
	app.saveArtifact(t, "review_get_by_id", body)
}

func registerClientAndLogin(t *testing.T, app *testApp, username string) string {
	t.Helper()

	status, body := app.request(t, http.MethodPost, "/api/v1/auth/register", "", map[string]any{
		"username":  username,
		"password":  "password123",
		"full_name": "Test Client",
		"phone":     "+79990002233",
		"email":     username + "@example.com",
		"role":      "CLIENT",
	})
	if status != http.StatusCreated {
		t.Fatalf("expected client register 201, got %d: %s", status, string(body))
	}

	return loginAndGetToken(t, app, username, "password123")
}

func seedBookableFixture(t *testing.T, app *testApp) uint {
	t.Helper()

	suffix := time.Now().UnixNano()
	material := models.Material{
		Name: fmt.Sprintf("Material-%d", suffix),
		Unit: "ml",
	}
	if err := app.db.Create(&material).Error; err != nil {
		t.Fatalf("failed to seed material: %v", err)
	}

	service := models.Service{
		Name:            fmt.Sprintf("Service-%d", suffix),
		Description:     "Test service for API e2e",
		Price:           1800,
		DurationMinutes: 75,
	}
	if err := app.db.Create(&service).Error; err != nil {
		t.Fatalf("failed to seed service: %v", err)
	}

	serviceMaterial := models.ServiceMaterial{
		ServiceID:      service.ID,
		MaterialID:     material.ID,
		QuantityPerUse: 1,
	}
	if err := app.db.Create(&serviceMaterial).Error; err != nil {
		t.Fatalf("failed to seed service materials: %v", err)
	}

	inventory := models.Inventory{
		SalonID:     app.baseSalonID,
		MaterialID:  material.ID,
		Quantity:    20,
		LastUpdated: time.Now(),
	}
	if err := app.db.Create(&inventory).Error; err != nil {
		t.Fatalf("failed to seed inventory: %v", err)
	}

	return service.ID
}

func bookingPayload(salonID, serviceID uint, start time.Time) map[string]any {
	return map[string]any{
		"start_time":  start.Format(time.RFC3339),
		"salon_id":    salonID,
		"service_ids": []uint{serviceID},
		"notes":       "E2E booking scenario",
	}
}

func createBookingAsClient(t *testing.T, app *testApp, clientToken string, payload map[string]any) (uint, []byte) {
	t.Helper()

	status, body := app.request(t, http.MethodPost, "/api/v1/bookings", clientToken, payload)
	if status != http.StatusCreated {
		t.Fatalf("expected booking create 201, got %d: %s", status, string(body))
	}

	var resp struct {
		ID uint `json:"id"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("failed to decode booking response: %v", err)
	}

	return resp.ID, body
}

func markBookingCompleted(t *testing.T, app *testApp, bookingID uint) {
	t.Helper()

	if err := app.db.Model(&models.Booking{}).
		Where("id = ?", bookingID).
		Update("status", models.BookingCompleted).Error; err != nil {
		t.Fatalf("failed to mark booking as completed: %v", err)
	}
}

func createReviewAsClient(t *testing.T, app *testApp, clientToken string) uint {
	t.Helper()

	status, body := app.request(t, http.MethodPost, "/api/v1/reviews", clientToken, map[string]any{
		"text":   "Отличный тестовый отзыв",
		"rating": 5,
	})
	if status != http.StatusCreated {
		t.Fatalf("expected review create 201, got %d: %s", status, string(body))
	}

	var resp struct {
		ID uint `json:"id"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("failed to decode review create response: %v", err)
	}

	app.saveArtifact(t, "review_create", body)
	return resp.ID
}
