package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/AlexGuo43/clans/clan-service/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClanRepository struct {
	db *pgxpool.Pool
}

func NewClanRepository(db *pgxpool.Pool) *ClanRepository {
	return &ClanRepository{
		db: db,
	}
}

func (r *ClanRepository) Create(ctx context.Context, clan *models.ClanRequest, userID int) (*models.Clan, error) {
	query := `
		INSERT INTO clans (name, display_name, description, owner_id, is_public)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	var id int
	var createdAt, updatedAt sql.NullTime

	err := r.db.QueryRow(ctx, query, clan.Name, clan.DisplayName, clan.Description, userID, clan.IsPublic).
		Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create clan: %w", err)
	}

	membershipQuery := `
		INSERT INTO clan_memberships (clan_id, user_id, role)
		VALUES ($1, $2, $3)`
	
	_, err = r.db.Exec(ctx, membershipQuery, id, userID, models.RoleOwner)
	if err != nil {
		return nil, fmt.Errorf("failed to create owner membership: %w", err)
	}

	return &models.Clan{
		ID:          id,
		Name:        clan.Name,
		DisplayName: clan.DisplayName,
		Description: clan.Description,
		OwnerID:     userID,
		MemberCount: 1,
		PostCount:   0,
		IsPublic:    clan.IsPublic,
		CreatedAt:   createdAt.Time,
		UpdatedAt:   updatedAt.Time,
	}, nil
}

func (r *ClanRepository) GetByID(ctx context.Context, id int) (*models.Clan, error) {
	query := `
		SELECT c.id, c.name, c.display_name, c.description, c.owner_id, u.username as owner_name,
			   COUNT(DISTINCT cm.id) as member_count,
			   COUNT(DISTINCT p.id) as post_count,
			   c.is_public, c.created_at, c.updated_at
		FROM clans c
		LEFT JOIN users u ON c.owner_id = u.id
		LEFT JOIN clan_memberships cm ON c.id = cm.clan_id
		LEFT JOIN posts p ON c.id = p.clan_id
		WHERE c.id = $1
		GROUP BY c.id, u.username`

	var clan models.Clan
	err := r.db.QueryRow(ctx, query, id).Scan(
		&clan.ID, &clan.Name, &clan.DisplayName, &clan.Description, &clan.OwnerID, &clan.OwnerName,
		&clan.MemberCount, &clan.PostCount, &clan.IsPublic, &clan.CreatedAt, &clan.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get clan: %w", err)
	}

	return &clan, nil
}

func (r *ClanRepository) GetByName(ctx context.Context, name string) (*models.Clan, error) {
	query := `
		SELECT c.id, c.name, c.display_name, c.description, c.owner_id, u.username as owner_name,
			   COUNT(DISTINCT cm.id) as member_count,
			   COUNT(DISTINCT p.id) as post_count,
			   c.is_public, c.created_at, c.updated_at
		FROM clans c
		LEFT JOIN users u ON c.owner_id = u.id
		LEFT JOIN clan_memberships cm ON c.id = cm.clan_id
		LEFT JOIN posts p ON c.id = p.clan_id
		WHERE c.name = $1
		GROUP BY c.id, u.username`

	var clan models.Clan
	err := r.db.QueryRow(ctx, query, name).Scan(
		&clan.ID, &clan.Name, &clan.DisplayName, &clan.Description, &clan.OwnerID, &clan.OwnerName,
		&clan.MemberCount, &clan.PostCount, &clan.IsPublic, &clan.CreatedAt, &clan.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get clan: %w", err)
	}

	return &clan, nil
}

func (r *ClanRepository) GetAll(ctx context.Context, limit, offset int) ([]models.Clan, error) {
	query := `
		SELECT c.id, c.name, c.display_name, c.description, c.owner_id, u.username as owner_name,
			   COUNT(DISTINCT cm.id) as member_count,
			   COUNT(DISTINCT p.id) as post_count,
			   c.is_public, c.created_at, c.updated_at
		FROM clans c
		LEFT JOIN users u ON c.owner_id = u.id
		LEFT JOIN clan_memberships cm ON c.id = cm.clan_id
		LEFT JOIN posts p ON c.id = p.clan_id
		WHERE c.is_public = true
		GROUP BY c.id, u.username
		ORDER BY c.created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get clans: %w", err)
	}
	defer rows.Close()

	var clans []models.Clan
	for rows.Next() {
		var clan models.Clan
		err := rows.Scan(
			&clan.ID, &clan.Name, &clan.DisplayName, &clan.Description, &clan.OwnerID, &clan.OwnerName,
			&clan.MemberCount, &clan.PostCount, &clan.IsPublic, &clan.CreatedAt, &clan.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan clan: %w", err)
		}
		clans = append(clans, clan)
	}

	return clans, nil
}

func (r *ClanRepository) Update(ctx context.Context, id int, clan *models.ClanRequest) (*models.Clan, error) {
	query := `
		UPDATE clans 
		SET display_name = $1, description = $2, is_public = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING name, owner_id, created_at, updated_at`

	var name string
	var ownerID int
	var createdAt, updatedAt sql.NullTime

	err := r.db.QueryRow(ctx, query, clan.DisplayName, clan.Description, clan.IsPublic, id).
		Scan(&name, &ownerID, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update clan: %w", err)
	}

	return r.GetByID(ctx, id)
}

func (r *ClanRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM clans WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete clan: %w", err)
	}
	return nil
}

func (r *ClanRepository) JoinClan(ctx context.Context, clanID, userID int) error {
	query := `
		INSERT INTO clan_memberships (clan_id, user_id, role)
		VALUES ($1, $2, $3)`
	
	_, err := r.db.Exec(ctx, query, clanID, userID, models.RoleMember)
	if err != nil {
		return fmt.Errorf("failed to join clan: %w", err)
	}
	return nil
}

func (r *ClanRepository) LeaveClan(ctx context.Context, clanID, userID int) error {
	query := `DELETE FROM clan_memberships WHERE clan_id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, clanID, userID)
	if err != nil {
		return fmt.Errorf("failed to leave clan: %w", err)
	}
	return nil
}

func (r *ClanRepository) GetMembership(ctx context.Context, clanID, userID int) (*models.ClanMembership, error) {
	query := `
		SELECT cm.id, cm.clan_id, cm.user_id, u.username, cm.role, cm.joined_at
		FROM clan_memberships cm
		JOIN users u ON cm.user_id = u.id
		WHERE cm.clan_id = $1 AND cm.user_id = $2`

	var membership models.ClanMembership
	err := r.db.QueryRow(ctx, query, clanID, userID).Scan(
		&membership.ID, &membership.ClanID, &membership.UserID,
		&membership.Username, &membership.Role, &membership.JoinedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get membership: %w", err)
	}

	return &membership, nil
}

func (r *ClanRepository) GetMembers(ctx context.Context, clanID int, limit, offset int) ([]models.ClanMembership, error) {
	query := `
		SELECT cm.id, cm.clan_id, cm.user_id, u.username, cm.role, cm.joined_at
		FROM clan_memberships cm
		JOIN users u ON cm.user_id = u.id
		WHERE cm.clan_id = $1
		ORDER BY cm.joined_at ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, clanID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get members: %w", err)
	}
	defer rows.Close()

	var members []models.ClanMembership
	for rows.Next() {
		var member models.ClanMembership
		err := rows.Scan(
			&member.ID, &member.ClanID, &member.UserID,
			&member.Username, &member.Role, &member.JoinedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan member: %w", err)
		}
		members = append(members, member)
	}

	return members, nil
}

func (r *ClanRepository) UpdateMemberRole(ctx context.Context, clanID, userID int, role models.ClanMembershipRole) error {
	query := `UPDATE clan_memberships SET role = $1 WHERE clan_id = $2 AND user_id = $3`
	_, err := r.db.Exec(ctx, query, role, clanID, userID)
	if err != nil {
		return fmt.Errorf("failed to update member role: %w", err)
	}
	return nil
}

func (r *ClanRepository) GetUserClans(ctx context.Context, userID int) ([]models.Clan, error) {
	query := `
		SELECT c.id, c.name, c.display_name, c.description, c.owner_id, u.username as owner_name,
			   COUNT(DISTINCT cm2.id) as member_count,
			   COUNT(DISTINCT p.id) as post_count,
			   c.is_public, c.created_at, c.updated_at
		FROM clans c
		LEFT JOIN users u ON c.owner_id = u.id
		LEFT JOIN clan_memberships cm ON c.id = cm.clan_id AND cm.user_id = $1
		LEFT JOIN clan_memberships cm2 ON c.id = cm2.clan_id
		LEFT JOIN posts p ON c.id = p.clan_id
		WHERE cm.user_id = $1
		GROUP BY c.id, u.username
		ORDER BY c.created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user clans: %w", err)
	}
	defer rows.Close()

	var clans []models.Clan
	for rows.Next() {
		var clan models.Clan
		err := rows.Scan(
			&clan.ID, &clan.Name, &clan.DisplayName, &clan.Description, &clan.OwnerID, &clan.OwnerName,
			&clan.MemberCount, &clan.PostCount, &clan.IsPublic, &clan.CreatedAt, &clan.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan clan: %w", err)
		}
		clans = append(clans, clan)
	}

	return clans, nil
}