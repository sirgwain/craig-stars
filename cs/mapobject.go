package cs

import (
	"encoding/json"
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
	TagTarget  = "target"
)

type Target[T PlayerMessageTargetType | MapObjectType] struct {
	TargetType      T      `json:"targetType,omitempty"`
	TargetName      string `json:"targetName,omitempty"`
	TargetNum       int    `json:"targetNum,omitempty"`
	TargetPlayerNum int    `json:"targetPlayerNum,omitempty"`
}

type MapObjectTarget Target[MapObjectType]

func (mo *MapObject) String() string {
	return fmt.Sprintf("GameID: %5d, ID: %5d, Num: %3d %s", mo.GameID, mo.ID, mo.Num, mo.Name)
}

func (mo *MapObject) Target() MapObjectTarget {
	return MapObjectTarget{TargetType: mo.Type, TargetName: mo.Name, TargetNum: mo.Num, TargetPlayerNum: mo.PlayerNum}
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

func (mo *MapObject) RemoveTag(key string) {
	delete(mo.Tags, key)
}

func (target MapObjectTarget) IsNone() bool {
	return target.TargetType == MapObjectTypeNone
}

func (target MapObjectTarget) String() string {
	json, _ := json.Marshal(target)
	return string(json)
}

func TargetFromString(value string) MapObjectTarget {
	target := MapObjectTarget{}
	json.Unmarshal([]byte(value), &target)
	return target
}
