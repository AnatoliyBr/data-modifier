package entity

func TestUser() *User {
	return &User{
		DisplayName: "Иванов Семен Петрович",
		Email:       "petrovich@mail.ru",
		MobilePhone: "+71234567890",
		WorkPhone:   "1234",
		ID:          123,
	}
}

func TestCredentials() *Credentials {
	return &Credentials{
		IP:         "127.0.0.1",
		Port:       ":8080",
		Login:      "admin",
		Password:   "password",
		AbsenceURL: "https://127.0.0.1:8080/Portal/springApi/api/absences",
		AuthURL:    "https://127.0.0.1:8080/tokens",
	}
}
