package handlers

import (
	"fmt"
	"net/http"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "hello")

}
