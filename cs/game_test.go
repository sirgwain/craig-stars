package cs

type testPlayerGetter struct {
	players []*Player
}

type testPlanetGetter struct {
	planets []*Planet
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

func newTestPlanetGetter(planets ...*Planet) planetGetter {
	return &testPlanetGetter{planets}
}

func (pg *testPlanetGetter) getPlanet(num int) *Planet {
	for _, p := range pg.planets {
		if p.Num == num {
			return p
		}
	}
	return nil
}
