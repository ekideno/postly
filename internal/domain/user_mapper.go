package domain

func ToPublicUserDTO(u *User) PublicUserDTO {
	return PublicUserDTO{
		ID:        u.ID,
		Username:  u.Username,
		Bio:       u.Bio,
		AvatarURL: u.AvatarURL,
		BannerURL: u.BannerURL,
	}
}

func ToPrivateUserDTO(u *User) PrivateUserDTO {
	return PrivateUserDTO{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Bio:       u.Bio,
		AvatarURL: u.AvatarURL,
		BannerURL: u.BannerURL,
	}
}

func FromRegisterRequest(req *RegisterRequest) *User {
	return &User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}
}

func MapToPublicUserDTOs(users []User) []PublicUserDTO {
	publicUsers := make([]PublicUserDTO, len(users))
	for i := range users {
		publicUsers[i] = ToPublicUserDTO(&users[i])
	}
	return publicUsers
}
