package db

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {

	type args struct {
		c    *client
		user *cs.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &cs.User{Username: "test"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := *tt.args.user
			got, err := tt.args.c.CreateUser(t.Context(), tt.args.user)

			// make sure auto generated fields are equal before compare
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("CreateUser() did not return error when expected")
				} else {
					t.Fatalf("CreateUser() errored unexpectedly; err = \n%v", err)
				}
			}

			want.DBObject = got.DBObject
			test.CompareAsJSON(t, got, want)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	c := connectTestDB()

	var user *cs.User
	var err error
	user, err = c.CreateUser(t.Context(), &cs.User{Username: "Test"})
	if err != nil {
		t.Errorf("create user %s", err)
		return
	}

	user.Username = "Test2"
	user.Password = "newpassword"
	user.Role = cs.RoleAdmin
	if err := c.UpdateUser(t.Context(), user); err != nil {
		t.Errorf("update user %s", err)
		return
	}

	updated, err := c.GetUser(t.Context(), user.ID)

	if err != nil {
		t.Errorf("get user %s", err)
		return
	}

	assert.Equal(t, user.Username, updated.Username)
	assert.Equal(t, user.Password, updated.Password)
	assert.Equal(t, user.Role, updated.Role)

}

func TestGetUser(t *testing.T) {
	c := connectTestDB()

	var user *cs.User
	var err error
	user, err = c.CreateUser(t.Context(), &cs.User{Username: "Test"})
	if err != nil {
		t.Errorf("create user %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *cs.User
		wantErr bool
	}{
		{"No results", args{id: 0}, nil, false},
		{"Got user", args{id: user.ID}, user, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetUser(t.Context(), tt.args.id)
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("GetUser() did not return error when expected")
				} else {
					t.Fatalf("GetUser() errored unexpectedly; err = \n%v", err)
				}
			}
			if got != nil {
				tt.want.UpdatedAt = got.UpdatedAt
				tt.want.CreatedAt = got.CreatedAt
			}
			test.CompareAsJSON(t, got, tt.want)
		})
	}
}

func TestGetUsers(t *testing.T) {
	c := connectTestDB()

	// start with 1 user from connectTestDB
	result, err := c.GetUsers(t.Context())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

	_, err = c.CreateUser(t.Context(), &cs.User{Username: "Test"})
	if err != nil {
		t.Errorf("create user %s", err)
		return
	}

	result, err = c.GetUsers(t.Context())
	assert.Nil(t, err)
	assert.Equal(t, 2, len(result))

}

func TestDeleteUsers(t *testing.T) {
	c := connectTestDB()

	_, err := c.GetUsers(t.Context())
	assert.Nil(t, err)

	var user *cs.User
	user, err = c.CreateUser(t.Context(), &cs.User{Username: "Test"})
	if err != nil {
		t.Errorf("create user %s", err)
		return
	}

	// should have our user in the db
	result, err := c.GetUsers(t.Context())
	assert.Nil(t, err)
	assert.Equal(t, 2, len(result))

	if err := c.DeleteUser(t.Context(), user.ID); err != nil {
		t.Errorf("delete user %s", err)
		return
	}

	// should be no users left in db
	result, err = c.GetUsers(t.Context())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))
}
