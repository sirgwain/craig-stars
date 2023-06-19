// Code generated by github.com/jmattheis/goverter, DO NOT EDIT.

package dbsqlx

import game "github.com/sirgwain/craig-stars/game"

type GameConverter struct{}

func (c *GameConverter) ConvertFleet(source *Fleet) *game.Fleet {
	var pGameFleet *game.Fleet
	if source != nil {
		gameFleet := c.dbsqlxFleetToGameFleet(*source)
		pGameFleet = &gameFleet
	}
	return pGameFleet
}
func (c *GameConverter) ConvertGame(source Game) game.Game {
	var gameGame game.Game
	gameGame.ID = source.ID
	gameGame.CreatedAt = TimeToTime(source.CreatedAt)
	gameGame.UpdatedAt = TimeToTime(source.UpdatedAt)
	gameGame.Name = source.Name
	gameGame.HostID = source.HostID
	gameGame.QuickStartTurns = source.QuickStartTurns
	gameGame.Size = game.Size(source.Size)
	gameGame.Density = game.Density(source.Density)
	gameGame.PlayerPositions = game.PlayerPositions(source.PlayerPositions)
	gameGame.RandomEvents = source.RandomEvents
	gameGame.ComputerPlayersFormAlliances = source.ComputerPlayersFormAlliances
	gameGame.PublicPlayerScores = source.PublicPlayerScores
	gameGame.StartMode = game.GameStartMode(source.StartMode)
	gameGame.Year = source.Year
	gameGame.State = game.GameState(source.State)
	gameGame.OpenPlayerSlots = source.OpenPlayerSlots
	gameGame.NumPlayers = source.NumPlayers
	gameGame.VictoryConditions = ExtendVictoryConditions(source)
	gameGame.VictorDeclared = source.VictorDeclared
	gameGame.Seed = source.Seed
	gameGame.Rules = ExtendDefaultRules(source)
	gameGame.Area = ExtendArea(source)
	return gameGame
}
func (c *GameConverter) ConvertGameFleet(source *game.Fleet) *Fleet {
	var pDbsqlxFleet *Fleet
	if source != nil {
		dbsqlxFleet := c.gameFleetToDbsqlxFleet(*source)
		pDbsqlxFleet = &dbsqlxFleet
	}
	return pDbsqlxFleet
}
func (c *GameConverter) ConvertGameGame(source *game.Game) *Game {
	var pDbsqlxGame *Game
	if source != nil {
		dbsqlxGame := c.gameGameToDbsqlxGame(*source)
		pDbsqlxGame = &dbsqlxGame
	}
	return pDbsqlxGame
}
func (c *GameConverter) ConvertGamePlanet(source *game.Planet) *Planet {
	var pDbsqlxPlanet *Planet
	if source != nil {
		dbsqlxPlanet := c.gamePlanetToDbsqlxPlanet(*source)
		pDbsqlxPlanet = &dbsqlxPlanet
	}
	return pDbsqlxPlanet
}
func (c *GameConverter) ConvertGamePlanetIntel(source *game.PlanetIntel) *PlanetIntel {
	var pDbsqlxPlanetIntel *PlanetIntel
	if source != nil {
		dbsqlxPlanetIntel := c.gamePlanetIntelToDbsqlxPlanetIntel(*source)
		pDbsqlxPlanetIntel = &dbsqlxPlanetIntel
	}
	return pDbsqlxPlanetIntel
}
func (c *GameConverter) ConvertGamePlayer(source *game.Player) *Player {
	var pDbsqlxPlayer *Player
	if source != nil {
		dbsqlxPlayer := c.gamePlayerToDbsqlxPlayer(*source)
		pDbsqlxPlayer = &dbsqlxPlayer
	}
	return pDbsqlxPlayer
}
func (c *GameConverter) ConvertGameRace(source *game.Race) *Race {
	var pDbsqlxRace *Race
	if source != nil {
		dbsqlxRace := c.gameRaceToDbsqlxRace(*source)
		pDbsqlxRace = &dbsqlxRace
	}
	return pDbsqlxRace
}
func (c *GameConverter) ConvertGameShipDesign(source *game.ShipDesign) *ShipDesign {
	var pDbsqlxShipDesign *ShipDesign
	if source != nil {
		dbsqlxShipDesign := c.gameShipDesignToDbsqlxShipDesign(*source)
		pDbsqlxShipDesign = &dbsqlxShipDesign
	}
	return pDbsqlxShipDesign
}
func (c *GameConverter) ConvertGameShipToken(source game.ShipToken) ShipToken {
	var dbsqlxShipToken ShipToken
	dbsqlxShipToken.ID = source.ID
	dbsqlxShipToken.CreatedAt = TimeToNullTime(source.CreatedAt)
	dbsqlxShipToken.UpdatedAt = TimeToNullTime(source.UpdatedAt)
	dbsqlxShipToken.FleetID = source.FleetID
	dbsqlxShipToken.DesignUUID = UUIDToUUID(source.DesignUUID)
	dbsqlxShipToken.Quantity = source.Quantity
	dbsqlxShipToken.Damage = source.Damage
	dbsqlxShipToken.QuantityDamaged = source.QuantityDamaged
	return dbsqlxShipToken
}
func (c *GameConverter) ConvertGameUser(source *game.User) *User {
	var pDbsqlxUser *User
	if source != nil {
		dbsqlxUser := c.gameUserToDbsqlxUser(*source)
		pDbsqlxUser = &dbsqlxUser
	}
	return pDbsqlxUser
}
func (c *GameConverter) ConvertGames(source []Game) []game.Game {
	gameGameList := make([]game.Game, len(source))
	for i := 0; i < len(source); i++ {
		gameGameList[i] = c.ConvertGame(source[i])
	}
	return gameGameList
}
func (c *GameConverter) ConvertPlanet(source *Planet) *game.Planet {
	var pGamePlanet *game.Planet
	if source != nil {
		gamePlanet := c.dbsqlxPlanetToGamePlanet(*source)
		pGamePlanet = &gamePlanet
	}
	return pGamePlanet
}
func (c *GameConverter) ConvertPlanetIntel(source *PlanetIntel) *game.PlanetIntel {
	var pGamePlanetIntel *game.PlanetIntel
	if source != nil {
		gamePlanetIntel := c.dbsqlxPlanetIntelToGamePlanetIntel(*source)
		pGamePlanetIntel = &gamePlanetIntel
	}
	return pGamePlanetIntel
}
func (c *GameConverter) ConvertPlayer(source Player) game.Player {
	var gamePlayer game.Player
	gamePlayer.ID = source.ID
	gamePlayer.CreatedAt = TimeToTime(source.CreatedAt)
	gamePlayer.UpdatedAt = TimeToTime(source.UpdatedAt)
	gamePlayer.GameID = source.GameID
	gamePlayer.UserID = source.UserID
	gamePlayer.Name = source.Name
	gamePlayer.Num = source.Num
	gamePlayer.Ready = source.Ready
	gamePlayer.AIControlled = source.AIControlled
	gamePlayer.SubmittedTurn = source.SubmittedTurn
	gamePlayer.Color = source.Color
	gamePlayer.DefaultHullSet = source.DefaultHullSet
	gamePlayer.Race = PlayerRaceToGameRace(source.Race)
	gamePlayer.TechLevels = ExtendTechLevels(source)
	gamePlayer.TechLevelsSpent = ExtendTechLevelsSpent(source)
	gamePlayer.ResearchAmount = source.ResearchAmount
	gamePlayer.ResearchSpentLastYear = source.ResearchSpentLastYear
	gamePlayer.NextResearchField = game.NextResearchField(source.NextResearchField)
	gamePlayer.Researching = game.TechField(source.Researching)
	gamePlayer.BattlePlans = BattlePlansToGameBattlePlans(source.BattlePlans)
	gamePlayer.ProductionPlans = ProductionPlansToGameProductionPlans(source.ProductionPlans)
	gamePlayer.TransportPlans = TransportPlansToGameTransportPlans(source.TransportPlans)
	gamePlayer.Stats = PlayerStatsToGamePlayerStats(source.Stats)
	gamePlayer.Spec = PlayerSpecToGamePlayerSpec(source.Spec)
	return gamePlayer
}
func (c *GameConverter) ConvertPlayers(source []Player) []game.Player {
	gamePlayerList := make([]game.Player, len(source))
	for i := 0; i < len(source); i++ {
		gamePlayerList[i] = c.ConvertPlayer(source[i])
	}
	return gamePlayerList
}
func (c *GameConverter) ConvertRace(source Race) game.Race {
	var gameRace game.Race
	gameRace.ID = source.ID
	gameRace.CreatedAt = TimeToTime(source.CreatedAt)
	gameRace.UpdatedAt = TimeToTime(source.UpdatedAt)
	gameRace.UserID = source.UserID
	var pInt64 *int64
	if source.PlayerID != nil {
		xint64 := *source.PlayerID
		pInt64 = &xint64
	}
	gameRace.PlayerID = pInt64
	gameRace.Name = source.Name
	gameRace.PluralName = source.PluralName
	gameRace.PRT = game.PRT(source.PRT)
	gameRace.LRTs = game.Bitmask(source.LRTs)
	gameRace.HabLow = ExtendHabLow(source)
	gameRace.HabHigh = ExtendHabHigh(source)
	gameRace.GrowthRate = source.GrowthRate
	gameRace.PopEfficiency = source.PopEfficiency
	gameRace.FactoryOutput = source.FactoryOutput
	gameRace.FactoryCost = source.FactoryCost
	gameRace.NumFactories = source.NumFactories
	gameRace.FactoriesCostLess = source.FactoriesCostLess
	gameRace.ImmuneGrav = source.ImmuneGrav
	gameRace.ImmuneTemp = source.ImmuneTemp
	gameRace.ImmuneRad = source.ImmuneRad
	gameRace.MineOutput = source.MineOutput
	gameRace.MineCost = source.MineCost
	gameRace.NumMines = source.NumMines
	gameRace.ResearchCost = ExtendResearchCost(source)
	gameRace.TechsStartHigh = source.TechsStartHigh
	gameRace.Spec = RaceSpecToGameRaceSpec(source.Spec)
	return gameRace
}
func (c *GameConverter) ConvertRaces(source []Race) []game.Race {
	gameRaceList := make([]game.Race, len(source))
	for i := 0; i < len(source); i++ {
		gameRaceList[i] = c.ConvertRace(source[i])
	}
	return gameRaceList
}
func (c *GameConverter) ConvertShipDesign(source *ShipDesign) *game.ShipDesign {
	var pGameShipDesign *game.ShipDesign
	if source != nil {
		gameShipDesign := c.dbsqlxShipDesignToGameShipDesign(*source)
		pGameShipDesign = &gameShipDesign
	}
	return pGameShipDesign
}
func (c *GameConverter) ConvertShipToken(source ShipToken) game.ShipToken {
	var gameShipToken game.ShipToken
	gameShipToken.ID = source.ID
	gameShipToken.CreatedAt = NullTimeToTime(source.CreatedAt)
	gameShipToken.UpdatedAt = NullTimeToTime(source.UpdatedAt)
	gameShipToken.FleetID = source.FleetID
	gameShipToken.DesignUUID = UUIDToUUID(source.DesignUUID)
	gameShipToken.Quantity = source.Quantity
	gameShipToken.Damage = source.Damage
	gameShipToken.QuantityDamaged = source.QuantityDamaged
	return gameShipToken
}
func (c *GameConverter) ConvertUser(source User) game.User {
	var gameUser game.User
	gameUser.ID = source.ID
	gameUser.CreatedAt = TimeToTime(source.CreatedAt)
	gameUser.UpdatedAt = TimeToTime(source.UpdatedAt)
	gameUser.Username = source.Username
	gameUser.Password = source.Password
	gameUser.Role = game.Role(source.Role)
	return gameUser
}
func (c *GameConverter) ConvertUsers(source []User) []game.User {
	gameUserList := make([]game.User, len(source))
	for i := 0; i < len(source); i++ {
		gameUserList[i] = c.ConvertUser(source[i])
	}
	return gameUserList
}
func (c *GameConverter) dbsqlxFleetToGameFleet(source Fleet) game.Fleet {
	var gameFleet game.Fleet
	gameFleet.MapObject = ExtendFleetMapObject(source)
	gameFleet.FleetOrders = ExtendFleetFleetOrders(source)
	gameFleet.PlanetID = source.PlanetID
	gameFleet.BaseName = source.BaseName
	gameFleet.Cargo = ExtendFleetCargo(source)
	gameFleet.Fuel = source.Fuel
	gameFleet.Damage = source.Damage
	gameFleet.BattlePlanName = source.BattlePlanName
	gameFleet.Heading = ExtendFleetHeading(source)
	gameFleet.WarpSpeed = source.WarpSpeed
	gameFleet.PreviousPosition = ExtendFleetPreviousPosition(source)
	gameFleet.OrbitingPlanetNum = source.OrbitingPlanetNum
	gameFleet.Spec = FleetSpecToGameFleetSpec(source.Spec)
	return gameFleet
}
func (c *GameConverter) dbsqlxPlanetIntelToGamePlanetIntel(source PlanetIntel) game.PlanetIntel {
	var gamePlanetIntel game.PlanetIntel
	gamePlanetIntel.MapObjectIntel = ExtendPlanetIntelMapObjectIntel(source)
	gamePlanetIntel.Hab = ExtendPlanetIntelHab(source)
	gamePlanetIntel.MineralConcentration = ExtendPlanetIntelMineralConcentration(source)
	gamePlanetIntel.Population = source.Population
	gamePlanetIntel.Cargo = ExtendPlanetIntelCargo(source)
	gamePlanetIntel.CargoDiscovered = source.CargoDiscovered
	return gamePlanetIntel
}
func (c *GameConverter) dbsqlxPlanetToGamePlanet(source Planet) game.Planet {
	var gamePlanet game.Planet
	gamePlanet.MapObject = ExtendPlanetMapObject(source)
	gamePlanet.Hab = ExtendHab(source)
	gamePlanet.BaseHab = ExtendBaseHab(source)
	gamePlanet.TerraformedAmount = ExtendTerraformedAmount(source)
	gamePlanet.MineralConcentration = ExtendMineralConcentration(source)
	gamePlanet.MineYears = ExtendMineYears(source)
	gamePlanet.Cargo = ExtendPlanetCargo(source)
	gamePlanet.Mines = source.Mines
	gamePlanet.Factories = source.Factories
	gamePlanet.Defenses = source.Defenses
	gamePlanet.Homeworld = source.Homeworld
	gamePlanet.ContributesOnlyLeftoverToResearch = source.ContributesOnlyLeftoverToResearch
	gamePlanet.Scanner = source.Scanner
	gamePlanet.PacketSpeed = source.PacketSpeed
	gamePlanet.BonusResources = source.BonusResources
	gamePlanet.ProductionQueue = ProductionQueueItemsToGameProductionQueueItems(source.ProductionQueue)
	gamePlanet.Spec = PlanetSpecToGamePlanetSpec(source.Spec)
	return gamePlanet
}
func (c *GameConverter) dbsqlxShipDesignToGameShipDesign(source ShipDesign) game.ShipDesign {
	var gameShipDesign game.ShipDesign
	gameShipDesign.ID = source.ID
	gameShipDesign.CreatedAt = TimeToTime(source.CreatedAt)
	gameShipDesign.UpdatedAt = TimeToTime(source.UpdatedAt)
	gameShipDesign.PlayerID = source.PlayerID
	gameShipDesign.PlayerNum = source.PlayerNum
	gameShipDesign.UUID = UUIDToUUID(source.UUID)
	gameShipDesign.Name = source.Name
	gameShipDesign.Version = source.Version
	gameShipDesign.Hull = source.Hull
	gameShipDesign.HullSetNumber = source.HullSetNumber
	gameShipDesign.CanDelete = source.CanDelete
	gameShipDesign.Slots = ShipDesignSlotsToGameShipDesignSlots(source.Slots)
	gameShipDesign.Purpose = game.ShipDesignPurpose(source.Purpose)
	gameShipDesign.Spec = ShipDesignSpecToGameShipDesignSpec(source.Spec)
	return gameShipDesign
}
func (c *GameConverter) gameFleetToDbsqlxFleet(source game.Fleet) Fleet {
	var dbsqlxFleet Fleet
	dbsqlxFleet.ID = source.MapObject.ID
	dbsqlxFleet.GameID = source.MapObject.GameID
	dbsqlxFleet.CreatedAt = TimeToTime(source.MapObject.CreatedAt)
	dbsqlxFleet.UpdatedAt = TimeToTime(source.MapObject.UpdatedAt)
	dbsqlxFleet.PlayerID = source.MapObject.PlayerID
	dbsqlxFleet.X = source.MapObject.Position.X
	dbsqlxFleet.Y = source.MapObject.Position.Y
	dbsqlxFleet.Name = source.MapObject.Name
	dbsqlxFleet.Num = source.MapObject.Num
	dbsqlxFleet.PlayerNum = source.MapObject.PlayerNum
	dbsqlxFleet.Waypoints = GameWaypointsToWaypoints(source.FleetOrders.Waypoints)
	dbsqlxFleet.RepeatOrders = source.FleetOrders.RepeatOrders
	dbsqlxFleet.PlanetID = source.PlanetID
	dbsqlxFleet.BaseName = source.BaseName
	dbsqlxFleet.Ironium = source.Cargo.Ironium
	dbsqlxFleet.Boranium = source.Cargo.Boranium
	dbsqlxFleet.Germanium = source.Cargo.Germanium
	dbsqlxFleet.Colonists = source.Cargo.Colonists
	dbsqlxFleet.Fuel = source.Fuel
	dbsqlxFleet.Damage = source.Damage
	dbsqlxFleet.BattlePlanName = source.BattlePlanName
	dbsqlxFleet.HeadingX = source.Heading.X
	dbsqlxFleet.HeadingY = source.Heading.Y
	dbsqlxFleet.WarpSpeed = source.WarpSpeed
	var pFloat64 *float64
	if source.PreviousPosition != nil {
		pFloat64 = &source.PreviousPosition.X
	}
	var pFloat642 *float64
	if pFloat64 != nil {
		xfloat64 := *pFloat64
		pFloat642 = &xfloat64
	}
	dbsqlxFleet.PreviousPositionX = pFloat642
	var pFloat643 *float64
	if source.PreviousPosition != nil {
		pFloat643 = &source.PreviousPosition.Y
	}
	var pFloat644 *float64
	if pFloat643 != nil {
		xfloat642 := *pFloat643
		pFloat644 = &xfloat642
	}
	dbsqlxFleet.PreviousPositionY = pFloat644
	dbsqlxFleet.OrbitingPlanetNum = source.OrbitingPlanetNum
	dbsqlxFleet.Starbase = source.Starbase
	dbsqlxFleet.Spec = GameFleetSpecToFleetSpec(source.Spec)
	return dbsqlxFleet
}
func (c *GameConverter) gameGameToDbsqlxGame(source game.Game) Game {
	var dbsqlxGame Game
	dbsqlxGame.ID = source.ID
	dbsqlxGame.CreatedAt = TimeToTime(source.CreatedAt)
	dbsqlxGame.UpdatedAt = TimeToTime(source.UpdatedAt)
	dbsqlxGame.Name = source.Name
	dbsqlxGame.HostID = source.HostID
	dbsqlxGame.QuickStartTurns = source.QuickStartTurns
	dbsqlxGame.Size = game.Size(source.Size)
	dbsqlxGame.Density = game.Density(source.Density)
	dbsqlxGame.PlayerPositions = game.PlayerPositions(source.PlayerPositions)
	dbsqlxGame.RandomEvents = source.RandomEvents
	dbsqlxGame.ComputerPlayersFormAlliances = source.ComputerPlayersFormAlliances
	dbsqlxGame.PublicPlayerScores = source.PublicPlayerScores
	dbsqlxGame.StartMode = game.GameStartMode(source.StartMode)
	dbsqlxGame.Year = source.Year
	dbsqlxGame.State = game.GameState(source.State)
	dbsqlxGame.OpenPlayerSlots = source.OpenPlayerSlots
	dbsqlxGame.NumPlayers = source.NumPlayers
	dbsqlxVictoryConditions := c.gameVictoryConditionListToDbsqlxVictoryConditions(source.VictoryConditions.Conditions)
	dbsqlxGame.VictoryConditionsConditions = &dbsqlxVictoryConditions
	dbsqlxGame.VictoryConditionsNumCriteriaRequired = source.VictoryConditions.NumCriteriaRequired
	dbsqlxGame.VictoryConditionsYearsPassed = source.VictoryConditions.YearsPassed
	dbsqlxGame.VictoryConditionsOwnPlanets = source.VictoryConditions.OwnPlanets
	dbsqlxGame.VictoryConditionsAttainTechLevel = source.VictoryConditions.AttainTechLevel
	dbsqlxGame.VictoryConditionsAttainTechLevelNumFields = source.VictoryConditions.AttainTechLevelNumFields
	dbsqlxGame.VictoryConditionsExceedsScore = source.VictoryConditions.ExceedsScore
	dbsqlxGame.VictoryConditionsExceedsSecondPlaceScore = source.VictoryConditions.ExceedsSecondPlaceScore
	dbsqlxGame.VictoryConditionsProductionCapacity = source.VictoryConditions.ProductionCapacity
	dbsqlxGame.VictoryConditionsOwnCapitalShips = source.VictoryConditions.OwnCapitalShips
	dbsqlxGame.VictoryConditionsHighestScoreAfterYears = source.VictoryConditions.HighestScoreAfterYears
	dbsqlxGame.VictorDeclared = source.VictorDeclared
	dbsqlxGame.Seed = source.Seed
	dbsqlxGame.Rules = GameRulesToRules(source.Rules)
	dbsqlxGame.AreaX = source.Area.X
	dbsqlxGame.AreaY = source.Area.Y
	return dbsqlxGame
}
func (c *GameConverter) gamePlanetIntelToDbsqlxPlanetIntel(source game.PlanetIntel) PlanetIntel {
	var dbsqlxPlanetIntel PlanetIntel
	dbsqlxPlanetIntel.ID = source.MapObjectIntel.Intel.ID
	dbsqlxPlanetIntel.CreatedAt = TimeToTime(source.MapObjectIntel.Intel.CreatedAt)
	dbsqlxPlanetIntel.UpdatedAt = TimeToTime(source.MapObjectIntel.Intel.UpdatedAt)
	dbsqlxPlanetIntel.Dirty = source.MapObjectIntel.Intel.Dirty
	dbsqlxPlanetIntel.PlayerID = source.MapObjectIntel.Intel.PlayerID
	dbsqlxPlanetIntel.Name = source.MapObjectIntel.Intel.Name
	dbsqlxPlanetIntel.Num = source.MapObjectIntel.Intel.Num
	dbsqlxPlanetIntel.PlayerNum = source.MapObjectIntel.Intel.PlayerNum
	dbsqlxPlanetIntel.ReportAge = source.MapObjectIntel.Intel.ReportAge
	dbsqlxPlanetIntel.Type = game.MapObjectType(source.MapObjectIntel.Type)
	dbsqlxPlanetIntel.X = source.MapObjectIntel.Position.X
	dbsqlxPlanetIntel.Y = source.MapObjectIntel.Position.Y
	dbsqlxPlanetIntel.Grav = source.Hab.Grav
	dbsqlxPlanetIntel.Temp = source.Hab.Temp
	dbsqlxPlanetIntel.Rad = source.Hab.Rad
	dbsqlxPlanetIntel.MineralConcIronium = source.MineralConcentration.Ironium
	dbsqlxPlanetIntel.MineralConcBoranium = source.MineralConcentration.Boranium
	dbsqlxPlanetIntel.MineralConcGermanium = source.MineralConcentration.Germanium
	dbsqlxPlanetIntel.Population = source.Population
	dbsqlxPlanetIntel.Ironium = source.Cargo.Ironium
	dbsqlxPlanetIntel.Boranium = source.Cargo.Boranium
	dbsqlxPlanetIntel.Germanium = source.Cargo.Germanium
	dbsqlxPlanetIntel.Colonists = source.Cargo.Colonists
	dbsqlxPlanetIntel.CargoDiscovered = source.CargoDiscovered
	return dbsqlxPlanetIntel
}
func (c *GameConverter) gamePlanetToDbsqlxPlanet(source game.Planet) Planet {
	var dbsqlxPlanet Planet
	dbsqlxPlanet.ID = source.MapObject.ID
	dbsqlxPlanet.GameID = source.MapObject.GameID
	dbsqlxPlanet.CreatedAt = TimeToTime(source.MapObject.CreatedAt)
	dbsqlxPlanet.UpdatedAt = TimeToTime(source.MapObject.UpdatedAt)
	dbsqlxPlanet.PlayerID = source.MapObject.PlayerID
	dbsqlxPlanet.X = source.MapObject.Position.X
	dbsqlxPlanet.Y = source.MapObject.Position.Y
	dbsqlxPlanet.Name = source.MapObject.Name
	dbsqlxPlanet.Num = source.MapObject.Num
	dbsqlxPlanet.PlayerNum = source.MapObject.PlayerNum
	dbsqlxPlanet.Grav = source.Hab.Grav
	dbsqlxPlanet.Temp = source.Hab.Temp
	dbsqlxPlanet.Rad = source.Hab.Rad
	dbsqlxPlanet.BaseGrav = source.BaseHab.Grav
	dbsqlxPlanet.BaseTemp = source.BaseHab.Temp
	dbsqlxPlanet.BaseRad = source.BaseHab.Rad
	dbsqlxPlanet.TerraformedAmountGrav = source.TerraformedAmount.Grav
	dbsqlxPlanet.TerraformedAmountTemp = source.TerraformedAmount.Temp
	dbsqlxPlanet.TerraformedAmountRad = source.TerraformedAmount.Rad
	dbsqlxPlanet.MineralConcIronium = source.MineralConcentration.Ironium
	dbsqlxPlanet.MineralConcBoranium = source.MineralConcentration.Boranium
	dbsqlxPlanet.MineralConcGermanium = source.MineralConcentration.Germanium
	dbsqlxPlanet.MineYearsIronium = source.MineYears.Ironium
	dbsqlxPlanet.MineYearsBoranium = source.MineYears.Boranium
	dbsqlxPlanet.MineYearsGermanium = source.MineYears.Germanium
	dbsqlxPlanet.Ironium = source.Cargo.Ironium
	dbsqlxPlanet.Boranium = source.Cargo.Boranium
	dbsqlxPlanet.Germanium = source.Cargo.Germanium
	dbsqlxPlanet.Colonists = source.Cargo.Colonists
	dbsqlxPlanet.Mines = source.Mines
	dbsqlxPlanet.Factories = source.Factories
	dbsqlxPlanet.Defenses = source.Defenses
	dbsqlxPlanet.Homeworld = source.Homeworld
	dbsqlxPlanet.ContributesOnlyLeftoverToResearch = source.ContributesOnlyLeftoverToResearch
	dbsqlxPlanet.Scanner = source.Scanner
	dbsqlxPlanet.PacketSpeed = source.PacketSpeed
	dbsqlxPlanet.BonusResources = source.BonusResources
	dbsqlxPlanet.ProductionQueue = GameProductionQueueItemsToProductionQueueItems(source.ProductionQueue)
	dbsqlxPlanet.Spec = GamePlanetSpecToPlanetSpec(source.Spec)
	return dbsqlxPlanet
}
func (c *GameConverter) gamePlayerToDbsqlxPlayer(source game.Player) Player {
	var dbsqlxPlayer Player
	dbsqlxPlayer.ID = source.ID
	dbsqlxPlayer.CreatedAt = TimeToTime(source.CreatedAt)
	dbsqlxPlayer.UpdatedAt = TimeToTime(source.UpdatedAt)
	dbsqlxPlayer.GameID = source.GameID
	dbsqlxPlayer.UserID = source.UserID
	dbsqlxPlayer.Name = source.Name
	dbsqlxPlayer.Num = source.Num
	dbsqlxPlayer.Ready = source.Ready
	dbsqlxPlayer.AIControlled = source.AIControlled
	dbsqlxPlayer.SubmittedTurn = source.SubmittedTurn
	dbsqlxPlayer.Color = source.Color
	dbsqlxPlayer.DefaultHullSet = source.DefaultHullSet
	dbsqlxPlayer.TechLevelsEnergy = source.TechLevels.Energy
	dbsqlxPlayer.TechLevelsWeapons = source.TechLevels.Weapons
	dbsqlxPlayer.TechLevelsPropulsion = source.TechLevels.Propulsion
	dbsqlxPlayer.TechLevelsConstruction = source.TechLevels.Construction
	dbsqlxPlayer.TechLevelsElectronics = source.TechLevels.Electronics
	dbsqlxPlayer.TechLevelsBiotechnology = source.TechLevels.Biotechnology
	dbsqlxPlayer.TechLevelsSpentEnergy = source.TechLevelsSpent.Energy
	dbsqlxPlayer.TechLevelsSpentWeapons = source.TechLevelsSpent.Weapons
	dbsqlxPlayer.TechLevelsSpentPropulsion = source.TechLevelsSpent.Propulsion
	dbsqlxPlayer.TechLevelsSpentConstruction = source.TechLevelsSpent.Construction
	dbsqlxPlayer.TechLevelsSpentElectronics = source.TechLevelsSpent.Electronics
	dbsqlxPlayer.TechLevelsSpentBiotechnology = source.TechLevelsSpent.Biotechnology
	dbsqlxPlayer.ResearchAmount = source.ResearchAmount
	dbsqlxPlayer.ResearchSpentLastYear = source.ResearchSpentLastYear
	dbsqlxPlayer.NextResearchField = game.NextResearchField(source.NextResearchField)
	dbsqlxPlayer.Researching = game.TechField(source.Researching)
	dbsqlxPlayer.BattlePlans = GameBattlePlansToBattlePlans(source.BattlePlans)
	dbsqlxPlayer.ProductionPlans = GameProductionPlansToProductionPlans(source.ProductionPlans)
	dbsqlxPlayer.TransportPlans = GameTransportPlansToTransportPlans(source.TransportPlans)
	dbsqlxPlayer.Race = GameRaceToPlayerRace(source.Race)
	dbsqlxPlayer.Stats = GamePlayerStatsToPlayerStats(source.Stats)
	dbsqlxPlayer.Spec = GamePlayerSpecToPlayerSpec(source.Spec)
	return dbsqlxPlayer
}
func (c *GameConverter) gameRaceToDbsqlxRace(source game.Race) Race {
	var dbsqlxRace Race
	dbsqlxRace.ID = source.ID
	dbsqlxRace.CreatedAt = TimeToTime(source.CreatedAt)
	dbsqlxRace.UpdatedAt = TimeToTime(source.UpdatedAt)
	dbsqlxRace.UserID = source.UserID
	var pInt64 *int64
	if source.PlayerID != nil {
		xint64 := *source.PlayerID
		pInt64 = &xint64
	}
	dbsqlxRace.PlayerID = pInt64
	dbsqlxRace.Name = source.Name
	dbsqlxRace.PluralName = source.PluralName
	dbsqlxRace.PRT = game.PRT(source.PRT)
	dbsqlxRace.LRTs = game.Bitmask(source.LRTs)
	dbsqlxRace.HabLowGrav = source.HabLow.Grav
	dbsqlxRace.HabLowTemp = source.HabLow.Temp
	dbsqlxRace.HabLowRad = source.HabLow.Rad
	dbsqlxRace.HabHighGrav = source.HabHigh.Grav
	dbsqlxRace.HabHighTemp = source.HabHigh.Temp
	dbsqlxRace.HabHighRad = source.HabHigh.Rad
	dbsqlxRace.GrowthRate = source.GrowthRate
	dbsqlxRace.PopEfficiency = source.PopEfficiency
	dbsqlxRace.FactoryOutput = source.FactoryOutput
	dbsqlxRace.FactoryCost = source.FactoryCost
	dbsqlxRace.NumFactories = source.NumFactories
	dbsqlxRace.FactoriesCostLess = source.FactoriesCostLess
	dbsqlxRace.ImmuneGrav = source.ImmuneGrav
	dbsqlxRace.ImmuneTemp = source.ImmuneTemp
	dbsqlxRace.ImmuneRad = source.ImmuneRad
	dbsqlxRace.MineOutput = source.MineOutput
	dbsqlxRace.MineCost = source.MineCost
	dbsqlxRace.NumMines = source.NumMines
	dbsqlxRace.ResearchCostEnergy = game.ResearchCostLevel(source.ResearchCost.Energy)
	dbsqlxRace.ResearchCostWeapons = game.ResearchCostLevel(source.ResearchCost.Weapons)
	dbsqlxRace.ResearchCostPropulsion = game.ResearchCostLevel(source.ResearchCost.Propulsion)
	dbsqlxRace.ResearchCostConstruction = game.ResearchCostLevel(source.ResearchCost.Construction)
	dbsqlxRace.ResearchCostElectronics = game.ResearchCostLevel(source.ResearchCost.Electronics)
	dbsqlxRace.ResearchCostBiotechnology = game.ResearchCostLevel(source.ResearchCost.Biotechnology)
	dbsqlxRace.TechsStartHigh = source.TechsStartHigh
	dbsqlxRace.Spec = GameRaceSpecToRaceSpec(source.Spec)
	return dbsqlxRace
}
func (c *GameConverter) gameShipDesignToDbsqlxShipDesign(source game.ShipDesign) ShipDesign {
	var dbsqlxShipDesign ShipDesign
	dbsqlxShipDesign.ID = source.ID
	dbsqlxShipDesign.CreatedAt = TimeToTime(source.CreatedAt)
	dbsqlxShipDesign.UpdatedAt = TimeToTime(source.UpdatedAt)
	dbsqlxShipDesign.PlayerID = source.PlayerID
	dbsqlxShipDesign.PlayerNum = source.PlayerNum
	dbsqlxShipDesign.UUID = UUIDToUUID(source.UUID)
	dbsqlxShipDesign.Name = source.Name
	dbsqlxShipDesign.Version = source.Version
	dbsqlxShipDesign.Hull = source.Hull
	dbsqlxShipDesign.HullSetNumber = source.HullSetNumber
	dbsqlxShipDesign.CanDelete = source.CanDelete
	dbsqlxShipDesign.Slots = GameShipDesignSlotsToShipDesignSlots(source.Slots)
	dbsqlxShipDesign.Purpose = game.ShipDesignPurpose(source.Purpose)
	dbsqlxShipDesign.Spec = GameShipDesignSpecToShipDesignSpec(source.Spec)
	return dbsqlxShipDesign
}
func (c *GameConverter) gameUserToDbsqlxUser(source game.User) User {
	var dbsqlxUser User
	dbsqlxUser.ID = source.ID
	dbsqlxUser.CreatedAt = TimeToTime(source.CreatedAt)
	dbsqlxUser.UpdatedAt = TimeToTime(source.UpdatedAt)
	dbsqlxUser.Username = source.Username
	dbsqlxUser.Password = source.Password
	dbsqlxUser.Role = string(source.Role)
	return dbsqlxUser
}
func (c *GameConverter) gameVictoryConditionListToDbsqlxVictoryConditions(source []game.VictoryCondition) VictoryConditions {
	dbsqlxVictoryConditions := make(VictoryConditions, len(source))
	for i := 0; i < len(source); i++ {
		dbsqlxVictoryConditions[i] = game.VictoryCondition(source[i])
	}
	return dbsqlxVictoryConditions
}
