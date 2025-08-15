package models

import "time"

type Comment struct {
	ID           int       `json:"id"`
	Content      string    `json:"content"`
	PostID       int       `json:"post_id"`
	UserID       int       `json:"user_id"`
	Username     string    `json:"username"`
	ParentID     *int      `json:"parent_id,omitempty"`
	VoteCount    int       `json:"vote_count"`
	ReplyCount   int       `json:"reply_count"`
	Depth        int       `json:"depth"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Replies      []*Comment `json:"replies,omitempty"`
}

type CommentVote struct {
	ID        int  `json:"id"`
	CommentID int  `json:"comment_id"`
	UserID    int  `json:"user_id"`
	IsUpvote  bool `json:"is_upvote"`
}

type CommentRequest struct {
	Content  string `json:"content"`
	PostID   int    `json:"post_id"`
	ParentID *int   `json:"parent_id,omitempty"`
}

type CommentTree struct {
	Comment Comment   `json:"comment"`
	Replies []Comment `json:"replies"`
}