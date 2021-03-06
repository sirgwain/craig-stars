
using Godot;
using Newtonsoft.Json;
using System;
using System.Collections.Generic;
using System.Linq;

namespace CraigStars
{
    // [JsonObject(IsReference = true)]
    public class PublicPlayerInfo
    {
        public int NetworkId { get; set; }
        public int Num { get; set; }
        public string Name { get; set; }
        public virtual string RaceName { get; set; }
        public virtual string RacePluralName { get; set; }
        public Boolean Ready { get; set; } = false;
        public Boolean AIControlled { get; set; }
        public AIDifficulty AIDifficulty { get; set; } = AIDifficulty.Normal;
        public Boolean SubmittedTurn { get; set; }
        public Color Color { get; set; } = Colors.Black;
        public PlayerScore PublicScore { get; set; } = new();
        public bool Victor { get; set; }
        public HashSet<VictoryConditionType> AchievedVictoryConditions { get; set; } = new HashSet<VictoryConditionType>();

        /// <summary>
        /// Update our data from another player info (probably from a network call)
        /// </summary>
        /// <param name="playerInfo"></param>
        public void Update(PublicPlayerInfo playerInfo)
        {
            NetworkId = playerInfo.NetworkId;
            Num = playerInfo.Num;
            RaceName = playerInfo.RaceName;
            RacePluralName = playerInfo.RacePluralName;
            Name = playerInfo.Name;
            Ready = playerInfo.Ready;
            AIControlled = playerInfo.AIControlled;
            SubmittedTurn = playerInfo.SubmittedTurn;
            Color = playerInfo.Color;
        }

        public override string ToString()
        {
            var networkDescription = AIControlled ? "AI Controlled" : $"NetworkId: {NetworkId}";
            return $"Player {Num} {Name} ({networkDescription})";
        }

    }
}