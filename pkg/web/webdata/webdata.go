package webdata

type TorrentWeb struct {
}

type LoginData struct {
	UserNameOrEmail string
	Password        string
}

type User struct {
	UserID        int
	Username      string
	Password      string
	Email         string
	Cover         string
	Joined        string
	LoggedIn      bool
	Disabled      bool
	SessionUID    string
	SessionExpiry string
	Gender        string
}
