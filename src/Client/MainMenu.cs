using CraigStars.Singletons;
using Godot;
using System;

namespace CraigStars
{
    public class MainMenu : MarginContainer
    {
        WindowDialog hostWindow;
        WindowDialog joinWindow;
        LineEdit joinHostEdit;
        LineEdit hostPortEdit;
        LineEdit joinPortEdit;

        Control continueGameInfo;
        Button continueGameButton;
        Label continueGameNameLabel;
        SpinBox continueGameYearSpinBox;

        private bool joining = false;

        public override void _Ready()
        {
            hostWindow = GetNode<WindowDialog>("HostWindow");
            joinWindow = GetNode<WindowDialog>("JoinWindow");
            hostPortEdit = (LineEdit)hostWindow.FindNode("PortEdit");
            joinHostEdit = (LineEdit)joinWindow.FindNode("HostEdit");
            joinPortEdit = (LineEdit)joinWindow.FindNode("PortEdit");

            continueGameInfo = (Control)FindNode("ContinueGameInfo");
            continueGameButton = (Button)FindNode("ContinueGameButton");
            continueGameNameLabel = (Label)FindNode("ContinueGameNameLabel");
            continueGameYearSpinBox = (SpinBox)FindNode("ContinueGameYearSpinBox");

            hostPortEdit.Text = Settings.Instance.ServerPort.ToString();
            joinHostEdit.Text = Settings.Instance.ClientHost;
            joinPortEdit.Text = Settings.Instance.ClientPort.ToString();


            if (Settings.Instance.ContinueGame != null)
            {
                continueGameButton.Visible = continueGameInfo.Visible = true;
                continueGameNameLabel.Text = Settings.Instance.ContinueGame;
                continueGameYearSpinBox.Value = Settings.Instance.ContinueYear;
                continueGameYearSpinBox.MaxValue = Settings.Instance.ContinueYear;
                var minSizeRect = continueGameYearSpinBox.RectMinSize;
                minSizeRect.x = continueGameYearSpinBox.GetFont("").GetStringSize("2400").x;
                continueGameYearSpinBox.RectMinSize = minSizeRect;

                continueGameButton.Connect("pressed", this, nameof(OnContinueGameButtonPressed));
            }

            FindNode("ExitButton").Connect("pressed", this, nameof(OnExitButtonPressed));
            FindNode("SettingsButton").Connect("pressed", this, nameof(OnSettingsButtonPressed));
            FindNode("NewGameButton").Connect("pressed", this, nameof(OnNewGameButtonPressed));
            FindNode("LoadGameButton").Connect("pressed", this, nameof(OnLoadGameButtonPressed));
            FindNode("HostGameButton").Connect("pressed", this, nameof(OnHostGameButtonPressed));
            FindNode("JoinGameButton").Connect("pressed", this, nameof(OnJoinGameButtonPressed));
            FindNode("CustomRacesButton").Connect("pressed", this, nameof(OnCustomRacesButtonPressed));

            joinWindow.Connect("popup_hide", this, nameof(OnJoinWindoPopupHide));
            joinWindow.FindNode("CancelButton").Connect("pressed", this, nameof(OnJoinWindowCancelButtonPressed));
            joinWindow.FindNode("JoinButton").Connect("pressed", this, nameof(OnJoinWindowJoinButtonPressed));
            hostWindow.Connect("popup_hide", this, nameof(OnHostWindowPopupHide));
            hostWindow.FindNode("HostButton").Connect("pressed", this, nameof(OnHostWindowHostButtonPressed));

            Signals.PlayerUpdatedEvent += OnPlayerUpdated;
            GetTree().Connect("server_disconnected", this, nameof(OnServerDisconnected));
            GetTree().Connect("connection_failed", this, nameof(OnConnectionFailed));
        }

        public override void _ExitTree()
        {
            Signals.PlayerUpdatedEvent -= OnPlayerUpdated;
        }

        void OnJoinWindowCancelButtonPressed()
        {
            joining = false;
            NetworkClient.Instance.CloseConnection();
            ((Button)joinWindow.FindNode("CancelButton")).Disabled = true;
            ((Button)joinWindow.FindNode("JoinButton")).Text = "Join";
        }

        void OnJoinWindowJoinButtonPressed()
        {
            joining = true;
            ((Button)joinWindow.FindNode("CancelButton")).Disabled = false;
            ((Button)joinWindow.FindNode("JoinButton")).Text = "Joining...";
            var host = ((LineEdit)joinWindow.FindNode("HostEdit")).Text;
            var port = int.Parse(((LineEdit)joinWindow.FindNode("PortEdit")).Text);
            Settings.Instance.ClientHost = host;
            Settings.Instance.ClientPort = port;
            NetworkClient.Instance.JoinGame(host, port);
        }

        void OnExitButtonPressed()
        {
            GetTree().Quit();
        }

        void OnNewGameButtonPressed()
        {
            GetTree().ChangeScene("res://src/Client/MenuScreens/NewGameMenu.tscn");
        }

        void OnLoadGameButtonPressed()
        {
            GetTree().ChangeScene("res://src/Client/MenuScreens/LoadGameMenu.tscn");
        }

        void OnContinueGameButtonPressed()
        {
            // like a new game, but we continue
            Settings.Instance.ShouldContinueGame = true;
            Settings.Instance.ContinueYear = (int)continueGameYearSpinBox.Value;
            GetTree().ChangeScene("res://src/Client/ClientView.tscn");
        }

        void OnSettingsButtonPressed()
        {
            GetTree().ChangeScene("res://src/Client/MenuScreens/SettingsMenu.tscn");
        }

        void OnCustomRacesButtonPressed()
        {
            GetTree().ChangeScene("res://src/Client/MenuScreens/CustomRacesMenu.tscn");
        }

        void OnHostGameButtonPressed()
        {
            // skip the host view and just go to lobby
            // we'll move the port to settings
            OnHostWindowHostButtonPressed();
            // Hide();
            // hostWindow.PopupCentered();
        }

        void OnJoinGameButtonPressed()
        {
            Hide();
            joinWindow.PopupCentered();
        }

        void OnJoinWindoPopupHide()
        {
            Show();
        }

        void OnHostWindowPopupHide()
        {
            Show();
        }

        void OnHostWindowHostButtonPressed()
        {
            PlayersManager.Instance.Reset();
            PlayersManager.Instance.SetupPlayers();
            Settings.Instance.ServerPort = int.Parse(hostPortEdit.Text);
            Network.Instance.HostGame(Settings.Instance.ServerPort);
            Network.Instance.BeginGame();
            GetTree().ChangeScene("res://src/Client/MenuScreens/Lobby.tscn");
        }

        public void OnServerDisconnected()
        {
            joining = false;
        }

        public void OnConnectionFailed()
        {
            joining = false;
        }

        void OnPlayerUpdated(PublicPlayerInfo player)
        {
            if (joining && this.IsClient() && player.NetworkId == this.GetNetworkId())
            {
                GetTree().ChangeScene("res://src/Client/MenuScreens/Lobby.tscn");
            }
        }
    }
}