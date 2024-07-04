package groups

import "github.com/oskar13/mini-tracker/pkg/web/webdata"

type GroupPost struct {
	Group      GroupInfo
	PostID     int
	User       webdata.User
	Title      string
	Content    string
	Date       string
	Updated    string
	Sticky     bool
	LastPost   *string
	Replies    []GroupPostReply
	ReplyCount int
}

type GroupPostReply struct {
	ReplyID int
	PostID  int
	User    webdata.User
	Content string
	Date    string
	Updated string
}

type GroupInfo struct {
	GroupName        string
	GroupID          int
	GroupIcon        *string
	GroupRole        *string
	GroupVisibility  string
	GroupTagline     *string
	GroupJoinType    string
	GroupDescription *string
}
