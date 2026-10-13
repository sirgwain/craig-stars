package cs

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// the longest mail a player can write
const PlayerMailMaxLength = 1000

// the most mail a player can write in one year
const PlayerMailMaxPerYear = 50

// mail with this ToPlayerNum goes to every other player in the game
const PlayerMailToEverybody = 0

// A PlayerMail is a message one player writes to another player, or to everybody. Mail is written
// during a year and delivered when the next turn is generated, so the sender can change or delete it
// until then.
type PlayerMail struct {
	GameDBObject
	Year          int    `json:"year"`
	FromPlayerNum int    `json:"fromPlayerNum"`
	ToPlayerNum   int    `json:"toPlayerNum"`
	Text          string `json:"text"`
}

// Delivered reports whether this mail was sent to its recipients by a turn generated before year
func (m *PlayerMail) Delivered(year int) bool {
	return m.Year < year
}

// SentTo reports whether a player receives this mail. Mail to everybody doesn't come back to the sender.
func (m *PlayerMail) SentTo(playerNum int) bool {
	if m.ToPlayerNum == PlayerMailToEverybody {
		return m.FromPlayerNum != playerNum
	}
	return m.ToPlayerNum == playerNum
}

// Validate cleans up mail from a player and returns an error if it can't be saved
func (m *PlayerMail) Validate(numPlayers int) error {
	m.Text = strings.TrimSpace(m.Text)
	if m.Text == "" {
		return fmt.Errorf("message is empty")
	}
	if utf8.RuneCountInString(m.Text) > PlayerMailMaxLength {
		return fmt.Errorf("message is longer than %d characters", PlayerMailMaxLength)
	}
	if m.ToPlayerNum < PlayerMailToEverybody || m.ToPlayerNum > numPlayers {
		return fmt.Errorf("message recipient is invalid")
	}
	return nil
}
