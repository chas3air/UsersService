package profiles

import (
	"users-service/internal/domain/models"
	umv1 "users-service/proto/gen"

	"github.com/google/uuid"
)

func UserToProtoUser(user models.User) *umv1.User {
	return &umv1.User{
		Id:       user.Id.String(),
		Login:    user.Login,
		Password: user.Password,
	}
}

func ProtoUserToUser(user *umv1.User) models.User {
	id, _ := uuid.Parse(user.Id)

	return models.User{
		Id:       id,
		Login:    user.Login,
		Password: user.Password,
	}
}
