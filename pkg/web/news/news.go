package news

import (
	"fmt"

	"github.com/oskar13/mini-tracker/pkg/db"
)

func LoadNewsArticle(newsID int) (NewsArticle, error) {

	var theArticle NewsArticle

	q := "SELECT post_ID, title, date, posted_by, content, commenting FROM site_news WHERE post_ID = ?"
	err := db.DB.QueryRow(q, newsID).Scan(&theArticle.NewsID, &theArticle.Title, &theArticle.Date, &theArticle.Author, &theArticle.Content, &theArticle.Commenting)
	if err != nil {

		return NewsArticle{}, err
	}

	theComments, err := LoadNewsComments(newsID)

	if err != nil {
		return theArticle, err
	}

	theArticle.Comments = theComments

	return theArticle, nil
}

func LoadNewsList(limit int) ([]NewsArticle, error) {
	var newsList []NewsArticle

	q := "SELECT post_ID, title, date, posted_by, excerpt FROM site_news ORDER BY date DESC LIMIT ? "
	rows, err := db.DB.Query(q, limit)
	if err != nil {
		return []NewsArticle{}, nil
	}
	defer rows.Close()
	for rows.Next() {
		var newsItem NewsArticle
		err = rows.Scan(&newsItem.NewsID, &newsItem.Title, &newsItem.Date, &newsItem.Author, &newsItem.Excerpt)
		if err != nil {
			return []NewsArticle{}, nil
		}

		newsList = append(newsList, newsItem)
	}
	err = rows.Err()
	if err != nil {
		return []NewsArticle{}, nil
	}

	return newsList, nil
}

func LoadNewsComments(newsID int) ([]NewsComment, error) {
	var comments []NewsComment

	fmt.Println("Loading comments")

	q := "SELECT comments.comment_ID, comments.content, users.user_ID, users.username, users.profile_pic FROM comments LEFT JOIN users ON users.user_ID = comments.user_ID WHERE comments.post_ID = ?"
	rows, err := db.DB.Query(q, newsID)
	if err != nil {
		return []NewsComment{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment NewsComment
		err = rows.Scan(&comment.CommentID, &comment.Content, &comment.UserID, &comment.Username, &comment.Cover)
		if err != nil {

			return []NewsComment{}, err
		}

		comments = append(comments, comment)
	}
	err = rows.Err()
	if err != nil {
		return []NewsComment{}, err
	}

	return comments, nil
}
