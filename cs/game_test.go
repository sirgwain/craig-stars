package cs

type testPlayerGetter struct {
	players []*Player
}

func newTestPlayerGetter(players ...*Player) playerGetter {
	return &testPlayerGetter{players}
}

func (pg *testPlayerGetter) getPlayer(num int) *Player {
	for _, p := range pg.players {
		if p.Num == num {
			return p
		}
	}
	return nil
}
