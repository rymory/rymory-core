module action

replace authenticate => ./module

require (
	authenticate v0.0.0-00010101000000-000000000000
	github.com/lemoras/goutils/api v1.0.2
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/gorm v1.9.16 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lemoras/goutils/db v1.0.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/crypto v0.33.0 // indirect
// golang.org/x/crypto v0.41.0 // indirect
)

// go 1.25.0
go 1.20
