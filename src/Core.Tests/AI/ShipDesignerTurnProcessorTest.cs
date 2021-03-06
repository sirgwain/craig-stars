using Godot;
using System;
using System.Collections.Generic;
using NUnit.Framework;

using CraigStars;
using CraigStars.Singletons;

namespace CraigStars.Tests
{
    [TestFixture]
    public class ShipDesignerTurnProcessorTest
    {
        ShipDesignerTurnProcessor shipDesignerTurnProcessor = new();

        [Test]
        public void TestDesignColonyShip()
        {
            var player = new Player()
            {
                TechLevels = new TechLevel(3, 3, 3, 3, 3, 3)
            };

            var design = shipDesignerTurnProcessor.DesignColonizer(player);

            Assert.AreEqual(Techs.MediumFreighter, design.Hull);
            Assert.AreEqual(Techs.ColonizationModule, design.Slots[1].HullComponent);

        }
    }
}