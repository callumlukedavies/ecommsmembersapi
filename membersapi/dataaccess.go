package membersapi

import (
	"database/sql"
	"fmt"
)

type DataAccess struct {
	DB *sql.DB
}

func (dataaccess *DataAccess) GetUser(emailAddr string) (User, error) {

	rows, err := dataaccess.DB.Query("SELECT * FROM usersdb.users where EmailAddress = (?)", emailAddr)
	if err != nil {
		fmt.Print(err)
		return User{}, err
	}
	defer rows.Close()

	var user User

	for rows.Next() {
		if err := rows.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.EmailAddress, &user.DateOfBirth, &user.Hash); err != nil {
			if err == sql.ErrNoRows {
				fmt.Print(err)
				return User{}, sql.ErrNoRows
			}
		}
	}

	return user, nil
}

func (dataaccess *DataAccess) DeleteUser(userID int) error {
	_, err := dataaccess.DB.Exec("DELETE FROM usersdb.users WHERE UserID = ?", userID)
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

func (dataaccess *DataAccess) UpdateUserFirstName(userID int, firstName string) error {

	_, err := dataaccess.DB.Exec("UPDATE usersdb.users SET FirstName = (?) where UserID = (?)", firstName, userID)

	if err != nil {
		return err
	}

	return nil
}

func (dataaccess *DataAccess) UpdateUserLastName(userID int, lastName string) error {

	_, err := dataaccess.DB.Exec("UPDATE usersdb.users SET LastName = (?) where UserID = (?)", lastName, userID)

	if err != nil {
		return err
	}

	return nil
}

func (dataaccess *DataAccess) UpdateUserDateOfBirth(userID int, dateOfBirth string) error {

	_, err := dataaccess.DB.Exec("UPDATE usersdb.users SET DateOfBirth = (?) where UserID = (?)", dateOfBirth, userID)

	if err != nil {
		return err
	}

	return nil
}

func (dataaccess *DataAccess) UpdateUserEmail(userID int, emailAddress string) error {

	_, err := dataaccess.DB.Exec("UPDATE usersdb.users SET EmailAddress = (?) where UserID = (?)", emailAddress, userID)

	if err != nil {
		return err
	}

	return nil
}

func (dataaccess *DataAccess) UpdateUserPassword(userID int, hashedPassword string) error {

	_, err := dataaccess.DB.Exec("UPDATE usersdb.users SET HashedPassword = (?) where UserID = (?)", hashedPassword, userID)

	if err != nil {
		return err
	}

	return nil
}

func (dataaccess *DataAccess) CheckUserExists(emailAddress string) (bool, error) {

	var emailExists bool
	row := dataaccess.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM usersdb.users WHERE EmailAddress = (?))", emailAddress)

	if err := row.Scan(&emailExists); err != nil {
		return false, err
	}

	return emailExists, nil
}
