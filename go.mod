// Note: Go directive bumped to 1.25.0 (from plan's 1.21) due to
// gorm.io/driver/postgres v1.6.3 and pgx/v5 v5.10.0 requiring go 1.25.0
module github.com/akordium-id/waqfwise

go 1.25.0

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/mux v1.8.1
	github.com/lib/pq v1.10.9
	github.com/pquerna/otp v1.4.0
	github.com/prometheus/client_golang v1.17.0
	github.com/rs/cors v1.10.1
	golang.org/x/crypto v0.16.0
)

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.10.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	gorm.io/driver/postgres v1.6.3 // indirect
	gorm.io/gorm v1.31.2 // indirect
)
