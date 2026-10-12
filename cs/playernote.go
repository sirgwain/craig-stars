package cs

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// the longest note a player can write
const PlayerNoteMaxLength = 2000

// the most notes a player can have in a game
const PlayerNoteMaxCount = 1000

// the longest tag a player can make
const PlayerNoteTagMaxLength = 30

// the most tags a player can put on one map object
const PlayerNoteMaxTags = 20

// A PlayerNote is private text a player writes for themselves. A note with a target is attached to
// a map object, i.e. a reminder about a planet, and can also tag it. A note without a target is a
// label placed at a Position on the map.
type PlayerNote struct {
	GameDBObject
	PlayerNum       int           `json:"playerNum"`
	TargetType      MapObjectType `json:"targetType,omitempty"`
	TargetNum       int           `json:"targetNum,omitempty"`
	TargetPlayerNum int           `json:"targetPlayerNum,omitempty"`
	Position        Vector        `json:"position"`
	Text            string        `json:"text"`
	Tags            []string      `json:"tags,omitempty"`
	ShowOnMap       bool          `json:"showOnMap,omitempty"`
}

// HasTarget reports whether this note is attached to a map object rather than a position on the map
func (n *PlayerNote) HasTarget() bool {
	return n.TargetType != MapObjectTypeNone
}

// Validate cleans up a note from a player and returns an error if it can't be saved
func (n *PlayerNote) Validate(area Vector) error {
	n.Text = strings.TrimSpace(n.Text)
	if len(n.Text) > PlayerNoteMaxLength {
		return fmt.Errorf("note is longer than %d characters", PlayerNoteMaxLength)
	}

	if n.HasTarget() {
		if err := n.cleanTags(); err != nil {
			return err
		}
		if n.Text == "" && len(n.Tags) == 0 {
			return fmt.Errorf("note is empty")
		}
		if n.TargetNum < 1 {
			return fmt.Errorf("note target is invalid")
		}
		switch n.TargetType {
		case MapObjectTypePlanet:
			// planets are keyed by num alone
			n.TargetPlayerNum = Unowned
		case MapObjectTypeFleet:
			// a player's fleet nums are reused, so fleet notes are deleted along with the fleet.
			// We don't know when another player's fleet is gone, so those can't have notes.
			if n.TargetPlayerNum != n.PlayerNum {
				return fmt.Errorf("notes can only be attached to your own fleets")
			}
			// fleets move, their notes aren't drawn on the map
			n.ShowOnMap = false
		default:
			return fmt.Errorf("notes can't be attached to a %s", n.TargetType)
		}
		// only the text is drawn on the map
		n.ShowOnMap = n.ShowOnMap && n.Text != ""
		return nil
	}

	if n.Text == "" {
		return fmt.Errorf("note is empty")
	}
	if n.Position.X < 0 || n.Position.Y < 0 || n.Position.X > area.X || n.Position.Y > area.Y {
		return fmt.Errorf("note position %v is outside the universe", n.Position)
	}
	// labels are only useful on the map, and there is no map object to tag
	n.ShowOnMap = true
	n.Tags = nil
	return nil
}

// cleanTags trims the note's tags and removes blanks and duplicates. Tags that differ only by
// case are the same tag, the first spelling wins.
func (n *PlayerNote) cleanTags() error {
	tags := make([]string, 0, len(n.Tags))
	seen := make(map[string]bool, len(n.Tags))
	for _, tag := range n.Tags {
		tag = strings.Join(strings.Fields(tag), " ")
		if tag == "" || seen[strings.ToLower(tag)] {
			continue
		}
		if utf8.RuneCountInString(tag) > PlayerNoteTagMaxLength {
			return fmt.Errorf("tag %q is longer than %d characters", tag, PlayerNoteTagMaxLength)
		}
		seen[strings.ToLower(tag)] = true
		tags = append(tags, tag)
	}
	if len(tags) > PlayerNoteMaxTags {
		return fmt.Errorf("a note can have at most %d tags", PlayerNoteMaxTags)
	}
	if len(tags) == 0 {
		tags = nil
	}
	n.Tags = tags
	return nil
}
