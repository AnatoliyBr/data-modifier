package entity

func TestUser() *User {
	return &User{
		DisplayName: "Иванов Семен Петрович",
		Email:       "petrovich@mail.ru",
		MobilePhone: "+71234567890",
		WorkPhone:   "1234",
		ID:          1234,
	}
}

func TestCredentials() *Credentials {
	return &Credentials{
		IP:          "127.0.0.1",
		Port:        ":8080",
		Login:       "admin",
		Password:    "password",
		EmployeeURL: "https://127.0.0.1:8080/Portal/springApi/api/employees",
		AbsenceURL:  "https://127.0.0.1:8080/Portal/springApi/api/absences",
	}
}
