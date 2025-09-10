package services

import (
	"context"
	"time"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/repositories/v2"
)

type GeneralService interface {
	GetStatistik(ctx context.Context) (map[string]interface{}, error)
}

type generalService struct {
	repo repositories.GeneralRepository
}

func NewServiceGeneral(repo repositories.GeneralRepository) GeneralService {
	return &generalService{repo: repo}
}

func (s *generalService) GetStatistik(ctx context.Context) (map[string]interface{}, error) {
	currentYear := time.Now().Year()
	prevYear := currentYear - 1

	dokumenNow, _ := s.repo.GetTotalDokumenByYear(ctx, currentYear)
	dokumenPrev, _ := s.repo.GetTotalDokumenByYear(ctx, prevYear)

	pasienNow, _ := s.repo.GetTotalPasienByYear(ctx, currentYear)
	pasienPrev, _ := s.repo.GetTotalPasienByYear(ctx, prevYear)

	kasusNow, _ := s.repo.GetTotalKasusByYear(ctx, currentYear)
	kasusPrev, _ := s.repo.GetTotalKasusByYear(ctx, prevYear)

	aliMediaNow, _ := s.repo.GetTotalAlihMediaByYear(ctx, currentYear)
	aliMediaPrev, _ := s.repo.GetTotalAlihMediaByYear(ctx, prevYear)

	retensiNow, _ := s.repo.GetTotalRetensiByYear(ctx, currentYear)
	retensiPrev, _ := s.repo.GetTotalRetensiByYear(ctx, prevYear)

	pemusnahanNow, _ := s.repo.GetTotalPemusnahanByYear(ctx, currentYear)
	pemusnahanPrev, _ := s.repo.GetTotalPemusnahanByYear(ctx, prevYear)

	kasusTerbanyak, totalKasus, _ := s.repo.GetMostCommonKasus(ctx)

	kunjunganNow, _ := s.repo.GetTotalKunjunganByYear(ctx, currentYear)
	kunjunganPrev, _ := s.repo.GetTotalKunjunganByYear(ctx, prevYear)

	kasusList, _ := s.repo.GetKasusList(ctx, currentYear)

	percent := func(prev, now int) float64 {
		total := prev + now
		if total == 0 {
			return 0
		}
		return float64(now) / float64(total) * 100
	}

	return map[string]interface{}{
		"dokumen": map[string]interface{}{
			"total":    dokumenNow,
			"persen":   percent(dokumenPrev, dokumenNow),
			"previous": dokumenPrev,
		},
		"pasien": map[string]interface{}{
			"total":    pasienNow,
			"persen":   percent(pasienPrev, pasienNow),
			"previous": pasienPrev,
		},
		"kasus": map[string]interface{}{
			"total":    kasusNow,
			"persen":   percent(kasusPrev, kasusNow),
			"previous": kasusPrev,
			"list":     kasusList,
		},
		"kunjungan": map[string]interface{}{
			"total":    kunjunganNow,
			"persen":   percent(kunjunganPrev, kunjunganNow),
			"previous": kunjunganPrev,
		},
		"alih_media": map[string]interface{}{
			"total":    aliMediaNow,
			"persen":   percent(aliMediaPrev, aliMediaNow),
			"previous": aliMediaPrev,
		},
		"retensi": map[string]interface{}{
			"total":    retensiNow,
			"persen":   percent(retensiPrev, retensiNow),
			"previous": retensiPrev,
		},
		"pemusnahan": map[string]interface{}{
			"total":    pemusnahanNow,
			"persen":   percent(pemusnahanPrev, pemusnahanNow),
			"previous": pemusnahanPrev,
		},
		"kasus_terbanyak": map[string]interface{}{
			"nama":  kasusTerbanyak,
			"total": totalKasus,
		},
	}, nil
}
