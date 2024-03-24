package service

import (
	"context"
	"fmt"
	"time"

	"mygram/model"
	"mygram/pkg/helper"
	"mygram/repository"
)

type UserService interface {
	GetUsers(ctx context.Context) ([]model.User, error)
	GetUsersById(ctx context.Context, id uint64) (model.User, error)
	DeleteUsersById(ctx context.Context, id uint64) (model.User, error)
	UpdateUserByID(ctx context.Context, id uint64, user model.User) (model.User, error)
	GetUsersByUsername(ctx context.Context, username string) (model.User, error)

	// activity
	SignUp(ctx context.Context, userSignUp model.UserSignUp) (model.User, error)
	SignIn(ctx context.Context, userSignIn model.UserSignIn) (model.User, error)

	// misc
	GenerateUserAccessToken(ctx context.Context, user model.User) (token string, err error)
}

type userServiceImpl struct {
	repo repository.UserQuery
}

func NewUserService(repo repository.UserQuery) UserService {
	return &userServiceImpl{repo: repo}
}

func (u *userServiceImpl) GetUsersByUsername(ctx context.Context, email string) (model.User, error) {
	user, err := u.repo.GetUsersByUsername(ctx, email)
	if err != nil {
		return model.User{}, err
	}
	return user, err
}

func (u *userServiceImpl) GetUsers(ctx context.Context) ([]model.User, error) {
	users, err := u.repo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, err
}

func (u *userServiceImpl) GetUsersById(ctx context.Context, id uint64) (model.User, error) {
	user, err := u.repo.GetUsersByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	return user, err
}

func (u *userServiceImpl) DeleteUsersById(ctx context.Context, id uint64) (model.User, error) {
	user, err := u.repo.GetUsersByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	// if user doesn't exist, return
	if user.ID == 0 {
		return model.User{}, nil
	}

	// delete user by id
	err = u.repo.DeleteUsersByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	return user, err
}

func (u *userServiceImpl) UpdateUserByID(ctx context.Context, id uint64, user model.User) (model.User, error) {
	// Get user by ID
	existingUser, err := u.repo.GetUsersByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	// Check if user exists
	if existingUser.ID == 0 {
		return model.User{}, fmt.Errorf("user with ID %d not found", id)
	}

	// Update user fields
	existingUser.Username = user.Username
	existingUser.Email = user.Email
	existingUser.DoB = user.DoB

	// Save updated user to repository
	updatedUser, err := u.repo.UpdateUserByID(ctx, id, existingUser)
	if err != nil {
		return model.User{}, err
	}

	return updatedUser, nil
}

func (u *userServiceImpl) SignUp(ctx context.Context, userSignUp model.UserSignUp) (model.User, error) {
	// assumption: semua user adalah user baru
	user := model.User{
		Username: userSignUp.Username,
		Email:    userSignUp.Email,
		DoB:      userSignUp.DoB,
		// FirstName: userSignUp.FirstName,
		// LastName:  userSignUp.LastName,
	}

	// encryption password
	// hashing
	pass, err := helper.GenerateHash(userSignUp.Password)
	if err != nil {
		return model.User{}, err
	}
	user.Password = pass

	// store to db
	res, err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return model.User{}, err
	}
	return res, err
}

func (u *userServiceImpl) SignIn(ctx context.Context, userSignIn model.UserSignIn) (model.User, error) {
	// Get user by username
	user, err := u.repo.GetUsersByUsername(ctx, userSignIn.Email)
	if err != nil {
		return model.User{}, err
	}

	// Check if user exists
	if user.ID == 0 {
		return model.User{}, fmt.Errorf("user with username %s not found", userSignIn.Email)
	}

	// Check if password matches
	match, err := helper.CompareHash(userSignIn.Password, user.Password)
	if err != nil {
		return model.User{}, err
	}
	if !match {
		return model.User{}, fmt.Errorf("invalid password")
	}

	return user, nil
}

func (u *userServiceImpl) GenerateUserAccessToken(ctx context.Context, user model.User) (token string, err error) {
	// generate claim
	now := time.Now()

	claim := model.StandardClaim{
		Jti: fmt.Sprintf("%v", time.Now().UnixNano()),
		Iss: "project",
		Aud: "mygram",
		Sub: "access-token",
		Exp: uint64(now.Add(time.Hour).Unix()),
		Iat: uint64(now.Unix()),
		Nbf: uint64(now.Unix()),
	}

	userClaim := model.AccessClaim{
		StandardClaim: claim,
		UserID:        user.ID,
		Username:      user.Username,
		Dob:           user.DoB,
	}

	token, err = helper.GenerateToken(userClaim)
	return
}
