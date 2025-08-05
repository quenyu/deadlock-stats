package services

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/quenyu/deadlock-stats/internal/domain"
	"github.com/quenyu/deadlock-stats/internal/repositories"
)

type CrosshairService struct {
	repo *repositories.CrosshairRepository
}

func NewCrosshairService(repo *repositories.CrosshairRepository) *CrosshairService {
	return &CrosshairService{repo: repo}
}

type CreateCrosshairRequest struct {
	Title       string                   `json:"title"`
	Description string                   `json:"description"`
	Settings    domain.CrosshairSettings `json:"settings"`
	IsPublic    bool                     `json:"is_public"`
}

type CrosshairResponse struct {
	ID           uuid.UUID                `json:"id"`
	AuthorID     uuid.UUID                `json:"author_id"`
	AuthorName   string                   `json:"author_name"`
	AuthorAvatar string                   `json:"author_avatar"`
	Title        string                   `json:"title"`
	Description  string                   `json:"description"`
	Settings     domain.CrosshairSettings `json:"settings"`
	LikesCount   int                      `json:"likes_count"`
	IsPublic     bool                     `json:"is_public"`
	ViewCount    int                      `json:"view_count"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
}

func (s *CrosshairService) Create(authorID uuid.UUID, req *CreateCrosshairRequest) (*CrosshairResponse, error) {
	settingsJSON, err := json.Marshal(req.Settings)
	if err != nil {
		return nil, err
	}

	crosshair := &domain.Crosshair{
		AuthorID:    authorID,
		Title:       req.Title,
		Description: req.Description,
		Settings:    settingsJSON,
		IsPublic:    req.IsPublic,
		LikesCount:  0,
		ViewCount:   0,
	}
	if err := s.repo.Create(crosshair); err != nil {
		return nil, err
	}
	return toCrosshairResponse(crosshair), nil
}

func (s *CrosshairService) GetByID(id uuid.UUID) (*CrosshairResponse, error) {
	crosshair, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return toCrosshairResponse(crosshair), nil
}

func (s *CrosshairService) GetAll(page, limit int) ([]CrosshairResponse, error) {
	crosshairs, err := s.repo.GetAll(page, limit)
	if err != nil {
		return nil, err
	}
	resp := make([]CrosshairResponse, len(crosshairs))
	for i, c := range crosshairs {
		resp[i] = *toCrosshairResponse(&c)
	}
	return resp, nil
}

func (s *CrosshairService) Count() (int64, error) {
	return s.repo.Count()
}

func (s *CrosshairService) GetByAuthorID(authorID uuid.UUID, limit int) ([]CrosshairResponse, error) {
	crosshairs, err := s.repo.GetByAuthorID(authorID, limit)
	if err != nil {
		return nil, err
	}
	resp := make([]CrosshairResponse, len(crosshairs))
	for i, c := range crosshairs {
		resp[i] = *toCrosshairResponse(&c)
	}
	return resp, nil
}

func (s *CrosshairService) Like(crosshairID, userID uuid.UUID) error {
	return s.repo.Like(crosshairID, userID)
}

func (s *CrosshairService) Unlike(crosshairID, userID uuid.UUID) error {
	return s.repo.Unlike(crosshairID, userID)
}

func (s *CrosshairService) Delete(id, authorID uuid.UUID) error {
	return s.repo.Delete(id, authorID)
}

func toCrosshairResponse(c *domain.Crosshair) *CrosshairResponse {
	var settings domain.CrosshairSettings
	if c.Settings != nil {
		json.Unmarshal(c.Settings, &settings)
	}

	authorName := ""
	authorAvatar := ""
	if c.Author != nil {
		authorName = c.Author.Nickname
		authorAvatar = c.Author.AvatarURL
	}

	return &CrosshairResponse{
		ID:           c.ID,
		AuthorID:     c.AuthorID,
		AuthorName:   authorName,
		AuthorAvatar: authorAvatar,
		Title:        c.Title,
		Description:  c.Description,
		Settings:     settings,
		LikesCount:   c.LikesCount,
		IsPublic:     c.IsPublic,
		ViewCount:    c.ViewCount,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}
