package membersapi

type User struct {
	UserID       int    `json:"UserID"`
	FirstName    string `json:"FirstName"`
	LastName     string `json:"LastName"`
	EmailAddress string `json:"EmailAddress"`
	DateOfBirth  string `json:"DateOfBirth"`
	Hash         string `json:"HashedPassword"`
}
