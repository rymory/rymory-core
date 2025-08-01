module authenticate

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/google/uuid v1.6.0
	github.com/jinzhu/gorm v1.9.16
	github.com/lemoras/goutils/api v0.0.0-20250801074636-babf72f0ff35
	github.com/lemoras/goutils/db v0.0.0-20250801074636-babf72f0ff35
	golang.org/x/crypto v0.29.0
)

require github.com/joho/godotenv v1.5.1 // indirect

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
)

go 1.20
