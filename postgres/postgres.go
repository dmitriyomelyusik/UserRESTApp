package postgres

import (
	"database/sql"
	"encoding/json"

	"../entity"
)

//Postgres is a PostgreSQL database
type Postgres struct {
	DB *sql.DB
}

//Users returns all users in DB
func (p Postgres) Users() ([]entity.User, error) {
	var users []entity.User
	rows, err := p.DB.Query("SELECT * from users")
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var u entity.User
		var rawInfo []byte
		err = rows.Scan(&u.ID, &u.Email, &u.Name, &rawInfo)
		if err != nil {
			return users, err
		}
		err = json.Unmarshal(rawInfo, &u.Info)
		if err != nil {
			return users, err
		}
		users = append(users, u)
	}
	return users, nil
}

//NewDB returns new posgres DB with configuration conf
func NewDB(conf string) (*Postgres, error) {
	db, err := sql.Open("postgres", conf)
	if err != nil {
		return nil, err
	}
	return &Postgres{DB: db}, nil
}
