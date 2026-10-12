package cs

import (
	"fmt"
	"strings"
)

// the longest note a player can write
const PlayerNoteMaxLength = 2000

// the most notes a player can have in a game
const PlayerNoteMaxCount = 1000

// A PlayerNote is private text a player writes for themselves. A note with a target is attached to
// a map object, i.e. a reminder about a planet. A note without a target is a label placed at a
// Position on the map.
type PlayerNote struct {
	GameDBObject
	PlayerNum       int           `json:"playerNum"`
	TargetType      MapObjectType `json:"targetType,omitempty"`
	TargetNum       int           `json:"targetNum,omitempty"`
	TargetPlayerNum int           `json:"targetPlayerNum,omitempty"`
	Position        Vector        `json:"position"`
	Text            string        `json:"text"`
	ShowOnMap       bool          `json:"showOnMap,omitempty"`
}

// HasTarget reports whether this note is attached to a map object rather than a position on the map
func (n *PlayerNote) HasTarget() bool {
	return n.TargetType != MapObjectTypeNone
}

// Validate cleans up a note from a player and returns an error if it can't be saved
func (n *PlayerNote) Validate(area Vector) error {
	n.Text = strings.TrimSpace(n.Text)
	if n.Text == "" {
		return fmt.Errorf("note is empty")
	}
	if len(n.Text) > PlayerNoteMaxLength {
		return fmt.Errorf("note is longer than %d characters", PlayerNoteMaxLength)
	}

	if n.HasTarget() {
		// planet nums never change. Fleet and other map object nums are reused after the object
		// is gone, so a note on one would end up attached to something else.
		if n.TargetType != MapObjectTypePlanet {
			return fmt.Errorf("notes can't be attached to a %s", n.TargetType)
		}
		if n.TargetNum < 1 {
			return fmt.Errorf("note target is invalid")
		}
		// planets are keyed by num alone
		n.TargetPlayerNum = Unowned
		return nil
	}

	if n.Position.X < 0 || n.Position.Y < 0 || n.Position.X > area.X || n.Position.Y > area.Y {
		return fmt.Errorf("note position %v is outside the universe", n.Position)
	}
	// labels are only useful on the map
	n.ShowOnMap = true
	return nil
}
