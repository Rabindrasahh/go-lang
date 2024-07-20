// route/router.go
package route

import (
	"net/http"
	"rest-api/controller"
)

func RegisterRoutes(uc *controller.UserController) {
	http.HandleFunc("/users", uc.GetUserHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
}
