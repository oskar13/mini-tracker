package webdata

type TorrentWeb struct {
	TorrentID   int
	Created     string
	UserID      int
	Name        string
	Anonymous   bool
	AccessType  string
	GroupID     int
	UpVotes     int
	DownVotes   int
	Description string
	InfoHash    string
	Size        string
	Path        string
	Tags        []string
	TagIDs      []string
	TagColors   []string
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
	Gender        *string
	Tagline       *string
	Bio           *string
}
