module action

replace validation => ./module

require (
	gitlab.com/onxorg/goutils/api v0.0.0-20241123105102-cf00b6958c18
	validation v0.0.0-00010101000000-000000000000
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/google/uuid v1.6.0 // indirect
)

go 1.20
