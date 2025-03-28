package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oskar13/mini-tracker/pkg/web/accounts"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

// Handle login and logout URL endpoint
func LoginPage(c *gin.Context) {

	var pageStruct struct {
		Error         bool
		ErrorText     string
		Message       bool
		MessageText   string
		LoggedIn      bool
		LogoutMessage bool
		UserData      webdata.User
	}

	userData := accounts.GetUserData(c.Request)

	if userData.LoggedIn {

		if c.Request.URL.Path == "/logout" {
			//Log user out
			userData = accounts.LogOutUser(c.Writer, userData)
			pageStruct.LogoutMessage = true

		} else {
			//Offer log out
			pageStruct.LoggedIn = true
		}

	} else {
		if c.Request.Method == "POST" {
			// Check login details

			if err := c.Request.ParseForm(); err != nil {
				http.Error(c.Writer, "Error parsing the form", http.StatusInternalServerError)
				return
			}

			loginData := webdata.LoginData{
				UserNameOrEmail: c.Request.FormValue("username"),
				Password:        c.Request.FormValue("password"),
			}

			_, err := accounts.LoginUser(c.Writer, c.Request, loginData.UserNameOrEmail, loginData.Password)

			if err != nil {

				pageStruct.Error = true
				pageStruct.ErrorText = "Invalid username or password"

				c.Writer.WriteHeader(http.StatusUnauthorized)
				webutils.RenderTemplate(c.Writer, []string{"pkg/web/templates/login.html"}, pageStruct)
				return
			}

			http.Redirect(c.Writer, c.Request, "/", http.StatusSeeOther)

			return
		} else if c.Request.Method == "GET" {
			message := c.Request.URL.Query().Get("message")

			if message != "" {

				if message == "creation-success" {
					pageStruct.Message = true

					pageStruct.MessageText = "Account creation success. You can now log in."
				}

				if message == "creation-fail" {
					pageStruct.Error = true

					pageStruct.ErrorText = "Failed creating account"

					if c.Request.URL.Query().Get("reason") != "" {
						pageStruct.ErrorText = pageStruct.ErrorText + ": " + c.Request.URL.Query().Get("reason")
					}
				}

			}

		}
	}

	webutils.RenderTemplate(c.Writer, []string{"pkg/web/templates/login.html"}, pageStruct)

}
