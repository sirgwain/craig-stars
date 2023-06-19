package db

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/game"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {

	type args struct {
		c    *client
		user *game.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &game.User{Username: "test"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := *tt.args.user
			err := tt.args.c.CreateUser(tt.args.user)

			// id is automatically added
			want.ID = tt.args.user.ID
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.args.user, &want) {
				t.Errorf("CreateUser() = \n%v, want \n%v", tt.args.user, want)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	c := connectTestDB()
	user := game.User{Username: "Test"}
	if err := c.CreateUser(&user); err != nil {
		t.Errorf("create user %s", err)
		return
	}

	user.Username = "Test2"
	user.Password = "newpassword"
	user.Role = game.RoleAdmin
	if err := c.UpdateUser(&user); err != nil {
		t.Errorf("update user %s", err)
		return
	}

	updated, err := c.GetUser(user.ID)

	if err != nil {
		t.Errorf("get user %s", err)
		return
	}

	assert.Equal(t, user.Username, updated.Username)
	assert.Equal(t, user.Password, updated.Password)
	assert.Equal(t, user.Role, updated.Role)
	assert.Less(t, user.UpdatedAt, updated.UpdatedAt)

}

func TestGetUser(t *testing.T) {
	c := connectTestDB()
	user := game.User{Username: "Test"}
	if err := c.CreateUser(&user); err != nil {
		t.Errorf("create user %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *game.User
		wantErr bool
	}{
		{"No results", args{id: 0}, nil, false},
		{"Got user", args{id: user.ID}, &user, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetUser(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				tt.want.UpdatedAt = got.UpdatedAt
				tt.want.CreatedAt = got.CreatedAt
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUsers(t *testing.T) {
	c := connectTestDB()

	// start with 1 user from connectTestDB
	result, err := c.GetUsers()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

	user := game.User{Username: "Test"}
	if err := c.CreateUser(&user); err != nil {
		t.Errorf("create user %s", err)
		return
	}

	result, err = c.GetUsers()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(result))

}

func TestDeleteUsers(t *testing.T) {
	c := connectTestDB()

	_, err := c.GetUsers()
	assert.Nil(t, err)

	user := game.User{Username: "Test"}
	if err := c.CreateUser(&user); err != nil {
		t.Errorf("create user %s", err)
		return
	}

	// should have our user in the db
	result, err := c.GetUsers()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(result))

	if err := c.DeleteUser(user.ID); err != nil {
		t.Errorf("delete user %s", err)
		return
	}

	// should be no users left in db
	result, err = c.GetUsers()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))
}
