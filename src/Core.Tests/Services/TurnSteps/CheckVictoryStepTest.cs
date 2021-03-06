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
    public class CheckVictoryStepTest
    {

        [Test]
        public void TestCheckOwnPlanets()
        {
            PlayerIntel playerIntel = new();
            var game = TestUtils.GetSingleUnitGame();
            var player = game.Players[0];

            // we must own >= 51% of the planets
            game.VictoryConditions.OwnPlanets = 51;

            // this player owns all planets
            // Discover this planet so our score calculation is aware of it
            playerIntel.Discover(player, game.Planets[0]);
            var step = new CheckVictoryStep(game);
            step.CheckOwnPlanets(player);

            Assert.IsTrue(player.AchievedVictoryConditions.Contains(VictoryConditionType.OwnPlanets));

            // add a planet, now our player only owns half the planets
            game.Planets.Add(new Planet() { Name = "Planet 2" });
            player.AchievedVictoryConditions.Clear();
            step.CheckOwnPlanets(player);
            Assert.IsFalse(player.AchievedVictoryConditions.Contains(VictoryConditionType.OwnPlanets));
        }

        [Test]
        public void TestCheckExceedSecondPlaceScore()
        {
            var game = TestUtils.GetSingleUnitGame();
            var player1 = game.Players[0];
            var player2 = new Player();
            game.Players.Add(player2);


            player1.Score.Score = 201;
            player2.Score.Score = 100;

            // we must exceed the second player score by 100%
            game.VictoryConditions.ExceedScore = 100;

            // this player owns all planets
            var step = new CheckVictoryStep(game);
            step.CheckExceedSecondPlaceScore(player1);

            Assert.IsTrue(player1.AchievedVictoryConditions.Contains(VictoryConditionType.ExceedSecondPlaceScore));

            // we only exceed the score by 20%
            player1.AchievedVictoryConditions.Clear();
            player1.Score.Score = 120;
            step.CheckExceedSecondPlaceScore(player1);
            Assert.IsFalse(player1.AchievedVictoryConditions.Contains(VictoryConditionType.ExceedSecondPlaceScore));
        }

    }
}