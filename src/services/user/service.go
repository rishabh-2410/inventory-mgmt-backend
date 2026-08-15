package services

import (
	"errors"
	"inventory-backend/src/auth"
	dbmodel "inventory-backend/src/models/db_object"
	enum "inventory-backend/src/models/enums"
	reqmodels "inventory-backend/src/models/req_dto"
	resmodels "inventory-backend/src/models/res_dto"
	tokenrepo "inventory-backend/src/repository/token"
	userrepo "inventory-backend/src/repository/user"
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo  *userrepo.UserRepo
	tokenRepo *tokenrepo.TokenRepo
}

func New(userRepo *userrepo.UserRepo, tokenRepo *tokenrepo.TokenRepo) *Service {
	return &Service{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (us *Service) RegisterUser(req *reqmodels.RegisterUserRequestDTO) (*resmodels.RegisterUserResponseDTO, error) {
	slog.Debug("registering user...", "email", req.Email, "name", req.Name)

	existingUser, err := us.userRepo.FindByEmail(req.Email)
	if err != nil && err.Error() != "sql: no rows in result set" {
		slog.Error("Unable to check db for existing user ", req.Email, err)
		return nil, err
	}

	if existingUser != nil {
		if existingUser.BusinessID == nil {
			return &resmodels.RegisterUserResponseDTO{ID: existingUser.ID}, nil // user is registered but not yet added to a business, so we can redirect them to the business registration page
		} else {
			slog.Info("register user skipped ", "email", req.Email)
			return nil, errors.New("email already exists")
		}
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash passwd", "email", req.Email, "err", err)
		return nil, err
	}

	isActive := true
	role := enum.RoleOwner
	newUser := &dbmodel.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hashBytes),
		Role:         &role,
		IsActive:     isActive,
	}

	userID, err := us.userRepo.AddUser(newUser)
	if err != nil {
		slog.Error("failed to add user", "email", req.Email, "err", err)
		return nil, err
	}

	slog.Info("user registered successfully", "email", req.Email)

	return &resmodels.RegisterUserResponseDTO{
		ID: userID,
	}, nil

}

func (us *Service) LoginUser(req *reqmodels.LoginUserRequestDTO) (*resmodels.LoginUserResponseDTO, error) {
	slog.Debug("login user service start", "email", req.Email)
	user, err := us.userRepo.FindByEmail(req.Email)
	if err != nil {
		slog.Error("login user lookup failed", "email", req.Email, "err", err)
		return nil, err
	}

	if user.IsActive == false {
		slog.Error("login user is inactive", "user_id", user.ID)
		return nil, errors.New("user is inactive")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		slog.Error("login user password mismatch", "email", req.Email, "user_id", user.ID, "err", err)
		return nil, err
	}

	if user.BusinessID == nil {
		slog.Error("login user business id is nil", "user_id", user.ID)
		return nil, errors.New("business id is required")
	}

	if user.Role == nil {
		slog.Error("login user role is nil", "user_id", user.ID)
		return nil, errors.New("role is required")
	}

	accessToken, accessTokenExpiresAt, err := auth.GenerateAccessToken(user.ID, *user.BusinessID, *user.Role)
	if err != nil {
		slog.Error("login user access token generation failed", "user_id", user.ID, "err", err)
		return nil, err
	}
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		slog.Error("login user refresh token generation failed", "user_id", user.ID, "err", err)
		return nil, err
	}

	hashedRefreshToken := auth.HashToken(refreshToken)

	err = us.tokenRepo.SaveRefreshToken(user.ID, hashedRefreshToken)
	if err != nil {
		slog.Error("login user refresh token save failed", "user_id", user.ID, "err", err)
		return nil, err
	}

	slog.Info("login user service success", "user_id", user.ID, "email", user.Email)

	return &resmodels.LoginUserResponseDTO{
		ID:                   user.ID,
		Name:                 user.Name,
		Email:                user.Email,
		CreatedAt:            user.CreatedAt,
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenExpiresAt,
		RefreshToken:         refreshToken,
		IsActive:             user.IsActive,
		BusinessID:           user.BusinessID,
		Role:                 user.Role,
	}, nil

}

func (us *Service) LogoutUser(refreshToken string) error {
	slog.Debug("logout user service start", "refresh_token", refreshToken)
	refreshTokenHash := auth.HashToken(refreshToken)

	err := us.userRepo.RevokeUser(refreshTokenHash)
	if err != nil {
		slog.Error("logout user revoke failed err=%v", err)
		return err
	}

	slog.Info("logout user service success")

	return nil
}

func (us *Service) RefreshToken(existingrefreshToken string) (*resmodels.LoginUserResponseDTO, error) {
	slog.Debug("refresh token service start", "refresh_token", existingrefreshToken)
	refreshTokenHash := auth.HashToken(existingrefreshToken)

	userID, err := us.tokenRepo.MarkTokenExpired(refreshTokenHash)
	if err != nil {
		slog.Error("refresh token mark expired failed err=%v", err)
		return nil, err
	}

	user, err := us.userRepo.FindByID(userID)
	if err != nil {
		slog.Error("refresh token user lookup failed user_id=%s err=%v", userID, err)
		return nil, err
	}

	accessToken, accessTokenExpiresAt, err := auth.GenerateAccessToken(user.ID, *user.BusinessID, *user.Role)
	if err != nil {
		slog.Error("refresh token access token generation failed user_id=%s err=%v", userID, err)
		return nil, err
	}

	newRefreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		slog.Error("refresh token generation failed user_id=%s err=%v", userID, err)
		return nil, err
	}

	hashedRefreshToken := auth.HashToken(newRefreshToken)

	err = us.tokenRepo.SaveRefreshToken(userID, hashedRefreshToken)
	if err != nil {
		slog.Error("refresh token save failed user_id=%s err=%v", userID, err)
		return nil, err
	}

	slog.Info("refresh token service success user_id=%s email=%s", user.ID, user.Email)

	return &resmodels.LoginUserResponseDTO{
		ID:                   user.ID,
		Name:                 user.Name,
		Email:                user.Email,
		CreatedAt:            user.CreatedAt,
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenExpiresAt,
		RefreshToken:         newRefreshToken,
		IsActive:             user.IsActive,
		BusinessID:           user.BusinessID,
		Role:                 user.Role,
	}, nil
}

func (us *Service) CreateUser(businessID uuid.UUID, req *reqmodels.RegisterUserRequestDTO) (*resmodels.CreateUserResponseDTO, error) {
	slog.Debug("creating employee...", "email", req.Email, "name", req.Name)

	existingUser, err := us.userRepo.FindByEmail(req.Email)
	if err != nil && err.Error() != "sql: no rows in result set" {
		slog.Error("Unable to check db for existing user ", req.Email, err)
		return nil, err
	}

	if existingUser != nil {
		slog.Info("employee email already exists ", "email", req.Email)
		return nil, errors.New("employee already exists")
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash passwd", "email", req.Email, "err", err)
		return nil, err
	}

	isActive := true
	role := enum.RoleEmployee
	newUser := &dbmodel.User{
		Email:        req.Email,
		Name:         req.Name,
		BusinessID:   &businessID,
		PasswordHash: string(hashBytes),
		Role:         &role,
		IsActive:     isActive,
	}

	employee, err := us.userRepo.AddEmployee(newUser)
	if err != nil {
		slog.Error("failed to add employee", "email", req.Email, "err", err)
		return nil, err
	}

	slog.Info("employee created successfully", "email", req.Email)

	return &resmodels.CreateUserResponseDTO{
		ID:         employee.ID,
		Email:      employee.Email,
		Name:       employee.Name,
		BusinessID: employee.BusinessID,
		Role:       employee.Role,
		IsActive:   employee.IsActive,
		CreatedAt:  employee.CreatedAt,
		UpdatedAt:  employee.UpdatedAt,
	}, nil

}

func (us *Service) DeleteUser(employeeID string) error {
	slog.Debug("deactivating employee...", "employee_id", employeeID)

	err := us.userRepo.DeleteUser(employeeID)
	if err != nil {
		slog.Error("failed to deactivate employee", "employee_id", employeeID, "err", err)
		return err
	}

	slog.Info("employee deactivated successfully", "employee_id", employeeID)
	return nil
}

func (us *Service) GetUsers(businessID uuid.UUID) ([]resmodels.GetUsersResponseDTO, error) {
	slog.Debug("getting users...", "business_id", businessID)

	response := make([]resmodels.GetUsersResponseDTO, 0)

	users, err := us.userRepo.FindByBusinessIDAndRole(businessID, enum.RoleEmployee)
	if err != nil {
		slog.Error("failed to get users", "business_id", businessID, "err", err)
		return nil, err
	}

	for _, user := range users {
		response = append(response, resmodels.GetUsersResponseDTO{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			BusinessID: user.BusinessID,
			Role:       user.Role,
			IsActive:   user.IsActive,
			CreatedAt:  user.CreatedAt,
			UpdatedAt:  user.UpdatedAt,
		})
	}

	return response, nil
}
