package participant

import (
	"context"
	"net/http"

	internalhttp "github.com/yassine-manai/go_zr_sdk/internal/http"
	"github.com/yassine-manai/go_zr_sdk/internal/logger"
	"github.com/yassine-manai/go_zr_sdk/models"
)

// Service handles contract operations
type ParticipantService struct {
	httpClient *internalhttp.Client
	logger     logger.Logger
}

// New creates a new contract service
func NewParticipantService(httpClient *internalhttp.Client, log logger.Logger) *ParticipantService {
	return &ParticipantService{
		httpClient: httpClient,
		logger:     log,
	}
}

// CreateContract creates a new contract
func (s *ParticipantService) CreateParticipant(ctx context.Context, req models.ContractRequest) (*models.ContractDetail, error) {
	s.logger.Info("creating contract", logger.String("name", req.Name), logger.String("valid_from", req.ValidFrom), logger.String("valid_until", req.ValidUntil))

	// Convert to XML structure
	contractDetail := req.ToXML()

	var result models.ContractDetail

	// Execute with retry
	err := s.httpClient.DoXMLRequest(ctx, http.MethodPost, models.ContractCustomerMedia, &contractDetail, &result)

	if err != nil {
		s.logger.Error("failed to create contract", logger.String("name", req.Name), logger.Error(err))
		return nil, err
	}

	s.logger.Info("contract created successfully", logger.Int("id", *result.Contract.ID), logger.String("name", result.Contract.Name))

	return &result, nil
}
