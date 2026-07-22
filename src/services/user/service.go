package services

import (
	"errors"
	dbmodel "inventory-backend/src/models/db_object"
	enum "inventory-backend/src/models/enums"
	reqmodels "inventory-backend/src/models/req_dto"
	tokenrepo "inventory-backend/src/repository/token"
	userrepo "inventory-backend/src/repository/user"
	"log/slog"

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

func (us *Service) RegisterUser(req *reqmodels.RegisterUserRequestDTO) error {
	slog.Debug("registering user...", "email", req.Email, "name", req.Name)

	existingUser, err := us.userRepo.FindByEmail(req.Email)
	if err != nil && err.Error() != "sql: no rows in result set" {
		slog.Error("Unable to check db for existing user ", req.Email, err)
		return err
	}

	if existingUser != nil {
		slog.Info("register user skipped ", "email", req.Email)
		return errors.New("Email already exists")
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash passwd", "email", req.Email, "err", err)
		return err
	}

	role := enum.RoleOwner
	newUser := &dbmodel.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hashBytes),
		Role:         &role,
	}

	err = us.userRepo.AddUser(newUser)
	if err != nil {
		slog.Error("failed to add user", "email", req.Email, "err", err)
		return err
	}

	slog.Info("user registered successfully", "email", req.Email)

	return nil

}

// func (us *Service) LoginUser(req *reqmodels.LoginUserRequestDTO) (*resmodels.LoginUserResponseDTO, error) {
// 	logging.Debug("login user service start email=%s", req.Email)
// 	user, err := us.userRepo.FindByEmail(req.Email)
// 	if err != nil {
// 		logging.Error("login user lookup failed email=%s err=%v", req.Email, err)
// 		return nil, err
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
// 	if err != nil {
// 		logging.Error("login user password mismatch email=%s user_id=%s err=%v", req.Email, user.ID, err)
// 		return nil, err
// 	}

// 	accessToken, err := auth.GenerateAccessToken(user.ID)
// 	if err != nil {
// 		logging.Error("login user access token generation failed user_id=%s err=%v", user.ID, err)
// 		return nil, err
// 	}

// 	refreshToken, err := auth.GenerateRefreshToken()
// 	if err != nil {
// 		logging.Error("login user refresh token generation failed user_id=%s err=%v", user.ID, err)
// 		return nil, err
// 	}

// 	hashedRefreshToken := auth.HashToken(refreshToken)

// 	err = us.tokenRepo.SaveRefreshToken(user.ID, hashedRefreshToken)
// 	if err != nil {
// 		logging.Error("login user refresh token save failed user_id=%s err=%v", user.ID, err)
// 		return nil, err
// 	}

// 	logging.Info("login user service success user_id=%s email=%s", user.ID, user.Email)

// 	return &resmodels.LoginUserResponseDTO{
// 		ID:           user.ID,
// 		Name:         user.Name,
// 		Email:        user.Email,
// 		CreatedAt:    user.CreatedAt,
// 		AccessToken:  accessToken,
// 		RefreshToken: refreshToken,
// 	}, nil

// }

// func (us *Service) LogoutUser(refreshToken string) error {
// 	logging.Debug("logout user service start refresh_token=%s", logging.Redacted("refresh_token"))
// 	refreshTokenHash := auth.HashToken(refreshToken)

// 	err := us.userRepo.RevokeUser(refreshTokenHash)
// 	if err != nil {
// 		logging.Error("logout user revoke failed err=%v", err)
// 		return err
// 	}

// 	logging.Info("logout user service success")

// 	return nil
// }

// func (us *Service) RefreshToken(existingrefreshToken string) (*resmodels.LoginUserResponseDTO, error) {
// 	logging.Debug("refresh token service start refresh_token=%s", logging.Redacted("refresh_token"))
// 	refreshTokenHash := auth.HashToken(existingrefreshToken)

// 	userID, err := us.tokenRepo.MarkTokenExpired(refreshTokenHash)
// 	if err != nil {
// 		logging.Error("refresh token mark expired failed err=%v", err)
// 		return nil, err
// 	}

// 	user, err := us.userRepo.FindById(userID)
// 	if err != nil {
// 		logging.Error("refresh token user lookup failed user_id=%s err=%v", userID, err)
// 		return nil, err
// 	}

// 	accessToken, err := auth.GenerateAccessToken(userID)
// 	if err != nil {
// 		logging.Error("refresh token access token generation failed user_id=%s err=%v", userID, err)
// 		return nil, err
// 	}

// 	newRefreshToken, err := auth.GenerateRefreshToken()
// 	if err != nil {
// 		logging.Error("refresh token generation failed user_id=%s err=%v", userID, err)
// 		return nil, err
// 	}

// 	hashedRefreshToken := auth.HashToken(newRefreshToken)

// 	err = us.tokenRepo.SaveRefreshToken(userID, hashedRefreshToken)
// 	if err != nil {
// 		logging.Error("refresh token save failed user_id=%s err=%v", userID, err)
// 		return nil, err
// 	}

// 	logging.Info("refresh token service success user_id=%s email=%s", user.ID, user.Email)

// 	return &resmodels.LoginUserResponseDTO{
// 		ID:           user.ID,
// 		Name:         user.Name,
// 		Email:        user.Email,
// 		CreatedAt:    user.CreatedAt,
// 		AccessToken:  accessToken,
// 		RefreshToken: newRefreshToken,
// 	}, nil
// }
