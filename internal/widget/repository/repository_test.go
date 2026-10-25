package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Discard,
	})
	if err != nil {
		t.Fatalf("failed to open gorm: %v", err)
	}
	cleanup := func() { db.Close() }
	return gormDB, mock, cleanup
}

func TestGetMosqueBySlug_Found(t *testing.T) {
	gormDB, mock, cleanup := setupMockDB(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name", "slug", "primary_color", "is_active"}).
		AddRow("550e8400-e29b-41d4-a716-446655440000", "Masjid Al-Ikhlas", "masjid-al-ikhlas", "#1B5E20", true)

	mock.ExpectQuery(`SELECT .* FROM "mosques" WHERE slug = \$1 AND is_active = \$2`).
		WillReturnRows(rows)

	repo := New(gormDB)
	mosque, err := repo.GetMosqueBySlug(context.Background(), "masjid-al-ikhlas")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mosque.Name != "Masjid Al-Ikhlas" {
		t.Errorf("Name = %q, want %q", mosque.Name, "Masjid Al-Ikhlas")
	}
}

func TestGetMosqueBySlug_NotFound(t *testing.T) {
	gormDB, mock, cleanup := setupMockDB(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT .* FROM "mosques"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	repo := New(gormDB)
	_, err := repo.GetMosqueBySlug(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
}

func TestGetCampaignsForWidget_FiltersActiveOnly(t *testing.T) {
	gormDB, mock, cleanup := setupMockDB(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT count\(\*\) FROM "campaigns"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	rows := sqlmock.NewRows([]string{"id", "title", "status", "deleted_at"}).
		AddRow(int64(1), "Campaign 1", "active", nil).
		AddRow(int64(2), "Campaign 2", "active", nil)

	mock.ExpectQuery(`SELECT .* FROM "campaigns" WHERE status = \$1 AND deleted_at IS NULL`).
		WillReturnRows(rows)

	repo := New(gormDB)
	campaigns, total, err := repo.GetCampaignsForWidget(context.Background(), CampaignQueryParams{
		Limit: 6,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}
	if len(campaigns) != 2 {
		t.Errorf("len(campaigns) = %d, want 2", len(campaigns))
	}
}
