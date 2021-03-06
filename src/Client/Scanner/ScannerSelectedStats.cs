using Godot;
using System;
using CraigStars.Singletons;

namespace CraigStars
{
    public class ScannerSelectedStats : MarginContainer
    {
        Label idLabel;
        Label xLabel;
        Label yLabel;
        Label nameLabel;
        Label distanceLabel;

        MapObjectSprite highlightedMapObject;
        MapObjectSprite selectedMapObject;
        MapObjectSprite commandedMapObject;

        public override void _Ready()
        {
            idLabel = FindNode("IdLabel") as Label;
            xLabel = FindNode("XLabel") as Label;
            yLabel = FindNode("YLabel") as Label;
            nameLabel = FindNode("NameLabel") as Label;
            distanceLabel = FindNode("DistanceLabel") as Label;


            Signals.MapObjectHighlightedEvent += OnMapObjectHighlighted;
            Signals.MapObjectSelectedEvent += OnMapObjectSelected;
            Signals.MapObjectCommandedEvent += OnMapObjectCommanded;
        }

        public override void _ExitTree()
        {
            Signals.MapObjectHighlightedEvent -= OnMapObjectHighlighted;
            Signals.MapObjectSelectedEvent -= OnMapObjectSelected;
            Signals.MapObjectCommandedEvent -= OnMapObjectCommanded;
        }

        void OnMapObjectHighlighted(MapObjectSprite mapObject)
        {
            highlightedMapObject = mapObject;
            if (mapObject == null)
            {
                // if we don't have a highlighted map object, use the selected one
                highlightedMapObject = selectedMapObject;
            }
            if (mapObject != null)
            {
                if (mapObject is PlanetSprite planet)
                {
                    idLabel.Visible = true;
                    idLabel.Text = $"ID: #{planet.Planet.Id}";
                }
                else
                {
                    idLabel.Visible = false;
                }

                xLabel.Text = $"X: {mapObject.Position.x:0.#}";
                yLabel.Text = $"Y: {mapObject.Position.y:0.#}";
                nameLabel.Text = $"{mapObject.ObjectName}";

                if (highlightedMapObject != commandedMapObject && commandedMapObject != null)
                {
                    var dist = Math.Abs(highlightedMapObject.Position.DistanceTo(commandedMapObject.Position));
                    distanceLabel.Text = $"{dist:0.#} light years from {commandedMapObject.ObjectName}";
                }
                else
                {
                    distanceLabel.Text = $"";
                }
            }
        }

        void OnMapObjectSelected(MapObjectSprite mapObject)
        {
            selectedMapObject = mapObject;

        }

        void OnMapObjectCommanded(MapObjectSprite mapObject)
        {
            commandedMapObject = mapObject;
        }
    }
}