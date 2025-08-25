module action

replace ticket => ./module

require (
	github.com/lemoras/goutils/api v1.0.0
	ticket v0.0.0-00010101000000-000000000000
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/google/uuid v1.6.0 // indirect
)

go 1.23.0

toolchain go1.24.4
