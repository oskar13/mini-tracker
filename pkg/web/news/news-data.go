package news

type NewsArticle struct {
	NewsID  int
	Title   string
	Content string
	Author  string
	Date    string
	Excerpt string
}

type NewsComment struct {
	CommentID int
	PostID    int
	UserID    int
	Username  string
	Cover     string
	Content   int
}