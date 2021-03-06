using System;
using System.Collections.Generic;
using Newtonsoft.Json;

namespace CraigStars
{
    /// <summary>
    /// An item in a ProductionQueue
    /// </summary>
    public class ProductionQueueItem
    {

        /// <summary>
        /// The type of production queue item to build
        /// </summary>
        /// <value></value>
        public QueueItemType Type { get; set; }

        /// <summary>
        /// The quantity requested to build
        /// For autobuild items this is the "Up to " Quantity
        /// </summary>
        /// <value></value>
        public int Quantity { get; set; }

        /// <summary>
        /// If this is a ship building item, this is the design to build
        /// </summary>
        /// <value></value>
        [JsonProperty(IsReference = true)]
        public ShipDesign Design { get; set; }

        /// <summary>
        /// The name of the fleet to place this item into
        /// </summary>
        /// <value></value>
        public string FleetName { get; set; }

        /// <summary>
        /// The amount allocated to this item
        /// </summary>
        /// <returns></returns>
        public Cost Allocated { get; set; } = new Cost();

        [JsonIgnore]
        public int yearsToBuildOne { get; set; } = -1;

        [JsonIgnore]
        public int yearsToBuildAll { get; set; } = -1;

        [JsonIgnore]
        public float percentComplete { get; set; } = 0;

        [JsonConstructor]
        public ProductionQueueItem(QueueItemType type, int quantity = 0, ShipDesign design = null, string fleetName = null, Cost allocated = new Cost())
        {
            Type = type;
            Quantity = quantity;
            Design = design;
            FleetName = fleetName;
            Allocated = allocated;
        }

        public ProductionQueueItem Clone()
        {
            return new ProductionQueueItem(Type, Quantity, Design, FleetName, Allocated);
        }

        public override string ToString()
        {
            return $"{Type} {Quantity}{(Design != null ? " " + Design.Name : "")}";
        }

        public override bool Equals(object obj) => this.Equals(obj as ProductionQueueItem);

        public bool Equals(ProductionQueueItem p)
        {
            if (p is null)
            {
                return false;
            }

            // Optimization for a common success case.
            if (Object.ReferenceEquals(this, p))
            {
                return true;
            }

            // If run-time types are not exactly the same, return false.
            if (this.GetType() != p.GetType())
            {
                return false;
            }

            // Return true if the fields match.
            // Note that the base class is not invoked because it is
            // System.Object, which defines Equals as reference equality.
            return (Type == p.Type) && (Quantity == p.Quantity) && (Design == p.Design) && (FleetName == p.FleetName) && (Allocated == p.Allocated);
        }

        public override int GetHashCode() => (Type, Quantity, Design, FleetName).GetHashCode();

        public static bool operator ==(ProductionQueueItem lhs, ProductionQueueItem rhs)
        {
            if (lhs is null)
            {
                if (rhs is null)
                {
                    return true;
                }

                // Only the left side is null.
                return false;
            }
            // Equals handles case of null on right side.
            return lhs.Equals(rhs);
        }

        public static bool operator !=(ProductionQueueItem lhs, ProductionQueueItem rhs) => !(lhs == rhs);


        public string ShortName
        {
            get
            {
                switch (Type)
                {
                    case QueueItemType.Starbase:
                    case QueueItemType.ShipToken:
                        return Design.Name;
                    case QueueItemType.TerraformEnvironment:
                        return "Terraform Environment";
                    case QueueItemType.AutoMines:
                        return "Mine (Auto)";
                    case QueueItemType.AutoFactories:
                        return "Factory (Auto)";
                    case QueueItemType.AutoDefenses:
                        return "Defenses (Auto)";
                    case QueueItemType.AutoMineralAlchemy:
                        return "Alchemy (Auto)";
                    case QueueItemType.AutoMaxTerraform:
                        return "Max Terraform (Auto)";
                    case QueueItemType.AutoMinTerraform:
                        return "Min Terraform (Auto)";
                    default:
                        return Type.ToString();
                }
            }
        }

        public string FullName
        {
            get
            {
                switch (Type)
                {
                    case QueueItemType.Starbase:
                    case QueueItemType.ShipToken:
                        return $"{Design.Name} v{Design.Version}";
                    case QueueItemType.AutoMineralAlchemy:
                        return "Alchemy (Auto Build)";
                    case QueueItemType.AutoMines:
                        return "Mine (Auto Build)";
                    case QueueItemType.AutoFactories:
                        return "Factory (Auto Build)";
                    case QueueItemType.AutoDefenses:
                        return "Defense (Auto Build)";
                    case QueueItemType.AutoMinTerraform:
                        return "Minimum Terraform";
                    case QueueItemType.AutoMaxTerraform:
                        return "Maximum Terraform";
                    default:
                        return Type.ToString();
                }
            }
        }

        /// <summary>
        /// True if this is a mineral packet
        /// </summary>
        /// <value></value>
        public bool IsPacket
        {
            get =>
                Type == QueueItemType.IroniumMineralPacket ||
                Type == QueueItemType.BoraniumMineralPacket ||
                Type == QueueItemType.GermaniumMineralPacket ||
                Type == QueueItemType.MixedMineralPacket ||
                Type == QueueItemType.AutoMineralPacket;
        }

        /// <summary>
        /// True if this is an auto build
        /// </summary>
        /// <value></value>
        public bool IsAuto
        {
            get =>
                Type == QueueItemType.AutoDefenses ||
                Type == QueueItemType.AutoFactories ||
                Type == QueueItemType.AutoMaxTerraform ||
                Type == QueueItemType.AutoMineralAlchemy ||
                Type == QueueItemType.AutoMineralPacket ||
                Type == QueueItemType.AutoMines ||
                Type == QueueItemType.AutoMinTerraform;
        }

        /// <summary>
        /// True if this is a mineral packet
        /// </summary>
        /// <value></value>
        public bool IsTerraform
        {
            get =>
                Type == QueueItemType.TerraformEnvironment ||
                Type == QueueItemType.AutoMaxTerraform ||
                Type == QueueItemType.AutoMinTerraform;
        }

        /// <summary>
        /// True if this is a mineral alchemy item
        /// </summary>
        /// <value></value>
        public bool IsMineralAlchemy
        {
            get =>
                Type == QueueItemType.MineralAlchemy ||
                Type == QueueItemType.AutoMineralAlchemy;
        }

        /// <summary>
        /// Get the cost of a single item in this ProductionQueueItem
        /// </summary>
        /// <param name="rules"></param>
        /// <param name="race"></param>
        /// <returns></returns>
        public Cost GetCostOfOne(Player player)
        {
            var rules = player.Rules;
            var race = player.Race;
            int resources = 0;
            int germanium = 0;
            if (Type == QueueItemType.Mine || Type == QueueItemType.AutoMines)
            {
                resources = race.MineCost;
            }
            else if (Type == QueueItemType.Factory || Type == QueueItemType.AutoFactories)
            {
                resources = race.FactoryCost;
                germanium = rules.FactoryCostGermanium;
                if (race.FactoriesCostLess)
                {
                    germanium = germanium - 1;
                }
            }
            else if (Type == QueueItemType.MineralAlchemy || Type == QueueItemType.AutoMineralAlchemy)
            {
                if (race.HasLRT(LRT.MA))
                {
                    resources = rules.MineralAlchemyLRTCost;
                }
                else
                {
                    resources = rules.MineralAlchemyCost;
                }
            }
            else if (Type == QueueItemType.Defenses || Type == QueueItemType.AutoDefenses)
            {
                return rules.DefenseCost;
            }
            else if (Type == QueueItemType.TerraformEnvironment || Type == QueueItemType.AutoMaxTerraform || Type == QueueItemType.AutoMinTerraform)
            {
                return player.TerraformCost;
            }
            else if (Type == QueueItemType.IroniumMineralPacket)
            {
                return new Cost(ironium: (int)(player.MineralsPerSingleMineralPacket * player.PacketCostFactor), resources: player.PacketResourceCost);
            }
            else if (Type == QueueItemType.BoraniumMineralPacket)
            {
                return new Cost(boranium: (int)(player.MineralsPerSingleMineralPacket * player.PacketCostFactor), resources: player.PacketResourceCost);
            }
            else if (Type == QueueItemType.GermaniumMineralPacket)
            {
                return new Cost(germanium: (int)(player.MineralsPerSingleMineralPacket * player.PacketCostFactor), resources: player.PacketResourceCost);
            }
            else if (Type == QueueItemType.MixedMineralPacket || Type == QueueItemType.AutoMineralPacket)
            {
                float packetCostFactor = player.PacketCostFactor;
                int mineralsPerMixedMineralPacket = player.MineralsPerMixedMineralPacket;
                int packetResourceCost = player.PacketResourceCost;

                return new Cost(
                    (int)(mineralsPerMixedMineralPacket * packetCostFactor),
                    (int)(mineralsPerMixedMineralPacket * packetCostFactor),
                    (int)(mineralsPerMixedMineralPacket * packetCostFactor),
                    player.PacketResourceCost
                );
            }
            else if (Type == QueueItemType.ShipToken || Type == QueueItemType.Starbase)
            {
                // ship designs have their own cost
                return Design.Aggregate.Cost;
            }

            return new Cost(germanium: germanium, resources: resources);
        }

    }
}