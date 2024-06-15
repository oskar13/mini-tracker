package news

import (
	"github.com/oskar13/mini-tracker/pkg/db"
)

func LoadNewsArticle(newsID int) NewsArticle {

	return NewsArticle{}
}

func LoadNewsList() []NewsArticle {
	var newsList []NewsArticle

	q := "SELECT post_ID, title, date, posted_by, excerpt FROM site_news LIMIT 100"
	rows, err := db.DB.Query(q)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var newsItem NewsArticle
		err = rows.Scan(&newsItem.NewsID, &newsItem.Title, &newsItem.Author, &newsItem.Excerpt)
		if err != nil {
			// handle this error
			panic(err)
		}

		newsList = append(newsList, newsItem)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return newsList
}

func LoadNewsComments(newsID int) []NewsComment {
	var comments []NewsComment

	q := "SELECT comments.comment_ID, comments.content, users.user_ID, users.username, users.profile_pic FROM comments LEFT JOIN users ON users.user_ID = comments.user_ID WHERE post_ID = ?"
	rows, err := db.DB.Query(q, newsID)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var comment NewsComment
		err = rows.Scan(&comment.CommentID, &comment.Content, &comment.UserID, &comment.Username, &comment.Cover)
		if err != nil {
			// handle this error
			panic(err)
		}

		comments = append(comments, comment)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return comments
}
