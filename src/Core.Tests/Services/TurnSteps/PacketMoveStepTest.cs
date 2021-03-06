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
    public class PacketMoveStepTest
    {
        static CSLog log = LogProvider.GetLogger(typeof(PacketMoveStepTest));


        [Test]
        public void TestCompleteMoveCaught()
        {
            var game = TestUtils.GetSingleUnitGame();
            PacketMoveStep step = new PacketMoveStep(game, 1);
            var player1 = game.Players[0];

            // create a starbase with a warp5 receiver
            var catchingStarbase = new ShipDesign()
            {
                Player = player1,
                Name = "Starbase",
                Hull = Techs.SpaceStation,
                HullSetNumber = 0,
                Slots = new List<ShipDesignSlot>()
                {
                    new ShipDesignSlot(Techs.UltraDriver10, 11, 1),
                }
            };
            var planet2 = new Planet()
            {
                Player = player1,
                Population = 250000,
                Defenses = 50,
                Starbase = new Starbase()
                {
                    Player = player1,
                    Tokens = new List<ShipToken>() {
                        new ShipToken(catchingStarbase, 1)
                    }
                }
            };

            // compute aggregate for this starbase so the receiver is up to date
            planet2.Starbase.ComputeAggregate();

            // create a 1000kT packet
            MineralPacket packet = new MineralPacket()
            {
                Player = player1,
                SafeWarpSpeed = 10,
                WarpFactor = 10,
                Cargo = new Cargo(1000),
                Target = planet2
            };

            // make landfall
            step.CompleteMove(packet);

            // starbase can accept this packet, should have same defenses/pop, and 1000kT more minerals
            Assert.AreEqual(250000, planet2.Population);
            Assert.AreEqual(50, planet2.Defenses);
            Assert.AreEqual(new Cargo(1000), planet2.Cargo.WithColonists(0));
        }

        [Test]
        public void TestCompleteMoveOverspeed()
        {
            var game = TestUtils.GetSingleUnitGame();
            PacketMoveStep step = new PacketMoveStep(game, 1);
            var player1 = game.Players[0];
            player1.TechLevels = new TechLevel(energy: 5);

            // create a starbase with a warp5 receiver
            var catchingStarbase = new ShipDesign()
            {
                Player = player1,
                Name = "Starbase",
                Hull = Techs.SpaceStation,
                HullSetNumber = 0,
                Slots = new List<ShipDesignSlot>()
                {
                    new ShipDesignSlot(Techs.MassDriver5, 11, 1),
                }
            };
            var planet2 = new Planet()
            {
                Player = player1,
                Population = 250000,
                Defenses = 50,
                Starbase = new Starbase()
                {
                    Player = player1,
                    Tokens = new List<ShipToken>() {
                        new ShipToken(catchingStarbase, 1)
                    }
                }
            };

            // compute aggregate for this starbase so the receiver is up to date
            planet2.Starbase.ComputeAggregate();

            // create a 1000kT packet
            MineralPacket packet = new MineralPacket()
            {
                Player = player1,
                SafeWarpSpeed = 10,
                WarpFactor = 10,
                Cargo = new Cargo(1000),
                Target = planet2
            };

            // make landfall
            step.CompleteMove(packet);

            // 42900 colonists destroyed, 8 defenses destroyed
            Assert.AreEqual(250000 - 42900, planet2.Population);
            Assert.AreEqual(50 - 8, planet2.Defenses);

        }

        [Test]
        public void TestCompleteMoveKillPlanet()
        {
            var game = TestUtils.GetSingleUnitGame();
            PacketMoveStep step = new PacketMoveStep(game, 1);
            var player = game.Players[0];

            var planet2 = new Planet()
            {
                Player = player,
                Population = 250000,
                Defenses = 0,
            };

            // create a 1000kT packet
            MineralPacket packet = new MineralPacket()
            {
                Player = player,
                SafeWarpSpeed = 13,
                WarpFactor = 13,
                Cargo = new Cargo(1000),
                Target = planet2
            };

            // make landfall
            step.CompleteMove(packet);

            // planet is destroyed
            Assert.AreEqual(0, planet2.Population);
            Assert.AreEqual(0, planet2.Defenses);

        }
    }
}