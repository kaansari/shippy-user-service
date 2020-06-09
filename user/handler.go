package user

import (
	"errors"
	"fmt"
	"log"

	pb "github.com/kaansari/shippy-user-service/proto/auth"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

const topic = "user.created"

type Service struct {
	Repo         Repository
	TokenService Authable
}

func (srv *Service) Get(ctx context.Context, req *pb.User) (*pb.Response, error) {
	user, err := srv.Repo.Get(req.Id)

	res := &pb.Response{}
	res.User = user
	return res, err
}

func (srv *Service) GetAll(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	users, err := srv.Repo.GetAll()
	res := &pb.Response{}
	res.Users = users
	return res, err
}

func (srv *Service) Auth(ctx context.Context, req *pb.User) (*pb.Token, error) {
	pbToken := &pb.Token{}
	if len(req.Email) != 0 {
		log.Println("Logging in with:", req.Email, req.Password)
		user, err := srv.Repo.GetByEmail(req.Email)
		log.Println(user, err)

		// Compares our given password against the hashed password
		// stored in the database
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {

		}
		token, err := srv.TokenService.Encode(user)

		pbToken.Token = token
		pbToken.Valid = true

	}
	if len(req.Id) != 0 {

		log.Println("Logging in with:", req.Id)
		user, err := srv.Repo.Get(req.Id)
		log.Println(user, err)
		token, err := srv.TokenService.Encode(user)

		pbToken.Token = token
		pbToken.Valid = true

	}

	return pbToken, nil
}

func (srv *Service) Create(ctx context.Context, req *pb.User) (*pb.Response, error) {

	log.Println("Creating user: ", req)

	// Generates a hashed version of our password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		err = errors.New(fmt.Sprintf("error hashing password: %v", err))
	}

	req.Password = string(hashedPass)
	if err := srv.Repo.Create(req); err != nil {
		err = errors.New(fmt.Sprintf("error creating user: %v", err))
	}

	token, err := srv.TokenService.Encode(req)
	res := &pb.Response{}
	res.User = req
	res.Token = &pb.Token{Token: token}

	/*
		if err := srv.Publisher.Publish(ctx, req); err != nil {
			return errors.New(fmt.Sprintf("error publishing event: %v", err))
		}*/

	return res, err
}

func (srv *Service) ValidateToken(ctx context.Context, req *pb.Token) (*pb.Token, error) {

	// Decode token
	claims, err := srv.TokenService.Decode(req.Token)

	if claims.User.Id == "" {
		err = errors.New("invalid user")
	}

	validToken := &pb.Token{}
	if err != nil {
		validToken.Token = req.Token
		validToken.Valid = true
	}

	return validToken, err
}
