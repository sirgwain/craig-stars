using System;
using System.Collections.Generic;
using System.ComponentModel;
using Newtonsoft.Json;

namespace CraigStars
{
    public class Rules
    {
        /// <summary>
        /// Allow setting of the random seed used
        /// </summary>
        /// <returns></returns>
        [JsonIgnore]
        public Random Random { get; set; } = new Random();

        [DefaultValue(2400)]
        public int StartingYear { get; set; } = 2400;

        /// <summary>
        /// If PublicPlayerScores is true, this is how many years to wait before you discover other player's scores.
        /// </summary>
        /// <value></value>
        [DefaultValue(20)]
        public int ShowPublicScoresAfterYears { get; set; } = 1;

        [DefaultValue(15)]
        public int PlanetMinDistance { get; } = 15;

        [DefaultValue(180)]
        public int MaxExtraWorldDistance { get; set; } = 180;

        [DefaultValue(130)]
        public int MinExtraWorldDistance { get; set; } = 130;

        // Mineral rules that we don't currently modify
        [DefaultValue(30)]
        public int MinHomeworldMineralConcentration { get; set; } = 30;

        [DefaultValue(30)]
        public int MinExtraPlanetMineralConcentration { get; set; } = 30;

        [DefaultValue(1)]
        public int MinMineralConcentration { get; set; } = 1;

        [DefaultValue(3)]
        public int MinStartingMineralConcentration { get; set; } = 3;

        [DefaultValue(100)]
        public int MaxStartingMineralConcentration { get; set; } = 100;

        [DefaultValue(1000)]
        public int MaxStartingMineralSurface { get; set; } = 1000;

        [DefaultValue(300)]
        public int MinStartingMineralSurface { get; set; } = 300;

        [DefaultValue(1500000)]
        public int MineralDecayFactor { get; set; } = 1500000;

        // Population rules
        [DefaultValue(25000)]
        public int StartingPopulation { get; set; } = 25000;
        [DefaultValue(20000)]
        public int StartingPopulationWithExtraPlanet { get; set; } = 20000;
        [DefaultValue(10000)]
        public int StartingPopulationExtraPlanet { get; set; } = 10000;
        [DefaultValue(.7f)]
        public float LowStartingPopulationFactor { get; set; } = .7f;

        // Bulding rules
        [DefaultValue(10)]
        public int StartingMines { get; set; } = 10;
        [DefaultValue(10)]
        public int StartingFactories { get; set; } = 10;
        [DefaultValue(10)]
        public int StartingDefenses { get; set; } = 10;

        // Race rules
        [DefaultValue(1650)]
        public int RaceStartingPoints { get; set; } = 1650;

        [DefaultValue(4)]
        public int FactoryCostGermanium { get; set; } = 4;
        public Cost DefenseCost { get; set; } = new Cost(5, 5, 5, 15);

        [DefaultValue(100)]
        public int MineralAlchemyCost { get; set; } = 100;

        [DefaultValue(25)]
        public int MineralAlchemyLRTCost { get; set; } = 25;

        public Cost TerraformCost { get; set; } = new Cost(0, 0, 0, 100);
        public Cost TotalTerraformCost { get; set; } = new Cost(0, 0, 0, 70);

        [DefaultValue(10)]
        public int PacketResourceCost { get; set; } = 10;
        [DefaultValue(5)]
        public int PacketResourceCostPP { get; set; } = 5;

        [DefaultValue(1.1f)]
        public float PacketMineralCostFactor { get; set; } = 1.1f;
        [DefaultValue(1f)]
        public float PacketMineralCostFactorPP { get; set; } = 1f;
        [DefaultValue(1.2f)]
        public float PacketMineralCostFactorIT { get; set; } = 1.2f;

        /// <summary>
        /// The decay rate for packets for the amount over the safe warp speed
        /// </summary>
        public Dictionary<int, float> PacketDecayRate { get; set; } = new Dictionary<int, float>() {
            { 1, .1f },
            { 2, .25f },
            { 3, .5f },
        };

        /// <summary>
        /// The amount of minerals in a single ironium, boranium, or germanium packet
        /// </summary>
        /// <value></value>
        [DefaultValue(100)]
        public int MineralsPerSingleMineralPacket { get; set; } = 100;
        [DefaultValue(70)]
        public int MineralsPerSingleMineralPacketPP { get; set; } = 70;


        /// <summary>
        /// The amount of minerals in a mixed (all minerals at once) packet
        /// </summary>
        /// <value></value>
        [DefaultValue(40)]
        public int MineralsPerMixedMineralPacket { get; set; } = 40;
        [DefaultValue(25)]
        public int MineralsPerMixedMineralPacketPP { get; set; } = 25;


        [DefaultValue(20)]
        public int BuiltInScannerJoaTMultiplier = 20;

        /// <summary>
        /// SS races have 300 cloak units built in (or 75% cloaking)
        /// </summary>
        [DefaultValue(300)]
        public int BuiltInSSCloakUnits = 300;

        [DefaultValue(5)]
        public int TachyonCloakReduction = 5;

        // Game rules
        [DefaultValue(1000000)]
        public int MaxPopulation = 1000000;

        /// <summary>
        /// The amount scanned populations differ from the actual population on the planet
        /// </summary>
        [DefaultValue(.2f)]
        public float PopulationScannerError = .2f;

        /// <summary>
        /// The factor that defenses coverage is multiplied by against smart bombs (i.e. halved)
        /// </summary>
        [DefaultValue(.5f)]
        public float SmartDefenseCoverageFactor = .5f;

        /// <summary>
        /// The factor that defenses coverage is multiplied by against invasions
        /// </summary>
        [DefaultValue(.5f)]
        public float InvasionDefenseCoverageFactor = .75f;

        /// <summary>
        /// The number of battle rounds per battle
        /// </summary>
        [DefaultValue(16)]
        public int NumBattleRounds = 16;

        /// <summary>
        /// The number of moves required before a token can run away
        /// </summary>
        [DefaultValue(7)]
        public int MovesToRunAway = 7;

        /// <summary>
        /// The maximum range penalty for beam weapons, i.e. 10%
        /// </summary>
        [DefaultValue(.1f)]
        public float BeamRangeDropoff = .1f;

        /// <summary>
        /// The number of moves required before a token can run away
        /// </summary>
        [DefaultValue(1 / 8f)]
        public float TorpedoSplashDamage = 1 / 8f;

        [DefaultValue(.1f)]
        public float SalvageDecayRate = .1f;
        [DefaultValue(10)]
        public int SalvageDecayMin = 10;


        /// <summary>
        /// MineFields are cloaked to 75% until spotted
        /// </summary>
        [DefaultValue(75)]
        public int MineFieldCloak = 75;

        /// <summary>
        /// Space Demolition fleets can travel 2 warp speeds faster through minefields
        /// </summary>
        [DefaultValue(2)]
        public int SDSafeWarpBonus = 2;
        [DefaultValue(.25f)]
        public float SDMinDecayFactor = .25f;
        [DefaultValue(.02f)]
        public float MineFieldBaseDecayRate = .02f;
        [DefaultValue(.04f)]
        public float MineFieldPlanetDecayRate = .04f;
        [DefaultValue(.25f)]
        public float MineFieldDetonateDecayRate = .25f;
        [DefaultValue(.5f)]
        public float MineFieldMaxDecayRate = .5f;

        /// <summary>
        /// The maximum factor for safe range, i.e. 5x safe range
        /// is the max range of a stargate
        /// </summary>
        [DefaultValue(5)]
        public int StargateMaxRangeFactor = 5;

        /// <summary>
        /// The maximum factor for safe hullmass, i.e. 5x safe hull mass
        /// is the max mass of a stargate
        /// </summary>
        [DefaultValue(5)]
        public int StargateMaxHullMassFactor = 5;

        /// <summary>
        /// MineFields are cloaked to 75% until spotted
        /// </summary>
        [DefaultValue(75)]
        public int WormholeCloak = 75;

        /// <summary>
        /// The minimum distance between a wormhole and any other object
        /// </summary>
        [DefaultValue(20)]
        public int WormholeMinDistance = 30;

        /// <summary>
        /// Every year wormholes degrage, jiggle, and sometimes jump
        /// </summary>
        public Dictionary<WormholeStability, WormholeStats> WormholeStatsByStability = new Dictionary<WormholeStability, WormholeStats>() {
            { WormholeStability.RockSolid, new WormholeStats(10, 0) },
            { WormholeStability.Stable, new WormholeStats(5, .005f) },
            { WormholeStability.MostlyStable, new WormholeStats(5, .02f) },
            { WormholeStability.Average, new WormholeStats(5, .04f) },
            { WormholeStability.SlightlyVolatile, new WormholeStats(5, .03f) },
            { WormholeStability.Volatile, new WormholeStats(5, .06f) },
            { WormholeStability.ExtremelyVolatile, new WormholeStats(int.MaxValue, .04f) },
        };

        /// <summary>
        /// The number of wormhole pairs for the size of the universe
        /// </summary>
        public Dictionary<Size, int> WormholePairsForSize = new Dictionary<Size, int>() {
            { Size.Tiny, 1},
            { Size.Small, 3},
            { Size.Medium, 4},
            { Size.Large, 5},
            { Size.Huge, 6},
        };

        /// <summary>
        /// Each minefield type has stats
        /// </summary>
        public Dictionary<MineFieldType, MineFieldStats> MineFieldStatsByType = new Dictionary<MineFieldType, MineFieldStats>() {
            {
                MineFieldType.Standard,
                new MineFieldStats() {
                    MinDamagePerFleetRS = 600,
                    DamagePerEngineRS = 125,
                    MaxSpeed = 4,
                    ChanceOfHit = .003f,
                    MinDamagePerFleet = 500,
                    DamagePerEngine = 100,
                    CanDetonate = true,
                    MinDecay = 10,
            } },
            {
                MineFieldType.Heavy,
                new MineFieldStats() {
                    MinDamagePerFleetRS = 2500,
                    DamagePerEngineRS = 600,
                    MaxSpeed = 6,
                    ChanceOfHit = .01f,
                    MinDamagePerFleet = 2000,
                    DamagePerEngine = 500,
                    MinDecay = 10,
            } },
            {
                MineFieldType.SpeedBump,
                new MineFieldStats() {
                    MaxSpeed = 5,
                    ChanceOfHit = .035f,
                    SweepFactor = 1f/3, // speed hump minds are swept at 1/3rd the normal rate
            } },
        };

        /// <summary>
        /// Get the Area of the universe
        /// </summary>
        /// <value></value>
        public int GetArea(Size size) => size switch
        {
            Size.Tiny => 400,
            Size.Small => 800,
            Size.Medium => 1200,
            Size.Large => 1600,
            Size.Huge => 2000,
            _ => throw new System.ArgumentException("Unknown Size=>  " + size),
        };

        /// <summary>
        /// Get the number of planets based on the size and density
        /// </summary>
        /// <value></value>
        public int GetNumPlanets(Size size, Density density) => size switch
        {
            Size.Huge => density switch
            {
                Density.Sparse => 600,
                Density.Normal => 800,
                Density.Dense => 940,
                Density.Packed => 945,
                _ => throw new System.ArgumentException($"Unknown Size {size} or Density {density}")
            },
            Size.Large => density switch
            {
                Density.Sparse => 384,
                Density.Normal => 512,
                Density.Dense => 640,
                Density.Packed => 910,
                _ => throw new System.ArgumentException($"Unknown Size {size} or Density {density}")
            },
            Size.Medium => density switch
            {
                Density.Sparse => 216,
                Density.Normal => 288,
                Density.Dense => 360,
                Density.Packed => 540,
                _ => throw new System.ArgumentException($"Unknown Size {size} or Density {density}")
            },
            Size.Small => density switch
            {
                Density.Sparse => 96,
                Density.Normal => 128,
                Density.Dense => 160,
                Density.Packed => 240,
                _ => throw new System.ArgumentException($"Unknown Size {size} or Density {density}")
            },
            Size.Tiny => density switch
            {

                Density.Sparse => 24,
                Density.Normal => 32,
                Density.Dense => 40,
                Density.Packed => 60,
                _ => throw new System.ArgumentException($"Unknown Size {size} or Density {density}")
            },
            _ => throw new System.ArgumentException($"Unknown Size {size} or Density {density}")
        };

        /// <summary>
        /// The maximum level achievable in any tech
        /// Note, if this is increased, additional fields must be added to TechBaseCost
        /// </summary>
        /// <value></value>
        [DefaultValue(26)]
        public int MaxTechLevel { get; set; } = 26;


        // The base cost for each tech level
        public int[] TechBaseCost { get; set; } = {
            0, // everyone starts on level 0
            50,
            80,
            130,
            210,
            340,
            550,
            890,
            1440,
            2330,
            3770,
            6100,
            9870,
            13850,
            18040,
            22440,
            27050,
            31870,
            36900,
            42140,
            47590,
            53250,
            59120,
            65200,
            71490,
            77990,
            84700
        };

        public Rules() { }

        /// <summary>
        /// Create a new Rules object with a specific random seed. This is used for unit tests but could also be used for
        /// world generation
        /// </summary>
        /// <param name="seed"></param>
        public Rules(int seed)
        {
            Random = new Random(seed);
        }
    }


}