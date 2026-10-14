package cs

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlayerMail_Validate(t *testing.T) {
	tests := []struct {
		name    string
		mail    PlayerMail
		want    string
		wantErr bool
	}{
		{name: "text is trimmed", mail: PlayerMail{FromPlayerNum: 1, ToPlayerNum: 2, Text: "  truce? \n"}, want: "truce?"},
		{name: "to everybody", mail: PlayerMail{FromPlayerNum: 1, ToPlayerNum: PlayerMailToEverybody, Text: "hi all"}, want: "hi all"},
		{name: "to yourself", mail: PlayerMail{FromPlayerNum: 1, ToPlayerNum: 1, Text: "remember"}, want: "remember"},
		{name: "longest", mail: PlayerMail{FromPlayerNum: 1, ToPlayerNum: 2, Text: strings.Repeat("é", PlayerMailMaxLength)}, want: strings.Repeat("é", PlayerMailMaxLength)},
		{name: "empty text", mail: PlayerMail{FromPlayerNum: 1, ToPlayerNum: 2, Text: "  "}, wantErr: true},
		{name: "too long", mail: PlayerMail{FromPlayerNum: 1, ToPlayerNum: 2, Text: strings.Repeat("a", PlayerMailMaxLength+1)}, wantErr: true},
		{name: "no such player", mail: PlayerMail{FromPlayerNum: 1, ToPlayerNum: 4, Text: "hi"}, wantErr: true},
		{name: "negative player", mail: PlayerMail{FromPlayerNum: 1, ToPlayerNum: -1, Text: "hi"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.mail.Validate(3)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, tt.mail.Text)
		})
	}
}

func TestPlayerMail_SentTo(t *testing.T) {
	direct := PlayerMail{FromPlayerNum: 1, ToPlayerNum: 2}
	assert.True(t, direct.SentTo(2))
	assert.False(t, direct.SentTo(1))
	assert.False(t, direct.SentTo(3))

	everybody := PlayerMail{FromPlayerNum: 1, ToPlayerNum: PlayerMailToEverybody}
	assert.False(t, everybody.SentTo(1))
	assert.True(t, everybody.SentTo(2))
	assert.True(t, everybody.SentTo(3))

	self := PlayerMail{FromPlayerNum: 1, ToPlayerNum: 1}
	assert.True(t, self.SentTo(1))
}

func Test_turn_deliverMail(t *testing.T) {
	s := TestScenario{Players: []ScenarioPlayer{
		{Player: NewPlayer(1, NewRace())},
		{Player: NewPlayer(2, NewRace())},
		{Player: NewPlayer(3, NewRace())},
	}}
	u := newTestUniverse(t, s)
	u.Player(1).PendingMail = []*PlayerMail{
		{FromPlayerNum: 1, ToPlayerNum: 2, Text: "truce?"},
		{FromPlayerNum: 1, ToPlayerNum: PlayerMailToEverybody, Text: "hi all"},
	}
	u.Player(3).PendingMail = []*PlayerMail{{FromPlayerNum: 3, ToPlayerNum: 3, Text: "remember"}}
	u.Run((*turnGenerator).deliverMail)

	// the sender doesn't get their own mail to everybody
	assert.Empty(t, u.Messages(1, PlayerMessagePlayerMail))

	mail := u.Messages(2, PlayerMessagePlayerMail)
	require.Len(t, mail, 2)
	assert.Equal(t, "truce?", mail[0].Spec.Note)
	assert.Equal(t, 1, mail[0].Spec.SourcePlayerNum)
	assert.Equal(t, 2, mail[0].Spec.DestPlayerNum)
	assert.Equal(t, "hi all", mail[1].Spec.Note)
	assert.Equal(t, PlayerMailToEverybody, mail[1].Spec.DestPlayerNum)

	mail = u.Messages(3, PlayerMessagePlayerMail)
	require.Len(t, mail, 2)
	assert.Equal(t, "hi all", mail[0].Spec.Note)
	assert.Equal(t, "remember", mail[1].Spec.Note)

	// mail from a player you haven't met tells you who they are
	assert.True(t, u.Player(2).PlayerIntels[0].Seen)
	assert.True(t, u.Player(3).PlayerIntels[0].Seen)
	assert.Len(t, u.Messages(2, PlayerMessagePlayerDiscovery), 1)
	assert.False(t, u.Player(1).PlayerIntels[1].Seen)

	assert.Empty(t, u.Player(1).PendingMail)
	assert.Empty(t, u.Player(3).PendingMail)
}
