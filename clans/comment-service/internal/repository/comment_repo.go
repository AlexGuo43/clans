package repository

import (
	"context"

	"github.com/AlexGuo43/clans/comment-service/internal/models"
	"github.com/jackc/pgx/v5"
)

type CommentRepository struct {
	db *pgx.Conn
}

func NewCommentRepository(db *pgx.Conn) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) CreateComment(comment *models.Comment) error {
	query := `
		INSERT INTO comments (content, post_id, user_id, parent_id, depth, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) 
		RETURNING id, created_at, updated_at`

	depth := 0
	if comment.ParentID != nil {
		// Calculate depth based on parent
		depthQuery := `SELECT depth FROM comments WHERE id = $1`
		var parentDepth int
		err := r.db.QueryRow(context.Background(), depthQuery, *comment.ParentID).Scan(&parentDepth)
		if err != nil {
			return err
		}
		depth = parentDepth + 1
	}

	err := r.db.QueryRow(context.Background(), query, 
		comment.Content, comment.PostID, comment.UserID, comment.ParentID, depth).
		Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
	
	return err
}

func (r *CommentRepository) GetCommentByID(id int) (*models.Comment, error) {
	query := `
		SELECT c.id, c.content, c.post_id, c.user_id, u.username, c.parent_id, c.depth,
			   COALESCE(SUM(CASE WHEN cv.is_upvote THEN 1 ELSE -1 END), 0) as vote_count,
			   (SELECT COUNT(*) FROM comments WHERE parent_id = c.id) as reply_count,
			   c.created_at, c.updated_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		LEFT JOIN comment_votes cv ON c.id = cv.comment_id
		WHERE c.id = $1
		GROUP BY c.id, u.username`

	comment := &models.Comment{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&comment.ID, &comment.Content, &comment.PostID, &comment.UserID, &comment.Username,
		&comment.ParentID, &comment.Depth, &comment.VoteCount, &comment.ReplyCount,
		&comment.CreatedAt, &comment.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *CommentRepository) GetCommentsByPost(postID int, limit, offset int) ([]*models.Comment, error) {
	query := `
		SELECT c.id, c.content, c.post_id, c.user_id, u.username, c.parent_id, c.depth,
			   COALESCE(SUM(CASE WHEN cv.is_upvote THEN 1 ELSE -1 END), 0) as vote_count,
			   (SELECT COUNT(*) FROM comments WHERE parent_id = c.id) as reply_count,
			   c.created_at, c.updated_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		LEFT JOIN comment_votes cv ON c.id = cv.comment_id
		WHERE c.post_id = $1
		GROUP BY c.id, u.username
		ORDER BY c.depth ASC, c.created_at ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(context.Background(), query, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(
			&comment.ID, &comment.Content, &comment.PostID, &comment.UserID, &comment.Username,
			&comment.ParentID, &comment.Depth, &comment.VoteCount, &comment.ReplyCount,
			&comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *CommentRepository) GetReplies(parentID int, limit, offset int) ([]*models.Comment, error) {
	query := `
		SELECT c.id, c.content, c.post_id, c.user_id, u.username, c.parent_id, c.depth,
			   COALESCE(SUM(CASE WHEN cv.is_upvote THEN 1 ELSE -1 END), 0) as vote_count,
			   (SELECT COUNT(*) FROM comments WHERE parent_id = c.id) as reply_count,
			   c.created_at, c.updated_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		LEFT JOIN comment_votes cv ON c.id = cv.comment_id
		WHERE c.parent_id = $1
		GROUP BY c.id, u.username
		ORDER BY c.created_at ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(context.Background(), query, parentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(
			&comment.ID, &comment.Content, &comment.PostID, &comment.UserID, &comment.Username,
			&comment.ParentID, &comment.Depth, &comment.VoteCount, &comment.ReplyCount,
			&comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *CommentRepository) UpdateComment(comment *models.Comment) error {
	query := `UPDATE comments SET content = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(context.Background(), query, comment.Content, comment.ID)
	return err
}

func (r *CommentRepository) DeleteComment(id int) error {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

func (r *CommentRepository) VoteComment(userID, commentID int, isUpvote bool) error {
	query := `
		INSERT INTO comment_votes (comment_id, user_id, is_upvote) 
		VALUES ($1, $2, $3)
		ON CONFLICT (comment_id, user_id) 
		DO UPDATE SET is_upvote = $3`
	
	_, err := r.db.Exec(context.Background(), query, commentID, userID, isUpvote)
	return err
}

func (r *CommentRepository) RemoveVote(userID, commentID int) error {
	query := `DELETE FROM comment_votes WHERE comment_id = $1 AND user_id = $2`
	_, err := r.db.Exec(context.Background(), query, commentID, userID)
	return err
}