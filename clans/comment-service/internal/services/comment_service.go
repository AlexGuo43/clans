package services

import (
	"errors"

	"github.com/AlexGuo43/clans/comment-service/internal/models"
	"github.com/AlexGuo43/clans/comment-service/internal/repository"
)

type CommentService struct {
	Repo *repository.CommentRepository
}

func NewCommentService(repo *repository.CommentRepository) *CommentService {
	return &CommentService{Repo: repo}
}

func (s *CommentService) CreateComment(content string, postID, userID int, parentID *int) (*models.Comment, error) {
	if content == "" {
		return nil, errors.New("content is required")
	}

	// Validate parent comment exists and belongs to same post
	if parentID != nil {
		parentComment, err := s.Repo.GetCommentByID(*parentID)
		if err != nil {
			return nil, errors.New("parent comment not found")
		}
		if parentComment.PostID != postID {
			return nil, errors.New("parent comment must be on the same post")
		}
	}

	comment := &models.Comment{
		Content:  content,
		PostID:   postID,
		UserID:   userID,
		ParentID: parentID,
	}

	err := s.Repo.CreateComment(comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) GetComment(id int) (*models.Comment, error) {
	return s.Repo.GetCommentByID(id)
}

func (s *CommentService) GetCommentsByPost(postID, page, limit int) ([]*models.Comment, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	flatComments, err := s.Repo.GetCommentsByPost(postID, limit, offset)
	if err != nil {
		return nil, err
	}

	return s.buildCommentTree(flatComments), nil
}

func (s *CommentService) buildCommentTree(flatComments []*models.Comment) []*models.Comment {
	commentMap := make(map[int]*models.Comment)
	var rootComments []*models.Comment

	// First pass: create map and initialize replies slice
	for _, comment := range flatComments {
		comment.Replies = make([]*models.Comment, 0)
		commentMap[comment.ID] = comment
	}

	// Second pass: build tree structure
	for _, comment := range flatComments {
		if comment.ParentID != nil {
			// This is a reply
			parentID := *comment.ParentID
			if parent, exists := commentMap[parentID]; exists {
				parent.Replies = append(parent.Replies, comment)
			}
		} else {
			// This is a top-level comment
			rootComments = append(rootComments, comment)
		}
	}

	return rootComments
}

func (s *CommentService) GetReplies(parentID, page, limit int) ([]*models.Comment, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	return s.Repo.GetReplies(parentID, limit, offset)
}

func (s *CommentService) UpdateComment(id int, content string, userID int) error {
	comment, err := s.Repo.GetCommentByID(id)
	if err != nil {
		return err
	}

	if comment.UserID != userID {
		return errors.New("unauthorized: can only edit your own comments")
	}

	if content != "" {
		comment.Content = content
	}

	return s.Repo.UpdateComment(comment)
}

func (s *CommentService) DeleteComment(id, userID int) error {
	comment, err := s.Repo.GetCommentByID(id)
	if err != nil {
		return err
	}

	if comment.UserID != userID {
		return errors.New("unauthorized: can only delete your own comments")
	}

	return s.Repo.DeleteComment(id)
}

func (s *CommentService) VoteComment(userID, commentID int, isUpvote bool) error {
	_, err := s.Repo.GetCommentByID(commentID)
	if err != nil {
		return errors.New("comment not found")
	}

	return s.Repo.VoteComment(userID, commentID, isUpvote)
}

func (s *CommentService) RemoveVote(userID, commentID int) error {
	return s.Repo.RemoveVote(userID, commentID)
}

