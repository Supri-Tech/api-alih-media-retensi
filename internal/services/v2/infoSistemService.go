package services

import (
	"context"
	"errors"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/repositories/v2"
)

type InfoSistemService interface {
	GetAllInfoSistem(ctx context.Context, infoSistem models.InfoSistem) (*models.InfoSistem, error)
	GetInfoSistem(ctx context.Context, id int) (*models.InfoSistem, error)
	CreateInfoSistem(ctx context.Context, infoSistem models.InfoSistem) (*models.InfoSistem, error)
	UpdateInfoSistem(ctx context.Context, infoSistem models.InfoSistem) (*models.InfoSistem, error)
}

type infoSistemService struct {
	repo repositories.InfoSistemRepository
}

func NewServiceInfoSistem(repo repositories.InfoSistemRepository) InfoSistemService {
	return &infoSistemService{repo: repo}
}

func (svc *infoSistemService) GetAllInfoSistem(ctx context.Context, infoSistem models.InfoSistem) (*models.InfoSistem, error) {
	info, err := svc.repo.GetAllInfoSistem(ctx, infoSistem)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (svc *infoSistemService) GetInfoSistem(ctx context.Context, id int) (*models.InfoSistem, error) {
	info, err := svc.repo.GetInfoSistem(ctx, id)
	if err != nil {
		return nil, err
	}

	if info == nil {
		return nil, errors.New("Info not found")
	}

	return info, nil
}

func (svc *infoSistemService) CreateInfoSistem(ctx context.Context, infoSistem models.InfoSistem) (*models.InfoSistem, error) {
	newInfo, err := svc.repo.CreateInfoSistem(ctx, infoSistem)
	if err != nil {
		return nil, err
	}

	return newInfo, nil
}

func (svc *infoSistemService) UpdateInfoSistem(ctx context.Context, infoSistem models.InfoSistem) (*models.InfoSistem, error) {
	existing, err := svc.repo.GetInfoSistem(ctx, infoSistem.ID)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("Info not found")
	}

	newInfo, err := svc.repo.UpdateInfoSistem(ctx, infoSistem)
	if err != nil {
		return nil, err
	}

	return newInfo, nil
}
