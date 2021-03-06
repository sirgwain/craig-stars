using System;
using System.Collections.Generic;
using System.Linq;
using CraigStars.Utils;
using log4net;

namespace CraigStars
{
    /// <summary>
    /// Build and deploy scout ships
    /// </summary>
    public class ColonyTurnProcessor : TurnProcessor
    {
        static CSLog log = LogProvider.GetLogger(typeof(ColonyTurnProcessor));

        // the required population density required of a planet in order to suck people off of it
        // setting this to .25 because we don't want to suck people off a planet until it's reached the
        // max of its growth rate (over 1/4 crowded)
        private const float PopulationDensityRequired = .25f;

        ShipDesignerTurnProcessor shipDesignerTurnProcessor = new ShipDesignerTurnProcessor();

        public ColonyTurnProcessor() : base("Colonizer") { }

        /// <summary>
        /// a new turn! build some ships
        /// </summary>
        public override void Process(Player player)
        {
            // find the first colony ship design
            ShipDesign colonyShip = shipDesignerTurnProcessor.DesignColonizer(player);

            var colonizablePlanets = player.AllPlanets
            .Where(planet => planet.Explored && planet.Uninhabited && player.Race.GetPlanetHabitability(planet.BaseHab.Value) > 0)
            .OrderByDescending(planet => player.Race.GetPlanetHabitability(planet.BaseHab.Value))
            .ToList();
            var buildablePlanets = player.Planets
                .Where(planet => planet.CanBuild(player, colonyShip.Aggregate.Mass) && planet.PopulationDensity >= PopulationDensityRequired)
                .ToList();
            var colonizerFleets = player.Fleets.Where(fleet => fleet.Aggregate.Purposes.Contains(ShipDesignPurpose.Colonizer));

            foreach (var fleet in colonizerFleets)
            {
                if (fleet.Waypoints.Count > 1 && fleet.Waypoints[1].Target is Planet planet)
                {
                    // don't send a colony ship to this planet
                    colonizablePlanets.Remove(planet);
                }
            }

            // go through each unassigned colonizer fleet and find it a new planet to colonize
            foreach (var fleet in colonizerFleets.Where(
                f => f.Waypoints.Count == 1 &&
                f.Orbiting?.Player == player &&
                f.Orbiting.GetPopulationDensity(f.Orbiting.Population - f.AvailableCapacity) >= PopulationDensityRequired)
            )
            {
                var planetToColonize = ClosestPlanet(fleet, colonizablePlanets);
                if (planetToColonize != null)
                {
                    var sourcePlanet = fleet.Orbiting;
                    Cargo colonists = new Cargo(colonists: fleet.AvailableCapacity);

                    if (sourcePlanet.AttemptTransfer(-colonists) && fleet.AttemptTransfer(colonists))
                    {
                        CargoTransferUtils.CreateCargoTransferOrder(fleet.Player, colonists, fleet, sourcePlanet);

                        fleet.Waypoints.Add(Waypoint.TargetWaypoint(planetToColonize, fleet.GetDefaultWarpFactor(), WaypointTask.Colonize));
                        fleet.Waypoints[1].WarpFactor = fleet.GetBestWarpFactor(fleet.Waypoints[0], fleet.Waypoints[1]);

                        // remove this planet from our list
                        colonizablePlanets.Remove(planetToColonize);
                    }
                    else
                    {
                        log.Error($"{player.Game.Year}: {player} Failed to transfer {fleet.AvailableCapacity}kT colonists from {sourcePlanet.Name} to {fleet.Name} for colonization.");
                    }

                }
            }

            BuildFleets(buildablePlanets, colonizablePlanets.Count, colonyShip);

        }

        /// <summary>
        /// Build any new colonizer fleets
        /// </summary>
        /// <param name="buildablePlanets"></param>
        /// <param name="numShipsNeeded"></param>
        /// <param name="colonyShip"></param>
        void BuildFleets(List<Planet> buildablePlanets, int numShipsNeeded, ShipDesign colonyShip)
        {
            if (colonyShip != null)
            {

                int queuedToBeBuilt = 0;

                List<Planet> planetsToBuildOn = new List<Planet>();
                foreach (Planet planet in buildablePlanets)
                {
                    bool isBuilding = false;
                    foreach (ProductionQueueItem item in planet.ProductionQueue?.Items)
                    {
                        if (item.Design != null && item.Design.Aggregate.Colonizer)
                        {
                            isBuilding = true;
                            queuedToBeBuilt++;
                            log.Debug($"{planet.Player.Game.Year}: {planet.Name} is already building a colony ship");
                        }
                    }

                    // if this planet isn't already building a colony ship, build one
                    if (!isBuilding)
                    {
                        planetsToBuildOn.Add(planet);
                    }
                }

                foreach (Planet planet in planetsToBuildOn)
                {
                    // if this planet isn't building a colony ship already, add
                    // one to the queue
                    if (queuedToBeBuilt < numShipsNeeded)
                    {
                        planet.ProductionQueue.AddAfter(new ProductionQueueItem(QueueItemType.ShipToken, 1, colonyShip), QueueItemType.AutoFactories);
                        log.Debug($"Added scout ship to planet queue: {planet.Name}");
                        queuedToBeBuilt++;
                    }
                }

            }

        }
    }
}