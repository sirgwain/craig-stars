package server

import (
	"context"
	"encoding/json"
	"testing"

	"connectrpc.com/connect"
	mcpruntime "github.com/redpanda-data/protoc-gen-go-mcp/pkg/runtime"
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
	"github.com/sirgwain/craig-stars/proto/converter"
	v1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1/craig_starsv1mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestUpdateGame(t *testing.T) {
	mask := func(paths ...string) *fieldmaskpb.FieldMask { return &fieldmaskpb.FieldMask{Paths: paths} }
	scoreRequest := func() *v1.UpdateGameRequest {
		return &v1.UpdateGameRequest{
			Settings:   &v1.GameSettings{PublicPlayerScores: true},
			UpdateMask: mask("public_player_scores"),
		}
	}
	type testCase struct {
		name    string
		prepare func(*cs.GameWithPlayers)
		request func() *v1.UpdateGameRequest
		want    func(*cs.Game)
		code    connect.Code
		nonHost bool
	}
	tests := []testCase{
		{name: "enable public scores only", request: scoreRequest, want: func(g *cs.Game) { g.PublicPlayerScores = true }},
		{name: "disable public scores explicitly", prepare: func(g *cs.GameWithPlayers) { g.PublicPlayerScores = true },
			request: func() *v1.UpdateGameRequest {
				return &v1.UpdateGameRequest{Settings: &v1.GameSettings{}, UpdateMask: mask("public_player_scores")}
			},
			want: func(g *cs.Game) { g.PublicPlayerScores = false }},
		{name: "rename and disable random events", request: func() *v1.UpdateGameRequest {
			return &v1.UpdateGameRequest{Settings: &v1.GameSettings{Name: "New name"}, UpdateMask: mask("name", "random_events")}
		}, want: func(g *cs.Game) { g.Name = "New name"; g.RandomEvents = false }},
		{name: "enable random events", prepare: func(g *cs.GameWithPlayers) { g.RandomEvents = false }, request: func() *v1.UpdateGameRequest {
			return &v1.UpdateGameRequest{Settings: &v1.GameSettings{RandomEvents: true}, UpdateMask: mask("random_events")}
		}, want: func(g *cs.Game) { g.RandomEvents = true }},
		{name: "replace victory conditions", request: func() *v1.UpdateGameRequest {
			return &v1.UpdateGameRequest{Settings: &v1.GameSettings{VictoryConditions: &v1.VictoryConditions{Conditions: 1, NumCriteriaRequired: 1, OwnPlanets: 90, YearsPassed: 60}}, UpdateMask: mask("victory_conditions")}
		}, want: func(g *cs.Game) {
			g.VictoryConditions = cs.VictoryConditions{Conditions: 1, NumCriteriaRequired: 1, OwnPlanets: 90, YearsPassed: 60}
		}},
		{name: "unmasked settings ignored", request: func() *v1.UpdateGameRequest {
			r := scoreRequest()
			r.Settings.Public = true
			r.Settings.Size = v1.Size_SIZE_HUGE
			return r
		}, want: func(g *cs.Game) { g.PublicPlayerScores = true }},
		{name: "missing settings", request: func() *v1.UpdateGameRequest { return &v1.UpdateGameRequest{UpdateMask: mask("name")} }, code: connect.CodeInvalidArgument},
		{name: "missing mask after setup", request: func() *v1.UpdateGameRequest {
			return &v1.UpdateGameRequest{Settings: &v1.GameSettings{PublicPlayerScores: true}}
		}, code: connect.CodeInvalidArgument},
		{name: "empty mask after setup", request: func() *v1.UpdateGameRequest { r := scoreRequest(); r.UpdateMask = mask(); return r }, code: connect.CodeInvalidArgument},
		{name: "unknown field rejects entire update", request: func() *v1.UpdateGameRequest {
			r := scoreRequest()
			r.UpdateMask.Paths = append(r.UpdateMask.Paths, "unknown")
			return r
		}, code: connect.CodeInvalidArgument},
		{name: "nested victory path rejected", request: func() *v1.UpdateGameRequest {
			r := scoreRequest()
			r.UpdateMask = mask("victory_conditions.own_planets")
			return r
		}, code: connect.CodeInvalidArgument},
		{name: "wildcard rejected", request: func() *v1.UpdateGameRequest { r := scoreRequest(); r.UpdateMask = mask("*"); return r }, code: connect.CodeInvalidArgument},
		{name: "players rejected", request: func() *v1.UpdateGameRequest { r := scoreRequest(); r.UpdateMask = mask("players"); return r }, code: connect.CodeInvalidArgument},
		{name: "non host rejected", request: scoreRequest, nonHost: true, code: connect.CodePermissionDenied},
		{name: "multiplayer rejected", prepare: func(g *cs.GameWithPlayers) { g.Players[1].AIControlled = false; g.Players[1].UserID = 999 }, request: scoreRequest, code: connect.CodeFailedPrecondition},
		{name: "unclaimed human slot rejected", prepare: func(g *cs.GameWithPlayers) { g.Players[1].AIControlled = false }, request: scoreRequest, code: connect.CodeFailedPrecondition},
		{name: "same user hot seat allowed", prepare: func(g *cs.GameWithPlayers) { g.Players[1].AIControlled = false; g.Players[1].UserID = g.HostID }, request: scoreRequest, want: func(g *cs.Game) { g.PublicPlayerScores = true }},
		{name: "victory conditions locked after victory", prepare: func(g *cs.GameWithPlayers) { g.VictorDeclared = true }, request: func() *v1.UpdateGameRequest {
			r := scoreRequest()
			r.UpdateMask.Paths = append(r.UpdateMask.Paths, "victory_conditions")
			return r
		}, code: connect.CodeFailedPrecondition},
		{name: "other settings allowed after victory", prepare: func(g *cs.GameWithPlayers) { g.VictorDeclared = true }, request: scoreRequest, want: func(g *cs.Game) { g.PublicPlayerScores = true }},
		{name: "setup full replacement remains supported", prepare: func(g *cs.GameWithPlayers) {
			g.State = cs.GameStateSetup
			g.Players[1].AIControlled = false
			g.Players[1].UserID = 999
		}, request: func() *v1.UpdateGameRequest {
			return &v1.UpdateGameRequest{Settings: &v1.GameSettings{Name: "Setup", Public: true, Size: v1.Size_SIZE_HUGE, Density: v1.Density_DENSITY_PACKED, PlayerPositions: v1.PlayerPositions_PLAYER_POSITIONS_DISTANT, ComputerPlayersFormAlliances: true, MaxMinerals: true, GalaxyClumping: true, StartMode: v1.GameStartMode_GAME_START_MODE_ACC_BBS, QuickStartTurns: 10}}
		}, want: func(g *cs.Game) {
			g.WithSettings(cs.GameSettings{Name: "Setup", Public: true, Size: cs.SizeHuge, Density: cs.DensityPacked, PlayerPositions: cs.PlayerPositionsDistant, ComputerPlayersFormAlliances: true, MaxMinerals: true, GalaxyClumping: true, StartMode: cs.GameStartModeAccBBS, QuickStartTurns: 10})
		}},
		{name: "setup masked update preserves other settings", prepare: func(g *cs.GameWithPlayers) { g.State = cs.GameStateSetup }, request: func() *v1.UpdateGameRequest {
			return &v1.UpdateGameRequest{Settings: &v1.GameSettings{Public: true}, UpdateMask: mask("public")}
		}, want: func(g *cs.Game) { g.Public = true }},
	}
	for _, path := range []string{"public", "size", "density", "player_positions", "computer_players_form_alliances", "max_minerals", "galaxy_clumping", "start_mode", "quick_start_turns"} {
		tests = append(tests, testCase{name: "locked setting " + path, request: func() *v1.UpdateGameRequest {
			r := scoreRequest()
			r.UpdateMask.Paths = append(r.UpdateMask.Paths, path)
			return r
		}, code: connect.CodeInvalidArgument})
	}
	for _, state := range []cs.GameState{cs.GameStateGeneratingUniverse, cs.GameStateGeneratingTurn, cs.GameStateGeneratingTurnError, cs.GameStateUnspecified} {
		tests = append(tests, testCase{name: "state " + string(state), prepare: func(g *cs.GameWithPlayers) { g.State = state }, request: scoreRequest, code: connect.CodeFailedPrecondition})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			conn := db.NewConn()
			cfg := config.Config{}
			cfg.Database.Filename = ":memory:"
			require.NoError(t, conn.Connect(ctx, &cfg))
			t.Cleanup(func() { conn.Close() })
			client := conn.NewReadWriteClient()
			host, err := client.CreateUser(ctx, cs.NewUser("host", "", "", cs.RoleUser))
			require.NoError(t, err)
			game := &cs.GameWithPlayers{Game: *cs.NewGame(), Players: []cs.GamePlayer{{UserID: host.ID, Num: 1}, {AIControlled: true, Num: 2}}}
			game.HostID = host.ID
			game.State = cs.GameStateWaitingForPlayers
			game.Year = 2429
			game.NumPlayers = 2
			if tt.prepare != nil {
				tt.prepare(game)
			}
			require.NoError(t, client.SaveGame(ctx, &game.Game))
			for _, p := range game.Players {
				player := cs.NewPlayer(p.UserID, cs.NewRace()).WithNum(p.Num).WithAIControlled(p.AIControlled)
				player.GameID = game.ID
				require.NoError(t, client.SavePlayer(ctx, player))
			}
			before, err := client.GetGame(ctx, game.ID)
			require.NoError(t, err)
			// Deliberately stale request context: updates must use persisted state,
			// settings and player ownership, not the earlier interceptor snapshot.
			stale := *before
			stale.State = cs.GameStateWaitingForPlayers
			stale.Year--
			stale.Players = []cs.GamePlayer{{UserID: host.ID, Num: 1}, {AIControlled: true, Num: 2}}
			ctx = context.WithValue(ctx, keyGame, &stale)
			ctx = context.WithValue(ctx, keyDbRead, conn.NewReadClient())
			ctx = context.WithValue(ctx, keyDbWrite, client)
			userID := host.ID
			if tt.nonHost {
				userID++
			}
			ctx = context.WithValue(ctx, keyUserSession, userSession{ID: userID})
			request := tt.request()
			request.GameId = game.ID
			response, err := NewGameServiceHandler(conn, cfg, nil).UpdateGame(ctx, connect.NewRequest(request))
			if tt.code != 0 {
				require.Error(t, err)
				assert.Equal(t, tt.code, connect.CodeOf(err))
			} else {
				require.NoError(t, err)
			}
			after, err := client.GetGame(ctx, game.ID)
			require.NoError(t, err)
			if response != nil {
				assert.True(t, proto.Equal(converter.C.ConvertCSGameWithPlayers(after), response.Msg.Game))
			}
			want := before.Game
			if tt.want != nil {
				tt.want(&want)
			}
			// GetGame constructs fresh default rules (including a random source)
			// on every read; these are not settings persisted by UpdateGame.
			after.Rules = want.Rules
			// Ignore database-managed timestamps when comparing persisted settings.
			after.Game.DBObject = want.DBObject
			assert.Equal(t, want, after.Game)
			assert.Equal(t, before.Players, after.Players)
		})
	}
}

func TestUpdateGameFieldMaskJSON(t *testing.T) {
	// FieldMask uses a comma-separated lowerCamelCase string in protobuf JSON.
	// This is the representation consumed by MCP's protojson decoder.
	var request v1.UpdateGameRequest
	var schema struct {
		Properties map[string]struct{ Type string }
	}
	require.NoError(t, json.Unmarshal(craig_starsv1mcp.GameService_UpdateGameTool.RawInputSchema, &schema))
	assert.Equal(t, "string", schema.Properties["update_mask"].Type)
	args := map[string]any{
		"game_id":     "123",
		"settings":    map[string]any{"public_player_scores": false, "random_events": false},
		"update_mask": "publicPlayerScores,randomEvents",
	}
	require.NoError(t, mcpruntime.DecodeArguments(request.ProtoReflect().Descriptor(), args))
	payload, err := json.Marshal(args)
	require.NoError(t, err)
	require.NoError(t, protojson.Unmarshal(payload, &request))
	assert.Equal(t, []string{"public_player_scores", "random_events"}, request.UpdateMask.Paths)
	game := &cs.GameWithPlayers{Game: *cs.NewGame(), Players: []cs.GamePlayer{{UserID: 1}, {AIControlled: true}}}
	game.State = cs.GameStateWaitingForPlayers
	game.PublicPlayerScores = true
	require.NoError(t, applyGameSettingsUpdate(game, &request))
	assert.False(t, game.PublicPlayerScores)
	assert.False(t, game.RandomEvents)
}
