package domain

func ToPublicUserDTO(u *User) PublicUserDTO {
	return PublicUserDTO{
		ID:       u.ID,
		Username: u.Username,
		Bio:      u.Bio,
	}
}

func ToPrivateUserDTO(u *User) PrivateUserDTO {
	return PrivateUserDTO{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Bio:      u.Bio,
	}
}

func FromRegisterRequest(req *RegisterRequest) *User {
	return &User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}
}
