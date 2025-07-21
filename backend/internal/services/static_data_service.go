package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/quenyu/deadlock-stats/internal/clients/deadlockapi"
	"go.uber.org/zap"
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
	var wg sync.WaitGroup
	errs := make(chan error, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := s.loadRanks(); err != nil {
			errs <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := s.loadHeroes(); err != nil {
			errs <- err
		}
	}()

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *StaticDataService) loadRanks() error {
	resp, err := http.Get(ranksURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var ranks []deadlockapi.RankV2
	if err := json.NewDecoder(resp.Body).Decode(&ranks); err != nil {
		return err
	}

	s.mx.Lock()
	defer s.mx.Unlock()
	for _, rank := range ranks {
		s.Ranks[rank.Tier] = rank
	}
	s.logger.Info("Successfully loaded ranks", zap.Int("count", len(s.Ranks)))
	return nil
}

func (s *StaticDataService) loadHeroes() error {
	resp, err := http.Get(heroesURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var heroes []deadlockapi.HeroV2
	if err := json.NewDecoder(resp.Body).Decode(&heroes); err != nil {
		return err
	}

	s.mx.Lock()
	defer s.mx.Unlock()
	for _, hero := range heroes {
		s.Heroes[hero.ClassName] = hero
		s.HeroesByHeroID[hero.ID] = hero
		key := fmt.Sprintf("HeroID_%d", hero.ID)
		s.Heroes[key] = hero
	}
	s.logger.Info("Successfully loaded heroes", zap.Int("count", len(s.Heroes)))
	return nil
}
