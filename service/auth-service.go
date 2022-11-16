package service

import (
	"github.com/mashingan/smapping"
	"github.com/thanhpk/randstr"
	"golang-api/config"
	"golang-api/dto"
	"golang-api/entity"
	"golang-api/helper"
	"golang-api/repository"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type AuthService interface {
	VerifyCredential(email string, password string) interface{}
	CreateUser(user dto.RegisterDTO) entity.User
	FindByEmail(email string) entity.User
	IsDuplicateEmail(email string) bool
	VerifyEmail(verificationCode string) (entity.User, error)
	IsUserVerified(email string) bool
}

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRep repository.UserRepository) AuthService {
	return &authService{
		userRepository: userRep,
	}
}

func (service *authService) VerifyCredential(email string, password string) interface{} {
	res := service.userRepository.VerifyCredential(email, password)
	if v, ok := res.(entity.User); ok {
		comparedPassword := comparePassword(v.Password, []byte(password))
		if v.Email == email && comparedPassword {
			return res
		}
		return false
	}
	return false
}

func (service *authService) CreateUser(user dto.RegisterDTO) entity.User {
	userToCreate := entity.User{}
	err := smapping.FillStruct(&userToCreate, smapping.MapFields(&user))
	code := randstr.String(20)
	verificationCode := helper.Encode(code)
	userToCreate.VerificationCode = verificationCode
	env := config.LoadEnv()
	emailData := helper.EmailData{
		URL:       env.CLIENT_URL + "/verify/" + verificationCode,
		FirstName: userToCreate.Name,
		Subject:   "Your account verification code",
	}
	helper.SendEmail(&userToCreate, &emailData)
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}
	res := service.userRepository.InsertUser(userToCreate)
	return res
}

func (service *authService) VerifyEmail(verificationCode string) (entity.User, error) {
	user, err := service.userRepository.VerifyEmail(verificationCode)
	return user, err
}

func (service *authService) IsUserVerified(email string) bool {
	user := service.userRepository.FindByEmail(email)
	if user.VerificationCode == "" && user.Verified {
		return true
	}
	return false
}

func (service *authService) FindByEmail(email string) entity.User {
	return service.userRepository.FindByEmail(email)
}

func (service *authService) IsDuplicateEmail(email string) bool {
	res := service.userRepository.IsDuplicateEmail(email)
	return !(res.Error == nil)
}

func comparePassword(hashedPwd string, plainPassword []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
