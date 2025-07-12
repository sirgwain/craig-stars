package db

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateRace(t *testing.T) {

	type args struct {
		c    *client
		race *cs.Race
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &cs.Race{UserID: 1, Name: "test", PluralName: "testers"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := *tt.args.race
			got, err := tt.args.c.CreateRace(t.Context(), tt.args.race)

			// id is automatically added
			want.DBObject = got.DBObject
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("CreateRace() did not return error when expected")
				} else {
					t.Fatalf("CreateRace() errored unexpectedly; err = \n%v", err)
				}
			}
			test.CompareAsJSON(t, got, want)
		})
	}
}

func TestUpdateRace(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	race, err := c.CreateRace(t.Context(), &cs.Race{UserID: 1, Name: "Test"})
	if err != nil {
		t.Errorf("create race %s", err)
		return
	}

	race.Name = "Test2"
	race.PluralName = "Testers"
	if err := c.UpdateRace(t.Context(), race); err != nil {
		t.Errorf("update race %s", err)
		return
	}

	updated, err := c.GetRace(t.Context(), race.ID)

	if err != nil {
		t.Errorf("get race %s", err)
		return
	}

	assert.Equal(t, race.Name, updated.Name)
	assert.Equal(t, race.PluralName, updated.PluralName)

}

func TestGetRace(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	rules := cs.NewRules()

	race := &cs.Race{UserID: 1, Name: "Test", PluralName: "testers"}
	race = race.WithSpec(&rules)
	var err error
	if race, err = c.CreateRace(t.Context(), race); err != nil {
		t.Errorf("create race %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *cs.Race
		wantErr bool
	}{
		{"No results", args{id: 0}, nil, false},
		{"Got race", args{id: race.ID}, race, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetRace(t.Context(), tt.args.id)
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("GetRace() did not return error when expected")
				} else {
					t.Fatalf("GetRace() errored unexpectedly; err = \n%v", err)
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

func TestGetRaces(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	// start with 1 race from connectTestDB
	result, err := c.GetRaces(t.Context())
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))

	if _, err = c.CreateRace(t.Context(), &cs.Race{UserID: 1, Name: "Test", PluralName: "testers"}); err != nil {
		t.Errorf("create race %s", err)
		return
	}

	result, err = c.GetRaces(t.Context())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestDeleteRaces(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	result, err := c.GetRaces(t.Context())
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))

	race, err := c.CreateRace(t.Context(), &cs.Race{UserID: 1, Name: "Test", PluralName: "Testers"})
	if err != nil {
		t.Errorf("create race %s", err)
		return
	}

	// should have our race in the db
	result, err = c.GetRaces(t.Context())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

	if err := c.DeleteRace(t.Context(), race.ID); err != nil {
		t.Errorf("delete race %s", err)
		return
	}

	// should be no races left in db
	result, err = c.GetRaces(t.Context())
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))
}
