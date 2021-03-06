using Godot;
using System;

namespace CraigStars
{
    public class RemoteMiningWaypointTaskContainer : VBoxContainer
    {
        Label remoteMiningLabel;

        Container remoteMiningSummaryContainer;
        Label ironium;
        Label boranium;
        Label germanium;

        public Planet Planet
        {
            get => planet;
            set
            {
                planet = value;
                UpdateControls();
            }
        }
        Planet planet;

        public Fleet Fleet
        {
            get => fleet;
            set
            {
                fleet = value;
                UpdateControls();
            }
        }
        Fleet fleet;

        public override void _Ready()
        {
            remoteMiningLabel = GetNode<Label>("RemoteMiningLabel");

            remoteMiningSummaryContainer = GetNode<Container>("RemoteMiningSummaryContainer");
            ironium = GetNode<Label>("RemoteMiningSummaryContainer/Ironium");
            boranium = GetNode<Label>("RemoteMiningSummaryContainer/Boranium");
            germanium = GetNode<Label>("RemoteMiningSummaryContainer/Germanium");
        }

        void UpdateControls()
        {
            remoteMiningSummaryContainer.Visible = false;
            remoteMiningLabel.Modulate = Colors.Red;

            // update our mining stats
            if (Fleet != null && Planet != null && fleet.Aggregate.MiningRate > 0 && Planet.Explored && Planet.Owner == null)
            {
                remoteMiningSummaryContainer.Visible = true;

                remoteMiningLabel.Modulate = Colors.White;
                remoteMiningLabel.Text = "Mining Rate per Year:";
                Mineral output = Planet.GetMineralOutput(Fleet.Aggregate.MiningRate);
                ironium.Text = output.Ironium.ToString();
                boranium.Text = output.Boranium.ToString();
                germanium.Text = output.Germanium.ToString();
            }
            else
            {
                if (Planet != null && Planet.Owner != null)
                {
                    remoteMiningLabel.Text = "Note: You can only remote mine unoccupied planets.";
                }
                else if (Planet != null && !Planet.Explored)
                {
                    remoteMiningLabel.Text = "Warning: This planet is unexplored. We have no way of knowing if we can mine it.";
                }
                else if (!(fleet?.Aggregate?.MiningRate > 0))
                {
                    remoteMiningLabel.Text = "Warning: This fleet contains no ships with remote mining modules.";
                }
                else
                {
                    remoteMiningLabel.Text = "Warning: Something went wrong.";
                }
            }

        }

    }
}