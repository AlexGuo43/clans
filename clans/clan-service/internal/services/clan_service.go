package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/AlexGuo43/clans/clan-service/internal/models"
	"github.com/AlexGuo43/clans/clan-service/internal/repository"
)

type ClanService struct {
	clanRepo *repository.ClanRepository
}

func NewClanService(clanRepo *repository.ClanRepository) *ClanService {
	return &ClanService{
		clanRepo: clanRepo,
	}
}

func (s *ClanService) CreateClan(ctx context.Context, req *models.ClanRequest, userID int) (*models.Clan, error) {
	req.Name = strings.ToLower(strings.TrimSpace(req.Name))
	req.DisplayName = strings.TrimSpace(req.DisplayName)
	req.Description = strings.TrimSpace(req.Description)

	if req.Name == "" {
		return nil, fmt.Errorf("clan name is required")
	}
	if req.DisplayName == "" {
		return nil, fmt.Errorf("clan display name is required")
	}
	if len(req.Name) < 3 || len(req.Name) > 20 {
		return nil, fmt.Errorf("clan name must be between 3 and 20 characters")
	}
	if len(req.DisplayName) > 50 {
		return nil, fmt.Errorf("clan display name must be 50 characters or less")
	}
	if len(req.Description) > 500 {
		return nil, fmt.Errorf("clan description must be 500 characters or less")
	}

	if !isValidClanName(req.Name) {
		return nil, fmt.Errorf("clan name can only contain letters, numbers, and underscores")
	}

	existing, _ := s.clanRepo.GetByName(ctx, req.Name)
	if existing != nil {
		return nil, fmt.Errorf("clan name already exists")
	}

	return s.clanRepo.Create(ctx, req, userID)
}

func (s *ClanService) GetClan(ctx context.Context, id int) (*models.Clan, error) {
	return s.clanRepo.GetByID(ctx, id)
}

func (s *ClanService) GetClanByName(ctx context.Context, name string) (*models.Clan, error) {
	return s.clanRepo.GetByName(ctx, name)
}

func (s *ClanService) GetClans(ctx context.Context, limit, offset int) ([]models.Clan, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	return s.clanRepo.GetAll(ctx, limit, offset)
}

func (s *ClanService) UpdateClan(ctx context.Context, id int, req *models.ClanRequest, userID int) (*models.Clan, error) {
	clan, err := s.clanRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("clan not found")
	}

	if clan.OwnerID != userID {
		membership, err := s.clanRepo.GetMembership(ctx, id, userID)
		if err != nil || membership.Role != models.RoleModerator {
			return nil, fmt.Errorf("insufficient permissions")
		}
	}

	req.DisplayName = strings.TrimSpace(req.DisplayName)
	req.Description = strings.TrimSpace(req.Description)

	if req.DisplayName == "" {
		return nil, fmt.Errorf("clan display name is required")
	}
	if len(req.DisplayName) > 50 {
		return nil, fmt.Errorf("clan display name must be 50 characters or less")
	}
	if len(req.Description) > 500 {
		return nil, fmt.Errorf("clan description must be 500 characters or less")
	}

	return s.clanRepo.Update(ctx, id, req)
}

func (s *ClanService) DeleteClan(ctx context.Context, id, userID int) error {
	clan, err := s.clanRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("clan not found")
	}

	if clan.OwnerID != userID {
		return fmt.Errorf("only the clan owner can delete the clan")
	}

	return s.clanRepo.Delete(ctx, id)
}

func (s *ClanService) JoinClan(ctx context.Context, clanID, userID int) error {
	clan, err := s.clanRepo.GetByID(ctx, clanID)
	if err != nil {
		return fmt.Errorf("clan not found")
	}

	if !clan.IsPublic {
		return fmt.Errorf("clan is private")
	}

	existing, _ := s.clanRepo.GetMembership(ctx, clanID, userID)
	if existing != nil {
		return fmt.Errorf("already a member of this clan")
	}

	return s.clanRepo.JoinClan(ctx, clanID, userID)
}

func (s *ClanService) LeaveClan(ctx context.Context, clanID, userID int) error {
	clan, err := s.clanRepo.GetByID(ctx, clanID)
	if err != nil {
		return fmt.Errorf("clan not found")
	}

	if clan.OwnerID == userID {
		return fmt.Errorf("clan owner cannot leave the clan")
	}

	membership, err := s.clanRepo.GetMembership(ctx, clanID, userID)
	if err != nil {
		return fmt.Errorf("not a member of this clan")
	}

	if membership == nil {
		return fmt.Errorf("not a member of this clan")
	}

	return s.clanRepo.LeaveClan(ctx, clanID, userID)
}

func (s *ClanService) GetMembers(ctx context.Context, clanID int, limit, offset int) ([]models.ClanMembership, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	_, err := s.clanRepo.GetByID(ctx, clanID)
	if err != nil {
		return nil, fmt.Errorf("clan not found")
	}

	return s.clanRepo.GetMembers(ctx, clanID, limit, offset)
}

func (s *ClanService) UpdateMemberRole(ctx context.Context, clanID, targetUserID, userID int, role models.ClanMembershipRole) error {
	clan, err := s.clanRepo.GetByID(ctx, clanID)
	if err != nil {
		return fmt.Errorf("clan not found")
	}

	if clan.OwnerID != userID {
		membership, err := s.clanRepo.GetMembership(ctx, clanID, userID)
		if err != nil || membership.Role != models.RoleModerator {
			return fmt.Errorf("insufficient permissions")
		}
	}

	if targetUserID == clan.OwnerID && role != models.RoleOwner {
		return fmt.Errorf("cannot change the role of the clan owner")
	}

	if role != models.RoleMember && role != models.RoleModerator && role != models.RoleOwner {
		return fmt.Errorf("invalid role")
	}

	targetMembership, err := s.clanRepo.GetMembership(ctx, clanID, targetUserID)
	if err != nil || targetMembership == nil {
		return fmt.Errorf("user is not a member of this clan")
	}

	if role == models.RoleOwner {
		if clan.OwnerID != userID {
			return fmt.Errorf("only the current owner can transfer ownership")
		}
		return fmt.Errorf("ownership transfer not implemented")
	}

	return s.clanRepo.UpdateMemberRole(ctx, clanID, targetUserID, role)
}

func (s *ClanService) GetUserClans(ctx context.Context, userID int) ([]models.Clan, error) {
	return s.clanRepo.GetUserClans(ctx, userID)
}

func (s *ClanService) GetMembership(ctx context.Context, clanID, userID int) (*models.ClanMembership, error) {
	return s.clanRepo.GetMembership(ctx, clanID, userID)
}

func isValidClanName(name string) bool {
	if len(name) == 0 {
		return false
	}
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '_') {
			return false
		}
	}
	return true
}