package db

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db/generated"
)

// ensure we are using indexes for our complex queries
func Test_client_explainQuery(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	tests := []struct {
		name      string
		query     string
		queryArgs []any
		want      []string
	}{
		{
			name:      "GetPlanetByNum",
			query:     generated.GetPlanetByNum,
			queryArgs: []any{1, 1},
			want: []string{
				"SEARCH f USING INDEX idx_fleets_planet_num (gameId=? AND planetNum=?) LEFT-JOIN",
				"SEARCH p USING INDEX sqlite_autoindex_planets_1 (gameId=? AND num=?)",
			},
		},
		{
			name:      "GetPlanetsForGame",
			query:     generated.GetPlanetsForGame,
			queryArgs: []any{1},
			want:      []string{"SEARCH planets USING INDEX sqlite_autoindex_planets_1 (gameId=?)"},
		},
		{
			name:      "GetPlayersForGame",
			query:     generated.GetPlayersForGame,
			queryArgs: []any{1},
			want:      []string{"SEARCH players USING INDEX idx_players_gameid_num (gameId=?)"},
		},
		{
			name:      "GetPlayerForGame",
			query:     generated.GetPlayerForGame,
			queryArgs: []any{1, 1},
			want: []string{
				"SEARCH d USING INDEX sqlite_autoindex_shipDesigns_1 (gameId=? AND playerNum=?) LEFT-JOIN",
				"SEARCH p USING INDEX sqlite_autoindex_players_1 (gameId=? AND num=?)",
			},
		},
		{
			name:      "GetPlayerForGameAndUser",
			query:     generated.GetPlayerForGameAndUser,
			queryArgs: []any{1, 1, 1},
			want: []string{
				"SEARCH d USING INDEX idx_ship_designs_gameid_playernum (gameId=? AND playerNum=?) LEFT-JOIN",
				"SEARCH p USING INDEX idx_players_userid_gameid (userId=? AND gameId=?)",
				"USE TEMP B-TREE FOR ORDER BY",
			},
		},
		{
			name:      "GetPlayersWithDesignsForGame",
			query:     generated.GetPlayersWithDesignsForGame,
			queryArgs: []any{1},
			want: []string{
				"SEARCH d USING INDEX sqlite_autoindex_shipDesigns_1 (gameId=? AND playerNum=?) LEFT-JOIN",
				"SEARCH p USING INDEX sqlite_autoindex_players_1 (gameId=?)",
			},
		},
		{
			name:      "GetLightPlayerForGame",
			query:     generated.GetLightPlayerForGame,
			queryArgs: []any{1, 1, 1},
			want:      []string{"SEARCH players USING INDEX idx_players_gameid_num (gameId=?)"},
		},
		{
			name:      "GetPlanetsForPlayer",
			query:     generated.GetPlanetsForPlayer,
			queryArgs: []any{1, 1},
			want:      []string{"SEARCH planets USING INDEX sqlite_autoindex_planets_1 (gameId=?)"},
		},
		{
			name:      "GetPlayersForGame",
			query:     generated.GetPlayersForGame,
			queryArgs: []any{1},
			want:      []string{"SEARCH players USING INDEX idx_players_gameid_num (gameId=?)"},
		},
		{
			name:      "GetGamesWithPlayers",
			query:     generated.GetGamesWithPlayers,
			queryArgs: []any{cs.GameStateSetup, true, true},
			want: []string{
				// TODO: use index
				"SCAN g",
				"SEARCH players USING INDEX idx_players_gameid_num (gameId=?) LEFT-JOIN",
			},
		},
		{
			name:      "GetGamesWithPlayersForUser",
			query:     generated.GetGamesWithPlayersForUser,
			queryArgs: []any{1},
			want: []string{
				"INDEX 1",
				"INDEX 2",
				"LIST SUBQUERY 1",
				"MULTI-INDEX OR",
				"SEARCH g USING INDEX idx_games_hostid_gameid (hostId=?)",
				"SEARCH g USING INTEGER PRIMARY KEY (rowid=?)",
				"SEARCH p USING COVERING INDEX idx_players_userid_gameid (userId=?)",
				"SEARCH players USING INDEX idx_players_gameid_num (gameId=?) LEFT-JOIN",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.explainQuery(t.Context(), tt.query, tt.queryArgs...)
			if err != nil {
				t.Fatalf("client.explainQuery() error = %v", err)
			}

			gotDetails := lo.Map(got, func(i QueryPlan, _ int) string { return i.Detail })
			slices.Sort(gotDetails)
			gotDetails = lo.Uniq(gotDetails)

			if !reflect.DeepEqual(gotDetails, tt.want) {
				t.Errorf("explainQuery() got: \n%s\nwant: \n%s", strings.Join(gotDetails, "\n"), strings.Join(tt.want, "\n"))
			}
		})
	}
}
