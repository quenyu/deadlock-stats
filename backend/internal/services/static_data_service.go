package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/quenyu/deadlock-stats/internal/clients/deadlockapi"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	ranksURL  = "https://assets.deadlock-api.com/v2/ranks"
	heroesURL = "https://assets.deadlock-api.com/v2/heroes"
)

type StaticDataService struct {
	logger         *zap.Logger
	Ranks          map[int]deadlockapi.RankV2
	Heroes         map[string]deadlockapi.HeroV2
	HeroesByHeroID map[int]deadlockapi.HeroV2
	mx             sync.RWMutex
}

func NewStaticDataService(logger *zap.Logger) *StaticDataService {
	return &StaticDataService{
		logger:         logger,
		Ranks:          make(map[int]deadlockapi.RankV2),
		Heroes:         make(map[string]deadlockapi.HeroV2),
		HeroesByHeroID: make(map[int]deadlockapi.HeroV2),
	}
}

func (s *StaticDataService) LoadStaticData() error {
	var g errgroup.Group

	g.Go(func() error {
		return s.loadRanks()
	})

	g.Go(func() error {
		return s.loadHeroes()
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (s *StaticDataService) loadRanks() error {
	data, err := s.fetchDataFromURL(ranksURL)
	if err != nil {
		return fmt.Errorf("failed to fetch ranks: %w", err)
	}

	return s.parseAndStoreRanks(data)
}

func (s *StaticDataService) loadHeroes() error {
	data, err := s.fetchDataFromURL(heroesURL)
	if err != nil {
		return fmt.Errorf("failed to fetch heroes: %w", err)
	}

	return s.parseAndStoreHeroes(data)
}

func (s *StaticDataService) fetchDataFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status code %d from %s", resp.StatusCode, url)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from %s: %w", url, err)
	}

	return data, nil
}

func (s *StaticDataService) parseAndStoreRanks(data []byte) error {
	var ranks []deadlockapi.RankV2
	if err := json.Unmarshal(data, &ranks); err != nil {
		return fmt.Errorf("failed to decode ranks JSON: %w", err)
	}

	s.mx.Lock()
	defer s.mx.Unlock()

	for _, rank := range ranks {
		s.Ranks[rank.Tier] = rank
	}

	s.logger.Info("Successfully loaded ranks from API", zap.Int("count", len(s.Ranks)))
	return nil
}

func (s *StaticDataService) parseAndStoreHeroes(data []byte) error {
	var heroes []deadlockapi.HeroV2
	if err := json.Unmarshal(data, &heroes); err != nil {
		return fmt.Errorf("failed to decode heroes JSON: %w", err)
	}

	s.mx.Lock()
	defer s.mx.Unlock()

	for _, hero := range heroes {
		s.Heroes[hero.ClassName] = hero
		s.HeroesByHeroID[hero.ID] = hero
	}

	s.logger.Info("Successfully loaded heroes", zap.Int("count", len(s.Heroes)))
	return nil
}

func (s *StaticDataService) GetRanksHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.Ranks)
}
