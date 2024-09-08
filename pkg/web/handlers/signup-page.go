package handlers

import (
	"fmt"
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
		//Show user a message about being logged in, also remind that one account per lifetime
		pageStruct.Error = true
		pageStruct.ErrorText = "You already have an account. Reminder: Ban evasion is not allowed."
	} else {
		if r.Method == "POST" {
			// Handle signup form contents

			r.ParseForm()
			username := r.Form.Get("username")
			password := r.Form.Get("password")
			password2 := r.Form.Get("password2")
			ref := r.Form.Get("ref")

			if username == "" || password == "" || password2 == "" || ref == "" {

				fmt.Println(r.Form)

				pageStruct.ErrorText = "Data incomplete"
				pageStruct.Error = true

			} else {
				err := webutils.CreateUser(username, password, password2, ref, 99, false)

				if err != nil {
					pageStruct.Error = true
					fmt.Println(err)
					pageStruct.ErrorText = fmt.Sprintf("%v", err)
				} else {
					http.Redirect(w, r, "/login?message=creation-success", http.StatusSeeOther)

				}

			}

		} else if r.Method == "GET" {

			ref := r.URL.Query().Get("ref")

			//Show Form for signup
			if ref == "" {
				//Show ref code input from and go for page reload

			} else {
				//check ref code

				_, err := webutils.CheckRefCode(ref)

				if err != nil {
					// invalid code show error message
					pageStruct.Error = true
					pageStruct.ErrorText = "Invalid invite code"
				} else {
					// correct code show signup form
					pageStruct.RefCodeValid = true
					pageStruct.RefCode = ref
				}

			}

		}
	}

	webutils.RenderTemplate(w, []string{"pkg/web/templates/signup.html"}, pageStruct)

}
