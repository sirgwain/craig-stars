using CraigStarsTable;
using Godot;
using System;
using System.Collections.Generic;

namespace CraigStars
{
    public class TestScene : Node2D
    {
        LineEdit lineEdit;
        CSTable table;

        public override void _Ready()
        {
            lineEdit = (LineEdit)FindNode("LineEdit");
            table = (CSTable)FindNode("CSTable");

            // add two columns
            table.Data.AddColumns(
                "Item",
                "Quantity",
                "Col3",
                new Column("hidden", hidden: true),
                new Column("Surface Minerals", scene: "res://src/Client/Controls/MineralsCell.tscn")
            );

            // add some rows of data
            table.Data.AddRowAdvanced(metadata: new Cost(2, 3, 4, 5), color: Colors.White, italic: false, "B", 30, 1, "hidden", new Cell("Surface Minerals", 10, metadata: new Cargo(10)));
            table.Data.AddRowAdvanced(metadata: new Cost(1, 2, 3, 4), color: Colors.White, italic: false, "A", 40, 2, "hidden", new Cell("Surface Minerals", 20, metadata: new Cargo(20)));
            table.Data.AddRow("C", 20, 4, "hidden", new Cell("Surface Minerals", 30, metadata: new Cargo(30)));
            table.Data.AddRow("D", 10, 3, "hidden", new Cell("Surface Minerals", 40, metadata: new Cargo(40)));

            table.ResetTable();
            table.UpdateRows();

            lineEdit.Connect("text_changed", this, nameof(OnSearchLineEditTextChanged));
        }

        void OnSearchLineEditTextChanged(string newText)
        {
            table.Data.Filter(newText);
        }

    }
}