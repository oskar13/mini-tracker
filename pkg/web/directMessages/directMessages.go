package directmessages

import (
	"fmt"

	"github.com/oskar13/mini-tracker/pkg/db"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func LoadDMThread(threadID int) webdata.DMThread {
	var thread webdata.DMThread
	var messages []webdata.DM

	users := loadDMThreadUsers(threadID)

	q := "SELECT direct_messages.message_ID, direct_messages.sender_ID, direct_messages.content, direct_messages.date, dm_threads.thread_title FROM direct_messages LEFT JOIN dm_threads ON dm_threads.dm_thread_ID = direct_messages.dm_thread_ID WHERE direct_messages.dm_thread_ID = ?"
	rows, err := db.DB.Query(q, threadID)
	if err != nil {
		//Handle error
		fmt.Println(err)
		return webdata.DMThread{}
	}
	defer rows.Close()
	for rows.Next() {
		var message webdata.DM
		err = rows.Scan(&message.MessageID, &message.SenderID, &message.Content, &message.Date, &thread.ThreadTitle)
		if err != nil {
			// handle this error
			panic(err)
		}

		message.Sender = getUserByID(message.SenderID, users)

		if message.Sender != nil {
			messages = append(messages, message)
		}

	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	thread.Users = users
	thread.ThreadID = threadID
	thread.Messages = messages

	return thread

}

func getUserByID(userID int, users []webdata.User) *webdata.User {
	for _, v := range users {
		if v.UserID == userID {
			return &v
		}
	}
	return nil
}

func loadDMThreadUsers(threadID int) []webdata.User {
	var users []webdata.User

	q := "SELECT dm_thread_users.user_ID FROM dm_thread_users WHERE dm_thread_users.dm_thread_ID = ?"
	rows, err := db.DB.Query(q, threadID)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var userID int
		err = rows.Scan(&userID)
		if err != nil {
			fmt.Println(err)
		}

		user, err := webutils.LoadUserProfileData(userID)

		if err != nil {
			fmt.Println(err)
		} else {
			users = append(users, user)
		}

	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}

	return users
}

func LoadDMThreadList(userID int) []webdata.DMThreadListItem {

	var threads []webdata.DMThreadListItem

	q := "SELECT dm_threads.dm_thread_ID, dm_threads.thread_title, dm_threads.last_message_date, dm_threads.last_message FROM dm_thread_users LEFT JOIN dm_threads ON dm_thread_users.dm_thread_ID = dm_threads.dm_thread_ID WHERE dm_thread_users.user_ID = ?"
	rows, err := db.DB.Query(q, userID)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var thread webdata.DMThreadListItem
		err = rows.Scan(&thread.ThreadID, &thread.ThreadTitle, &thread.LastMessageDate, &thread.LastMessage)
		if err != nil {
			// handle this error
			panic(err)
		}

		threads = append(threads, thread)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return threads
}
