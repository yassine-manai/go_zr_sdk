package contract

import (
	"context"
	"fmt"
	"net/http"

	internalhttp "github.com/yassine-manai/go_zr_sdk/internal/http"
	"github.com/yassine-manai/go_zr_sdk/internal/logger"
	"github.com/yassine-manai/go_zr_sdk/models"
)

// Service handles contract operations
type ContractService struct {
	httpClient *internalhttp.Client
	logger     logger.Logger
}

// NewContractService creates a new contract service
func NewContractService(httpClient *internalhttp.Client, log logger.Logger) *ContractService {
	return &ContractService{
		httpClient: httpClient,
		logger:     log,
	}
}

// CreateContract creates a new contract
func (s *ContractService) CreateContract(ctx context.Context, req models.ContractRequest) (*models.ContractDetail, error) {
	s.logger.Info("creating contract", logger.String("name", req.Name), logger.String("valid_from", req.ValidFrom), logger.String("valid_until", req.ValidUntil))

	// Convert to XML structure
	contractDetail := req.ToXML()
	var result models.ContractDetail

	// Execute request
	err := s.httpClient.DoXMLRequest(
		ctx,
		http.MethodPost,
		models.ContractCustomerMedia,
		&contractDetail,
		&result,
	)

	if err != nil {
		s.logger.Error("failed to create contract", logger.String("name", req.Name), logger.Error(err))
		return nil, err
	}

	s.logger.Info("contract created successfully", logger.Int("id", *result.Contract.ID), logger.String("name", result.Contract.Name))

	return &result, nil
}

// GetContract retrieves a contract by ID
func (s *ContractService) GetContractById(ctx context.Context, contractID int) (*models.ContractDetail, error) {
	s.logger.Info("getting contract", logger.Int("contract_id", contractID))

	path := fmt.Sprintf(models.ContractCustomerMediaByID, contractID)
	var result models.ContractDetail

	err := s.httpClient.DoXMLRequest(
		ctx,
		http.MethodGet, path,
		nil, &result,
	)

	if err != nil {
		s.logger.Error("failed to get contract", logger.Int("contract_id", contractID), logger.Error(err))
		return nil, err
	}

	s.logger.Info("contract retrieved successfully", logger.Int("contract_id", contractID), logger.String("name", result.Contract.Name))

	return &result, nil
}

func (s *ContractService) GetContractList(ctx context.Context) (*models.Contracts, error) {
	s.logger.Info("getting all contract List from zr")

	var result models.Contracts

	err := s.httpClient.DoXMLRequest(
		ctx,
		http.MethodGet, models.ContractCustomerMedia,
		nil, &result,
	)

	if err != nil {
		s.logger.Error("failed to get contracts", logger.Error(err))
		return nil, err
	}

	s.logger.Info("contracts retrieved successfully", logger.Int("Count", len(result.Contract)))

	return &result, nil
}

// UpdateContract updates an existing contract
func (s *ContractService) UpdateContract(ctx context.Context, req models.ContractRequest) (*models.ContractDetail, error) {
	s.logger.Info("updating contract", logger.Int("contract_id", *req.ID), logger.String("name", req.Name))

	// Convert to XML structure
	contractDetail := req.ToXML()
	var result models.ContractDetail

	// Build path with contract ID
	path := fmt.Sprintf(models.ContractCustomerMediaDetail, req.ID)

	// Execute request
	err := s.httpClient.DoXMLRequest(
		ctx,
		http.MethodPut,
		path,
		&contractDetail,
		&result,
	)

	if err != nil {
		s.logger.Error("failed to update contract", logger.Int("contract_id", *req.ID), logger.Error(err))
		return nil, err
	}

	s.logger.Info("contract updated successfully", logger.Int("contract_id", *req.ID), logger.String("name", result.Contract.Name))

	return &result, nil
}

// DeleteContract deletes a contract by ID
func (s *ContractService) DeleteContract(ctx context.Context, contractID int) error {
	s.logger.Info("deleting contract", logger.Int("contract_id", contractID))

	path := fmt.Sprintf(models.ContractCustomerMediaByID, contractID)

	// Execute DELETE request (no response body expected on success)
	err := s.httpClient.DoXMLRequest(
		ctx,
		http.MethodDelete,
		path,
		nil,
		nil, // No result expected on 200 OK
	)

	if err != nil {
		s.logger.Error("failed to delete contract", logger.Int("contract_id", contractID), logger.Error(err))
		return err
	}

	s.logger.Info("contract deleted successfully", logger.Int("contract_id", contractID))

	return nil
}
