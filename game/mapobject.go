package game

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type MapObject struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	Dirty     bool           `json:"-" gorm:"-"`
	GameID    uint           `json:"gameId"`
	PlayerID  uint           `json:"playerId"`
	Position  Vector         `json:"position" gorm:"embedded"`
	Name      string         `json:"name"`
	Num       int            `json:"num"`
	PlayerNum *int           `json:"playerNum"`
	// Tags      Tags           `json:"tags" gorm:"serializer:json"`
}

func (mo *MapObject) String() string {
	return fmt.Sprintf("GameID: %5d, ID: %5d, Num: %3d %s", mo.GameID, mo.ID, mo.Num, mo.Name)
}

func (mo *MapObject) Owned() bool {
	return mo.PlayerNum != nil
}

func (mo *MapObject) OwnedBy(num int) bool {
	return mo.PlayerNum != nil && *mo.PlayerNum == num
}
