package entity_api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/request"
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/response"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type PartsSupplyRepository struct {
	http *http.Client
}

var _ interfaces.IPartsSupplyRepository = (*PartsSupplyRepository)(nil)

func NewPartsSupplyRepository() *PartsSupplyRepository {
	return &PartsSupplyRepository{http: http.DefaultClient}
}

func (s *PartsSupplyRepository) GetByID(ctx context.Context, id string) (*response.PartsSupplyResponse, error) {
	logger := logs.Logger()
	path := fmt.Sprintf(PARTS_SUPPLY_ID_ENDPOINT, id)
	resp, err := http.Get(path)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find vehicle by plate")
		return nil, err
	}
	defer resp.Body.Close()

	var partsSupply response.PartsSupplyResponse
	if err := json.NewDecoder(resp.Body).Decode(&partsSupply); err != nil {
		logger.Error().Err(err).Msg("failed to decode vehicle")
		return nil, err
	}
	return &partsSupply, nil
}

func (s *PartsSupplyRepository) Reserve(ctx context.Context, partsSupply []request.PartsSupplyRequest) error {
	logger := logs.Logger()
	if txn := newrelic.FromContext(ctx); txn != nil {
		startSegment := txn.StartSegment("PartsSupplyRepository.Reserve")
		defer startSegment.End()
	}

	payload, err := json.Marshal(partsSupply)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal parts supply reserve request")
		return err
	}

	body := io.NopCloser(bytes.NewReader(payload))
	resp, err := s.http.Post(PARTS_SUPPLY_RESERVE_ENDPOINT, "application/json", body)
	if err != nil {
		logger.Error().Err(err).Msg("failed to reserve parts supply")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logger.Error().Msgf("failed to reserve parts supply: status %d, response: %s", resp.StatusCode, string(bodyBytes))
		return fmt.Errorf("failed to reserve parts supply: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func (s *PartsSupplyRepository) Release(ctx context.Context, partsSupply []request.PartsSupplyRequest) error {
	logger := logs.Logger()
	if txn := newrelic.FromContext(ctx); txn != nil {
		startSegment := txn.StartSegment("PartsSupplyRepository.Release")
		defer startSegment.End()
	}

	payload, err := json.Marshal(partsSupply)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal parts supply release request")
		return err
	}

	body := io.NopCloser(bytes.NewReader(payload))
	resp, err := s.http.Post(PARTS_SUPPLY_RELEASE_ENDPOINT, "application/json", body)
	if err != nil {
		logger.Error().Err(err).Msg("failed to release parts supply")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logger.Error().Msgf("failed to release parts supply: status %d, response: %s", resp.StatusCode, string(bodyBytes))
		return fmt.Errorf("failed to release parts supply: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func (s *PartsSupplyRepository) WriteOff(ctx context.Context, partsSupply []request.PartsSupplyRequest) error {
	logger := logs.Logger()
	if txn := newrelic.FromContext(ctx); txn != nil {
		startSegment := txn.StartSegment("PartsSupplyRepository.WriteOff")
		defer startSegment.End()
	}

	payload, err := json.Marshal(partsSupply)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal parts supply writeoff request")
		return err
	}

	body := io.NopCloser(bytes.NewReader(payload))
	resp, err := s.http.Post(PARTS_SUPPLY_WRITEOFF_ENDPOINT, "application/json", body)
	if err != nil {
		logger.Error().Err(err).Msg("failed to writeoff parts supply")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logger.Error().Msgf("failed to writeoff parts supply: status %d, response: %s", resp.StatusCode, string(bodyBytes))
		return fmt.Errorf("failed to writeoff parts supply: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func (s *PartsSupplyRepository) AuthorizeReserve(ctx context.Context, partsSupply []request.PartsSupplyRequest) error {
	logger := logs.Logger()
	if txn := newrelic.FromContext(ctx); txn != nil {
		startSegment := txn.StartSegment("PartsSupplyRepository.AuthorizeReserve")
		defer startSegment.End()
	}

	payload, err := json.Marshal(partsSupply)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal parts supply authorize reserve request")
		return err
	}

	body := io.NopCloser(bytes.NewReader(payload))
	resp, err := s.http.Post(PARTS_SUPPLY_AUTHORIZE_RESERVE_ENDPOINT, "application/json", body)
	if err != nil {
		logger.Error().Err(err).Msg("failed to authorize reserve parts supply")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logger.Error().Msgf("failed to authorize reserve parts supply: status %d, response: %s", resp.StatusCode, string(bodyBytes))
		return fmt.Errorf("failed to authorize reserve parts supply: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
