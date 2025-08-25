package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)

	//r.Use(authMiddleware)

	r.HandleFunc("/security/authenticate", Authenticate) // ----> To request all groceries
	r.HandleFunc("/security/account", Account)
	r.HandleFunc("/security/role", Role)
	r.HandleFunc("/security/validation", Validation)
	r.HandleFunc("/system/init/{key}", Initialize) // ----> To request all groceries
	r.HandleFunc("/system/member", Member)
	r.HandleFunc("/system/project", Project)
	r.HandleFunc("/system/zombie", Zombie)

	r.HandleFunc("/security/ticket", Ticket)

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			// "Access-Control-Expose-Headers",
		}),
		handlers.MaxAge(3600),

		// handlers.ExposedHeaders([]string{"userId", "UserId"}),
	)

	handler := corsMiddleware(r)

	log.Fatal(http.ListenAndServe(":8080", handler))
	//http.Handle("/", r)
}

// func authMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Println(r.Method, r.URL.Path)

// 		//token := r.Header.Get("authentication") // authorization
// 		w.Header().Set("Content-Type", "application/json")

// 		next.ServeHTTP(w, r)
// 	})
// }
