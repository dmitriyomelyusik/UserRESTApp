// Package postgres has a set of tests needed to check correctness of wroten postgres methods.
package postgres

import (
	"log"
	"os"
	"testing"

	"github.com/UserRESTApp/entity"
	"github.com/UserRESTApp/errors"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

var (
	p  *Postgres
	ID []string
)

func TestMain(m *testing.M) {
	var err error
	p, err = NewDB("user=postgres dbname=postgres password=password sslmode=disable")
	if err != nil {
		log.Fatalf("Cannot open database: %v", err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestUser_DeleteUser(t *testing.T) {
	users := []entity.User{
		{Name: "deleteuser_test1", Email: "deleteuser_test1", Info: "deleteuser_test1"},
		{Name: "deleteuser_test2", Email: "deleteuser_test2", Info: "deleteuser_test2"},
		{Name: "deleteuser_test3", Email: "deleteuser_test3", Info: "deleteuser_test3"},
	}
	userIDs := []string{}
	for _, user := range users {
		id, err := p.CreateUser(user)
		if err != nil {
			t.Fatal(err)
		}
		userIDs = append(userIDs, id)
	}

	tt := []struct {
		name          string
		id            string
		expectedError error
	}{
		{
			name:          "not found",
			id:            "TestItem_test1",
			expectedError: errors.Error{Code: errors.UserNotFound, Message: "user is not found, id: TestItem_test1"},
		},
		{
			name:          "ok",
			id:            userIDs[0],
			expectedError: nil,
		},
		{
			name:          "user not found because deleted",
			id:            userIDs[0],
			expectedError: errors.Error{Code: errors.UserNotFound, Message: "user is not found, id: " + userIDs[0]},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := p.DeleteUser(tc.id)
			assert.Equal(t, tc.expectedError, err)
		})
	}

	for _, id := range userIDs {
		p.DeleteUser(id)
	}
}

func TestUser_UpdateUser(t *testing.T) {
	users := []entity.User{
		{Name: "updateuser_test1", Email: "updateuser_test1", Info: "updateuser_test1"},
		{Name: "updateuser_test2", Email: "updateuser_test2", Info: "updateuser_test2"},
		{Name: "updateuser_test3", Email: "updateuser_test3", Info: "updateuser_test3"},
	}
	userIDs := []string{}
	for _, user := range users {
		id, err := p.CreateUser(user)
		if err != nil {
			t.Fatal(err)
		}
		userIDs = append(userIDs, id)
	}

	tt := []struct {
		name             string
		user             entity.User
		expectedUser     entity.User
		expectedGetError error
		expectedError    error
	}{
		{
			name:             "not found",
			user:             entity.User{ID: "updateuser_wrongID", Name: "updateuser_wrongName", Email: "updateuser_wrongEmail", Info: "updateuser_wrongInfo"},
			expectedUser:     entity.User{},
			expectedGetError: errors.Error{Code: errors.UserNotFound, Message: "cannot find user in database with id updateuser_wrongID"},
			expectedError:    errors.Error{Code: errors.UserNotFound, Message: "user is not found, id: updateuser_wrongID"},
		},
		{
			name:             "ok",
			user:             entity.User{ID: userIDs[0], Name: "updateuser_normalName", Email: "updateuser_normalEmail", Info: "updateuser_normalInfo"},
			expectedUser:     entity.User{ID: userIDs[0], Name: "updateuser_normalName", Email: "updateuser_normalEmail", Info: "updateuser_normalInfo"},
			expectedGetError: nil,
			expectedError:    nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := p.UpdateUser(tc.user)
			assert.Equal(t, tc.expectedError, err)
			user, err := p.GetUser(tc.user.ID)
			assert.Equal(t, tc.expectedGetError, err)
			assert.Equal(t, tc.expectedUser, user)
		})
	}

	for _, id := range userIDs {
		p.DeleteUser(id)
	}
}

func TestUser_UpdateUserFields(t *testing.T) {
	users := []entity.User{
		{Name: "updateuserfields_test1", Email: "updateuserfields_test1", Info: "updateuserfields_test1"},
		{Name: "updateuserfields_test2", Email: "updateuserfields_test2", Info: "updateuserfields_test2"},
		{Name: "updateuserfields_test3", Email: "updateuserfields_test3", Info: "updateuserfields_test3"},
	}
	userIDs := []string{}
	for i, user := range users {
		id, err := p.CreateUser(user)
		if err != nil {
			t.Fatal(err)
		}
		users[i].ID = id
		userIDs = append(userIDs, id)
	}

	tt := []struct {
		name             string
		id               string
		changes          map[string]interface{}
		expectedUser     entity.User
		expectedGetError error
		expectedError    error
	}{
		{
			name:             "not found",
			id:               "updateuserfields_wrongID",
			changes:          map[string]interface{}{"name": "useless", "email": "useless"},
			expectedUser:     entity.User{},
			expectedGetError: errors.Error{Code: errors.UserNotFound, Message: "cannot find user in database with id updateuserfields_wrongID"},
			expectedError:    errors.Error{Code: errors.UserNotFound, Message: "user is not found, id: updateuserfields_wrongID"},
		},
		{
			name:             "ok",
			id:               userIDs[1],
			changes:          map[string]interface{}{"name": "updateuserfields_test2.2", "email": "updateuserfields_test2.2"},
			expectedUser:     entity.User{ID: userIDs[1], Email: "updateuserfields_test2.2", Name: "updateuserfields_test2.2", Info: "updateuserfields_test2"},
			expectedGetError: nil,
			expectedError:    nil,
		},
		{
			name:             "wrong fields",
			id:               userIDs[0],
			changes:          map[string]interface{}{"surname": "updateuserfields_wrongtest"},
			expectedUser:     users[0],
			expectedGetError: nil,
			expectedError:    errors.Error{Code: errors.DatabaseQueryError, Message: "bad query for database: UPDATE users SET surname='updateuserfields_wrongtest' WHERE id='" + userIDs[0] + "'"},
		},
		{
			name:             "empty changes",
			id:               userIDs[2],
			changes:          map[string]interface{}{},
			expectedUser:     users[2],
			expectedGetError: nil,
			expectedError:    errors.Error{Code: errors.DatabaseQueryError, Message: "bad query for database: UPDATE users SET  WHERE id='" + userIDs[2] + "'"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := p.UpdateUserFields(tc.changes, tc.id)
			assert.Equal(t, tc.expectedError, err)
			user, err := p.GetUser(tc.id)
			assert.Equal(t, tc.expectedGetError, err)
			assert.Equal(t, tc.expectedUser, user)
		})
	}

	for _, id := range userIDs {
		p.DeleteUser(id)
	}
}

func TestUser_GetUsers(t *testing.T) {
	users := []entity.User{
		{Name: "getusers_test1", Email: "getusers_test1", Info: "getusers_test1"},
		{Name: "getusers_test2", Email: "getusers_test2", Info: "getusers_test2"},
		{Name: "getusers_test3", Email: "getusers_test3", Info: "getusers_test3"},
	}
	userIDs := []string{}
	for _, user := range users {
		id, err := p.CreateUser(user)
		if err != nil {
			t.Fatal(err)
		}
		userIDs = append(userIDs, id)
	}

	tt := []struct {
		name          string
		expectedError error
	}{
		{
			name:          "ok",
			expectedError: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			u, err := p.GetUsers()
			assert.Equal(t, tc.expectedError, err)
			for i := range users {
				users[i].ID = userIDs[i]
			}
			numOfU := len(users)
			for _, v1 := range u {
				for _, v2 := range users {
					if v1 == v2 {
						numOfU--
					}
				}
			}
			assert.Equal(t, 0, numOfU)
		})
	}

	for _, id := range userIDs {
		p.DeleteUser(id)
	}
}

func TestUser_GetUser(t *testing.T) {
	users := []entity.User{
		{Name: "getuser_test1", Email: "getuser_test1", Info: "getuser_test1"},
		{Name: "getuser_test2", Email: "getuser_test2", Info: "getuser_test2"},
		{Name: "getuser_test3", Email: "getuser_test3", Info: "getuser_test3"},
	}
	userIDs := []string{}
	for i, user := range users {
		id, err := p.CreateUser(user)
		if err != nil {
			t.Fatal(err)
		}
		users[i].ID = id
		userIDs = append(userIDs, id)
	}

	tt := []struct {
		name          string
		id            string
		expectedUser  entity.User
		expectedError error
	}{
		{
			name:          "ok0",
			id:            userIDs[0],
			expectedUser:  users[0],
			expectedError: nil,
		},
		{
			name:          "ok1",
			id:            userIDs[1],
			expectedUser:  users[1],
			expectedError: nil,
		},
		{
			name:          "ok2",
			id:            userIDs[2],
			expectedUser:  users[2],
			expectedError: nil,
		},
		{
			name:          "wrong id",
			id:            "getuser_wrongID",
			expectedUser:  entity.User{},
			expectedError: errors.Error{Code: errors.UserNotFound, Message: "cannot find user in database with id getuser_wrongID"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			u, err := p.GetUser(tc.id)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedUser, u)
		})
	}

	for _, id := range userIDs {
		p.DeleteUser(id)
	}
}

func TestUser_CreateUser(t *testing.T) {
	users := []entity.User{
		{Name: "createuser_test1", Email: "createuser_test1", Info: "createuser_test1"},
		{Name: "createuser_test2", Email: "createuser_test2", Info: "createuser_test2"},
		{Name: "createuser_test3", Email: "createuser_test3", Info: "createuser_test3"},
	}
	userIDs := []string{}

	tt := []struct {
		name          string
		newUser       entity.User
		expectedError error
	}{
		{
			name:          "ok0",
			newUser:       users[0],
			expectedError: nil,
		},
		{
			name:          "ok1",
			newUser:       users[1],
			expectedError: nil,
		},
		{
			name:          "ok2",
			newUser:       users[2],
			expectedError: nil,
		},
		{
			name:          "adding similar user",
			newUser:       users[0],
			expectedError: nil,
		},
		{
			name:          "empty user",
			newUser:       entity.User{},
			expectedError: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			id, err := p.CreateUser(tc.newUser)
			assert.Equal(t, tc.expectedError, err)
			userIDs = append(userIDs, id)
			u, err := p.GetUser(id)
			assert.Equal(t, tc.expectedError, err)
			tc.newUser.ID = id
			assert.Equal(t, tc.newUser, u)
		})
	}

	for _, id := range userIDs {
		p.DeleteUser(id)
	}
}
