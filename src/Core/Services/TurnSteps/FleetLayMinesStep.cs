using System;
using System.Collections.Generic;
using System.Linq;
using CraigStars.Singletons;
using CraigStars.Utils;
using Godot;
using log4net;
using static CraigStars.Utils.Utils;

namespace CraigStars
{
    /// <summary>
    /// Lay mines
    /// </summary>
    public class FleetLayMinesStep : TurnGenerationStep
    {
        static CSLog log = LogProvider.GetLogger(typeof(FleetLayMinesStep));

        public FleetLayMinesStep(Game game) : base(game, TurnGenerationState.MineLaying) { }

        public override void Process()
        {

            // Separate our waypoint tasks into groups
            foreach (var fleet in Game.Fleets.Where(fleet =>
                fleet.Aggregate.CanLayMines &&
                fleet.Waypoints.Count > 0 &&
                fleet.Waypoints[0].Task == WaypointTask.LayMineField))
            {
                LayMineField(fleet);
            }
        }

        internal void LayMineField(Fleet fleet)
        {
            foreach (var entry in fleet.Aggregate.MineLayingRateByMineType)
            {
                // see if we are adding to an existing minefield
                var mineField = LocateExistingMineField(fleet.Player, fleet.Position, entry.Key);
                if (mineField == null)
                {
                    mineField = new MineField()
                    {
                        Name = $"{fleet.RaceName} {EnumUtils.GetLabelForMineFieldType(entry.Key)} Mine Field",
                        Player = fleet.Player,
                        Position = fleet.Position,
                        Type = entry.Key
                    };
                    EventManager.PublishMapObjectCreatedEvent(mineField);
                }

                // add to it!
                long currentMines = mineField.NumMines;
                int minesLaid = entry.Value;
                mineField.NumMines += minesLaid;
                Message.MinesLaid(fleet.Player, fleet, mineField, minesLaid);

                if (mineField.Position != fleet.Position)
                {
                    // move this minefield closer to us (in case it's not in our location)
                    // This was taken from the FreeStars codebase (like many other things)
                    mineField.Position = new Vector2(
                        minesLaid / mineField.NumMines * (fleet.Position.x - mineField.Position.x) + mineField.Position.x,
                        minesLaid / mineField.NumMines * (fleet.Position.y - mineField.Position.y) + mineField.Position.y
                    );
                }
            }
        }

        /// <summary>
        /// Find an existing minefield that contains this position and is owned by the player
        /// </summary>
        /// <param name="player"></param>
        /// <param name="position"></param>
        /// <param name="type"></param>
        /// <returns></returns>
        MineField LocateExistingMineField(Player player, Vector2 position, MineFieldType type)
        {
            foreach (var mineField in Game.MineFields.Where(mf => mf.Player == player && mf.Type == type))
            {
                if (IsPointInCircle(position, mineField.Position, mineField.Radius))
                {
                    return mineField;
                }
            }
            return null;
        }

    }
}
