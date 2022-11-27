package game

import (
	"fmt"
	"time"
)

type MapObject struct {
	ID        int64         `json:"id"`
	GameID    int64         `json:"gameId" boltholdIndex:"GameID"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
	Type      MapObjectType `json:"type"`
	PlayerID  int64         `json:"-"`
	Dirty     bool          `json:"-"`
	Delete    bool          `json:"-"`
	Position  Vector        `json:"position"`
	Name      string        `json:"name"`
	Num       int           `json:"num"`
	PlayerNum int           `json:"playerNum"`
	// Tags      Tags           `json:"tags"`
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

func (mo *MapObject) String() string {
	return fmt.Sprintf("GameID: %5d, ID: %5d, Num: %3d %s", mo.GameID, mo.ID, mo.Num, mo.Name)
}

func (mo *MapObject) Owned() bool {
	return mo.PlayerNum != Unowned
}

func (mo *MapObject) OwnedBy(num int) bool {
	return mo.PlayerNum != Unowned && mo.PlayerNum == num
}
