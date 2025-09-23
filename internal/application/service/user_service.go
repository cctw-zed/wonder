package service

import (
	"context"
	"time"

	"github.com/cctw-zed/wonder/internal/domain/user"
	"github.com/cctw-zed/wonder/pkg/errors"
	"github.com/cctw-zed/wonder/pkg/logger"
	"github.com/cctw-zed/wonder/pkg/snowflake/id"
)

type userService struct {
	repo  user.UserRepository
	idGen id.Generator
}

func NewUserService(repo user.UserRepository, idGen id.Generator) user.UserService {
	return &userService{
		repo:  repo,
		idGen: idGen,
	}
}

func (s *userService) Register(ctx context.Context, email, name string) (*user.User, error) {
	appLogger := logger.NewApplicationLogger(logger.NewLogger())
	startTime := time.Now()

	// Log use case start
	appLogger.Info(ctx, "Starting user registration use case",
		logger.String("email", email),
		logger.String("name", name),
	)

	// Business rule validation
	appLogger.LogServiceCall(ctx, "UserService", "validateEmail", time.Now())
	if err := s.validateEmail(ctx, email); err != nil {
		appLogger.LogUseCase(ctx, "RegisterUser", startTime, false,
			logger.String("email", email),
			logger.String("error", err.Error()),
			logger.String("phase", "validation"),
		)
		return nil, err // Return domain validation error directly
	}

	// Check if email already exists
	appLogger.Debug(ctx, "Checking if email already exists",
		logger.String("email", email),
	)
	existingUser, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		appLogger.LogUseCase(ctx, "RegisterUser", startTime, false,
			logger.String("email", email),
			logger.String("error", err.Error()),
			logger.String("phase", "email_check"),
		)
		// Repository layer should return proper infrastructure errors
		return nil, err
	}
	if existingUser != nil {
		appLogger.LogUseCase(ctx, "RegisterUser", startTime, false,
			logger.String("email", email),
			logger.String("existing_user_id", existingUser.ID),
			logger.String("phase", "duplicate_check"),
		)
		return nil, errors.NewDuplicateEntryError("user", "email", email, existingUser.ID)
	}

	// Create user aggregate
	userID := s.idGen.Generate()
	u := &user.User{
		ID:        userID,
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	appLogger.Debug(ctx, "Created user aggregate",
		logger.String("user_id", userID),
		logger.String("email", email),
		logger.String("name", name),
	)

	// Validate the aggregate before persisting
	if err := u.Validate(ctx); err != nil {
		appLogger.LogUseCase(ctx, "RegisterUser", startTime, false,
			logger.String("user_id", userID),
			logger.String("email", email),
			logger.String("error", err.Error()),
			logger.String("phase", "aggregate_validation"),
		)
		return nil, err // Return domain validation error directly
	}

	// Persist the user
	appLogger.Debug(ctx, "Persisting user to repository",
		logger.String("user_id", userID),
	)
	if err := s.repo.Create(ctx, u); err != nil {
		appLogger.LogUseCase(ctx, "RegisterUser", startTime, false,
			logger.String("user_id", userID),
			logger.String("email", email),
			logger.String("error", err.Error()),
			logger.String("phase", "persistence"),
		)
		// Repository layer should return proper infrastructure errors
		return nil, err
	}

	// Log successful registration
	appLogger.LogUseCase(ctx, "RegisterUser", startTime, true,
		logger.String("user_id", userID),
		logger.String("email", email),
		logger.String("name", name),
	)

	appLogger.Info(ctx, "User registration completed successfully",
		logger.String("user_id", userID),
		logger.String("email", email),
		logger.Duration("duration", time.Since(startTime)),
	)

	return u, nil
}

func (s *userService) validateEmail(ctx context.Context, email string) error {
	appLogger := logger.NewApplicationLogger(logger.NewLogger())

	appLogger.Debug(ctx, "Validating email format",
		logger.String("email", email),
	)

	var validationErrors []string

	if email == "" {
		validationErrors = append(validationErrors, "email is required")
		appLogger.LogValidation(ctx, "email_required", false, validationErrors,
			logger.String("field", "email"),
		)
		return errors.NewRequiredFieldError("email", email)
	}

	u := &user.User{Email: email}
	if !u.IsEmailValid() {
		validationErrors = append(validationErrors, "invalid email format")
		appLogger.LogValidation(ctx, "email_format", false, validationErrors,
			logger.String("field", "email"),
			logger.String("value", email),
		)
		return errors.NewInvalidFormatError("email", email, "valid email address")
	}

	appLogger.LogValidation(ctx, "email_validation", true, nil,
		logger.String("field", "email"),
		logger.String("value", email),
	)

	return nil
}
