module action

replace validation => ./module

require (
	github.com/lemoras/goutils/api v0.0.0-20250803100205-481cd7ccb67e
	validation v0.0.0-00010101000000-000000000000
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/google/uuid v1.6.0 // indirect
)

go 1.20
