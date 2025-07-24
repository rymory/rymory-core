module main

replace account => ./../packages/security/account/module

replace authenticate => ./../packages/security/authenticate/module

replace role => ./../packages/security/role/module

replace validation => ./../packages/security/validation/module

replace initialize => ./../packages/system/init/module

replace member => ./../packages/system/member/module

replace project => ./../packages/system/project/module

replace zombie => ./../packages/system/zombie/module

go 1.23.5

require role v0.0.0-00010101000000-000000000000

require github.com/felixge/httpsnoop v1.0.3 // indirect

require (
	account v0.0.0-00010101000000-000000000000 // indirect
	authenticate v0.0.0-00010101000000-000000000000 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/jinzhu/gorm v1.9.16 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	gitlab.com/onxorg/goutils/api v0.0.0-20241123105102-cf00b6958c18 // indirect
	gitlab.com/onxorg/goutils/db v0.0.0-20241123105102-cf00b6958c18 // indirect
	golang.org/x/crypto v0.30.0 // indirect
	initialize v0.0.0-00010101000000-000000000000 // indirect
	member v0.0.0-00010101000000-000000000000 // indirect
	project v0.0.0-00010101000000-000000000000 // indirect
	validation v0.0.0-00010101000000-000000000000 // indirect
	zombie v0.0.0-00010101000000-000000000000 // indirect
)
