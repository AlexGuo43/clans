package repository

import (
	"context"

	"github.com/AlexGuo43/clans/post-service/internal/models"
	"github.com/jackc/pgx/v5"
)

type PostRepository struct {
	db *pgx.Conn
}

func NewPostRepository(db *pgx.Conn) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) CreatePost(post *models.Post) error {
	query := `
		INSERT INTO posts (title, content, user_id, clan_id, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, NOW(), NOW()) 
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(context.Background(), query, 
		post.Title, post.Content, post.UserID, post.ClanID).
		Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	
	return err
}

func (r *PostRepository) GetPostByID(id int) (*models.Post, error) {
	query := `
		SELECT p.id, p.title, p.content, p.user_id, u.username, 
			   p.clan_id, c.name as clan_name, 
			   COALESCE(SUM(CASE WHEN pv.is_upvote THEN 1 ELSE -1 END), 0) as vote_count,
			   p.created_at, p.updated_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN clans c ON p.clan_id = c.id
		LEFT JOIN post_votes pv ON p.id = pv.post_id
		WHERE p.id = $1
		GROUP BY p.id, u.username, c.name`

	post := &models.Post{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&post.ID, &post.Title, &post.Content, &post.UserID, &post.Username,
		&post.ClanID, &post.ClanName, &post.VoteCount,
		&post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *PostRepository) GetPosts(limit, offset int) ([]*models.Post, error) {
	query := `
		SELECT p.id, p.title, p.content, p.user_id, u.username, 
			   p.clan_id, c.name as clan_name, 
			   COALESCE(SUM(CASE WHEN pv.is_upvote THEN 1 ELSE -1 END), 0) as vote_count,
			   p.created_at, p.updated_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN clans c ON p.clan_id = c.id
		LEFT JOIN post_votes pv ON p.id = pv.post_id
		GROUP BY p.id, u.username, c.name
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.UserID, &post.Username,
			&post.ClanID, &post.ClanName, &post.VoteCount,
			&post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) GetPostsByClan(clanID, limit, offset int) ([]*models.Post, error) {
	query := `
		SELECT p.id, p.title, p.content, p.user_id, u.username, 
			   p.clan_id, c.name as clan_name, 
			   COALESCE(SUM(CASE WHEN pv.is_upvote THEN 1 ELSE -1 END), 0) as vote_count,
			   p.created_at, p.updated_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN clans c ON p.clan_id = c.id
		LEFT JOIN post_votes pv ON p.id = pv.post_id
		WHERE p.clan_id = $1
		GROUP BY p.id, u.username, c.name
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(context.Background(), query, clanID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.UserID, &post.Username,
			&post.ClanID, &post.ClanName, &post.VoteCount,
			&post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) UpdatePost(post *models.Post) error {
	query := `UPDATE posts SET title = $1, content = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.Exec(context.Background(), query, post.Title, post.Content, post.ID)
	return err
}

func (r *PostRepository) DeletePost(id int) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

func (r *PostRepository) VotePost(userID, postID int, isUpvote bool) error {
	query := `
		INSERT INTO post_votes (post_id, user_id, is_upvote) 
		VALUES ($1, $2, $3)
		ON CONFLICT (post_id, user_id) 
		DO UPDATE SET is_upvote = $3`
	
	_, err := r.db.Exec(context.Background(), query, postID, userID, isUpvote)
	return err
}

func (r *PostRepository) RemoveVote(userID, postID int) error {
	query := `DELETE FROM post_votes WHERE post_id = $1 AND user_id = $2`
	_, err := r.db.Exec(context.Background(), query, postID, userID)
	return err
}