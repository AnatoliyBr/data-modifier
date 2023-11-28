package entity

import "time"

// TestUser returns the User entity initialized with valid data.
func TestUser() *User {
	return &User{
		DisplayName: "Иванов Семен Петрович",
		Email:       "petrovich@mail.ru",
		MobilePhone: "+71234567890",
		WorkPhone:   "1234",
		ID:          1234,
	}
}

// TestUserAbsenceData returns the UserAbsenceData entity
// initialized with valid data.
func TestUserAbsenceData() *UserAbsenceData {
	createdDate, _ := time.Parse("2006-01-02", "2023-08-14")
	dateFrom, _ := time.Parse("2006-01-02T15:04:05", "2023-08-12T00:00:00")
	dateTo, _ := time.Parse("2006-01-02T15:04:05", "2023-08-12T23:59:59")

	return &UserAbsenceData{
		CreatedDate: CustomDate{Time: createdDate},
		DateFrom:    CustomTime{Time: dateFrom},
		DateTo:      CustomTime{Time: dateTo},
		ID:          28246,
		PersonID:    1234,
		ReasonID:    1,
	}
}

// TestCredentials returns the Credentials entity
// initialized with valid data.
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
