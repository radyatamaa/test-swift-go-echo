package usecase

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/auth/user"
	guuid "github.com/google/uuid"
	"github.com/models"
	"golang.org/x/crypto/bcrypt"
	"io"
	"time"
)

type userUsecase struct {
	userRepo user.Repository
	contextTimeout time.Duration
}



// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewuserUsecase(timeout time.Duration,	userRepo user.Repository) user.Usecase {
	return &userUsecase{
		userRepo:userRepo,
		contextTimeout: timeout,
	}
}
var(
	key = []byte("TW4e87abf80afb7467eb89qwesad1234") // 32 bytes
)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func (a userUsecase) ValidateUser(c context.Context, email,password string) (*models.UserDto, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	//plaintext := []byte(password)
	//ciphertext, err := encrypt(key, plaintext)
	//if err != nil {
	//	return nil,err
	//}
	//fmt.Printf("%0x\n", ciphertext)
	user ,err := a.userRepo.ValidateUser(ctx,email)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}
	match := CheckPasswordHash(password, user.Password)
	if match == false{
		return nil,models.ErrUnAuthorize
	}
	password = user.Password
	if err != nil{
		return nil,models.ErrUnAuthorize
	}
	result := &models.UserDto{
		Id:        user.Id,
		UserEmail: user.UserEmail,
		Password:  password,
		Phone:     user.Phone,
	}

	return result,nil
}
func (a userUsecase) Create(c context.Context, user models.NewCommandUser) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	userM := models.User{
		Id:          guuid.New().String(),
		CreatedBy:   user.UserEmail,
		CreatedDate:  time.Time{},
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		UserEmail:    user.UserEmail,
		Password:    "",
		Phone:        user.Phone,
	}

	hash, _ := HashPassword(user.Password)
	userM.Password = hash

	insert ,err := a.userRepo.Create(ctx,userM)
	if err != nil {
		return nil,err
	}

	result := &models.ResponseDelete{
		Id:      *insert,
		Message: "succcess insert",
	}

	return result,nil
}

func encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}
