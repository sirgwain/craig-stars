using CraigStars;
using Godot;
using System;

[Tool]
public class TileContainer : Control
{
    [Export]
    public string Title
    {
        get => title;
        set
        {
            title = value;
            UpdateTitle();
        }
    }
    string title = "";

    [Export(PropertyHint.File, "*.tscn")]
    public string ControlsScene
    {
        get => controlsScene; set
        {
            controlsScene = value;
            UpdateControlsScene();
        }
    }
    string controlsScene;

    public Control Controls { get => controls; }

    Label titleLabel;
    Control controls;
    TextureRect disclosureButton;
    Control titleContainer;

    public override void _Ready()
    {
        titleContainer = GetNode<Container>("VBoxContainer/TitleContainer");
        titleLabel = GetNode<Label>("VBoxContainer/TitleContainer/TitleLabel");
        disclosureButton = GetNode<TextureRect>("VBoxContainer/TitleContainer/DisclosureButton");
        controls = GetNode<Control>("VBoxContainer/Controls");

        titleContainer.Connect("gui_input", this, nameof(OnTitleContainerGuiInput));

        UpdateControlsScene();
        UpdateTitle();
    }

    void OnTitleContainerGuiInput(InputEvent @event)
    {
        if (@event.IsActionPressed("ui_select"))
        {
            controls.Visible = !controls.Visible;
            disclosureButton.FlipV = !controls.Visible;
        }
    }

    void UpdateControlsScene()
    {
        if (ControlsScene != null && controls != null)
        {
            try
            {
                foreach (Node child in controls.GetChildren())
                {
                    if (child is ITileContent childTileContent)
                    {
                        childTileContent.UpdateTitleEvent -= OnUpdateTitle;
                        childTileContent.UpdateVisibilityEvent -= OnUpdateVisibility;
                    }

                    controls.RemoveChild(child);
                    child.QueueFree();
                }
                var scene = GD.Load<PackedScene>(ControlsScene);
                var instance = scene.Instance();

                if (instance is ITileContent titleContent)
                {
                    titleContent.UpdateTitleEvent += OnUpdateTitle;
                    titleContent.UpdateVisibilityEvent += OnUpdateVisibility;
                }
                controls.AddChild(instance);
            }
            catch (Exception e)
            {
                GD.PrintErr("Failed to instantiate controls scene: " + ControlsScene, e);
                GD.PrintStack();
            }
        }
    }

    protected virtual void UpdateTitle()
    {
        if (titleLabel != null && Title != null)
        {
            titleLabel.Text = Title;
        }
    }

    /// <summary>
    /// Called whenever a child control wants to update our title
    /// </summary>
    /// <param name="title"></param>
    void OnUpdateTitle(string title)
    {
        Title = title;
    }

    void OnUpdateVisibility(bool visible)
    {
        Visible = visible;
    }

}
