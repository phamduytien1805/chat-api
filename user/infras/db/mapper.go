package db

import "github.com/phamduytien1805/user/domain"

func mapToUser(u User) domain.User {
	return domain.User{
		ID:            u.ID,
		Username:      u.Username,
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
	}
}

func mapToUserCredential(uc UserCredential) domain.UserCredential {
	return domain.UserCredential{
		HashedPassword: uc.HashedPassword,
	}
}
