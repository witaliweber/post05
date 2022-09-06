package post05

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type Userdata struct {
	ID          int
	Username    string
	Name        string
	Surname     string
	Description string
}

var (
	Hostname = ""
	Port     = 2345
	Username = ""
	Password = ""
	Database = ""
)

func openConnection() (*sql.DB, error) {
	// connection string
	conn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Hostname,
		Port,
		Username,
		Password,
		Database,
	)

	// open database
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func exists(username string) int {
	username = strings.ToLower(username)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	userID := -1
	statement := fmt.Sprintf(
		`SELECT "id" FROM "users" WHERE username = '%s'`,
		username,
	)
	rows, err := db.Query(statement)
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println("Scan", err)
			return -1
		}
		userID = id
	}
	defer rows.Close()
	return userID
}

func ListUsers() ([]Userdata, error) {
	Data := []Userdata{}
	db, err := openConnection()
	if err != nil {
		return Data, err
	}
	defer db.Close()

	rows, err := db.Query(
		`SELECT "id", "username", "name", "description"
		 FROM "users", "userdata"
		 WHERE users.id = userdata.userid`)
	if err != nil {
		return Data, err
	}

	for rows.Next() {
		var id int
		var username string
		var name string
		var surname string
		var description string
		err = rows.Scan(&id, &username, &name, &surname, &description)
		temp := Userdata{
			ID:          id,
			Username:    username,
			Name:        name,
			Surname:     surname,
			Description: description,
		}
		Data = append(Data, temp)
		if err != nil {
			return Data, err
		}
	}
	defer rows.Close()
	return Data, nil
}
