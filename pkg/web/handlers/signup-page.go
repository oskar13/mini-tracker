package handlers

import (
	"net/http"

	db "github.com/oskar13/mini-tracker/pkg/db"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

// Handle login and logout URL endpoint
func SignupPage(w http.ResponseWriter, r *http.Request) {

	var pageStruct struct {
		Error         bool
		ErrorText     string
		RefCodeValid  bool
		RefCode       string
		LoggedIn      bool
		LogoutMessage bool
		UserData      webdata.User
	}

	userData := webutils.GetUserData(r, db.DB)

	if userData.LoggedIn {
		//Show user a message about being logged in, aslo remind that one account per lifetime

	} else {
		if r.Method == "POST" {
			// Handle sighnup form contents

		} else if r.Method == "GET" {

			ref := r.URL.Query().Get("ref")

			//Show Form for signup
			if ref == "" {
				//Show ref code input from and go for page reload

			} else {
				//check ref code

				if ref == "testcode" {
					// correct code show signup form
					pageStruct.RefCodeValid = true
					pageStruct.RefCode = "testcode"

				} else {
					// invalid code show error message
					pageStruct.Error = true
					pageStruct.ErrorText = "Invalid invite code"
				}

			}

		} else {

		}
	}

	webutils.RenderTemplate(w, []string{"pkg/web/templates/signup.html"}, pageStruct)

}
