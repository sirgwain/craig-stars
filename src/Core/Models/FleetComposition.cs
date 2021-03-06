using System;
using System.Collections.Generic;

namespace CraigStars
{
    /// <summary>
    /// Players (and AIs) use FleetCompositions to determine what they want a fleet to look like.
    /// </summary>
    public class FleetComposition
    {
        public Guid Guid { get; set; }
        public List<ShipToken> Tokens { get; set; } = new List<ShipToken>();
    }
}