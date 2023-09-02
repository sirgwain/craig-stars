package server

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
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

type GameRunner interface {
	HostGame(hostID int64, settings *cs.GameSettings) (*cs.FullGame, error)
	JoinGame(gameID int64, userID int64, raceID int64) error
	LeaveGame(gameID, userID int64) error
	KickPlayer(gameID int64, playerNum int) error
	AddOpenPlayerSlot(game *cs.GameWithPlayers) (*cs.Player, error)
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

		for i, playerSetting := range settings.Players {
			if playerSetting.Type == cs.NewGamePlayerTypeHost {

				// add the host player with the host race
				race, err := c.GetRace(playerSetting.RaceID)
				if err != nil {
					return err
				}
				if race.UserID != hostID {
					return fmt.Errorf("user %d does not own Race %d", hostID, race.ID)
				}
				log.Debug().Int64("hostID", hostID).Msg("Adding host to game")
				player := gr.client.NewPlayer(hostID, *race, &game.Rules)
				player.GameID = game.ID
				player.Num = i + 1
				player.Name = user.Username
				player.Color = playerSetting.Color
				player.DefaultHullSet = playerSetting.DefaultHullSet
				player.Ready = true
				players = append(players, player)
			} else if playerSetting.Type == cs.NewGamePlayerTypeAI {
				log.Debug().Int64("hostID", hostID).Msg("Adding ai player to game")
				race := *cs.NewRace()
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
				log.Debug().Uint("openPlayerSlots", game.OpenPlayerSlots).Msg("Added open player slot to game")
				race := cs.NewRace()
				player := gr.client.NewPlayer(0, *race, &game.Rules)
				player.GameID = game.ID
				player.Num = i + 1
				player.Name = "Open Slot"
				player.Color = playerSetting.Color
				player.DefaultHullSet = playerSetting.DefaultHullSet
				players = append(players, player)
				game.OpenPlayerSlots++
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
		if game.OpenPlayerSlots == 0 {
			err = gr.generateUniverse(fullGame)
			if err != nil {
				return err
			}
		}

		// save the game
		if err := c.UpdateFullGame(fullGame); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return fullGame, nil
}

// add a player to an existing game
func (gr *gameRunner) JoinGame(gameID int64, userID int64, raceID int64) error {
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

	race, err := readClient.GetRace(raceID)
	if err != nil {
		log.Error().Err(err).Int64("ID", raceID).Msg("get race from database")
		return fmt.Errorf("failed to get race from database")
	}

	// validate
	if race == nil {
		return fmt.Errorf("no race for id %d found. %w", raceID, errNotFound)
	}

	if race.UserID != userID {
		return fmt.Errorf("user %d does not control race %d", userID, raceID)
	}

	if fullGame == nil {
		return fmt.Errorf("no game for id %d found. %w", gameID, errNotFound)
	}

	if fullGame.OpenPlayerSlots == 0 {
		return fmt.Errorf("no slots left in this game")
	}

	player := gr.client.NewPlayer(userID, *race, &fullGame.Rules)
	player.GameID = fullGame.ID
	player.Name = user.Username
	player.Ready = true
	player.AIControlled = false

	if err := gr.dbConn.WrapInTransaction(func(c db.Client) error {

		// claim an open slot
		for i, p := range fullGame.Players {
			if !p.AIControlled && p.UserID == 0 {
				// take over this empty player
				player.Num = p.Num
				player.ID = p.ID
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

				fullGame.OpenPlayerSlots--
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

// add a player to an existing game
func (gr *gameRunner) LeaveGame(gameID, userID int64) error {

	game, err := gr.loadGame(gr.dbConn.NewReadClient(), gameID)
	if err != nil {
		return fmt.Errorf("unable to load game %d: %w", gameID, err)
	}

	if game == nil {
		return fmt.Errorf("no game for id %d found. %w", gameID, errNotFound)
	}

	if err := gr.dbConn.WrapInTransaction(func(c db.Client) error {
		for i, player := range game.Players {
			if player.UserID == userID {
				if err := c.DeletePlayer(player.ID); err != nil {
					return fmt.Errorf("delete open slot player %s from game %d: %w", player, gameID, err)
				}

				race := cs.NewRace()
				player := gr.client.NewPlayer(0, *race, &game.Rules)
				player.GameID = game.ID
				player.Num = i + 1
				player.Name = "Open Slot"
				game.Players[i] = player
				game.OpenPlayerSlots++

				if player.Num-1 < len(colors) {
					player.Color = colors[player.Num-1]
				} else {
					color := make([]byte, 3)
					rand.Read(color)
					player.Color = fmt.Sprintf("#%s", hex.EncodeToString(color))
				}

				if err := c.CreatePlayer(player); err != nil {
					return fmt.Errorf("update open slot player %s for game %d: %w", player, gameID, err)
				}

			}
		}

		if err := c.UpdateGame(game.Game); err != nil {
			return fmt.Errorf("save game %d: %w", gameID, err)
		}

		return nil
	}); err != nil {
		return err
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

	if err := gr.dbConn.WrapInTransaction(func(c db.Client) error {

		for i, player := range game.Players {
			if player.Num == playerNum {
				if err := c.DeletePlayer(player.ID); err != nil {
					return fmt.Errorf("delete open slot player %s from game %d: %w", player, gameID, err)
				}

				race := cs.NewRace()
				player := gr.client.NewPlayer(0, *race, &game.Rules)
				player.GameID = game.ID
				player.Num = i + 1
				player.Name = "Open Slot"
				game.Players[i] = player
				game.OpenPlayerSlots++

				if player.Num-1 < len(colors) {
					player.Color = colors[player.Num-1]
				} else {
					color := make([]byte, 3)
					rand.Read(color)
					player.Color = fmt.Sprintf("#%s", hex.EncodeToString(color))
				}

				if err := c.CreatePlayer(player); err != nil {
					return fmt.Errorf("update open slot player %s for game %d: %w", player, gameID, err)
				}

			}
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
					game.OpenPlayerSlots--
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
				// if the player was using a default color, use the previous one
				if len(colors) > i+1 && player.Color == colors[i+1] {
					player.Color = colors[i]
				}
				if player.AIControlled {
					player.Name = fmt.Sprintf("AI Player %d", aiPlayerNum+1)
				}

				if err := c.UpdatePlayer(player); err != nil {
					return fmt.Errorf("update player %s for game %d: %w", player.Name, gameID, err)
				}
			}
		}

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

	dbClient := gr.dbConn.NewReadWriteClient()
	if err := dbClient.CreatePlayer(player); err != nil {
		return nil, fmt.Errorf("added slot player %s for game %d: %w", player, game.ID, err)
	}

	log.Info().Int64("GameID", game.ID).Int("Num", player.Num).Msgf("added player slot %d %s", player.Num, game.Name)

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
		game.Rules.WithTechStore(&cs.StaticTechStore)
	} else {
		techs, err := readClient.GetTechStore(game.Rules.TechsID)
		if err != nil {
			return nil, nil, err
		}
		game.Rules.WithTechStore(techs)
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
