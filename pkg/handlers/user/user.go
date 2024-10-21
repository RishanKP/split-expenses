package user

import (
	"context"
	"net/http"

	"split-expenses/library/api"
	"split-expenses/library/jwt"
	"split-expenses/library/utils"
	"split-expenses/pkg/interfaces"
	"split-expenses/pkg/repository"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler interface {
	SignUp(c *gin.Context)
	Login(c *gin.Context)
	GetUserById(c *gin.Context)
}

type handler struct {
	repo repository.UserRepository
}

func (h handler) SignUp(c *gin.Context) {
	var req interfaces.UserCreationRequest
	err := c.BindJSON(&req)

	if err != nil {
		log.Err(err).Msg("unable to unmarshal request")
		api.NewClientError(c, http.StatusBadRequest, "invalid request")
		return
	}

	user, _ := h.repo.GetByEmail(context.TODO(), req.Email)
	if user.Email != "" {
		api.NewClientError(c, http.StatusBadRequest, "email id exists")
		return
	}

	user = req.AsUser()
	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		log.Err(err).Msg("error hashing password")
		api.NewInternalError(c, http.StatusInternalServerError, "server error")
		return
	}

	err = h.repo.Create(context.TODO(), user)
	if err != nil {
		log.Err(err).Msg("error creating user")
		api.NewInternalError(c, http.StatusInternalServerError, "error creating user")
		return
	}

	api.Result(c, http.StatusCreated, "success", struct{}{})
}

func (h handler) Login(c *gin.Context) {

	var req interfaces.LoginCredentials

	err := c.BindJSON(&req)
	if err != nil {
		api.NewClientError(c, http.StatusBadRequest, "invalid request")
	}

	user, err := h.repo.GetByEmail(context.TODO(), req.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			api.NewClientError(c, http.StatusBadRequest, "account not found")
			return
		}
		api.NewInternalError(c, http.StatusInternalServerError, "failed to fetch details")
		return
	}

	if !utils.ComparePassword(req.Password, user.Password) {
		api.NewClientError(c, http.StatusUnauthorized, "invalid password")
		return
	}

	token, err := jwt.CreateToken(jwt.Claims{
		Username: user.Email,
		UserId:   user.ID.Hex(),
		Role:     "user",
	})

	if err != nil {
		api.NewInternalError(c, http.StatusInternalServerError, "failed to generate token")
		return
	}

	res := interfaces.UserLoginResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Token:     token,
		Contact:   user.Contact,
		Id:        user.ID.Hex(),
	}

	api.Result(c, http.StatusOK, "authorized", res)
}

func (h handler) GetUserById(c *gin.Context) {
	user, err := h.repo.GetById(context.TODO(), c.Param("id"))
	if err != nil {
		log.Err(err).Msg("failed to get user details")
		api.NewClientError(c, http.StatusBadRequest, "failed")
		return
	}

	api.Result(c, http.StatusOK, "success", user)
	return
}

func Newhandler(repo repository.UserRepository) UserHandler {
	return handler{
		repo: repo,
	}
}
