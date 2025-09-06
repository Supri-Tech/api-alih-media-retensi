package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/repositories/v2"
)

type CronService interface {
	CheckAndProcessKunjungan(ctx context.Context) error
	ProcessKunjungan(ctx context.Context, id int) error
}

type cronService struct {
	kunjunganRepo repositories.KunjunganRepository
	kasusRepo     repositories.KasusRepository
	alihMediaRepo repositories.AlihMediaRepository
}

func NewCronService(
	kunjunganRepo repositories.KunjunganRepository,
	kasusRepo repositories.KasusRepository,
	alihMediaRepo repositories.AlihMediaRepository,
) CronService {
	return &cronService{
		kunjunganRepo: kunjunganRepo,
		kasusRepo:     kasusRepo,
		alihMediaRepo: alihMediaRepo,
	}
}

func (svc *cronService) CheckAndProcessKunjungan(ctx context.Context) error {
	log.Println("Starting cron job: Checking for kunjungan tidak aktif")

	activeKunjungan, err := svc.kunjunganRepo.GetActiveKunjungan(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	processedCount := 0

	for _, kunjungan := range activeKunjungan {
		kasus, err := svc.kasusRepo.GetKasusByID(ctx, kunjungan.IDKasus)
		if err != nil {
			log.Printf("Error getting kasus for kunjungan %d: %v", kunjungan.ID, err)
			continue
		}

		var masaInaktif int
		if kunjungan.JenisKunjungan == "RI" {
			masaInaktif = kasus.MasaInaktifRI
		} else if kunjungan.JenisKunjungan == "RJ" {
			masaInaktif = kasus.MasaInaktifRJ
		} else {
			continue
		}

		expired := kunjungan.TanggalMasuk.AddDate(masaInaktif, 0, 0)

		if now.After(expired) {
			_, err := svc.processInactiveKunjungan(ctx, kunjungan)
			if err != nil {
				log.Printf("Error processing kunjungan %d: %v", kunjungan.ID, err)
				continue
			}
			processedCount++
		}
	}

	log.Printf("Cron job completed. Processed %d kunjungen", processedCount)
	return nil
}

func (svc *cronService) ProcessKunjungan(ctx context.Context, id int) error {
	kunjungan, err := svc.kunjunganRepo.GetKunjunganBasicByID(ctx, id)
	if err != nil {
		return err
	}
	if kunjungan == nil {
		return fmt.Errorf("kunjungan with ID %d not found", id)
	}

	if kunjungan.Status == "tidak aktif" {
		log.Printf("Kunjungan %d is already inactive", id)
		return nil
	}

	kasus, err := svc.kasusRepo.GetKasusByID(ctx, kunjungan.IDKasus)
	if err != nil {
		return err
	}

	var masaInaktif int
	if kunjungan.JenisKunjungan == "RI" {
		masaInaktif = kasus.MasaInaktifRI
	} else if kunjungan.JenisKunjungan == "RJ" {
		masaInaktif = kasus.MasaInaktifRJ
	} else {
		return fmt.Errorf("unknown jenis kunjungan: %s", kunjungan.JenisKunjungan)
	}

	expired := kunjungan.TanggalMasuk.AddDate(masaInaktif, 0, 0)
	now := time.Now()

	if now.After(expired) {
		_, err := svc.processInactiveKunjungan(ctx, kunjungan)
		if err != nil {
			return err
		}
		log.Printf("Successfully processed kunjungan %d as inactive", id)
	} else {
		log.Printf("Kunjungan %d is not yet expired (expires on %v)", id, expired)
	}

	return nil
}

func (svc *cronService) processInactiveKunjungan(ctx context.Context, kunjungan *models.Kunjungan) (*models.AlihMedia, error) {
	err := svc.kunjunganRepo.UpdateKunjunganStatus(ctx, kunjungan.ID, "tidak aktif")
	if err != nil {
		return nil, err
	}

	alihMedia := &models.AlihMedia{
		ID:         kunjungan.ID,
		TglLaporan: nil,
		Status:     "belum di alih media",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	createdAlihMedia, err := svc.alihMediaRepo.CreateAlihMedia(ctx, alihMedia)
	if err != nil {
		return nil, err
	}

	return createdAlihMedia, nil
}
