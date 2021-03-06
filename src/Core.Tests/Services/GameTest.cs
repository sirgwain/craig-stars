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

namespace CraigStars.Tests
{
    [TestFixture]
    public class GameTest
    {
        static CSLog log = LogProvider.GetLogger(typeof(GameTest));



        [Test]
        public void TestGenerateUniverse()
        {
            var game = new Game() { SaveToDisk = false };
            var rules = new Rules(0);
            game.Init(new List<Player>() { new Player() }, rules, StaticTechStore.Instance, new TestGamesManager(), new TestTurnProcessorManager());
            game.GenerateUniverse();

            Assert.AreEqual(rules.GetNumPlanets(game.Size, game.Density), game.Planets.Count);
            Assert.AreEqual(rules.GetNumPlanets(game.Size, game.Density), game.Players[0].AllPlanets.ToList().Count);
            Assert.AreEqual(game.Fleets.Count, game.Players[0].Fleets.Count);
        }

        [Test]
        public async Task TestGenerateTurn()
        {
            // create a new game with universe
            var game = new Game() { SaveToDisk = false };
            var player = new Player();
            var rules = new Rules(0);
            game.Init(new List<Player>() { player }, rules, StaticTechStore.Instance, new TestGamesManager(), new TestTurnProcessorManager());
            game.GenerateUniverse();

            // submit the player
            game.SubmitTurn(player);

            // generate the turn
            await game.GenerateTurn();

            // make sure our turn was generated and the player's report was updated
            Assert.Greater(player.Homeworld.Population, rules.StartingPopulation);
            Assert.AreEqual(0, player.Homeworld.ReportAge);
        }

        /// <summary>
        /// Test generating multiple turns with an AI and Human player
        /// </summary>
        [Test]
        public async Task TestGenerateManyTurns()
        {
            // create a new game with universe
            var game = new Game()
            {
                SaveToDisk = false,
                Size = Size.Huge,
                Density = Density.Packed,
            };
            game.GameInfo.QuickStartTurns = 0;
            var player = new Player();
            var aiPlayer = new Player() { AIControlled = true };
            var rules = new Rules(0);
            game.Init(new List<Player>() { player, aiPlayer }, rules, StaticTechStore.Instance, new TestGamesManager(), new TestTurnProcessorManager());
            game.GenerateUniverse();

            // // turn off logging but for errors
            // var logger = (Logger)log.Logger;
            // logger.Hierarchy.Root.Level = Level.Error;

            // generate a thousand turns
            var stopwatch = new Stopwatch();
            stopwatch.Start();
            int numTurns = 10; //1000;

            // generate a bunch of turns
            for (int i = 0; i < numTurns; i++)
            {
                // submit the player
                game.SubmitTurn(player);

                // generate the turn
                await game.GenerateTurn();
            }
            stopwatch.Stop();

            // turn back on logging defaults
            // logger.Hierarchy.Root.Level = Level.All;
            log.Debug($"Generated {numTurns} turns in {stopwatch.ElapsedMilliseconds / 1000.0f} seconds");

            // make sure our turn was generated and the player's report was updated
            Assert.Greater(player.Homeworld.Population, rules.StartingPopulation);
            Assert.AreEqual(0, player.Homeworld.ReportAge);
            Assert.AreEqual(rules.StartingYear + numTurns, game.Year);
        }

    }

}