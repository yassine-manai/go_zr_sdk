package contract

import (
	"context"
	"net/http"

	internalhttp "github.com/yourcompany/thirdparty-sdk/internal/http"
	"github.com/yourcompany/thirdparty-sdk/internal/retry"
	"github.com/yourcompany/thirdparty-sdk/logger"
	"github.com/yourcompany/thirdparty-sdk/models"
)

// Service handles contract operations
type Service struct {
	httpClient *internalhttp.Client
	retryer    *retry.Retryer
	logger     logger.Logger
}

// New creates a new contract service
func New(httpClient *internalhttp.Client, retryer *retry.Retryer, log logger.Logger) *Service {
	return &Service{
		httpClient: httpClient,
		retryer:    retryer,
		logger:     log,
	}
}

// CreateContract creates a new contract
func (s *Service) CreateContract(ctx context.Context, req models.CreateContractRequest) (*models.ContractDetail, error) {
	s.logger.Info("creating contract",
		logger.String("name", req.Name),
		logger.String("valid_from", req.ValidFrom),
		logger.String("valid_until", req.ValidUntil),
	)

	// Convert to XML structure
	contractDetail := req.ToXML()

	var result models.ContractDetail

	// Execute with retry
	err := s.retryer.Do(ctx, func() error {
		return s.httpClient.DoXML(
			ctx,
			http.MethodPost,
			"/contracts",
			&contractDetail,
			&result,
		)
	})

	if err != nil {
		s.logger.Error("failed to create contract",
			logger.String("name", req.Name),
			logger.Error(err),
		)
		return nil, err
	}

	s.logger.Info("contract created successfully",
		logger.Int("id", *result.Contract.ID),
		logger.String("name", result.Contract.Name),
	)

	return &result, nil
}
