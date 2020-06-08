package main

import (
	"strings"
	"testing"

	pb "github.com/kaansari/shippy-user-service/proto/auth"
	"github.com/kaansari/shippy-user-service/user"
)

var (
	t_user = &pb.User{
		Id:    "abc123",
		Email: "ewan.valentine89@gmail.com",
	}
)

type MockRepo struct{}

func (repo *MockRepo) GetAll() ([]*pb.User, error) {
	var users []*pb.User
	return users, nil
}

func (repo *MockRepo) Get(id string) (*pb.User, error) {
	var user *pb.User
	return user, nil
}

func (repo *MockRepo) Create(user *pb.User) error {
	return nil
}

func (repo *MockRepo) GetByEmail(email string) (*pb.User, error) {
	var user *pb.User
	return user, nil
}

func newInstance() user.Authable {
	repo := &MockRepo{}
	return &user.TokenService{repo}
}

func TestCanCreateToken(t *testing.T) {
	srv := newInstance()
	token, err := srv.Encode(t_user)
	if err != nil {
		t.Fail()
	}

	if token == "" {
		t.Fail()
	}

	if len(strings.Split(token, ".")) != 3 {
		t.Fail()
	}
}

func TestCanDecodeToken(t *testing.T) {
	srv := newInstance()
	token, err := srv.Encode(t_user)
	t.Log(token)
	if err != nil {
		t.Fail()
	}
	claims, err := srv.Decode(token)
	t.Log(claims)
	if err != nil {
		t.Fail()
	}
	if claims.User == nil {
		t.Fail()
	}
	if claims.User.Email != "ewan.valentine89@gmail.com" {
		t.Fail()
	}
}

func TestThrowsErrorIfTokenInvalid(t *testing.T) {
	srv := newInstance()
	_, err := srv.Decode("nope.nope.nope")
	if err == nil {
		t.Fail()
	}
}
