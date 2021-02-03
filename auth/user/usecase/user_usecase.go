package usecase

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/auth/user"
	guuid "github.com/google/uuid"
	"github.com/models"
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
	key = []byte("a very very very very secret key") // 32 bytes
)
func (a userUsecase) ValidateUser(c context.Context, email,password string) (*models.UserDto, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	plaintext := []byte(password)
	ciphertext, err := decrypt(key, plaintext)
	if err != nil {
		return nil,err
	}
	password = fmt.Sprintf("%0x\n", ciphertext)
	user ,err := a.userRepo.ValidateUser(ctx,email,password)



	if err != nil{
		return nil,models.ErrUnAuthorize
	}
	result := &models.UserDto{
		Id:        user.Id,
		UserEmail: user.UserEmail,
		Password:  user.Password,
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

	plaintext := []byte(user.Password)
	ciphertext, err := encrypt(key, plaintext)
	if err != nil {
		return nil,err
	}
	userM.Password = fmt.Sprintf("%0x\n", ciphertext)

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
