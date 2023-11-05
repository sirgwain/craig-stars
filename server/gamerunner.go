package server

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/ai"
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
)

var errNotFound = errors.New("resource was not found")

type TurnGenerationCheckResult uint
type DBConnection db.DBConn
type DBClient db.Client

const (
	TurnNotGenerated TurnGenerationCheckResult = iota
	TurnGenerated
)

var colors = []string{
	"#0000FF",
	"#C33232",
	"#1F8BA7",
	"#43A43E",
	"#8D29CB",
	"#B88628",
	"#FF4500",
	"#FF8C00",
	"#008000",
	"#00FA9A",
	"#7FFFD4",
	"#8A2BE2",
	"#FF1493",
	"#D2691E",
	"#F0FFF0",
}

// The GameRunner handles hosting games, updating players of games before they are started, loading a full Player game from
// the database, and generating new turns.
type GameRunner interface {
	HostGame(hostID int64, settings *cs.GameSettings) (*cs.FullGame, error)
	JoinGame(gameID int64, userID int64, name string, race cs.Race) error
	LeaveGame(gameID, userID int64) error
	KickPlayer(gameID int64, playerNum int) error
	AddOpenPlayerSlot(game *cs.GameWithPlayers) (*cs.Player, error)
	AddGuestPlayer(game *cs.GameWithPlayers) (*cs.Player, error)
	AddAIPlayer(game *cs.GameWithPlayers) (*cs.Player, error)
	DeletePlayerSlot(gameID int64, playerNum int) error
	StartGame(game *cs.Game) error
	LoadPlayerGame(gameID int64, userID int64) (*cs.GameWithPlayers, *cs.FullPlayer, error)
	SubmitTurn(gameID int64, userID int64) error
	CheckAndGenerateTurn(gameID int64) (TurnGenerationCheckResult, error)
	GenerateTurn(gameID int64) (TurnGenerationCheckResult, error)
}

type gameRunner struct {
	dbConn DBConnection
	client cs.Gamer
	config config.Config
}

func NewGameRunner(dbConn DBConnection, config config.Config) GameRunner {
	return &gameRunner{dbConn, cs.NewGamer(), config}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

// extract the guest number from a guest username
func (gr *gameRunner) getGuestNum(u *cs.User) (int, error) {
	r := regexp.MustCompile(`\w+-\d+-(\d+)`)
	m := r.FindStringSubmatch(u.Username)

	if len(m) == 2 {
		return strconv.Atoi(m[1])
	} else {
		return 0, fmt.Errorf("unable to determine guest num from %s", u.Username)
	}
}

// host a new game
func (gr *gameRunner) HostGame(hostID int64, settings *cs.GameSettings) (*cs.FullGame, error) {
	game := gr.client.CreateGame(hostID, *settings)
	var fullGame *cs.FullGame

	if err := gr.dbConn.WrapInTransaction(func(c db.Client) error {
		user, err := c.GetUser(hostID)
		if err != nil {
			return err
		}

		// create the game in the database
		err = c.CreateGame(game)
		if err != nil {
			return err
		}

		// generate an invite hash for the game
		game.Hash = game.GenerateHash(gr.config.Game.InviteLinkSalt)
		if err := c.UpdateGame(game); err != nil {
			return err
		}

		game.NumPlayers = len(settings.Players)
		players := make([]*cs.Player, 0, len(settings.Players))
		guestNumber := 1

		for i, playerSetting := range settings.Players {
			if playerSetting.Type == cs.NewGamePlayerTypeHost {

				log.Debug().Int64("hostID", hostID).Msg("Adding host to game")
				player := gr.client.NewPlayer(hostID, playerSetting.Race, &game.Rules)
				player.GameID = game.ID
				player.Num = i + 1
				player.Name = user.Username
				player.Color = playerSetting.Color
				player.DefaultHullSet = playerSetting.DefaultHullSet
				player.Ready = true
				players = append(players, player)
			} else if playerSetting.Type == cs.NewGamePlayerTypeAI {
				log.Debug().Int64("hostID", hostID).Msg("Adding ai player to game")
				race := ai.GetRandomRace()
				if playerSetting.Race.Name != "" {
					race = playerSetting.Race
				}
				player := gr.client.NewPlayer(0, race, &game.Rules)
				player.GameID = game.ID
				player.Num = i + 1
				player.AIControlled = true
				player.Name = fmt.Sprintf("AI Player %d", player.Num)
				player.Color = playerSetting.Color
				player.DefaultHullSet = playerSetting.DefaultHullSet
				player.Ready = true
				players = append(players, player)
			} else if playerSetting.Type == cs.NewGamePlayerTypeOpen {
				log.Debug().Int("openPlayerSlots", game.OpenPlayerSlots).Msg("Added open player slot to game")
				race := cs.NewRace()
				player := gr.client.NewPlayer(0, *race, &game.Rules)
				player.GameID = game.ID
				player.Num = i + 1
				player.Name = "Open Slot"
				player.Color = playerSetting.Color
				player.DefaultHullSet = playerSetting.DefaultHullSet
				players = append(players, player)
				game.OpenPlayerSlots++
			} else if playerSetting.Type == cs.NewGamePlayerTypeGuest {
				// create a new guest user for this game
				playerNum := i + 1
				// username is based on game/number
				username := fmt.Sprintf("guest-%d-%d", game.ID, guestNumber)
				guestUser := cs.NewGuestUser(username, game.ID, playerNum)
				if err := c.CreateUser(guestUser); err != nil {
					return fmt.Errorf("failed to create guest user %w", err)
				}
				guestUser.GenerateHash(gr.config.Game.InviteLinkSalt)
				if err := c.UpdateUser(guestUser); err != nil {
					return fmt.Errorf("failed to update guest user %w", err)
				}

				race := cs.NewRace()
				player := gr.client.NewPlayer(guestUser.ID, *race, &game.Rules)
				player.GameID = game.ID
				player.Num = playerNum
				player.Name = "Guest"
				player.Guest = true
				player.UserID = guestUser.ID
				player.Color = playerSetting.Color
				player.DefaultHullSet = playerSetting.DefaultHullSet
				players = append(players, player)
				guestNumber++
			}
		}

		for _, player := range players {
			if player.Color == "" || (player.Num != 1 && player.Color == "#0000FF") {
				if player.Num-1 < len(colors) {
					player.Color = colors[player.Num-1]
				} else {
					color := make([]byte, 3)
					rand.Read(color)
					player.Color = fmt.Sprintf("#%s", hex.EncodeToString(color))
				}
			}
		}

		// load the tech store for this game
		techStore, err := gr.loadTechStore(c, game)
		if err != nil {
			return err
		}

		universe := cs.NewUniverse(&game.Rules)
		fullGame = &cs.FullGame{
			Game:      game,
			Universe:  &universe,
			TechStore: techStore,
			Players:   players,
		}

		// generate the universe if this game is ready
		if fullGame.IsSinglePlayer() {
			err = gr.generateUniverse(fullGame)
			if err != nil {
				return err
			}
		}

		// save the game
		if err := c.UpdateFullGame(fullGame); err != nil {
			return fmt.Errorf("update full game %w", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return fullGame, nil
}

// add a player to an existing game
func (gr *gameRunner) JoinGame(gameID int64, userID int64, name string, race cs.Race) error {
	readClient := gr.dbConn.NewReadClient()
	user, err := readClient.GetUser(userID)
	if err != nil {
		return fmt.Errorf("unable to load user %d: %w", userID, err)
	}
	if user == nil {
		return fmt.Errorf("no user for id %d found. %w", userID, errNotFound)
	}

	fullGame, err := gr.loadGame(readClient, gameID)
	if err != nil {
		return fmt.Errorf("unable to load game %d: %w", gameID, err)
	}

	if fullGame == nil {
		return fmt.Errorf("no game for id %d found. %w", gameID, errNotFound)
	}

	invitedGuest := false
	if user.Role == cs.RoleGuest {
		for _, player := range fullGame.Players {
			if player.UserID == user.ID {
				invitedGuest = true
			}
		}

		if !invitedGuest {
			return fmt.Errorf("guests can't join games they aren't invited to")
		}
	}

	if !invitedGuest && fullGame.OpenPlayerSlots == 0 {
		return fmt.Errorf("no slots left in this game")
	}

	player := gr.client.NewPlayer(userID, race, &fullGame.Rules)
	player.GameID = fullGame.ID
	player.Ready = true
	player.AIControlled = false
	player.Guest = invitedGuest
	player.Name = name

	if err := gr.dbConn.WrapInTransaction(func(c db.Client) error {

		// claim an open slot
		for i, p := range fullGame.Players {
			// if this is an invited guest, take over the invite user
			// otherwise, take over an open slot
			if (invitedGuest && p.UserID == user.ID) || (!invitedGuest && !p.AIControlled && p.UserID == 0) {
				// take over this empty player
				player.Num = p.Num
				player.ID = p.ID
				player.UserID = user.ID
				if player.Num-1 < len(colors) {
					player.Color = colors[player.Num-1]
				} else {
					color := make([]byte, 3)
					rand.Read(color)
					player.Color = fmt.Sprintf("#%s", hex.EncodeToString(color))
				}

				if err := c.UpdatePlayer(player); err != nil {
					return fmt.Errorf("update open slot player %s for game %d: %w", p, gameID, err)
				}

				fullGame.Players[i] = player

				if !invitedGuest {
					fullGame.OpenPlayerSlots--
				}

				if err := c.UpdateGame(fullGame.Game); err != nil {
					return fmt.Errorf("save game %d: %w", gameID, err)
				}

				break
			}
		}
		return nil
	}); err != nil {
		return err
	}

	log.Info().Int64("GameID", gameID).Int64("UserID", userID).Str("Race", player.Race.Name).Msgf("Joined game %s", fullGame.Name)
	return nil
}

// delete a guest user and any single player games or races they've created
func (gr *gameRunner) deleteGuestUser(c DBClient, user *cs.User) error {
	if err := c.DeleteUser(user.ID); err != nil {
		return fmt.Errorf("delete guest user %d %w", user.ID, err)
	}
	if err := c.DeleteUserGames(user.ID); err != nil {
		return fmt.Errorf("delete guest user %d games %w", user.ID, err)
	}
	if err := c.DeleteUserRaces(user.ID); err != nil {
		return fmt.Errorf("delete guest user %d races %w", user.ID, err)
	}
	return nil
}

// reset the player colors if players are kicked or leave
func (gr *gameRunner) resetPlayerColors(c DBClient, game *cs.FullGame) error {
	for _, player := range game.Players {
		if player.Num-1 < len(colors) {
			player.Color = colors[player.Num-1]
		} else {
			color := make([]byte, 3)
			rand.Read(color)
			player.Color = fmt.Sprintf("#%s", hex.EncodeToString(color))
		}
		if err := c.UpdatePlayer(player); err != nil {
			return fmt.Errorf("update player %s color for game %d: %w", player, game.ID, err)
		}
	}
	return nil
}

// add a player to an existing game
func (gr *gameRunner) LeaveGame(gameID, userID int64) error {

	game, err := gr.loadGame(gr.dbConn.NewReadClient(), gameID)
	if err != nil {
		return fmt.Errorf("unable to load game %d: %w", gameID, err)
	}

	if game == nil {
		return fmt.Errorf("no game for id %d found. %w", gameID, errNotFound)
	}

	readClient := gr.dbConn.NewReadClient()
	user, err := readClient.GetUser(userID)
	if err != nil {
		return fmt.Errorf("no user for %d %w", userID, err)
	}

	if user.IsGuest() {
		// guests can leave and join again
		if err := gr.dbConn.WrapInTransaction(func(c db.Client) error {
			for _, player := range game.Players {
				if player.UserID == userID {
					player.Ready = false
					if err := c.UpdatePlayer(player); err != nil {
						return fmt.Errorf("update open slot player %s for game %d: %w", player, gameID, err)
					}
				}
			}
			return nil
		}); err != nil {
			return err
		}
	} else {
		// open player slots delete the player and recreate it as an open player
		if err := gr.dbConn.WrapInTransaction(func(c db.Client) error {
			for i, player := range game.Players {
				if player.UserID == userID {
					if err := c.DeletePlayer(player.ID); err != nil {
						return fmt.Errorf("delete open slot player %s from game %d: %w", player, gameID, err)
					}

					race := cs.NewRace()
					player := gr.client.NewPlayer(0, *race, &game.Rules)
					player.GameID = game.ID
					player.UserID = 0
					player.Num = i + 1
					player.Name = "Open Slot"
					player.Guest = false
					game.Players[i] = player
					game.OpenPlayerSlots++

					if err := c.CreatePlayer(player); err != nil {
						return fmt.Errorf("update open slot player %s for game %d: %w", player, gameID, err)
					}

				}
			}

			if err := gr.resetPlayerColors(c, game); err != nil {
				return fmt.Errorf("update player colors %d: %w", gameID, err)
			}
			if err := c.UpdateGame(game.Game); err != nil {
				return fmt.Errorf("save game %d: %w", gameID, err)
			}

			return nil
		}); err != nil {
			return err
		}
	}

	log.Info().Int64("GameID", gameID).Int64("UserID", userID).Msgf("Left game %s", game.Name)

	return nil
}

// add a player to an existing game
func (gr *gameRunner) KickPlayer(gameID int64, playerNum int) error {

	game, err := gr.loadGame(gr.dbConn.NewReadClient(), gameID)
	if err != nil {
		return fmt.Errorf("unable to load game %d: %w", gameID, err)
	}

	if game == nil {
		return fmt.Errorf("no game for id %d found. %w", gameID, errNotFound)
	}

	readClient := gr.dbConn.NewReadClient()
	player, err := readClient.GetPlayerByNum(gameID, playerNum)
	if err != nil {
		return fmt.Errorf("no player %d for %d %w", playerNum, gameID, err)
	}
	var user *cs.User
	if player.UserID != 0 {
		user, err = readClient.GetUser(player.UserID)
		if err != nil {
			return fmt.Errorf("no user for %d %w", player.UserID, err)
		}
	}

	if err := gr.dbConn.WrapInTransaction(func(c db.Client) error {

		for i, player := range game.Players {
			if player.Num == playerNum {
				if err := c.DeletePlayer(player.ID); err != nil {
					return fmt.Errorf("delete open slot player %s from game %d: %w", player, gameID, err)
				}

				if user != nil && user.IsGuest() {
					if err := gr.deleteGuestUser(c, user); err != nil {
						return err
					}
				}

				race := cs.NewRace()
				player := gr.client.NewPlayer(0, *race, &game.Rules)
				player.GameID = game.ID
				player.UserID = 0
				player.Num = i + 1
				player.Name = "Open Slot"
				player.Guest = false
				game.Players[i] = player
				game.OpenPlayerSlots++
				game.NumPlayers--

				if err := c.CreatePlayer(player); err != nil {
					return fmt.Errorf("update open slot player %s for game %d: %w", player, gameID, err)
				}		
			}
		}

		if err := gr.resetPlayerColors(c, game); err != nil {
			return fmt.Errorf("update player colors %d: %w", gameID, err)
		}

		if err := c.UpdateGame(game.Game); err != nil {
			return fmt.Errorf("save game %d: %w", gameID, err)
		}

		return nil
	}); err != nil {
		return err
	}

	log.Info().Int64("GameID", gameID).Int("PlayerNum", playerNum).Msgf("kicked player %d %s", playerNum, game.Name)

	return nil
}

// delete an entire player slot
func (gr *gameRunner) DeletePlayerSlot(gameID int64, playerNum int) error {

	game, err := gr.loadGame(gr.dbConn.NewReadClient(), gameID)
	if err != nil {
		return fmt.Errorf("unable to load game %d: %w", gameID, err)
	}

	if game == nil {
		return fmt.Errorf("no game for id %d found. %w", gameID, errNotFound)
	}

	if err := gr.dbConn.WrapInTransaction(func(c db.Client) error {
		deleted := false
		aiPlayerNum := 0
		for i, player := range game.Players {
			if player.Num == playerNum {
				if !player.AIControlled {
					if player.Guest {
						user, err := c.GetUser(player.UserID)
						if err != nil {
							return fmt.Errorf("no user for %d %w", player.UserID, err)
						}
						if err := gr.deleteGuestUser(c, user); err != nil {
							return err
						}
					} else {
						game.OpenPlayerSlots--
					}
				}

				if err := c.DeletePlayer(player.ID); err != nil {
					return fmt.Errorf("delete open slot player %s from game %d: %w", player, gameID, err)
				}
				deleted = true

				continue
			}
			if player.AIControlled {
				aiPlayerNum++
			}
			// move the other player's up
			// if we delete player 2 (i = 1), make player 3 player 2
			if deleted {
				player.Num = i
				player.UserID = 0
				// if the player was using a default color, use the previous one
				if len(colors) > i+1 && player.Color == colors[i+1] {
					player.Color = colors[i]
				}
				if player.AIControlled {
					player.Name = fmt.Sprintf("AI Player %d", aiPlayerNum+1)
				} else if player.Guest {
					// update this guest user with a new playerNum
					user, err := c.GetUser(player.UserID)
					if err != nil {
						return fmt.Errorf("no user for %d %w", player.UserID, err)
					}
					// update the guest user
					user.PlayerNum = player.Num
					if err := c.UpdateUser(user); err != nil {
						return fmt.Errorf("update guest user %d with new playerNum %w", user.ID, err)
					}
				}

				if err := c.UpdatePlayer(player); err != nil {
					return fmt.Errorf("update player %s for game %d: %w", player.Name, gameID, err)
				}
			}
		}

		if err := gr.resetPlayerColors(c, game); err != nil {
			return fmt.Errorf("update player colors %d: %w", gameID, err)
		}

		game.NumPlayers--
		if err := c.UpdateGame(game.Game); err != nil {
			return fmt.Errorf("save game %d: %w", gameID, err)
		}
		return nil
	}); err != nil {
		return err
	}

	log.Info().Int64("GameID", gameID).Int("PlayerNum", playerNum).Msgf("deleted player slot %d %s", playerNum, game.Name)

	return nil
}

// delete an entire player slot
func (gr *gameRunner) AddOpenPlayerSlot(game *cs.GameWithPlayers) (*cs.Player, error) {

	race := *cs.NewRace()
	player := gr.client.NewPlayer(0, race, &game.Rules)
	player.GameID = game.ID
	player.Num = len(game.Players) + 1
	player.Name = "Open Player"
	if len(colors) > player.Num {
		player.Color = colors[player.Num-1]
	}

	if err := gr.dbConn.WrapInTransaction(func(c db.Client) error {
		if err := c.CreatePlayer(player); err != nil {
			return fmt.Errorf("added slot player %s for game %d: %w", player, game.ID, err)
		}

		game.OpenPlayerSlots++
		game.NumPlayers++
		if err := c.UpdateGame(&game.Game); err != nil {
			return fmt.Errorf("updating open player slots for game %d: %w", game.ID, err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	log.Info().Int64("GameID", game.ID).Int("Num", player.Num).Msgf("added player slot %d %s", player.Num, game.Name)

	return player, nil
}

// delete an entire player slot
func (gr *gameRunner) AddGuestPlayer(game *cs.GameWithPlayers) (*cs.Player, error) {

	race := *cs.NewRace()
	player := gr.client.NewPlayer(0, race, &game.Rules)
	player.GameID = game.ID
	player.Num = len(game.Players) + 1
	player.Name = "Guest"
	player.Guest = true
	if len(colors) > player.Num {
		player.Color = colors[player.Num-1]
	}

	readClient := gr.dbConn.NewReadClient()
	users, err := readClient.GetGuestUsersForGame(game.ID)
	if err != nil {
		return nil, fmt.Errorf("load guest users for game %d %w", game.ID, err)
	}
	guestNumber := 1
	for _, user := range users {
		userGuestNumber, err := gr.getGuestNum(&user)
		if err != nil {
			return nil, err
		}
		// make sure our guestNumber is larger than the highest guestNum
		guestNumber = int(math.Max(float64(userGuestNumber+1), float64(guestNumber)))
	}

	if err := gr.dbConn.WrapInTransaction(func(c db.Client) error {

		// username is based on game/number
		username := fmt.Sprintf("guest-%d-%d", game.ID, guestNumber)
		guestUser := cs.NewGuestUser(username, game.ID, player.Num)
		if err := c.CreateUser(guestUser); err != nil {
			return fmt.Errorf("failed to create guest user %w", err)
		}
		guestUser.GenerateHash(gr.config.Game.InviteLinkSalt)
		if err := c.UpdateUser(guestUser); err != nil {
			return fmt.Errorf("failed to update guest user %w", err)
		}

		player.UserID = guestUser.ID

		if err := c.CreatePlayer(player); err != nil {
			return fmt.Errorf("added slot player %s for game %d: %w", player, game.ID, err)
		}

		game.NumPlayers++
		if err := c.UpdateGame(&game.Game); err != nil {
			return fmt.Errorf("updating open player slots for game %d: %w", game.ID, err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	log.Info().Int64("GameID", game.ID).Int("Num", player.Num).Msgf("added player slot %d %s", player.Num, game.Name)

	return player, nil
}

// delete an entire player slot
func (gr *gameRunner) AddAIPlayer(game *cs.GameWithPlayers) (*cs.Player, error) {

	race := *cs.NewRace()
	player := gr.client.NewPlayer(0, race, &game.Rules)
	player.GameID = game.ID
	player.Num = len(game.Players) + 1
	player.AIControlled = true
	player.Name = fmt.Sprintf("AI Player %d", player.Num)
	if len(colors) > player.Num {
		player.Color = colors[player.Num-1]
	}
	player.DefaultHullSet = player.Num % 4
	player.Ready = true

	if err := gr.dbConn.WrapInTransaction(func(c db.Client) error {

		if err := c.CreatePlayer(player); err != nil {
			return fmt.Errorf("added slot player %s for game %d: %w", player, game.ID, err)
		}

		game.NumPlayers++
		if err := c.UpdateGame(&game.Game); err != nil {
			return fmt.Errorf("updating open player slots for game %d: %w", game.ID, err)
		}

		log.Info().Int64("GameID", game.ID).Int("Num", player.Num).Msgf("added player slot %d %s", player.Num, game.Name)
		return nil
	}); err != nil {
		return nil, err
	}
	return player, nil
}

func (gr *gameRunner) StartGame(game *cs.Game) error {
	defer timeTrack(time.Now(), "GenerateUniverse")
	fullGame, err := gr.loadGame(gr.dbConn.NewReadClient(), game.ID)
	if err != nil {
		return fmt.Errorf("load game %d: %w", game.ID, err)
	}

	// generate the universe
	if err := gr.generateUniverse(fullGame); err != nil {
		return fmt.Errorf("generate universe %w", err)
	}

	// save the game to the db
	if err := gr.dbConn.WrapInTransaction(func(c db.Client) error {
		if err := c.UpdateFullGame(fullGame); err != nil {
			return fmt.Errorf("update full game failed: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("save game %d: %w", fullGame.ID, err)
	}

	return nil
}

// load a player and the light version of the player game
func (gr *gameRunner) LoadPlayerGame(gameID int64, userID int64) (*cs.GameWithPlayers, *cs.FullPlayer, error) {

	readClient := gr.dbConn.NewReadClient()
	game, err := readClient.GetGame(gameID)

	if err != nil {
		return nil, nil, err
	}

	if game.Rules.TechsID == 0 {
		game.Rules.SetTechStore(&cs.StaticTechStore)
	} else {
		techs, err := readClient.GetTechStore(game.Rules.TechsID)
		if err != nil {
			return nil, nil, err
		}
		game.Rules.SetTechStore(techs)
	}

	player, err := readClient.GetFullPlayerForGame(gameID, userID)

	if err != nil {
		return nil, nil, err
	}

	return game, player, nil
}

// submit a turn for a player
func (gr *gameRunner) SubmitTurn(gameID int64, userID int64) error {
	client := gr.dbConn.NewReadWriteClient()
	player, err := client.GetLightPlayerForGame(gameID, userID)
	if err != nil {
		return fmt.Errorf("find player for user %d, game %d: %w", userID, gameID, err)
	}

	if player == nil {
		return fmt.Errorf("player for user %d, game %d not found", userID, gameID)
	}

	if err := client.SubmitPlayerTurn(gameID, player.Num, true); err != nil {
		return fmt.Errorf("submitting player turn %w", err)
	}
	return nil
}

// CheckAndGenerateTurn check a game for every player submitted and generate a turn
func (gr *gameRunner) CheckAndGenerateTurn(gameID int64) (TurnGenerationCheckResult, error) {
	defer timeTrack(time.Now(), "CheckAndGenerateTurn")
	fullGame, err := gr.loadGame(gr.dbConn.NewReadClient(), gameID)

	if err != nil {
		return TurnNotGenerated, err
	}

	// if everyone submitted their turn, generate a new turn
	if gr.client.CheckAllPlayersSubmitted(fullGame.Players) {
		readWriteClient := gr.dbConn.NewReadWriteClient()
		if err := readWriteClient.UpdateGameState(fullGame.ID, cs.GameStateGeneratingTurn); err != nil {
			return TurnNotGenerated, err
		}

		if err := gr.generateTurn(readWriteClient, fullGame); err != nil {
			return TurnNotGenerated, err
		} else {
			return TurnGenerated, nil
		}
	}

	return TurnNotGenerated, nil
}

// GenerateTurn generate a new turn, regardless of whether player's have all submitted their turns
func (gr *gameRunner) GenerateTurn(gameID int64) (TurnGenerationCheckResult, error) {
	defer timeTrack(time.Now(), "GenerateTurn")

	readWriteClient := gr.dbConn.NewReadWriteClient()

	// update the state so no one can make calls against it
	if err := readWriteClient.UpdateGameState(gameID, cs.GameStateGeneratingTurn); err != nil {
		return TurnNotGenerated, err
	}

	fullGame, err := gr.loadGame(readWriteClient, gameID)

	if err != nil {
		return TurnNotGenerated, err
	}

	if err := gr.generateTurn(readWriteClient, fullGame); err != nil {
		return TurnNotGenerated, err
	}

	return TurnGenerated, nil
}

// load a full game
func (gr *gameRunner) loadGame(db DBClient, gameID int64) (*cs.FullGame, error) {
	defer timeTrack(time.Now(), "loadGame")

	fullGame, err := db.GetFullGame(gameID)

	if err != nil {
		return nil, fmt.Errorf("load game %d: %w", gameID, err)
	}

	return fullGame, nil
}

// generate a new universe and save the game in the database
func (gr *gameRunner) generateUniverse(fullGame *cs.FullGame) error {

	universe, err := gr.client.GenerateUniverse(fullGame.Game, fullGame.Players)
	if err != nil {
		return fmt.Errorf("generate universe for game %d: %w", fullGame.ID, err)
	}

	// ai processing
	gr.processAITurns(fullGame)

	fullGame.State = cs.GameStateWaitingForPlayers
	fullGame.Universe = universe

	return nil
}

// process an the ai player's turns
func (gr *gameRunner) processAITurns(fullGame *cs.FullGame) {
	for _, player := range fullGame.Players {
		if !player.AIControlled {
			continue
		}
		// TODO: make this use copies to ensure the ai only updates orders?
		// TODO: ai only ai processing
		pmo := fullGame.Universe.GetPlayerMapObjects(player.Num)
		ai := ai.NewAIPlayer(fullGame.Game, fullGame.TechStore, player, pmo)
		if err := ai.ProcessTurn(); err != nil {
			log.Error().Err(err).Int64("GameID", fullGame.ID).Int("PlayerNum", player.Num).Msgf("ai process turn")
		}

		if player.AIControlled {
			player.SubmittedTurn = true
		}
	}
}

// generate a turn for a game
func (gr *gameRunner) generateTurn(readWriteClient DBClient, fullGame *cs.FullGame) error {

	defer func() {
		// if we panic, update the game state to fail
		if r := recover(); r != nil {
			// update the state so it's clear we got an error during turn generation
			if err := readWriteClient.UpdateGameState(fullGame.ID, cs.GameStateGeneratingTurnError); err != nil {
				log.Error().Err(err).Msgf("failed to update game state")
			}

			// middleware.PrintPrettyStack(r)
			log.Panic().Msgf("failed to generate turn %v", r)
		}
	}()

	// if everyone submitted their turn, generate a new turn
	if err := gr.client.GenerateTurn(fullGame.Game, fullGame.Universe, fullGame.Players); err != nil {
		// update the state so it's clear we got an error during turn generation
		if err := readWriteClient.UpdateGameState(fullGame.ID, cs.GameStateGeneratingTurnError); err != nil {
			log.Error().Err(err).Msgf("failed to update game state")
		}

		return fmt.Errorf("generate turn -> %w", err)
	}

	// ai processing
	gr.processAITurns(fullGame)

	// save the game and update the game state
	if err := gr.dbConn.WrapInTransaction(func(c db.Client) error {
		// for UI debugging
		// if fullGame.ID == 29 {
		// 	time.Sleep(10 * time.Second)
		// }

		for _, player := range fullGame.Players {
			if !player.AIControlled {
				continue
			}

			for _, design := range player.Designs {
				if design.Delete && design.ID != 0 && !design.CannotDelete {
					log.Debug().Msgf("player %d deleting design %d - id %d", player.Num, design.Num, design.ID)
					if err := c.DeleteShipDesign(design.ID); err != nil {
						return fmt.Errorf("delete AI player design %s %v", design.Name, err)
					}
				}
			}
		}

		if err := c.UpdateFullGame(fullGame); err != nil {

			if err := readWriteClient.UpdateGameState(fullGame.ID, cs.GameStateGeneratingTurnError); err != nil {
				log.Error().Err(err).Msgf("failed to update game state")
			}

			return fmt.Errorf("save game after turn generation -> %w", err)
		}
		return nil
	}); err != nil {
		// update the state of the game to error
		if err := readWriteClient.UpdateGameState(fullGame.ID, cs.GameStateGeneratingTurnError); err != nil {
			log.Error().Err(err).Msgf("failed to update game state after failing to save game during turn generation")
		}

		return err
	}

	return nil

}

// load the tech store into a game
func (gr *gameRunner) loadTechStore(db DBClient, game *cs.Game) (*cs.TechStore, error) {
	if game.Rules.TechsID == 0 {
		return &cs.StaticTechStore, nil
	}
	techs, err := db.GetTechStore(game.Rules.TechsID)
	if err != nil {
		return nil, err
	}
	return techs, nil
}
