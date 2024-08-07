package cs

import (
	"fmt"
	"time"
)

// Every object stored in the database has an ID and a create/update timestamp.
// Though the cs package doesn't deal with the database, they are still part of the models
type DBObject struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// A GameObject is a database object that is associated with a game
type GameDBObject struct {
	ID        int64     `json:"id"`
	GameID    int64     `json:"gameId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Each object in the universe is a MapObject. MapObjects have a unique Num (and often a PlayerNum for player owned
// map objects), as well as a Position in space.
type MapObject struct {
	GameDBObject
	Type      MapObjectType `json:"type"`
	Dirty     bool          `json:"-"`
	Delete    bool          `json:"-"`
	Position  Vector        `json:"position"`
	Num       int           `json:"num"`
	PlayerNum int           `json:"playerNum"`
	Name      string        `json:"name"`
	Tags      Tags          `json:"tags"`
}

type MapObjectType string

const (
	MapObjectTypeNone          MapObjectType = ""
	MapObjectTypePlanet        MapObjectType = "Planet"
	MapObjectTypeFleet         MapObjectType = "Fleet"
	MapObjectTypeWormhole      MapObjectType = "Wormhole"
	MapObjectTypeMineField     MapObjectType = "MineField"
	MapObjectTypeMysteryTrader MapObjectType = "MysteryTrader"
	MapObjectTypeSalvage       MapObjectType = "Salvage"
	MapObjectTypeMineralPacket MapObjectType = "MineralPacket"
)

const (
	TagPurpose = "purpose"
)

func (mo *MapObject) String() string {
	return fmt.Sprintf("GameID: %5d, ID: %5d, Num: %3d %s", mo.GameID, mo.ID, mo.Num, mo.Name)
}

func (mo *MapObject) Owned() bool {
	return mo.PlayerNum != Unowned
}

func (mo *MapObject) OwnedBy(num int) bool {
	return mo.PlayerNum != Unowned && mo.PlayerNum == num
}

func (mo *MapObject) MarkDirty() {
	mo.Dirty = true
}

func (mo *MapObject) GetTag(key string) string {
	return mo.Tags[key]
}

func (mo *MapObject) SetTag(key, value string) {
	if mo.Tags == nil {
		mo.Tags = make(Tags)
	}
	mo.Tags[key] = value
}
