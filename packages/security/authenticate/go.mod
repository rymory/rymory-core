module action

replace authenticate => ./module

require (
	authenticate v0.0.0-00010101000000-000000000000
	gitlab.com/onxorg/goutils/api v0.0.0-20241123105102-cf00b6958c18
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/gorm v1.9.16 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	gitlab.com/onxorg/goutils/db v0.0.0-20241123105102-cf00b6958c18 // indirect
	golang.org/x/crypto v0.29.0 // indirect
)

go 1.20
