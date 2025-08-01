module action

go 1.20

replace account => ./module

require (
	account v0.0.0-00010101000000-000000000000
	github.com/lemoras/goutils/api v0.0.0-20250801074636-babf72f0ff35
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/gorm v1.9.16 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lemoras/goutils/db v0.0.0-20250801074636-babf72f0ff35 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/crypto v0.30.0 // indirect
)
