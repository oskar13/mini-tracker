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

		} else {

			//Show Form for signup
			if r.URL.Path == "/ref/" {
				//Handle ref code in url

			} else {
				//Otherwise handle ref code in form

			}

		}
	}

	webutils.RenderTemplate(w, []string{"pkg/web/templates/signup.html"}, pageStruct)

}
