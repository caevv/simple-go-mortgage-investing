package repository

import (
	"context"

	"github.com/caevv/simple-go-prepaid-card/data"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type CardRepository interface {
	TopUp(ctx context.Context, id string, amount int64) (*data.Card, error)
}

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) (CardRepository, error) {
	return &Repository{db: db}, nil
}

func (r *Repository) TopUp(ctx context.Context, id string, amount int64) (*data.Card, error) {
	var card data.Card

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&card).Where("id = ?", id).Error; err != nil {
			return errors.Wrap(err, "failed to fetch card")
		}

		card.Amount += amount

		if err := tx.Save(card).Error; err != nil {
			return errors.Wrap(err, "failed to update amount")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &card, nil
}
