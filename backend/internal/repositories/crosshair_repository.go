package repositories

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/quenyu/deadlock-stats/internal/domain"
	"gorm.io/gorm"
)

type CrosshairRepository struct {
	db *gorm.DB
}

func NewCrosshairRepository(db *gorm.DB) *CrosshairRepository {
	return &CrosshairRepository{db: db}
}

func (r *CrosshairRepository) Create(crosshair *domain.Crosshair) error {
	crosshair.ID = uuid.New()
	crosshair.CreatedAt = time.Now()
	crosshair.UpdatedAt = time.Now()
	return r.db.Create(crosshair).Error
}

func (r *CrosshairRepository) GetByID(id uuid.UUID) (*domain.Crosshair, error) {
	var crosshair domain.Crosshair
	err := r.db.First(&crosshair, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("crosshair not found")
		}
		return nil, err
	}
	return &crosshair, nil
}

func (r *CrosshairRepository) GetAll(page, limit int) ([]domain.Crosshair, error) {
	var crosshairs []domain.Crosshair
	offset := (page - 1) * limit
	err := r.db.Preload("Author").Order("created_at DESC").Limit(limit).Offset(offset).Find(&crosshairs).Error
	return crosshairs, err
}

func (r *CrosshairRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Crosshair{}).Count(&count).Error
	return count, err
}

func (r *CrosshairRepository) GetByAuthorID(authorID uuid.UUID, limit int) ([]domain.Crosshair, error) {
	var crosshairs []domain.Crosshair
	err := r.db.Where("author_id = ?", authorID).Order("created_at DESC").Limit(limit).Find(&crosshairs).Error
	return crosshairs, err
}

func (r *CrosshairRepository) Like(crosshairID, userID uuid.UUID) error {
	// TODO: Implement
	return nil
}

func (r *CrosshairRepository) Unlike(crosshairID, userID uuid.UUID) error {
	// TODO: Implement
	return nil
}

func (r *CrosshairRepository) Delete(id, authorID uuid.UUID) error {
	result := r.db.Where("id = ? AND author_id = ?", id, authorID).Delete(&domain.Crosshair{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("crosshair not found or not owned by user")
	}
	return nil
}
