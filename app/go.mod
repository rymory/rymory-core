module main

replace account => ./../packages/security/account/module

replace authenticate => ./../packages/security/authenticate/module

replace role => ./../packages/security/role/module

replace validation => ./../packages/security/validation/module

replace initialize => ./../packages/system/init/module

replace member => ./../packages/system/member/module

replace project => ./../packages/system/project/module

replace zombie => ./../packages/system/zombie/module

replace ticket => ./../packages/security/ticket/module

go 1.25.0

require (
	github.com/lemoras/goutils/api v1.0.3
	role v0.0.0-00010101000000-000000000000
)

require (
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/lemoras/goutils/db v1.0.1 // indirect
)

require (
	account v0.0.0-00010101000000-000000000000
	authenticate v0.0.0-00010101000000-000000000000
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1
	github.com/jinzhu/gorm v1.9.16 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	initialize v0.0.0-00010101000000-000000000000
	member v0.0.0-00010101000000-000000000000
	project v0.0.0-00010101000000-000000000000
	ticket v0.0.0-00010101000000-000000000000
	validation v0.0.0-00010101000000-000000000000
	zombie v0.0.0-00010101000000-000000000000
)
