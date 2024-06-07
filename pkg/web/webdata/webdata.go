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
	UserID           int
	Username         string
	Password         string
	Email            string
	Cover            string
	Banner           *string
	Joined           string
	LoggedIn         bool
	Disabled         bool
	SessionUID       string
	SessionExpiry    string
	Gender           *string
	Tagline          *string
	Bio              *string
	InvitationStatus int
	Blocked          bool
	UserBadges       *[]Badges
	UserBadgesBlob   *string
}

type Badges struct {
	BadgeTitle string `json:"badgeTitle"`
	Color      string `json:"color"`
}

type DMThreadListItem struct {
	ThreadID        int
	ThreadTitle     string
	LastMessage     string
	LastMessageDate string
}

type DMThread struct {
	Users    []User
	Messages []DM
}

// Sender ID is used temporarily to hold userID till the user pointer gets tied to a correct entry
type DM struct {
	MessageID int
	SenderID  int
	Sender    *User
	Content   string
	Date      string
}
