package repositories

import (
	"errors"

	"github.com/quenyu/deadlock-stats/internal/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindBySteamID(steamID string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("steam_id = ?", steamID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) FindOrCreate(user *domain.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		query := `
			INSERT INTO users (id, steam_id, nickname, avatar_url, profile_url, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7) 
			ON CONFLICT (steam_id) 
			DO UPDATE SET nickname = $3, avatar_url = $4, profile_url = $5, updated_at = $7 
			RETURNING id
		`
		if err := tx.Raw(query, user.ID, user.SteamID, user.Nickname, user.AvatarURL, user.ProfileURL, user.CreatedAt, user.UpdatedAt).Scan(user).Error; err != nil {
			return err
		}

		statsQuery := `
			INSERT INTO player_stats (user_id) 
			VALUES ($1) 
			ON CONFLICT (user_id) 
			DO NOTHING
		`
		if err := tx.Exec(statsQuery, user.ID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *UserRepository) FindByID(id string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, steam_id, nickname, avatar_url, profile_url, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.Raw(query, id).Scan(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
