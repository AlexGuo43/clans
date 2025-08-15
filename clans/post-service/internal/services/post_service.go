package services

import (
	"errors"

	"github.com/AlexGuo43/clans/post-service/internal/models"
	"github.com/AlexGuo43/clans/post-service/internal/repository"
)

type PostService struct {
	Repo *repository.PostRepository
}

func NewPostService(repo *repository.PostRepository) *PostService {
	return &PostService{Repo: repo}
}

func (s *PostService) CreatePost(title, content string, userID int, clanID *int) (*models.Post, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	if content == "" {
		return nil, errors.New("content is required")
	}

	post := &models.Post{
		Title:   title,
		Content: content,
		UserID:  userID,
		ClanID:  clanID,
	}

	err := s.Repo.CreatePost(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) GetPost(id int) (*models.Post, error) {
	return s.Repo.GetPostByID(id)
}

func (s *PostService) GetPosts(page, limit int) ([]*models.Post, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	return s.Repo.GetPosts(limit, offset)
}

func (s *PostService) GetPostsByClan(clanID, page, limit int) ([]*models.Post, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	return s.Repo.GetPostsByClan(clanID, limit, offset)
}

func (s *PostService) UpdatePost(id int, title, content string, userID int) error {
	post, err := s.Repo.GetPostByID(id)
	if err != nil {
		return err
	}

	if post.UserID != userID {
		return errors.New("unauthorized: can only edit your own posts")
	}

	if title != "" {
		post.Title = title
	}
	if content != "" {
		post.Content = content
	}

	return s.Repo.UpdatePost(post)
}

func (s *PostService) DeletePost(id, userID int) error {
	post, err := s.Repo.GetPostByID(id)
	if err != nil {
		return err
	}

	if post.UserID != userID {
		return errors.New("unauthorized: can only delete your own posts")
	}

	return s.Repo.DeletePost(id)
}

func (s *PostService) VotePost(userID, postID int, isUpvote bool) error {
	_, err := s.Repo.GetPostByID(postID)
	if err != nil {
		return errors.New("post not found")
	}

	return s.Repo.VotePost(userID, postID, isUpvote)
}

func (s *PostService) RemoveVote(userID, postID int) error {
	return s.Repo.RemoveVote(userID, postID)
}