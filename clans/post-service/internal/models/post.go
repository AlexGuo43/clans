package models

import "time"

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	ClanID    *int      `json:"clan_id,omitempty"`
	ClanName  *string   `json:"clan_name,omitempty"`
	VoteCount int       `json:"vote_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostVote struct {
	ID      int  `json:"id"`
	PostID  int  `json:"post_id"`
	UserID  int  `json:"user_id"`
	IsUpvote bool `json:"is_upvote"`
}