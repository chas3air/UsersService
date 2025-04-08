package profiles

import (
	"users-service/internal/domain/models"
	umv1 "users-service/proto/gen"

	"github.com/google/uuid"
)

func UserToProtoUser(user models.User) (umv1.User, error) {
	return umv1.User{
		Id:       user.Id.String(),
		Login:    user.Login,
		Password: user.Password,
	}, nil
}

func ProtoUserToUser(user *umv1.User) (models.User, error) {
	id, err := uuid.Parse(user.Id)
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		Id:       id,
		Login:    user.Login,
		Password: user.Password,
	}, nil
}
