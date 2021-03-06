using Godot;
using System;
using System.Collections.Generic;
using NUnit.Framework;

using CraigStars.Singletons;
using log4net;
using System.Diagnostics;
using log4net.Core;
using log4net.Repository.Hierarchy;
using System.Threading.Tasks;
using System.Linq;
using CraigStars.UniverseGeneration;

namespace CraigStars.Tests
{
    [TestFixture]
    public class PlayerScanStepTest
    {
        static CSLog log = LogProvider.GetLogger(typeof(PlayerScanStepTest));

        [Test]
        public void TestScan()
        {
            var game = TestUtils.GetSingleUnitGame();
            game.Planets[0].Population = 120000;

            var scanStep = new PlayerScanStep(game);
            scanStep.Execute(new TurnGenerationContext(), game.OwnedPlanets.ToList());

            // our player should know about the planet updates
            Assert.AreEqual(game.Planets[0].Population, game.Players[0].Planets[0].Population);
            Assert.IsTrue(game.Players[0].Fleets.Count > 0);
        }

        [Test]
        public void TestScanFleets()
        {
            PlayerIntel playerIntel = new();
            var game = TestUtils.GetTwoPlayerGame();
            var player1 = game.Players[0];

            // our scan won't work unless we know about our own designs
            playerIntel.Discover(player1, game.Designs[0]);

            var fleet2 = game.Fleets[1];

            // create a pen scanner on top of our other fleet (which is behind a planet)
            var scanners = new List<PlayerScanStep.Scanner>() {
                new PlayerScanStep.Scanner(fleet2.Position, rangePen: 1)
            };

            // we should discover this fleet
            var scanStep = new PlayerScanStep(game);
            scanStep.ScanFleets(player1, scanners);

            Assert.AreEqual(1, player1.ForeignFleets.Count);
            Assert.AreEqual(fleet2.Guid, player1.ForeignFleets[0].Guid);
            Assert.AreEqual(1, player1.ForeignDesigns.Count);
            Assert.AreEqual(fleet2.Tokens[0].Design.Guid, player1.ForeignDesigns[0].Guid);
        }

        // Disabled because this is long running
        // [Test]
        public void TestScanManyFleets()
        {
            PlayerIntel playerIntel = new();
            var game = TestUtils.GetTwoPlayerGame();
            var player1 = game.Players[0];
            var player2 = game.Players[1];

            // our scan won't work unless we know about our own designs
            playerIntel.Discover(player1, game.Designs[0]);

            // make many fleets
            player1.TechLevels = new TechLevel(electronics: 16);
            var planet = game.Planets[1];
            var random = game.Rules.Random;
            var scanners = new List<PlayerScanStep.Scanner>();
            for (int i = 0; i < 1000; i++)
            {
                var isOrbiting = random.Next() % 2 > 0;

                var position1 = isOrbiting ? planet.Position : new Vector2(random.Next(1000), random.Next(1000));
                var design1 = ShipDesigns.LongRangeScount.Clone(player1);
                design1.Name = "Design " + i;
                var fleet1 = new Fleet()
                {
                    Player = player1,
                    Name = "New Fleet " + i,
                    Position = position1,
                    Tokens = new List<ShipToken>() {
                        new ShipToken(design1, 1),
                    },
                    BattlePlan = player1.BattlePlans[0],
                    Orbiting = isOrbiting ? planet : null,
                    Waypoints = new List<Waypoint>() {
                        isOrbiting ? Waypoint.TargetWaypoint(planet) : Waypoint.PositionWaypoint(position1)
                    }
                };
                game.Fleets.Add(fleet1);
                game.Designs.Add(design1);

                var position2 = isOrbiting ? planet.Position : new Vector2(random.Next(1000), random.Next(1000));
                var design2 = ShipDesigns.Teamster.Clone(player2);
                design2.Name = "Design " + i;
                var fleet2 = new Fleet()
                {
                    Player = player2,
                    Name = "New Fleet " + i,
                    Position = position2,
                    Tokens = new List<ShipToken>() {
                        new ShipToken(design2, 1),
                    },
                    BattlePlan = player2.BattlePlans[0],
                    Orbiting = isOrbiting ? planet : null,
                    Waypoints = new List<Waypoint>() {
                        isOrbiting ? Waypoint.TargetWaypoint(planet) : Waypoint.PositionWaypoint(position2)
                    }
                };
                game.Fleets.Add(fleet2);
                game.Designs.Add(design2);


            }

            game.UpdateDictionaries();
            game.ComputeAggregates();

            // we should discover this fleet
            var scanStep = new PlayerScanStep(game);

            Stopwatch stopwatch = new Stopwatch();
            stopwatch.Reset();
            stopwatch.Start();

            for (int i = 0; i < 500; i++)
            {
                scanStep.Execute(new TurnGenerationContext(), game.OwnedPlanets.ToList());
            }
            stopwatch.Stop();
            log.Debug($"Completed All Scans: ({TimeSpan.FromMilliseconds(stopwatch.ElapsedMilliseconds):c})");

        }
    }

}