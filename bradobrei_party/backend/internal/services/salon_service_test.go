package services

import (
	"testing"

	"bradobrei/backend/internal/models"
)

func TestNormalizePoint(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "lat lon pair", input: "58.0141, 56.2230", want: "POINT(56.223 58.0141)"},
		{name: "semicolon pair", input: "58.0141;56.2230", want: "POINT(56.223 58.0141)"},
		{name: "wkt point", input: "POINT(56.2230 58.0141)", want: "POINT(56.223 58.0141)"},
		{name: "wkt point with comma inside", input: "POINT(58.0141, 56.2230)", want: "POINT(56.223 58.0141)"},
		{name: "invalid latitude", input: "120, 56.2230", wantErr: true},
		{name: "invalid format", input: "POINT(56.2230)", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizePoint(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestNormalizeSalonLocation(t *testing.T) {
	t.Run("empty string becomes nil", func(t *testing.T) {
		raw := "   "
		salon := &models.Salon{Location: &raw}

		if err := normalizeSalonLocation(salon); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if salon.Location != nil {
			t.Fatal("expected nil location after normalization")
		}
	})

	t.Run("coordinates become wkt", func(t *testing.T) {
		raw := "58.0141, 56.2230"
		salon := &models.Salon{Location: &raw}

		if err := normalizeSalonLocation(salon); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if salon.Location == nil || *salon.Location != "POINT(56.223 58.0141)" {
			t.Fatalf("unexpected normalized location: %#v", salon.Location)
		}
	})
}

func TestValidateLatLon(t *testing.T) {
	if err := validateLatLon(58.0141, 56.2230); err != nil {
		t.Fatalf("expected valid coords, got error: %v", err)
	}

	if err := validateLatLon(-91, 56.2230); err == nil {
		t.Fatal("expected error for invalid latitude")
	}

	if err := validateLatLon(58.0141, 181); err == nil {
		t.Fatal("expected error for invalid longitude")
	}
}
