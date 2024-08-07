package webdata

// Data object to hold torrent data when loading data from and to database, also to be used in the templates
type TorrentWeb struct {
	TorrentID        int
	Announce         string //unique url if torrent is not public
	AnnounceList     []string
	AnnounceListJSON *string
	Name             string
	Description      string //comment added on web page
	Comment          string //file metadata comment
	Type             string //category it was posted in
	PieceLength      *int64
	Pieces           *[]byte
	GroupID          *int
	GroupName        *string
	User             User //Uploader
	Uploaded         string
	InfoHash         string
	Encoding         string
	PathJSON         *string
	Size             string
	Seeders          int
	Leechers         int
	FilesCount       int
	Discussion       []TorrentComment
	Tags             []string
	TagIDs           []string
	TagColors        []string
	UpVotes          int
	DownVotes        int
	Anonymous        bool
	AccessType       string
	InfoField        []byte
	Uuid             string
	CategoryID       int
	ParentCategoryID int
	Category         string
	ParentCategory   string
	KeepTrackers     bool
}

type TorrentComment struct {
	CommentID int
	User      User
	Date      string
	Content   string
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

type Strike struct {
	StrikeID    int
	UserID      int
	Heading     string
	Description string
	Date        string
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
	ThreadTitle string
	ThreadID    int
	Users       []User
	Messages    []DM
}

// Sender ID is used temporarily to hold userID till the user pointer gets tied to a correct entry
type DM struct {
	MessageID int
	SenderID  int
	Sender    *User
	Content   string
	Date      string
}

type FriendRequests struct {
	Incoming []FriendRequest
	Outgoing []FriendRequest
}

type FriendRequest struct {
	User    User
	Message *string
	Date    string
}

var SiteName string = "Mini Tracker"
