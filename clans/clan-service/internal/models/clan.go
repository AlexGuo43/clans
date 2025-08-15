package models

import "time"

type Clan struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	OwnerID     int       `json:"owner_id"`
	OwnerName   string    `json:"owner_name"`
	MemberCount int       `json:"member_count"`
	PostCount   int       `json:"post_count"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ClanMembership struct {
	ID       int                `json:"id"`
	ClanID   int                `json:"clan_id"`
	UserID   int                `json:"user_id"`
	Username string             `json:"username"`
	Role     ClanMembershipRole `json:"role"`
	JoinedAt time.Time          `json:"joined_at"`
}

type ClanMembershipRole string

const (
	RoleMember    ClanMembershipRole = "member"
	RoleModerator ClanMembershipRole = "moderator"
	RoleOwner     ClanMembershipRole = "owner"
)

type ClanRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

type ClanStats struct {
	MemberCount int `json:"member_count"`
	PostCount   int `json:"post_count"`
}

type ClanJoinRequest struct {
	ClanID int `json:"clan_id"`
}

type ClanMembershipUpdate struct {
	Role ClanMembershipRole `json:"role"`
}