package membersapi

type User struct {
	UserID       int    `json:"UserID"`
	EmailAddress string `json:"Username"`
	Password     string `json:"Password"`
	DateOfBirth  string `json:"DateOfBirth"`
	FirstName    string `json:"FirstName"`
	Surname      string `json:"LastName"`
}
