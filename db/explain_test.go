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
				"SEARCH f USING INDEX idx_fleets_gameid (game_id=?) LEFT-JOIN",
				"SEARCH p USING INDEX sqlite_autoindex_planets_1 (game_id=? AND num=?)",
			},
		},
		{
			name:      "GetPlanetsForGame",
			query:     generated.GetPlanetsForGame,
			queryArgs: []any{1},
			want:      []string{"SEARCH planets USING INDEX sqlite_autoindex_planets_1 (game_id=?)"},
		},
		{
			name:      "GetPlayersForGame",
			query:     generated.GetPlayersForGame,
			queryArgs: []any{1},
			want:      []string{"SEARCH players USING INDEX idx_players_gameid_num (game_id=?)"},
		},
		{
			name:      "GetPlayerForGame",
			query:     generated.GetPlayerForGame,
			queryArgs: []any{1, 1},
			want: []string{
				"SEARCH d USING INDEX sqlite_autoindex_ship_designs_1 (game_id=? AND player_num=?) LEFT-JOIN",
				"SEARCH p USING INDEX sqlite_autoindex_players_1 (game_id=? AND num=?)",
			},
		},
		{
			name:      "GetPlayersWithDesignsForGame",
			query:     generated.GetPlayersWithDesignsForGame,
			queryArgs: []any{1},
			want: []string{
				"SEARCH d USING INDEX sqlite_autoindex_ship_designs_1 (game_id=? AND player_num=?) LEFT-JOIN",
				"SEARCH p USING INDEX sqlite_autoindex_players_1 (game_id=?)",
			},
		},
		{
			name:      "GetLightPlayerForGame",
			query:     generated.GetLightPlayerForGame,
			queryArgs: []any{1, 1, 1},
			want:      []string{"SEARCH players USING INDEX idx_players_gameid_num (game_id=?)"},
		},
		{
			name:      "GetPlanetsForPlayer",
			query:     generated.GetPlanetsForPlayer,
			queryArgs: []any{1, 1},
			want:      []string{"SEARCH planets USING INDEX sqlite_autoindex_planets_1 (game_id=?)"},
		},
		{
			name:      "GetPlayersForGame",
			query:     generated.GetPlayersForGame,
			queryArgs: []any{1},
			want:      []string{"SEARCH players USING INDEX idx_players_gameid_num (game_id=?)"},
		},
		{
			name:      "GetGamesWithPlayers",
			query:     generated.GetGamesWithPlayers,
			queryArgs: []any{cs.GameStateSetup, true, true},
			want: []string{
				// TODO: use index
				"SCAN g",
				"SEARCH players USING INDEX idx_players_gameid_num (game_id=?) LEFT-JOIN",
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
				"REUSE LIST SUBQUERY 1",
				"SEARCH g USING INDEX idx_games_hostid_gameid (host_id=?)",
				"SEARCH g USING INTEGER PRIMARY KEY (rowid=?)",
				"SEARCH p USING COVERING INDEX idx_players_userid_gameid (user_id=?)",
				"SEARCH players USING INDEX idx_players_gameid_num (game_id=?) LEFT-JOIN",
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
