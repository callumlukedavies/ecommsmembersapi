package membersapi

import (
	"database/sql"
	"fmt"
)

type DataAccess struct {
	DB *sql.DB
}

func (dataaccess *DataAccess) GetUser(emailAddr, password string) (User, error) {

	rows, err := dataaccess.DB.Query("SELECT * FROM usersdb.users where Username = (?) AND Password = (?)", emailAddr, password)
	if err != nil {
		fmt.Print(err)
		return User{}, err
	}
	defer rows.Close()

	var user User

	for rows.Next() {
		if err := rows.Scan(&user.UserID, &user.EmailAddress, &user.Password, &user.DateOfBirth, &user.FirstName, &user.Surname); err != nil {
			if err == sql.ErrNoRows {
				fmt.Print(err)
				return User{}, sql.ErrNoRows
			}
		}
	}

	return user, nil
}

func (dataaccess *DataAccess) DeleteUser(userID int64) error {
	_, err := dataaccess.DB.Exec("DELETE FROM usersdb.users WHERE userID = ?", userID)
	return err
}

func (dataaccess *DataAccess) CreateUser(firstName, lastname, email, dob, password string) error {

	_, err := dataaccess.DB.Exec("INSERT INTO usersdb.users"+
		"(FirstName, LastName, EmailAddress, DateOfBirth, HashedPassword)"+
		"VALUES (?, ?, ?, ?, ?)", firstName, lastname, email, dob, password)

	if err != nil {
		return err
	}

	return nil
}

func (dataaccess *DataAccess) UpdateUserData(userID int64, updateKey, updateValue string) error {

	_, err := dataaccess.DB.Exec("UPDATE usersdb.users SET %s = (?) where userID = (?)", updateKey, updateValue, userID)

	if err != nil {
		return err
	}

	return nil
}
