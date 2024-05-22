package handlers

import (
	"fmt"
	"net/http"

	db "github.com/oskar13/mini-tracker/pkg/db"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

// Handle login and logout URL endpoint
func LoginPage(w http.ResponseWriter, r *http.Request) {

	var pageStruct struct {
		Error         bool
		ErrorText     string
		LoggedIn      bool
		LogoutMessage bool
		UserData      webdata.User
	}

	userData := webutils.GetUserData(r, db.DB)

	if userData.LoggedIn {

		if r.URL.Path == "/logout" {
			//Log user out
			userData = webutils.LogOutUser(w, userData)
			pageStruct.LogoutMessage = true

		} else {
			//Offer log out
			pageStruct.LoggedIn = true
		}

	} else {
		if r.Method == "POST" {
			// Check login details

			if err := r.ParseForm(); err != nil {
				http.Error(w, "Error parsing the form", http.StatusInternalServerError)
				return
			}

			loginData := webdata.LoginData{
				UserNameOrEmail: r.FormValue("username"),
				Password:        r.FormValue("password"),
			}

			err, resultUserData := webutils.LoginUser(w, loginData.UserNameOrEmail, loginData.Password)

			if err != nil {

				pageStruct.Error = true
				pageStruct.ErrorText = "Invalid username or password"

				w.WriteHeader(http.StatusUnauthorized)
				webutils.RenderTemplate(w, []string{"pkg/web/templates/login.html"}, pageStruct)
				return
			}

			fmt.Println(resultUserData)

			http.Redirect(w, r, "/", http.StatusSeeOther)

			return
		}
	}

	webutils.RenderTemplate(w, []string{"pkg/web/templates/login.html"}, pageStruct)

}
