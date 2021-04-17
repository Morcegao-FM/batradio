namespace BatRadio.UI
{
    partial class frmMain
    {
        /// <summary>
        /// Required designer variable.
        /// </summary>
        private System.ComponentModel.IContainer components = null;

        /// <summary>
        /// Clean up any resources being used.
        /// </summary>
        /// <param name="disposing">true if managed resources should be disposed; otherwise, false.</param>
        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }
            base.Dispose(disposing);
        }

        #region Windows Form Designer generated code

        /// <summary>
        /// Required method for Designer support - do not modify
        /// the contents of this method with the code editor.
        /// </summary>
        private void InitializeComponent()
        {
            this.components = new System.ComponentModel.Container();
            this.styleManager = new MetroFramework.Components.MetroStyleManager(this.components);
            this.panelMusicStatus = new MetroFramework.Controls.MetroPanel();
            this.trackCurrentMusic = new MetroFramework.Controls.MetroTrackBar();
            this.labelNowPlaying = new MetroFramework.Controls.MetroLabel();
            this.timerUpdate = new System.Windows.Forms.Timer(this.components);
            this.tabMain = new MetroFramework.Controls.MetroTabControl();
            this.tabPlaying = new MetroFramework.Controls.MetroTabPage();
            this.metroPanel2 = new MetroFramework.Controls.MetroPanel();
            this.labelShuffle = new MetroFramework.Controls.MetroLabel();
            this.toggleShuffle = new MetroFramework.Controls.MetroToggle();
            this.labelIsPlaying = new MetroFramework.Controls.MetroLabel();
            this.toggleIsPlaying = new MetroFramework.Controls.MetroToggle();
            this.metroTabPage2 = new MetroFramework.Controls.MetroTabPage();
            this.tabConfiguration = new MetroFramework.Controls.MetroTabPage();
            this.buttonUpdateMusicDatabase = new MetroFramework.Controls.MetroButton();
            this.textboxServerAPIKey = new MetroFramework.Controls.MetroTextBox();
            this.textboxServerPort = new MetroFramework.Controls.MetroTextBox();
            this.textboxServerAddress = new MetroFramework.Controls.MetroTextBox();
            this.metroLabel1 = new MetroFramework.Controls.MetroLabel();
            this.labelServerPort = new MetroFramework.Controls.MetroLabel();
            this.labelServerAddress = new MetroFramework.Controls.MetroLabel();
            this.buttonConfigurationSave = new MetroFramework.Controls.MetroButton();
            this.buttonConfigurationCancel = new MetroFramework.Controls.MetroButton();
            this.labelRepeat = new MetroFramework.Controls.MetroLabel();
            this.toggleRepeatPlaylist = new MetroFramework.Controls.MetroToggle();
            ((System.ComponentModel.ISupportInitialize)(this.styleManager)).BeginInit();
            this.panelMusicStatus.SuspendLayout();
            this.tabMain.SuspendLayout();
            this.tabPlaying.SuspendLayout();
            this.metroPanel2.SuspendLayout();
            this.tabConfiguration.SuspendLayout();
            this.SuspendLayout();
            // 
            // styleManager
            // 
            this.styleManager.Owner = this;
            this.styleManager.Style = MetroFramework.MetroColorStyle.Red;
            this.styleManager.Theme = MetroFramework.MetroThemeStyle.Dark;
            // 
            // panelMusicStatus
            // 
            this.panelMusicStatus.Controls.Add(this.trackCurrentMusic);
            this.panelMusicStatus.Controls.Add(this.labelNowPlaying);
            this.panelMusicStatus.Dock = System.Windows.Forms.DockStyle.Bottom;
            this.panelMusicStatus.HorizontalScrollbarBarColor = true;
            this.panelMusicStatus.HorizontalScrollbarHighlightOnWheel = false;
            this.panelMusicStatus.HorizontalScrollbarSize = 10;
            this.panelMusicStatus.Location = new System.Drawing.Point(20, 701);
            this.panelMusicStatus.Name = "panelMusicStatus";
            this.panelMusicStatus.Size = new System.Drawing.Size(1239, 57);
            this.panelMusicStatus.TabIndex = 0;
            this.panelMusicStatus.Theme = MetroFramework.MetroThemeStyle.Dark;
            this.panelMusicStatus.VerticalScrollbarBarColor = true;
            this.panelMusicStatus.VerticalScrollbarHighlightOnWheel = false;
            this.panelMusicStatus.VerticalScrollbarSize = 10;
            // 
            // trackCurrentMusic
            // 
            this.trackCurrentMusic.BackColor = System.Drawing.Color.Transparent;
            this.trackCurrentMusic.Dock = System.Windows.Forms.DockStyle.Bottom;
            this.trackCurrentMusic.Location = new System.Drawing.Point(0, 30);
            this.trackCurrentMusic.Name = "trackCurrentMusic";
            this.trackCurrentMusic.Size = new System.Drawing.Size(1239, 27);
            this.trackCurrentMusic.TabIndex = 3;
            this.trackCurrentMusic.Text = "trackMusicPosition";
            this.trackCurrentMusic.Theme = MetroFramework.MetroThemeStyle.Dark;
            // 
            // labelNowPlaying
            // 
            this.labelNowPlaying.Dock = System.Windows.Forms.DockStyle.Top;
            this.labelNowPlaying.Location = new System.Drawing.Point(0, 0);
            this.labelNowPlaying.Name = "labelNowPlaying";
            this.labelNowPlaying.Size = new System.Drawing.Size(1239, 19);
            this.labelNowPlaying.TabIndex = 2;
            this.labelNowPlaying.Text = "Não conectado - Por favor, aguarde ...";
            this.labelNowPlaying.Theme = MetroFramework.MetroThemeStyle.Dark;
            // 
            // timerUpdate
            // 
            this.timerUpdate.Interval = 5000;
            this.timerUpdate.Tick += new System.EventHandler(this.timeUpdate_Tick);
            // 
            // tabMain
            // 
            this.tabMain.Controls.Add(this.tabPlaying);
            this.tabMain.Controls.Add(this.metroTabPage2);
            this.tabMain.Controls.Add(this.tabConfiguration);
            this.tabMain.Dock = System.Windows.Forms.DockStyle.Fill;
            this.tabMain.Location = new System.Drawing.Point(20, 60);
            this.tabMain.Name = "tabMain";
            this.tabMain.SelectedIndex = 0;
            this.tabMain.Size = new System.Drawing.Size(1239, 641);
            this.tabMain.Style = MetroFramework.MetroColorStyle.Red;
            this.tabMain.TabIndex = 1;
            this.tabMain.UseSelectable = true;
            // 
            // tabPlaying
            // 
            this.tabPlaying.Controls.Add(this.metroPanel2);
            this.tabPlaying.HorizontalScrollbarBarColor = true;
            this.tabPlaying.HorizontalScrollbarHighlightOnWheel = false;
            this.tabPlaying.HorizontalScrollbarSize = 10;
            this.tabPlaying.Location = new System.Drawing.Point(4, 38);
            this.tabPlaying.Name = "tabPlaying";
            this.tabPlaying.Size = new System.Drawing.Size(1231, 599);
            this.tabPlaying.TabIndex = 0;
            this.tabPlaying.Text = "Tocando Agora";
            this.tabPlaying.VerticalScrollbarBarColor = true;
            this.tabPlaying.VerticalScrollbarHighlightOnWheel = false;
            this.tabPlaying.VerticalScrollbarSize = 10;
            // 
            // metroPanel2
            // 
            this.metroPanel2.Controls.Add(this.labelRepeat);
            this.metroPanel2.Controls.Add(this.toggleRepeatPlaylist);
            this.metroPanel2.Controls.Add(this.labelShuffle);
            this.metroPanel2.Controls.Add(this.toggleShuffle);
            this.metroPanel2.Controls.Add(this.labelIsPlaying);
            this.metroPanel2.Controls.Add(this.toggleIsPlaying);
            this.metroPanel2.Dock = System.Windows.Forms.DockStyle.Top;
            this.metroPanel2.HorizontalScrollbarBarColor = true;
            this.metroPanel2.HorizontalScrollbarHighlightOnWheel = false;
            this.metroPanel2.HorizontalScrollbarSize = 10;
            this.metroPanel2.Location = new System.Drawing.Point(0, 0);
            this.metroPanel2.Name = "metroPanel2";
            this.metroPanel2.Size = new System.Drawing.Size(1231, 46);
            this.metroPanel2.TabIndex = 2;
            this.metroPanel2.VerticalScrollbarBarColor = true;
            this.metroPanel2.VerticalScrollbarHighlightOnWheel = false;
            this.metroPanel2.VerticalScrollbarSize = 10;
            // 
            // labelShuffle
            // 
            this.labelShuffle.AutoSize = true;
            this.labelShuffle.Location = new System.Drawing.Point(238, 14);
            this.labelShuffle.Name = "labelShuffle";
            this.labelShuffle.Size = new System.Drawing.Size(48, 19);
            this.labelShuffle.TabIndex = 5;
            this.labelShuffle.Text = "Shuffle";
            // 
            // toggleShuffle
            // 
            this.toggleShuffle.AutoSize = true;
            this.toggleShuffle.Location = new System.Drawing.Point(292, 16);
            this.toggleShuffle.Name = "toggleShuffle";
            this.toggleShuffle.Size = new System.Drawing.Size(80, 17);
            this.toggleShuffle.TabIndex = 4;
            this.toggleShuffle.Text = "Off";
            this.toggleShuffle.UseSelectable = true;
            // 
            // labelIsPlaying
            // 
            this.labelIsPlaying.AutoSize = true;
            this.labelIsPlaying.Location = new System.Drawing.Point(7, 14);
            this.labelIsPlaying.Name = "labelIsPlaying";
            this.labelIsPlaying.Size = new System.Drawing.Size(83, 19);
            this.labelIsPlaying.TabIndex = 3;
            this.labelIsPlaying.Text = "Tocar Playlist";
            // 
            // toggleIsPlaying
            // 
            this.toggleIsPlaying.AutoSize = true;
            this.toggleIsPlaying.Location = new System.Drawing.Point(98, 16);
            this.toggleIsPlaying.Name = "toggleIsPlaying";
            this.toggleIsPlaying.Size = new System.Drawing.Size(80, 17);
            this.toggleIsPlaying.TabIndex = 2;
            this.toggleIsPlaying.Text = "Off";
            this.toggleIsPlaying.UseSelectable = true;
            // 
            // metroTabPage2
            // 
            this.metroTabPage2.HorizontalScrollbarBarColor = true;
            this.metroTabPage2.HorizontalScrollbarHighlightOnWheel = false;
            this.metroTabPage2.HorizontalScrollbarSize = 10;
            this.metroTabPage2.Location = new System.Drawing.Point(4, 38);
            this.metroTabPage2.Name = "metroTabPage2";
            this.metroTabPage2.Size = new System.Drawing.Size(1231, 599);
            this.metroTabPage2.TabIndex = 1;
            this.metroTabPage2.Text = "Playlist e Arquivos";
            this.metroTabPage2.VerticalScrollbarBarColor = true;
            this.metroTabPage2.VerticalScrollbarHighlightOnWheel = false;
            this.metroTabPage2.VerticalScrollbarSize = 10;
            // 
            // tabConfiguration
            // 
            this.tabConfiguration.Controls.Add(this.buttonUpdateMusicDatabase);
            this.tabConfiguration.Controls.Add(this.textboxServerAPIKey);
            this.tabConfiguration.Controls.Add(this.textboxServerPort);
            this.tabConfiguration.Controls.Add(this.textboxServerAddress);
            this.tabConfiguration.Controls.Add(this.metroLabel1);
            this.tabConfiguration.Controls.Add(this.labelServerPort);
            this.tabConfiguration.Controls.Add(this.labelServerAddress);
            this.tabConfiguration.Controls.Add(this.buttonConfigurationSave);
            this.tabConfiguration.Controls.Add(this.buttonConfigurationCancel);
            this.tabConfiguration.HorizontalScrollbarBarColor = true;
            this.tabConfiguration.HorizontalScrollbarHighlightOnWheel = false;
            this.tabConfiguration.HorizontalScrollbarSize = 10;
            this.tabConfiguration.Location = new System.Drawing.Point(4, 38);
            this.tabConfiguration.Name = "tabConfiguration";
            this.tabConfiguration.Size = new System.Drawing.Size(1231, 599);
            this.tabConfiguration.TabIndex = 2;
            this.tabConfiguration.Text = "Configurações";
            this.tabConfiguration.VerticalScrollbarBarColor = true;
            this.tabConfiguration.VerticalScrollbarHighlightOnWheel = false;
            this.tabConfiguration.VerticalScrollbarSize = 10;
            // 
            // buttonUpdateMusicDatabase
            // 
            this.buttonUpdateMusicDatabase.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Bottom | System.Windows.Forms.AnchorStyles.Left)));
            this.buttonUpdateMusicDatabase.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(128)))), ((int)(((byte)(128)))), ((int)(((byte)(255)))));
            this.buttonUpdateMusicDatabase.FontSize = MetroFramework.MetroButtonSize.Tall;
            this.buttonUpdateMusicDatabase.ForeColor = System.Drawing.Color.Black;
            this.buttonUpdateMusicDatabase.Location = new System.Drawing.Point(15, 301);
            this.buttonUpdateMusicDatabase.Name = "buttonUpdateMusicDatabase";
            this.buttonUpdateMusicDatabase.Size = new System.Drawing.Size(284, 53);
            this.buttonUpdateMusicDatabase.TabIndex = 10;
            this.buttonUpdateMusicDatabase.Text = "&Atualizar dados do Servidor";
            this.buttonUpdateMusicDatabase.UseCustomBackColor = true;
            this.buttonUpdateMusicDatabase.UseSelectable = true;
            this.buttonUpdateMusicDatabase.Click += new System.EventHandler(this.buttonUpdateMusicDatabase_Click);
            // 
            // textboxServerAPIKey
            // 
            // 
            // 
            // 
            this.textboxServerAPIKey.CustomButton.Image = null;
            this.textboxServerAPIKey.CustomButton.Location = new System.Drawing.Point(486, 1);
            this.textboxServerAPIKey.CustomButton.Name = "";
            this.textboxServerAPIKey.CustomButton.Size = new System.Drawing.Size(21, 21);
            this.textboxServerAPIKey.CustomButton.Style = MetroFramework.MetroColorStyle.Blue;
            this.textboxServerAPIKey.CustomButton.TabIndex = 1;
            this.textboxServerAPIKey.CustomButton.Theme = MetroFramework.MetroThemeStyle.Light;
            this.textboxServerAPIKey.CustomButton.UseSelectable = true;
            this.textboxServerAPIKey.CustomButton.Visible = false;
            this.textboxServerAPIKey.Lines = new string[] {
        "my_api_key"};
            this.textboxServerAPIKey.Location = new System.Drawing.Point(179, 137);
            this.textboxServerAPIKey.MaxLength = 32767;
            this.textboxServerAPIKey.Name = "textboxServerAPIKey";
            this.textboxServerAPIKey.PasswordChar = '\0';
            this.textboxServerAPIKey.ScrollBars = System.Windows.Forms.ScrollBars.None;
            this.textboxServerAPIKey.SelectedText = "";
            this.textboxServerAPIKey.SelectionLength = 0;
            this.textboxServerAPIKey.SelectionStart = 0;
            this.textboxServerAPIKey.ShortcutsEnabled = true;
            this.textboxServerAPIKey.Size = new System.Drawing.Size(508, 23);
            this.textboxServerAPIKey.TabIndex = 9;
            this.textboxServerAPIKey.Text = "my_api_key";
            this.textboxServerAPIKey.UseSelectable = true;
            this.textboxServerAPIKey.WaterMarkColor = System.Drawing.Color.FromArgb(((int)(((byte)(109)))), ((int)(((byte)(109)))), ((int)(((byte)(109)))));
            this.textboxServerAPIKey.WaterMarkFont = new System.Drawing.Font("Segoe UI", 12F, System.Drawing.FontStyle.Italic, System.Drawing.GraphicsUnit.Pixel);
            // 
            // textboxServerPort
            // 
            // 
            // 
            // 
            this.textboxServerPort.CustomButton.Image = null;
            this.textboxServerPort.CustomButton.Location = new System.Drawing.Point(79, 1);
            this.textboxServerPort.CustomButton.Name = "";
            this.textboxServerPort.CustomButton.Size = new System.Drawing.Size(21, 21);
            this.textboxServerPort.CustomButton.Style = MetroFramework.MetroColorStyle.Blue;
            this.textboxServerPort.CustomButton.TabIndex = 1;
            this.textboxServerPort.CustomButton.Theme = MetroFramework.MetroThemeStyle.Light;
            this.textboxServerPort.CustomButton.UseSelectable = true;
            this.textboxServerPort.CustomButton.Visible = false;
            this.textboxServerPort.Lines = new string[] {
        "9320"};
            this.textboxServerPort.Location = new System.Drawing.Point(179, 84);
            this.textboxServerPort.MaxLength = 32767;
            this.textboxServerPort.Name = "textboxServerPort";
            this.textboxServerPort.PasswordChar = '\0';
            this.textboxServerPort.ScrollBars = System.Windows.Forms.ScrollBars.None;
            this.textboxServerPort.SelectedText = "";
            this.textboxServerPort.SelectionLength = 0;
            this.textboxServerPort.SelectionStart = 0;
            this.textboxServerPort.ShortcutsEnabled = true;
            this.textboxServerPort.Size = new System.Drawing.Size(101, 23);
            this.textboxServerPort.TabIndex = 8;
            this.textboxServerPort.Text = "9320";
            this.textboxServerPort.UseSelectable = true;
            this.textboxServerPort.WaterMarkColor = System.Drawing.Color.FromArgb(((int)(((byte)(109)))), ((int)(((byte)(109)))), ((int)(((byte)(109)))));
            this.textboxServerPort.WaterMarkFont = new System.Drawing.Font("Segoe UI", 12F, System.Drawing.FontStyle.Italic, System.Drawing.GraphicsUnit.Pixel);
            // 
            // textboxServerAddress
            // 
            // 
            // 
            // 
            this.textboxServerAddress.CustomButton.Image = null;
            this.textboxServerAddress.CustomButton.Location = new System.Drawing.Point(486, 1);
            this.textboxServerAddress.CustomButton.Name = "";
            this.textboxServerAddress.CustomButton.Size = new System.Drawing.Size(21, 21);
            this.textboxServerAddress.CustomButton.Style = MetroFramework.MetroColorStyle.Blue;
            this.textboxServerAddress.CustomButton.TabIndex = 1;
            this.textboxServerAddress.CustomButton.Theme = MetroFramework.MetroThemeStyle.Light;
            this.textboxServerAddress.CustomButton.UseSelectable = true;
            this.textboxServerAddress.CustomButton.Visible = false;
            this.textboxServerAddress.Lines = new string[] {
        "localhost"};
            this.textboxServerAddress.Location = new System.Drawing.Point(179, 38);
            this.textboxServerAddress.MaxLength = 32767;
            this.textboxServerAddress.Name = "textboxServerAddress";
            this.textboxServerAddress.PasswordChar = '\0';
            this.textboxServerAddress.ScrollBars = System.Windows.Forms.ScrollBars.None;
            this.textboxServerAddress.SelectedText = "";
            this.textboxServerAddress.SelectionLength = 0;
            this.textboxServerAddress.SelectionStart = 0;
            this.textboxServerAddress.ShortcutsEnabled = true;
            this.textboxServerAddress.Size = new System.Drawing.Size(508, 23);
            this.textboxServerAddress.TabIndex = 7;
            this.textboxServerAddress.Text = "localhost";
            this.textboxServerAddress.UseSelectable = true;
            this.textboxServerAddress.WaterMarkColor = System.Drawing.Color.FromArgb(((int)(((byte)(109)))), ((int)(((byte)(109)))), ((int)(((byte)(109)))));
            this.textboxServerAddress.WaterMarkFont = new System.Drawing.Font("Segoe UI", 12F, System.Drawing.FontStyle.Italic, System.Drawing.GraphicsUnit.Pixel);
            // 
            // metroLabel1
            // 
            this.metroLabel1.AutoSize = true;
            this.metroLabel1.Location = new System.Drawing.Point(47, 141);
            this.metroLabel1.Name = "metroLabel1";
            this.metroLabel1.Size = new System.Drawing.Size(91, 19);
            this.metroLabel1.TabIndex = 6;
            this.metroLabel1.Text = "&Chave de API:";
            // 
            // labelServerPort
            // 
            this.labelServerPort.AutoSize = true;
            this.labelServerPort.Location = new System.Drawing.Point(47, 88);
            this.labelServerPort.Name = "labelServerPort";
            this.labelServerPort.Size = new System.Drawing.Size(119, 19);
            this.labelServerPort.TabIndex = 5;
            this.labelServerPort.Text = "&Porta do Servidor:";
            // 
            // labelServerAddress
            // 
            this.labelServerAddress.AutoSize = true;
            this.labelServerAddress.Location = new System.Drawing.Point(24, 38);
            this.labelServerAddress.Name = "labelServerAddress";
            this.labelServerAddress.Size = new System.Drawing.Size(142, 19);
            this.labelServerAddress.TabIndex = 4;
            this.labelServerAddress.Text = "&Endereço do Servidor:";
            // 
            // buttonConfigurationSave
            // 
            this.buttonConfigurationSave.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Bottom | System.Windows.Forms.AnchorStyles.Right)));
            this.buttonConfigurationSave.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(128)))), ((int)(((byte)(255)))), ((int)(((byte)(128)))));
            this.buttonConfigurationSave.FontSize = MetroFramework.MetroButtonSize.Tall;
            this.buttonConfigurationSave.ForeColor = System.Drawing.Color.Black;
            this.buttonConfigurationSave.Location = new System.Drawing.Point(931, 301);
            this.buttonConfigurationSave.Name = "buttonConfigurationSave";
            this.buttonConfigurationSave.Size = new System.Drawing.Size(150, 53);
            this.buttonConfigurationSave.TabIndex = 3;
            this.buttonConfigurationSave.Text = "&Salvar";
            this.buttonConfigurationSave.UseCustomBackColor = true;
            this.buttonConfigurationSave.UseCustomForeColor = true;
            this.buttonConfigurationSave.UseSelectable = true;
            this.buttonConfigurationSave.Click += new System.EventHandler(this.buttonConfigurationSave_Click);
            // 
            // buttonConfigurationCancel
            // 
            this.buttonConfigurationCancel.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Bottom | System.Windows.Forms.AnchorStyles.Right)));
            this.buttonConfigurationCancel.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(255)))), ((int)(((byte)(128)))), ((int)(((byte)(128)))));
            this.buttonConfigurationCancel.FontSize = MetroFramework.MetroButtonSize.Medium;
            this.buttonConfigurationCancel.ForeColor = System.Drawing.Color.Black;
            this.buttonConfigurationCancel.Location = new System.Drawing.Point(758, 301);
            this.buttonConfigurationCancel.Name = "buttonConfigurationCancel";
            this.buttonConfigurationCancel.Size = new System.Drawing.Size(150, 53);
            this.buttonConfigurationCancel.TabIndex = 2;
            this.buttonConfigurationCancel.Text = "&Cancelar";
            this.buttonConfigurationCancel.UseCustomBackColor = true;
            this.buttonConfigurationCancel.UseCustomForeColor = true;
            this.buttonConfigurationCancel.UseSelectable = true;
            this.buttonConfigurationCancel.Click += new System.EventHandler(this.buttonConfigurationCancel_Click);
            // 
            // labelRepeat
            // 
            this.labelRepeat.AutoSize = true;
            this.labelRepeat.Location = new System.Drawing.Point(424, 14);
            this.labelRepeat.Name = "labelRepeat";
            this.labelRepeat.Size = new System.Drawing.Size(166, 19);
            this.labelRepeat.TabIndex = 7;
            this.labelRepeat.Text = "Repetir playlist ao término:";
            // 
            // toggleRepeatPlaylist
            // 
            this.toggleRepeatPlaylist.AutoSize = true;
            this.toggleRepeatPlaylist.Location = new System.Drawing.Point(596, 16);
            this.toggleRepeatPlaylist.Name = "toggleRepeatPlaylist";
            this.toggleRepeatPlaylist.Size = new System.Drawing.Size(80, 17);
            this.toggleRepeatPlaylist.TabIndex = 6;
            this.toggleRepeatPlaylist.Text = "Off";
            this.toggleRepeatPlaylist.UseSelectable = true;
            // 
            // frmMain
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 13F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(1279, 778);
            this.Controls.Add(this.tabMain);
            this.Controls.Add(this.panelMusicStatus);
            this.Name = "frmMain";
            this.Style = MetroFramework.MetroColorStyle.Red;
            this.Text = "BAT Radio by MorcegaoFM";
            this.Theme = MetroFramework.MetroThemeStyle.Dark;
            this.Load += new System.EventHandler(this.frmMain_Load);
            ((System.ComponentModel.ISupportInitialize)(this.styleManager)).EndInit();
            this.panelMusicStatus.ResumeLayout(false);
            this.tabMain.ResumeLayout(false);
            this.tabPlaying.ResumeLayout(false);
            this.metroPanel2.ResumeLayout(false);
            this.metroPanel2.PerformLayout();
            this.tabConfiguration.ResumeLayout(false);
            this.tabConfiguration.PerformLayout();
            this.ResumeLayout(false);

        }

        #endregion

        private MetroFramework.Components.MetroStyleManager styleManager;
        private MetroFramework.Controls.MetroPanel panelMusicStatus;
        private MetroFramework.Controls.MetroTrackBar trackCurrentMusic;
        private MetroFramework.Controls.MetroLabel labelNowPlaying;
        private System.Windows.Forms.Timer timerUpdate;
        private MetroFramework.Controls.MetroTabControl tabMain;
        private MetroFramework.Controls.MetroTabPage tabPlaying;
        private MetroFramework.Controls.MetroTabPage metroTabPage2;
        private MetroFramework.Controls.MetroTabPage tabConfiguration;
        private MetroFramework.Controls.MetroButton buttonUpdateMusicDatabase;
        private MetroFramework.Controls.MetroTextBox textboxServerAPIKey;
        private MetroFramework.Controls.MetroTextBox textboxServerPort;
        private MetroFramework.Controls.MetroTextBox textboxServerAddress;
        private MetroFramework.Controls.MetroLabel metroLabel1;
        private MetroFramework.Controls.MetroLabel labelServerPort;
        private MetroFramework.Controls.MetroLabel labelServerAddress;
        private MetroFramework.Controls.MetroButton buttonConfigurationSave;
        private MetroFramework.Controls.MetroButton buttonConfigurationCancel;
        private MetroFramework.Controls.MetroPanel metroPanel2;
        private MetroFramework.Controls.MetroLabel labelShuffle;
        private MetroFramework.Controls.MetroToggle toggleShuffle;
        private MetroFramework.Controls.MetroLabel labelIsPlaying;
        private MetroFramework.Controls.MetroToggle toggleIsPlaying;
        private MetroFramework.Controls.MetroLabel labelRepeat;
        private MetroFramework.Controls.MetroToggle toggleRepeatPlaylist;
    }
}