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
            System.Windows.Forms.DataGridViewCellStyle dataGridViewCellStyle16 = new System.Windows.Forms.DataGridViewCellStyle();
            System.Windows.Forms.DataGridViewCellStyle dataGridViewCellStyle17 = new System.Windows.Forms.DataGridViewCellStyle();
            System.Windows.Forms.DataGridViewCellStyle dataGridViewCellStyle18 = new System.Windows.Forms.DataGridViewCellStyle();
            System.Windows.Forms.DataGridViewCellStyle dataGridViewCellStyle19 = new System.Windows.Forms.DataGridViewCellStyle();
            System.Windows.Forms.DataGridViewCellStyle dataGridViewCellStyle20 = new System.Windows.Forms.DataGridViewCellStyle();
            System.Windows.Forms.DataGridViewCellStyle dataGridViewCellStyle1 = new System.Windows.Forms.DataGridViewCellStyle();
            System.Windows.Forms.DataGridViewCellStyle dataGridViewCellStyle2 = new System.Windows.Forms.DataGridViewCellStyle();
            System.Windows.Forms.DataGridViewCellStyle dataGridViewCellStyle3 = new System.Windows.Forms.DataGridViewCellStyle();
            System.Windows.Forms.DataGridViewCellStyle dataGridViewCellStyle4 = new System.Windows.Forms.DataGridViewCellStyle();
            System.Windows.Forms.DataGridViewCellStyle dataGridViewCellStyle5 = new System.Windows.Forms.DataGridViewCellStyle();
            System.Windows.Forms.DataGridViewCellStyle dataGridViewCellStyle6 = new System.Windows.Forms.DataGridViewCellStyle();
            System.Windows.Forms.DataGridViewCellStyle dataGridViewCellStyle7 = new System.Windows.Forms.DataGridViewCellStyle();
            System.Windows.Forms.DataGridViewCellStyle dataGridViewCellStyle8 = new System.Windows.Forms.DataGridViewCellStyle();
            System.Windows.Forms.DataGridViewCellStyle dataGridViewCellStyle9 = new System.Windows.Forms.DataGridViewCellStyle();
            System.Windows.Forms.DataGridViewCellStyle dataGridViewCellStyle10 = new System.Windows.Forms.DataGridViewCellStyle();
            this.styleManager = new MetroFramework.Components.MetroStyleManager(this.components);
            this.panelMusicStatus = new MetroFramework.Controls.MetroPanel();
            this.trackCurrentMusic = new MetroFramework.Controls.MetroTrackBar();
            this.labelNowPlaying = new MetroFramework.Controls.MetroLabel();
            this.timerUpdate = new System.Windows.Forms.Timer(this.components);
            this.tabMain = new MetroFramework.Controls.MetroTabControl();
            this.tabPlaying = new MetroFramework.Controls.MetroTabPage();
            this.metroPanel1 = new MetroFramework.Controls.MetroPanel();
            this.buttonCurrentMusic = new MetroFramework.Controls.MetroButton();
            this.textboxSearchPlaylist = new MetroFramework.Controls.MetroTextBox();
            this.labelSearchPlaylist = new MetroFramework.Controls.MetroLabel();
            this.gridPlaylist = new MetroFramework.Controls.MetroGrid();
            this.datagridPlaylistTime = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.datagridPlaylistFile = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.datagridPlaylistArtist = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.datagridPlaylistTitle = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.datagridPlaylistTotalDuration = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.datagridPlaylistAlbum = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.gridPlaylistPosition = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.panSearch = new MetroFramework.Controls.MetroPanel();
            this.metroButton2 = new MetroFramework.Controls.MetroButton();
            this.metroButton1 = new MetroFramework.Controls.MetroButton();
            this.textboxSearch = new MetroFramework.Controls.MetroTextBox();
            this.labelSearch = new MetroFramework.Controls.MetroLabel();
            this.gridSearch = new MetroFramework.Controls.MetroGrid();
            this.gridSearchFile = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.gridSearchArtist = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.gridSearchTitle = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.gridSearchDuration = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.gridSearchAlbum = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.metroTabPage2 = new MetroFramework.Controls.MetroTabPage();
            this.panelPlaylists = new MetroFramework.Controls.MetroPanel();
            this.labelInfo = new MetroFramework.Controls.MetroLabel();
            this.metroGrid1 = new MetroFramework.Controls.MetroGrid();
            this.dataGridViewTextBoxColumn1 = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.dataGridViewTextBoxColumn2 = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.dataGridViewTextBoxColumn3 = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.dataGridViewTextBoxColumn4 = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.dataGridViewTextBoxColumn5 = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.dataGridViewTextBoxColumn6 = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.dataGridViewTextBoxColumn7 = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.metroComboBox1 = new MetroFramework.Controls.MetroComboBox();
            this.labelPlaylists = new MetroFramework.Controls.MetroLabel();
            this.tabConfiguration = new MetroFramework.Controls.MetroTabPage();
            this.metroPanel2 = new MetroFramework.Controls.MetroPanel();
            this.labelMusicFade = new MetroFramework.Controls.MetroLabel();
            this.toggleFadeIn = new MetroFramework.Controls.MetroToggle();
            this.labelRepeat = new MetroFramework.Controls.MetroLabel();
            this.toggleRepeatPlaylist = new MetroFramework.Controls.MetroToggle();
            this.labelShuffle = new MetroFramework.Controls.MetroLabel();
            this.toggleShuffle = new MetroFramework.Controls.MetroToggle();
            this.labelIsPlaying = new MetroFramework.Controls.MetroLabel();
            this.toggleIsPlaying = new MetroFramework.Controls.MetroToggle();
            this.buttonUpdateMusicDatabase = new MetroFramework.Controls.MetroButton();
            this.textboxServerAPIKey = new MetroFramework.Controls.MetroTextBox();
            this.textboxServerPort = new MetroFramework.Controls.MetroTextBox();
            this.textboxServerAddress = new MetroFramework.Controls.MetroTextBox();
            this.metroLabel1 = new MetroFramework.Controls.MetroLabel();
            this.labelServerPort = new MetroFramework.Controls.MetroLabel();
            this.labelServerAddress = new MetroFramework.Controls.MetroLabel();
            this.buttonConfigurationSave = new MetroFramework.Controls.MetroButton();
            this.buttonConfigurationCancel = new MetroFramework.Controls.MetroButton();
            this.contextSearchGrid = new System.Windows.Forms.ContextMenuStrip(this.components);
            this.menuAddAbovePlaylist = new System.Windows.Forms.ToolStripMenuItem();
            this.menuAddBellowPlaylist = new System.Windows.Forms.ToolStripMenuItem();
            this.toolStripSeparator3 = new System.Windows.Forms.ToolStripSeparator();
            this.inserirMúsicaVáriasVezesNaProgramaçãoToolStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.contextMenuPlaylist = new System.Windows.Forms.ContextMenuStrip(this.components);
            this.menuRemoveMusic = new System.Windows.Forms.ToolStripMenuItem();
            this.menuMoveAbove = new System.Windows.Forms.ToolStripMenuItem();
            this.menuMoveBellow = new System.Windows.Forms.ToolStripMenuItem();
            this.menuSwap = new System.Windows.Forms.ToolStripMenuItem();
            this.mnuMerge = new System.Windows.Forms.ToolStripMenuItem();
            this.toolStripSeparator1 = new System.Windows.Forms.ToolStripSeparator();
            this.menuGoToCurrent = new System.Windows.Forms.ToolStripMenuItem();
            this.toolStripSeparator2 = new System.Windows.Forms.ToolStripSeparator();
            this.menuPlaySelectedMusic = new System.Windows.Forms.ToolStripMenuItem();
            ((System.ComponentModel.ISupportInitialize)(this.styleManager)).BeginInit();
            this.panelMusicStatus.SuspendLayout();
            this.tabMain.SuspendLayout();
            this.tabPlaying.SuspendLayout();
            this.metroPanel1.SuspendLayout();
            ((System.ComponentModel.ISupportInitialize)(this.gridPlaylist)).BeginInit();
            this.panSearch.SuspendLayout();
            ((System.ComponentModel.ISupportInitialize)(this.gridSearch)).BeginInit();
            this.metroTabPage2.SuspendLayout();
            this.panelPlaylists.SuspendLayout();
            ((System.ComponentModel.ISupportInitialize)(this.metroGrid1)).BeginInit();
            this.tabConfiguration.SuspendLayout();
            this.metroPanel2.SuspendLayout();
            this.contextSearchGrid.SuspendLayout();
            this.contextMenuPlaylist.SuspendLayout();
            this.SuspendLayout();
            // 
            // styleManager
            // 
            this.styleManager.Owner = this;
            this.styleManager.Style = MetroFramework.MetroColorStyle.Red;
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
            this.panelMusicStatus.Size = new System.Drawing.Size(1285, 57);
            this.panelMusicStatus.TabIndex = 0;
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
            this.trackCurrentMusic.Size = new System.Drawing.Size(1285, 27);
            this.trackCurrentMusic.TabIndex = 3;
            this.trackCurrentMusic.Text = "trackMusicPosition";
            this.trackCurrentMusic.Scroll += new System.Windows.Forms.ScrollEventHandler(this.trackCurrentMusic_Scroll);
            // 
            // labelNowPlaying
            // 
            this.labelNowPlaying.Dock = System.Windows.Forms.DockStyle.Top;
            this.labelNowPlaying.Location = new System.Drawing.Point(0, 0);
            this.labelNowPlaying.Name = "labelNowPlaying";
            this.labelNowPlaying.Size = new System.Drawing.Size(1285, 19);
            this.labelNowPlaying.TabIndex = 2;
            this.labelNowPlaying.Text = "Não conectado - Por favor, aguarde ...";
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
            this.tabMain.SelectedIndex = 1;
            this.tabMain.Size = new System.Drawing.Size(1285, 641);
            this.tabMain.Style = MetroFramework.MetroColorStyle.Red;
            this.tabMain.TabIndex = 1;
            this.tabMain.UseSelectable = true;
            // 
            // tabPlaying
            // 
            this.tabPlaying.Controls.Add(this.metroPanel1);
            this.tabPlaying.Controls.Add(this.panSearch);
            this.tabPlaying.HorizontalScrollbarBarColor = true;
            this.tabPlaying.HorizontalScrollbarHighlightOnWheel = false;
            this.tabPlaying.HorizontalScrollbarSize = 10;
            this.tabPlaying.Location = new System.Drawing.Point(4, 38);
            this.tabPlaying.Name = "tabPlaying";
            this.tabPlaying.Size = new System.Drawing.Size(1277, 599);
            this.tabPlaying.TabIndex = 0;
            this.tabPlaying.Text = "Tocando Agora";
            this.tabPlaying.VerticalScrollbarBarColor = true;
            this.tabPlaying.VerticalScrollbarHighlightOnWheel = false;
            this.tabPlaying.VerticalScrollbarSize = 10;
            // 
            // metroPanel1
            // 
            this.metroPanel1.Controls.Add(this.buttonCurrentMusic);
            this.metroPanel1.Controls.Add(this.textboxSearchPlaylist);
            this.metroPanel1.Controls.Add(this.labelSearchPlaylist);
            this.metroPanel1.Controls.Add(this.gridPlaylist);
            this.metroPanel1.Dock = System.Windows.Forms.DockStyle.Fill;
            this.metroPanel1.HorizontalScrollbarBarColor = true;
            this.metroPanel1.HorizontalScrollbarHighlightOnWheel = false;
            this.metroPanel1.HorizontalScrollbarSize = 10;
            this.metroPanel1.Location = new System.Drawing.Point(0, 304);
            this.metroPanel1.Name = "metroPanel1";
            this.metroPanel1.Size = new System.Drawing.Size(1277, 295);
            this.metroPanel1.Style = MetroFramework.MetroColorStyle.Purple;
            this.metroPanel1.TabIndex = 4;
            this.metroPanel1.VerticalScrollbarBarColor = true;
            this.metroPanel1.VerticalScrollbarHighlightOnWheel = false;
            this.metroPanel1.VerticalScrollbarSize = 10;
            // 
            // buttonCurrentMusic
            // 
            this.buttonCurrentMusic.Location = new System.Drawing.Point(555, 13);
            this.buttonCurrentMusic.Name = "buttonCurrentMusic";
            this.buttonCurrentMusic.Size = new System.Drawing.Size(146, 19);
            this.buttonCurrentMusic.TabIndex = 8;
            this.buttonCurrentMusic.Text = "Música Atual";
            this.buttonCurrentMusic.UseSelectable = true;
            this.buttonCurrentMusic.Click += new System.EventHandler(this.buttonCurrentMusic_Click_1);
            // 
            // textboxSearchPlaylist
            // 
            // 
            // 
            // 
            this.textboxSearchPlaylist.CustomButton.Image = null;
            this.textboxSearchPlaylist.CustomButton.Location = new System.Drawing.Point(334, 2);
            this.textboxSearchPlaylist.CustomButton.Name = "";
            this.textboxSearchPlaylist.CustomButton.Size = new System.Drawing.Size(15, 15);
            this.textboxSearchPlaylist.CustomButton.Style = MetroFramework.MetroColorStyle.Blue;
            this.textboxSearchPlaylist.CustomButton.TabIndex = 1;
            this.textboxSearchPlaylist.CustomButton.Theme = MetroFramework.MetroThemeStyle.Light;
            this.textboxSearchPlaylist.CustomButton.UseSelectable = true;
            this.textboxSearchPlaylist.CustomButton.Visible = false;
            this.textboxSearchPlaylist.Lines = new string[0];
            this.textboxSearchPlaylist.Location = new System.Drawing.Point(167, 13);
            this.textboxSearchPlaylist.MaxLength = 32767;
            this.textboxSearchPlaylist.Name = "textboxSearchPlaylist";
            this.textboxSearchPlaylist.PasswordChar = '\0';
            this.textboxSearchPlaylist.PromptText = "Digite a música que deseja pesquisar";
            this.textboxSearchPlaylist.ScrollBars = System.Windows.Forms.ScrollBars.None;
            this.textboxSearchPlaylist.SelectedText = "";
            this.textboxSearchPlaylist.SelectionLength = 0;
            this.textboxSearchPlaylist.SelectionStart = 0;
            this.textboxSearchPlaylist.ShortcutsEnabled = true;
            this.textboxSearchPlaylist.Size = new System.Drawing.Size(352, 20);
            this.textboxSearchPlaylist.TabIndex = 6;
            this.textboxSearchPlaylist.UseSelectable = true;
            this.textboxSearchPlaylist.WaterMark = "Digite a música que deseja pesquisar";
            this.textboxSearchPlaylist.WaterMarkColor = System.Drawing.Color.FromArgb(((int)(((byte)(109)))), ((int)(((byte)(109)))), ((int)(((byte)(109)))));
            this.textboxSearchPlaylist.WaterMarkFont = new System.Drawing.Font("Segoe UI", 12F, System.Drawing.FontStyle.Italic, System.Drawing.GraphicsUnit.Pixel);
            this.textboxSearchPlaylist.KeyPress += new System.Windows.Forms.KeyPressEventHandler(this.textboxSearchPlaylist_KeyPress);
            // 
            // labelSearchPlaylist
            // 
            this.labelSearchPlaylist.AutoSize = true;
            this.labelSearchPlaylist.FontWeight = MetroFramework.MetroLabelWeight.Bold;
            this.labelSearchPlaylist.Location = new System.Drawing.Point(12, 13);
            this.labelSearchPlaylist.Name = "labelSearchPlaylist";
            this.labelSearchPlaylist.Size = new System.Drawing.Size(149, 19);
            this.labelSearchPlaylist.TabIndex = 5;
            this.labelSearchPlaylist.Text = "Pesquisar na Playlist:";
            // 
            // gridPlaylist
            // 
            this.gridPlaylist.AllowUserToAddRows = false;
            this.gridPlaylist.AllowUserToDeleteRows = false;
            this.gridPlaylist.AllowUserToOrderColumns = true;
            this.gridPlaylist.AllowUserToResizeRows = false;
            dataGridViewCellStyle16.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(255)))), ((int)(((byte)(192)))), ((int)(((byte)(192)))));
            dataGridViewCellStyle16.Font = new System.Drawing.Font("Microsoft Sans Serif", 11.25F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            dataGridViewCellStyle16.ForeColor = System.Drawing.Color.Black;
            dataGridViewCellStyle16.SelectionBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(192)))), ((int)(((byte)(0)))), ((int)(((byte)(0)))));
            dataGridViewCellStyle16.SelectionForeColor = System.Drawing.Color.White;
            this.gridPlaylist.AlternatingRowsDefaultCellStyle = dataGridViewCellStyle16;
            this.gridPlaylist.Anchor = ((System.Windows.Forms.AnchorStyles)((((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.gridPlaylist.AutoSizeColumnsMode = System.Windows.Forms.DataGridViewAutoSizeColumnsMode.Fill;
            this.gridPlaylist.BackgroundColor = System.Drawing.Color.FromArgb(((int)(((byte)(255)))), ((int)(((byte)(255)))), ((int)(((byte)(255)))));
            this.gridPlaylist.BorderStyle = System.Windows.Forms.BorderStyle.None;
            this.gridPlaylist.CellBorderStyle = System.Windows.Forms.DataGridViewCellBorderStyle.None;
            this.gridPlaylist.ColumnHeadersBorderStyle = System.Windows.Forms.DataGridViewHeaderBorderStyle.None;
            dataGridViewCellStyle17.Alignment = System.Windows.Forms.DataGridViewContentAlignment.MiddleLeft;
            dataGridViewCellStyle17.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(209)))), ((int)(((byte)(17)))), ((int)(((byte)(65)))));
            dataGridViewCellStyle17.Font = new System.Drawing.Font("Segoe UI", 11F, System.Drawing.FontStyle.Bold, System.Drawing.GraphicsUnit.Pixel);
            dataGridViewCellStyle17.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(255)))), ((int)(((byte)(255)))), ((int)(((byte)(255)))));
            dataGridViewCellStyle17.SelectionBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(234)))), ((int)(((byte)(19)))), ((int)(((byte)(73)))));
            dataGridViewCellStyle17.SelectionForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(17)))), ((int)(((byte)(17)))), ((int)(((byte)(17)))));
            dataGridViewCellStyle17.WrapMode = System.Windows.Forms.DataGridViewTriState.True;
            this.gridPlaylist.ColumnHeadersDefaultCellStyle = dataGridViewCellStyle17;
            this.gridPlaylist.ColumnHeadersHeightSizeMode = System.Windows.Forms.DataGridViewColumnHeadersHeightSizeMode.AutoSize;
            this.gridPlaylist.Columns.AddRange(new System.Windows.Forms.DataGridViewColumn[] {
            this.datagridPlaylistTime,
            this.datagridPlaylistFile,
            this.datagridPlaylistArtist,
            this.datagridPlaylistTitle,
            this.datagridPlaylistTotalDuration,
            this.datagridPlaylistAlbum,
            this.gridPlaylistPosition});
            this.gridPlaylist.ContextMenuStrip = this.contextMenuPlaylist;
            dataGridViewCellStyle18.Alignment = System.Windows.Forms.DataGridViewContentAlignment.MiddleLeft;
            dataGridViewCellStyle18.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(255)))), ((int)(((byte)(255)))), ((int)(((byte)(255)))));
            dataGridViewCellStyle18.Font = new System.Drawing.Font("Segoe UI", 11F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Pixel);
            dataGridViewCellStyle18.ForeColor = System.Drawing.Color.Black;
            dataGridViewCellStyle18.SelectionBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(234)))), ((int)(((byte)(19)))), ((int)(((byte)(73)))));
            dataGridViewCellStyle18.SelectionForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(17)))), ((int)(((byte)(17)))), ((int)(((byte)(17)))));
            dataGridViewCellStyle18.WrapMode = System.Windows.Forms.DataGridViewTriState.False;
            this.gridPlaylist.DefaultCellStyle = dataGridViewCellStyle18;
            this.gridPlaylist.EnableHeadersVisualStyles = false;
            this.gridPlaylist.Font = new System.Drawing.Font("Segoe UI", 11F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Pixel);
            this.gridPlaylist.GridColor = System.Drawing.Color.FromArgb(((int)(((byte)(255)))), ((int)(((byte)(255)))), ((int)(((byte)(255)))));
            this.gridPlaylist.Location = new System.Drawing.Point(0, 44);
            this.gridPlaylist.Name = "gridPlaylist";
            this.gridPlaylist.ReadOnly = true;
            this.gridPlaylist.RowHeadersBorderStyle = System.Windows.Forms.DataGridViewHeaderBorderStyle.None;
            dataGridViewCellStyle19.Alignment = System.Windows.Forms.DataGridViewContentAlignment.MiddleLeft;
            dataGridViewCellStyle19.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(209)))), ((int)(((byte)(17)))), ((int)(((byte)(65)))));
            dataGridViewCellStyle19.Font = new System.Drawing.Font("Segoe UI", 11F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Pixel);
            dataGridViewCellStyle19.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(255)))), ((int)(((byte)(255)))), ((int)(((byte)(255)))));
            dataGridViewCellStyle19.SelectionBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(234)))), ((int)(((byte)(19)))), ((int)(((byte)(73)))));
            dataGridViewCellStyle19.SelectionForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(17)))), ((int)(((byte)(17)))), ((int)(((byte)(17)))));
            dataGridViewCellStyle19.WrapMode = System.Windows.Forms.DataGridViewTriState.True;
            this.gridPlaylist.RowHeadersDefaultCellStyle = dataGridViewCellStyle19;
            this.gridPlaylist.RowHeadersWidthSizeMode = System.Windows.Forms.DataGridViewRowHeadersWidthSizeMode.DisableResizing;
            dataGridViewCellStyle20.BackColor = System.Drawing.Color.White;
            dataGridViewCellStyle20.Font = new System.Drawing.Font("Microsoft Sans Serif", 11.25F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            dataGridViewCellStyle20.ForeColor = System.Drawing.Color.Black;
            dataGridViewCellStyle20.SelectionBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(192)))), ((int)(((byte)(0)))), ((int)(((byte)(0)))));
            dataGridViewCellStyle20.SelectionForeColor = System.Drawing.Color.White;
            this.gridPlaylist.RowsDefaultCellStyle = dataGridViewCellStyle20;
            this.gridPlaylist.SelectionMode = System.Windows.Forms.DataGridViewSelectionMode.FullRowSelect;
            this.gridPlaylist.ShowCellToolTips = false;
            this.gridPlaylist.ShowEditingIcon = false;
            this.gridPlaylist.Size = new System.Drawing.Size(1277, 205);
            this.gridPlaylist.Style = MetroFramework.MetroColorStyle.Red;
            this.gridPlaylist.TabIndex = 3;
            this.gridPlaylist.UseCustomBackColor = true;
            this.gridPlaylist.UseCustomForeColor = true;
            this.gridPlaylist.MouseDown += new System.Windows.Forms.MouseEventHandler(this.gridPlaylist_MouseClick);
            // 
            // datagridPlaylistTime
            // 
            this.datagridPlaylistTime.DataPropertyName = "NextPresentation";
            this.datagridPlaylistTime.HeaderText = "Horário prox. Reprodução";
            this.datagridPlaylistTime.Name = "datagridPlaylistTime";
            this.datagridPlaylistTime.ReadOnly = true;
            // 
            // datagridPlaylistFile
            // 
            this.datagridPlaylistFile.DataPropertyName = "file";
            this.datagridPlaylistFile.HeaderText = "Arquivo";
            this.datagridPlaylistFile.Name = "datagridPlaylistFile";
            this.datagridPlaylistFile.ReadOnly = true;
            // 
            // datagridPlaylistArtist
            // 
            this.datagridPlaylistArtist.DataPropertyName = "artist";
            this.datagridPlaylistArtist.HeaderText = "Artista";
            this.datagridPlaylistArtist.Name = "datagridPlaylistArtist";
            this.datagridPlaylistArtist.ReadOnly = true;
            // 
            // datagridPlaylistTitle
            // 
            this.datagridPlaylistTitle.DataPropertyName = "title";
            this.datagridPlaylistTitle.HeaderText = "Titúlo";
            this.datagridPlaylistTitle.Name = "datagridPlaylistTitle";
            this.datagridPlaylistTitle.ReadOnly = true;
            // 
            // datagridPlaylistTotalDuration
            // 
            this.datagridPlaylistTotalDuration.DataPropertyName = "TotalDuration";
            this.datagridPlaylistTotalDuration.HeaderText = "Duração";
            this.datagridPlaylistTotalDuration.Name = "datagridPlaylistTotalDuration";
            this.datagridPlaylistTotalDuration.ReadOnly = true;
            // 
            // datagridPlaylistAlbum
            // 
            this.datagridPlaylistAlbum.DataPropertyName = "album";
            this.datagridPlaylistAlbum.HeaderText = "Album";
            this.datagridPlaylistAlbum.Name = "datagridPlaylistAlbum";
            this.datagridPlaylistAlbum.ReadOnly = true;
            // 
            // gridPlaylistPosition
            // 
            this.gridPlaylistPosition.DataPropertyName = "Pos";
            this.gridPlaylistPosition.HeaderText = "Posição";
            this.gridPlaylistPosition.Name = "gridPlaylistPosition";
            this.gridPlaylistPosition.ReadOnly = true;
            // 
            // panSearch
            // 
            this.panSearch.BorderStyle = System.Windows.Forms.BorderStyle.FixedSingle;
            this.panSearch.Controls.Add(this.metroButton2);
            this.panSearch.Controls.Add(this.metroButton1);
            this.panSearch.Controls.Add(this.textboxSearch);
            this.panSearch.Controls.Add(this.labelSearch);
            this.panSearch.Controls.Add(this.gridSearch);
            this.panSearch.Dock = System.Windows.Forms.DockStyle.Top;
            this.panSearch.HorizontalScrollbarBarColor = true;
            this.panSearch.HorizontalScrollbarHighlightOnWheel = false;
            this.panSearch.HorizontalScrollbarSize = 10;
            this.panSearch.Location = new System.Drawing.Point(0, 0);
            this.panSearch.Name = "panSearch";
            this.panSearch.Size = new System.Drawing.Size(1277, 304);
            this.panSearch.Style = MetroFramework.MetroColorStyle.Silver;
            this.panSearch.TabIndex = 3;
            this.panSearch.VerticalScrollbarBarColor = true;
            this.panSearch.VerticalScrollbarHighlightOnWheel = false;
            this.panSearch.VerticalScrollbarSize = 10;
            // 
            // metroButton2
            // 
            this.metroButton2.Location = new System.Drawing.Point(913, 16);
            this.metroButton2.Name = "metroButton2";
            this.metroButton2.Size = new System.Drawing.Size(244, 19);
            this.metroButton2.TabIndex = 6;
            this.metroButton2.Text = "Adicionar acima da seleção";
            this.metroButton2.UseSelectable = true;
            this.metroButton2.Click += new System.EventHandler(this.buttonAddAbove_Click);
            // 
            // metroButton1
            // 
            this.metroButton1.Location = new System.Drawing.Point(646, 16);
            this.metroButton1.Name = "metroButton1";
            this.metroButton1.Size = new System.Drawing.Size(244, 19);
            this.metroButton1.TabIndex = 5;
            this.metroButton1.Text = "Adicionar abaixo da seleção";
            this.metroButton1.UseSelectable = true;
            this.metroButton1.Click += new System.EventHandler(this.buttonAddBellow_Click);
            // 
            // textboxSearch
            // 
            // 
            // 
            // 
            this.textboxSearch.CustomButton.Image = null;
            this.textboxSearch.CustomButton.Location = new System.Drawing.Point(456, 2);
            this.textboxSearch.CustomButton.Name = "";
            this.textboxSearch.CustomButton.Size = new System.Drawing.Size(15, 15);
            this.textboxSearch.CustomButton.Style = MetroFramework.MetroColorStyle.Blue;
            this.textboxSearch.CustomButton.TabIndex = 1;
            this.textboxSearch.CustomButton.Theme = MetroFramework.MetroThemeStyle.Light;
            this.textboxSearch.CustomButton.UseSelectable = true;
            this.textboxSearch.CustomButton.Visible = false;
            this.textboxSearch.Lines = new string[0];
            this.textboxSearch.Location = new System.Drawing.Point(142, 16);
            this.textboxSearch.MaxLength = 32767;
            this.textboxSearch.Name = "textboxSearch";
            this.textboxSearch.PasswordChar = '\0';
            this.textboxSearch.PromptText = "Digite a música que deseja pesquisar";
            this.textboxSearch.ScrollBars = System.Windows.Forms.ScrollBars.None;
            this.textboxSearch.SelectedText = "";
            this.textboxSearch.SelectionLength = 0;
            this.textboxSearch.SelectionStart = 0;
            this.textboxSearch.ShortcutsEnabled = true;
            this.textboxSearch.Size = new System.Drawing.Size(474, 20);
            this.textboxSearch.TabIndex = 4;
            this.textboxSearch.UseSelectable = true;
            this.textboxSearch.WaterMark = "Digite a música que deseja pesquisar";
            this.textboxSearch.WaterMarkColor = System.Drawing.Color.FromArgb(((int)(((byte)(109)))), ((int)(((byte)(109)))), ((int)(((byte)(109)))));
            this.textboxSearch.WaterMarkFont = new System.Drawing.Font("Segoe UI", 12F, System.Drawing.FontStyle.Italic, System.Drawing.GraphicsUnit.Pixel);
            this.textboxSearch.TextChanged += new System.EventHandler(this.TextboxSearch_TextChanged);
            // 
            // labelSearch
            // 
            this.labelSearch.AutoSize = true;
            this.labelSearch.FontWeight = MetroFramework.MetroLabelWeight.Bold;
            this.labelSearch.Location = new System.Drawing.Point(27, 16);
            this.labelSearch.Name = "labelSearch";
            this.labelSearch.Size = new System.Drawing.Size(78, 19);
            this.labelSearch.TabIndex = 3;
            this.labelSearch.Text = "Pesquisar:";
            // 
            // gridSearch
            // 
            this.gridSearch.AllowUserToAddRows = false;
            this.gridSearch.AllowUserToDeleteRows = false;
            this.gridSearch.AllowUserToOrderColumns = true;
            this.gridSearch.AllowUserToResizeRows = false;
            dataGridViewCellStyle1.BackColor = System.Drawing.Color.White;
            dataGridViewCellStyle1.ForeColor = System.Drawing.Color.Black;
            dataGridViewCellStyle1.SelectionBackColor = System.Drawing.Color.Red;
            this.gridSearch.AlternatingRowsDefaultCellStyle = dataGridViewCellStyle1;
            this.gridSearch.Anchor = ((System.Windows.Forms.AnchorStyles)((((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.gridSearch.AutoSizeColumnsMode = System.Windows.Forms.DataGridViewAutoSizeColumnsMode.Fill;
            this.gridSearch.BackgroundColor = System.Drawing.Color.FromArgb(((int)(((byte)(255)))), ((int)(((byte)(255)))), ((int)(((byte)(255)))));
            this.gridSearch.BorderStyle = System.Windows.Forms.BorderStyle.None;
            this.gridSearch.CellBorderStyle = System.Windows.Forms.DataGridViewCellBorderStyle.None;
            this.gridSearch.ColumnHeadersBorderStyle = System.Windows.Forms.DataGridViewHeaderBorderStyle.None;
            dataGridViewCellStyle2.Alignment = System.Windows.Forms.DataGridViewContentAlignment.MiddleLeft;
            dataGridViewCellStyle2.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(0)))), ((int)(((byte)(0)))), ((int)(((byte)(0)))));
            dataGridViewCellStyle2.Font = new System.Drawing.Font("Segoe UI", 11F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Pixel);
            dataGridViewCellStyle2.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(255)))), ((int)(((byte)(255)))), ((int)(((byte)(255)))));
            dataGridViewCellStyle2.SelectionBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(25)))), ((int)(((byte)(25)))), ((int)(((byte)(25)))));
            dataGridViewCellStyle2.SelectionForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(17)))), ((int)(((byte)(17)))), ((int)(((byte)(17)))));
            dataGridViewCellStyle2.WrapMode = System.Windows.Forms.DataGridViewTriState.True;
            this.gridSearch.ColumnHeadersDefaultCellStyle = dataGridViewCellStyle2;
            this.gridSearch.ColumnHeadersHeightSizeMode = System.Windows.Forms.DataGridViewColumnHeadersHeightSizeMode.AutoSize;
            this.gridSearch.Columns.AddRange(new System.Windows.Forms.DataGridViewColumn[] {
            this.gridSearchFile,
            this.gridSearchArtist,
            this.gridSearchTitle,
            this.gridSearchDuration,
            this.gridSearchAlbum});
            this.gridSearch.ContextMenuStrip = this.contextSearchGrid;
            dataGridViewCellStyle3.Alignment = System.Windows.Forms.DataGridViewContentAlignment.MiddleLeft;
            dataGridViewCellStyle3.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(255)))), ((int)(((byte)(255)))), ((int)(((byte)(255)))));
            dataGridViewCellStyle3.Font = new System.Drawing.Font("Segoe UI", 11F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Pixel);
            dataGridViewCellStyle3.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(136)))), ((int)(((byte)(136)))), ((int)(((byte)(136)))));
            dataGridViewCellStyle3.SelectionBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(25)))), ((int)(((byte)(25)))), ((int)(((byte)(25)))));
            dataGridViewCellStyle3.SelectionForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(17)))), ((int)(((byte)(17)))), ((int)(((byte)(17)))));
            dataGridViewCellStyle3.WrapMode = System.Windows.Forms.DataGridViewTriState.False;
            this.gridSearch.DefaultCellStyle = dataGridViewCellStyle3;
            this.gridSearch.EnableHeadersVisualStyles = false;
            this.gridSearch.Font = new System.Drawing.Font("Segoe UI", 11F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Pixel);
            this.gridSearch.GridColor = System.Drawing.Color.FromArgb(((int)(((byte)(255)))), ((int)(((byte)(255)))), ((int)(((byte)(255)))));
            this.gridSearch.Location = new System.Drawing.Point(0, 56);
            this.gridSearch.Name = "gridSearch";
            this.gridSearch.ReadOnly = true;
            this.gridSearch.RowHeadersBorderStyle = System.Windows.Forms.DataGridViewHeaderBorderStyle.None;
            dataGridViewCellStyle4.Alignment = System.Windows.Forms.DataGridViewContentAlignment.MiddleLeft;
            dataGridViewCellStyle4.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(0)))), ((int)(((byte)(0)))), ((int)(((byte)(0)))));
            dataGridViewCellStyle4.Font = new System.Drawing.Font("Segoe UI", 11F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Pixel);
            dataGridViewCellStyle4.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(255)))), ((int)(((byte)(255)))), ((int)(((byte)(255)))));
            dataGridViewCellStyle4.SelectionBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(25)))), ((int)(((byte)(25)))), ((int)(((byte)(25)))));
            dataGridViewCellStyle4.SelectionForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(17)))), ((int)(((byte)(17)))), ((int)(((byte)(17)))));
            dataGridViewCellStyle4.WrapMode = System.Windows.Forms.DataGridViewTriState.True;
            this.gridSearch.RowHeadersDefaultCellStyle = dataGridViewCellStyle4;
            this.gridSearch.RowHeadersWidthSizeMode = System.Windows.Forms.DataGridViewRowHeadersWidthSizeMode.DisableResizing;
            dataGridViewCellStyle5.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(224)))), ((int)(((byte)(224)))), ((int)(((byte)(224)))));
            dataGridViewCellStyle5.Font = new System.Drawing.Font("Microsoft Sans Serif", 11.25F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            dataGridViewCellStyle5.ForeColor = System.Drawing.Color.Black;
            dataGridViewCellStyle5.SelectionBackColor = System.Drawing.Color.Red;
            dataGridViewCellStyle5.SelectionForeColor = System.Drawing.Color.White;
            this.gridSearch.RowsDefaultCellStyle = dataGridViewCellStyle5;
            this.gridSearch.SelectionMode = System.Windows.Forms.DataGridViewSelectionMode.FullRowSelect;
            this.gridSearch.Size = new System.Drawing.Size(1275, 246);
            this.gridSearch.Style = MetroFramework.MetroColorStyle.Black;
            this.gridSearch.TabIndex = 2;
            this.gridSearch.MouseDown += new System.Windows.Forms.MouseEventHandler(this.gridSearch_MouseClick);
            // 
            // gridSearchFile
            // 
            this.gridSearchFile.DataPropertyName = "file";
            this.gridSearchFile.HeaderText = "Arquivo";
            this.gridSearchFile.Name = "gridSearchFile";
            this.gridSearchFile.ReadOnly = true;
            // 
            // gridSearchArtist
            // 
            this.gridSearchArtist.DataPropertyName = "artist";
            this.gridSearchArtist.HeaderText = "Artista";
            this.gridSearchArtist.Name = "gridSearchArtist";
            this.gridSearchArtist.ReadOnly = true;
            // 
            // gridSearchTitle
            // 
            this.gridSearchTitle.DataPropertyName = "title";
            this.gridSearchTitle.HeaderText = "Titúlo";
            this.gridSearchTitle.Name = "gridSearchTitle";
            this.gridSearchTitle.ReadOnly = true;
            // 
            // gridSearchDuration
            // 
            this.gridSearchDuration.DataPropertyName = "TotalDuration";
            this.gridSearchDuration.HeaderText = "Duração";
            this.gridSearchDuration.Name = "gridSearchDuration";
            this.gridSearchDuration.ReadOnly = true;
            // 
            // gridSearchAlbum
            // 
            this.gridSearchAlbum.DataPropertyName = "album";
            this.gridSearchAlbum.HeaderText = "Album";
            this.gridSearchAlbum.Name = "gridSearchAlbum";
            this.gridSearchAlbum.ReadOnly = true;
            // 
            // metroTabPage2
            // 
            this.metroTabPage2.Controls.Add(this.panelPlaylists);
            this.metroTabPage2.HorizontalScrollbarBarColor = true;
            this.metroTabPage2.HorizontalScrollbarHighlightOnWheel = false;
            this.metroTabPage2.HorizontalScrollbarSize = 10;
            this.metroTabPage2.Location = new System.Drawing.Point(4, 38);
            this.metroTabPage2.Name = "metroTabPage2";
            this.metroTabPage2.Size = new System.Drawing.Size(1277, 599);
            this.metroTabPage2.TabIndex = 1;
            this.metroTabPage2.Text = "Playlist e Arquivos";
            this.metroTabPage2.VerticalScrollbarBarColor = true;
            this.metroTabPage2.VerticalScrollbarHighlightOnWheel = false;
            this.metroTabPage2.VerticalScrollbarSize = 10;
            // 
            // panelPlaylists
            // 
            this.panelPlaylists.Controls.Add(this.labelInfo);
            this.panelPlaylists.Controls.Add(this.metroGrid1);
            this.panelPlaylists.Controls.Add(this.metroComboBox1);
            this.panelPlaylists.Controls.Add(this.labelPlaylists);
            this.panelPlaylists.Dock = System.Windows.Forms.DockStyle.Top;
            this.panelPlaylists.HorizontalScrollbarBarColor = true;
            this.panelPlaylists.HorizontalScrollbarHighlightOnWheel = false;
            this.panelPlaylists.HorizontalScrollbarSize = 10;
            this.panelPlaylists.Location = new System.Drawing.Point(0, 0);
            this.panelPlaylists.Name = "panelPlaylists";
            this.panelPlaylists.Size = new System.Drawing.Size(1277, 530);
            this.panelPlaylists.TabIndex = 2;
            this.panelPlaylists.VerticalScrollbarBarColor = true;
            this.panelPlaylists.VerticalScrollbarHighlightOnWheel = false;
            this.panelPlaylists.VerticalScrollbarSize = 10;
            this.panelPlaylists.Paint += new System.Windows.Forms.PaintEventHandler(this.metroPanel3_Paint);
            // 
            // labelInfo
            // 
            this.labelInfo.AutoSize = true;
            this.labelInfo.Location = new System.Drawing.Point(25, 266);
            this.labelInfo.Name = "labelInfo";
            this.labelInfo.Size = new System.Drawing.Size(775, 19);
            this.labelInfo.TabIndex = 5;
            this.labelInfo.Text = "O Horário da Próxima reprodução é exibido de acordo com o horário atual (simuland" +
    "o se colocasse a programação agora no ar)";
            // 
            // metroGrid1
            // 
            this.metroGrid1.AllowUserToAddRows = false;
            this.metroGrid1.AllowUserToDeleteRows = false;
            this.metroGrid1.AllowUserToOrderColumns = true;
            this.metroGrid1.AllowUserToResizeRows = false;
            dataGridViewCellStyle6.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(255)))), ((int)(((byte)(192)))), ((int)(((byte)(192)))));
            dataGridViewCellStyle6.Font = new System.Drawing.Font("Microsoft Sans Serif", 11.25F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            dataGridViewCellStyle6.ForeColor = System.Drawing.Color.Black;
            dataGridViewCellStyle6.SelectionBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(192)))), ((int)(((byte)(0)))), ((int)(((byte)(0)))));
            dataGridViewCellStyle6.SelectionForeColor = System.Drawing.Color.White;
            this.metroGrid1.AlternatingRowsDefaultCellStyle = dataGridViewCellStyle6;
            this.metroGrid1.Anchor = ((System.Windows.Forms.AnchorStyles)((((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.metroGrid1.AutoSizeColumnsMode = System.Windows.Forms.DataGridViewAutoSizeColumnsMode.Fill;
            this.metroGrid1.BackgroundColor = System.Drawing.Color.FromArgb(((int)(((byte)(255)))), ((int)(((byte)(255)))), ((int)(((byte)(255)))));
            this.metroGrid1.BorderStyle = System.Windows.Forms.BorderStyle.None;
            this.metroGrid1.CellBorderStyle = System.Windows.Forms.DataGridViewCellBorderStyle.None;
            this.metroGrid1.ColumnHeadersBorderStyle = System.Windows.Forms.DataGridViewHeaderBorderStyle.None;
            dataGridViewCellStyle7.Alignment = System.Windows.Forms.DataGridViewContentAlignment.MiddleLeft;
            dataGridViewCellStyle7.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(209)))), ((int)(((byte)(17)))), ((int)(((byte)(65)))));
            dataGridViewCellStyle7.Font = new System.Drawing.Font("Segoe UI", 11F, System.Drawing.FontStyle.Bold, System.Drawing.GraphicsUnit.Pixel);
            dataGridViewCellStyle7.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(255)))), ((int)(((byte)(255)))), ((int)(((byte)(255)))));
            dataGridViewCellStyle7.SelectionBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(234)))), ((int)(((byte)(19)))), ((int)(((byte)(73)))));
            dataGridViewCellStyle7.SelectionForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(17)))), ((int)(((byte)(17)))), ((int)(((byte)(17)))));
            dataGridViewCellStyle7.WrapMode = System.Windows.Forms.DataGridViewTriState.True;
            this.metroGrid1.ColumnHeadersDefaultCellStyle = dataGridViewCellStyle7;
            this.metroGrid1.ColumnHeadersHeightSizeMode = System.Windows.Forms.DataGridViewColumnHeadersHeightSizeMode.AutoSize;
            this.metroGrid1.Columns.AddRange(new System.Windows.Forms.DataGridViewColumn[] {
            this.dataGridViewTextBoxColumn1,
            this.dataGridViewTextBoxColumn2,
            this.dataGridViewTextBoxColumn3,
            this.dataGridViewTextBoxColumn4,
            this.dataGridViewTextBoxColumn5,
            this.dataGridViewTextBoxColumn6,
            this.dataGridViewTextBoxColumn7});
            dataGridViewCellStyle8.Alignment = System.Windows.Forms.DataGridViewContentAlignment.MiddleLeft;
            dataGridViewCellStyle8.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(255)))), ((int)(((byte)(255)))), ((int)(((byte)(255)))));
            dataGridViewCellStyle8.Font = new System.Drawing.Font("Segoe UI", 11F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Pixel);
            dataGridViewCellStyle8.ForeColor = System.Drawing.Color.Black;
            dataGridViewCellStyle8.SelectionBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(234)))), ((int)(((byte)(19)))), ((int)(((byte)(73)))));
            dataGridViewCellStyle8.SelectionForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(17)))), ((int)(((byte)(17)))), ((int)(((byte)(17)))));
            dataGridViewCellStyle8.WrapMode = System.Windows.Forms.DataGridViewTriState.False;
            this.metroGrid1.DefaultCellStyle = dataGridViewCellStyle8;
            this.metroGrid1.EnableHeadersVisualStyles = false;
            this.metroGrid1.Font = new System.Drawing.Font("Segoe UI", 11F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Pixel);
            this.metroGrid1.GridColor = System.Drawing.Color.FromArgb(((int)(((byte)(255)))), ((int)(((byte)(255)))), ((int)(((byte)(255)))));
            this.metroGrid1.Location = new System.Drawing.Point(25, 58);
            this.metroGrid1.Name = "metroGrid1";
            this.metroGrid1.ReadOnly = true;
            this.metroGrid1.RowHeadersBorderStyle = System.Windows.Forms.DataGridViewHeaderBorderStyle.None;
            dataGridViewCellStyle9.Alignment = System.Windows.Forms.DataGridViewContentAlignment.MiddleLeft;
            dataGridViewCellStyle9.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(209)))), ((int)(((byte)(17)))), ((int)(((byte)(65)))));
            dataGridViewCellStyle9.Font = new System.Drawing.Font("Segoe UI", 11F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Pixel);
            dataGridViewCellStyle9.ForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(255)))), ((int)(((byte)(255)))), ((int)(((byte)(255)))));
            dataGridViewCellStyle9.SelectionBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(234)))), ((int)(((byte)(19)))), ((int)(((byte)(73)))));
            dataGridViewCellStyle9.SelectionForeColor = System.Drawing.Color.FromArgb(((int)(((byte)(17)))), ((int)(((byte)(17)))), ((int)(((byte)(17)))));
            dataGridViewCellStyle9.WrapMode = System.Windows.Forms.DataGridViewTriState.True;
            this.metroGrid1.RowHeadersDefaultCellStyle = dataGridViewCellStyle9;
            this.metroGrid1.RowHeadersWidthSizeMode = System.Windows.Forms.DataGridViewRowHeadersWidthSizeMode.DisableResizing;
            dataGridViewCellStyle10.BackColor = System.Drawing.Color.White;
            dataGridViewCellStyle10.Font = new System.Drawing.Font("Microsoft Sans Serif", 11.25F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            dataGridViewCellStyle10.ForeColor = System.Drawing.Color.Black;
            dataGridViewCellStyle10.SelectionBackColor = System.Drawing.Color.FromArgb(((int)(((byte)(192)))), ((int)(((byte)(0)))), ((int)(((byte)(0)))));
            dataGridViewCellStyle10.SelectionForeColor = System.Drawing.Color.White;
            this.metroGrid1.RowsDefaultCellStyle = dataGridViewCellStyle10;
            this.metroGrid1.SelectionMode = System.Windows.Forms.DataGridViewSelectionMode.FullRowSelect;
            this.metroGrid1.ShowCellToolTips = false;
            this.metroGrid1.ShowEditingIcon = false;
            this.metroGrid1.Size = new System.Drawing.Size(1230, 205);
            this.metroGrid1.Style = MetroFramework.MetroColorStyle.Red;
            this.metroGrid1.TabIndex = 4;
            this.metroGrid1.UseCustomBackColor = true;
            this.metroGrid1.UseCustomForeColor = true;
            // 
            // dataGridViewTextBoxColumn1
            // 
            this.dataGridViewTextBoxColumn1.DataPropertyName = "NextPresentation";
            this.dataGridViewTextBoxColumn1.HeaderText = "Horário prox. Reprodução";
            this.dataGridViewTextBoxColumn1.Name = "dataGridViewTextBoxColumn1";
            this.dataGridViewTextBoxColumn1.ReadOnly = true;
            // 
            // dataGridViewTextBoxColumn2
            // 
            this.dataGridViewTextBoxColumn2.DataPropertyName = "file";
            this.dataGridViewTextBoxColumn2.HeaderText = "Arquivo";
            this.dataGridViewTextBoxColumn2.Name = "dataGridViewTextBoxColumn2";
            this.dataGridViewTextBoxColumn2.ReadOnly = true;
            // 
            // dataGridViewTextBoxColumn3
            // 
            this.dataGridViewTextBoxColumn3.DataPropertyName = "artist";
            this.dataGridViewTextBoxColumn3.HeaderText = "Artista";
            this.dataGridViewTextBoxColumn3.Name = "dataGridViewTextBoxColumn3";
            this.dataGridViewTextBoxColumn3.ReadOnly = true;
            // 
            // dataGridViewTextBoxColumn4
            // 
            this.dataGridViewTextBoxColumn4.DataPropertyName = "title";
            this.dataGridViewTextBoxColumn4.HeaderText = "Titúlo";
            this.dataGridViewTextBoxColumn4.Name = "dataGridViewTextBoxColumn4";
            this.dataGridViewTextBoxColumn4.ReadOnly = true;
            // 
            // dataGridViewTextBoxColumn5
            // 
            this.dataGridViewTextBoxColumn5.DataPropertyName = "TotalDuration";
            this.dataGridViewTextBoxColumn5.HeaderText = "Duração";
            this.dataGridViewTextBoxColumn5.Name = "dataGridViewTextBoxColumn5";
            this.dataGridViewTextBoxColumn5.ReadOnly = true;
            // 
            // dataGridViewTextBoxColumn6
            // 
            this.dataGridViewTextBoxColumn6.DataPropertyName = "album";
            this.dataGridViewTextBoxColumn6.HeaderText = "Album";
            this.dataGridViewTextBoxColumn6.Name = "dataGridViewTextBoxColumn6";
            this.dataGridViewTextBoxColumn6.ReadOnly = true;
            // 
            // dataGridViewTextBoxColumn7
            // 
            this.dataGridViewTextBoxColumn7.DataPropertyName = "Pos";
            this.dataGridViewTextBoxColumn7.HeaderText = "Posição";
            this.dataGridViewTextBoxColumn7.Name = "dataGridViewTextBoxColumn7";
            this.dataGridViewTextBoxColumn7.ReadOnly = true;
            // 
            // metroComboBox1
            // 
            this.metroComboBox1.FormattingEnabled = true;
            this.metroComboBox1.ItemHeight = 23;
            this.metroComboBox1.Location = new System.Drawing.Point(105, 13);
            this.metroComboBox1.Name = "metroComboBox1";
            this.metroComboBox1.Size = new System.Drawing.Size(365, 29);
            this.metroComboBox1.TabIndex = 3;
            this.metroComboBox1.UseSelectable = true;
            // 
            // labelPlaylists
            // 
            this.labelPlaylists.AutoSize = true;
            this.labelPlaylists.FontSize = MetroFramework.MetroLabelSize.Tall;
            this.labelPlaylists.Location = new System.Drawing.Point(25, 13);
            this.labelPlaylists.Name = "labelPlaylists";
            this.labelPlaylists.Size = new System.Drawing.Size(74, 25);
            this.labelPlaylists.TabIndex = 2;
            this.labelPlaylists.Text = "Playlists:";
            // 
            // tabConfiguration
            // 
            this.tabConfiguration.Controls.Add(this.metroPanel2);
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
            this.tabConfiguration.Size = new System.Drawing.Size(1277, 599);
            this.tabConfiguration.TabIndex = 2;
            this.tabConfiguration.Text = "Configurações";
            this.tabConfiguration.VerticalScrollbarBarColor = true;
            this.tabConfiguration.VerticalScrollbarHighlightOnWheel = false;
            this.tabConfiguration.VerticalScrollbarSize = 10;
            // 
            // metroPanel2
            // 
            this.metroPanel2.Controls.Add(this.labelMusicFade);
            this.metroPanel2.Controls.Add(this.toggleFadeIn);
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
            this.metroPanel2.Size = new System.Drawing.Size(1277, 46);
            this.metroPanel2.TabIndex = 11;
            this.metroPanel2.VerticalScrollbarBarColor = true;
            this.metroPanel2.VerticalScrollbarHighlightOnWheel = false;
            this.metroPanel2.VerticalScrollbarSize = 10;
            // 
            // labelMusicFade
            // 
            this.labelMusicFade.AutoSize = true;
            this.labelMusicFade.Location = new System.Drawing.Point(731, 14);
            this.labelMusicFade.Name = "labelMusicFade";
            this.labelMusicFade.Size = new System.Drawing.Size(145, 19);
            this.labelMusicFade.TabIndex = 9;
            this.labelMusicFade.Text = "Transição entre músicas";
            // 
            // toggleFadeIn
            // 
            this.toggleFadeIn.AutoSize = true;
            this.toggleFadeIn.Location = new System.Drawing.Point(903, 16);
            this.toggleFadeIn.Name = "toggleFadeIn";
            this.toggleFadeIn.Size = new System.Drawing.Size(80, 17);
            this.toggleFadeIn.TabIndex = 8;
            this.toggleFadeIn.Text = "Off";
            this.toggleFadeIn.UseSelectable = true;
            this.toggleFadeIn.CheckedChanged += new System.EventHandler(this.toggleFadeIn_CheckedChanged);
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
            this.toggleRepeatPlaylist.CheckedChanged += new System.EventHandler(this.toggleRepeatPlaylist_CheckedChanged);
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
            this.toggleShuffle.CheckedChanged += new System.EventHandler(this.toggleShuffle_CheckedChanged);
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
            this.toggleIsPlaying.CheckedChanged += new System.EventHandler(this.toggleIsPlaying_CheckedChanged);
            // 
            // buttonUpdateMusicDatabase
            // 
            this.buttonUpdateMusicDatabase.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Bottom | System.Windows.Forms.AnchorStyles.Left)));
            this.buttonUpdateMusicDatabase.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(128)))), ((int)(((byte)(128)))), ((int)(((byte)(255)))));
            this.buttonUpdateMusicDatabase.FontSize = MetroFramework.MetroButtonSize.Tall;
            this.buttonUpdateMusicDatabase.ForeColor = System.Drawing.Color.Black;
            this.buttonUpdateMusicDatabase.Location = new System.Drawing.Point(24, 532);
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
            this.textboxServerAPIKey.Location = new System.Drawing.Point(179, 192);
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
            this.textboxServerPort.Location = new System.Drawing.Point(179, 139);
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
            this.textboxServerAddress.Location = new System.Drawing.Point(179, 93);
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
            this.metroLabel1.Location = new System.Drawing.Point(47, 196);
            this.metroLabel1.Name = "metroLabel1";
            this.metroLabel1.Size = new System.Drawing.Size(91, 19);
            this.metroLabel1.TabIndex = 6;
            this.metroLabel1.Text = "&Chave de API:";
            // 
            // labelServerPort
            // 
            this.labelServerPort.AutoSize = true;
            this.labelServerPort.Location = new System.Drawing.Point(47, 143);
            this.labelServerPort.Name = "labelServerPort";
            this.labelServerPort.Size = new System.Drawing.Size(119, 19);
            this.labelServerPort.TabIndex = 5;
            this.labelServerPort.Text = "&Porta do Servidor:";
            // 
            // labelServerAddress
            // 
            this.labelServerAddress.AutoSize = true;
            this.labelServerAddress.Location = new System.Drawing.Point(24, 93);
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
            this.buttonConfigurationSave.Location = new System.Drawing.Point(1052, 532);
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
            this.buttonConfigurationCancel.Location = new System.Drawing.Point(887, 532);
            this.buttonConfigurationCancel.Name = "buttonConfigurationCancel";
            this.buttonConfigurationCancel.Size = new System.Drawing.Size(150, 53);
            this.buttonConfigurationCancel.TabIndex = 2;
            this.buttonConfigurationCancel.Text = "&Cancelar";
            this.buttonConfigurationCancel.UseCustomBackColor = true;
            this.buttonConfigurationCancel.UseCustomForeColor = true;
            this.buttonConfigurationCancel.UseSelectable = true;
            this.buttonConfigurationCancel.Click += new System.EventHandler(this.buttonConfigurationCancel_Click);
            // 
            // contextSearchGrid
            // 
            this.contextSearchGrid.Items.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.menuAddAbovePlaylist,
            this.menuAddBellowPlaylist,
            this.toolStripSeparator3,
            this.inserirMúsicaVáriasVezesNaProgramaçãoToolStripMenuItem});
            this.contextSearchGrid.Name = "contextSearchGrid";
            this.contextSearchGrid.Size = new System.Drawing.Size(302, 76);
            this.contextSearchGrid.Text = "Editar playlist";
            this.contextSearchGrid.Opening += new System.ComponentModel.CancelEventHandler(this.contextSearchGrid_Opening);
            // 
            // menuAddAbovePlaylist
            // 
            this.menuAddAbovePlaylist.Name = "menuAddAbovePlaylist";
            this.menuAddAbovePlaylist.Size = new System.Drawing.Size(301, 22);
            this.menuAddAbovePlaylist.Text = "Adicionar acima da Seleção da Playlist";
            this.menuAddAbovePlaylist.Click += new System.EventHandler(this.menuAddAbovePlaylist_Click);
            // 
            // menuAddBellowPlaylist
            // 
            this.menuAddBellowPlaylist.Name = "menuAddBellowPlaylist";
            this.menuAddBellowPlaylist.Size = new System.Drawing.Size(301, 22);
            this.menuAddBellowPlaylist.Text = "Adicionar abaixo da seleção da playlist";
            this.menuAddBellowPlaylist.Click += new System.EventHandler(this.menuAddBellowPlaylist_Click);
            // 
            // toolStripSeparator3
            // 
            this.toolStripSeparator3.Name = "toolStripSeparator3";
            this.toolStripSeparator3.Size = new System.Drawing.Size(298, 6);
            // 
            // inserirMúsicaVáriasVezesNaProgramaçãoToolStripMenuItem
            // 
            this.inserirMúsicaVáriasVezesNaProgramaçãoToolStripMenuItem.Name = "inserirMúsicaVáriasVezesNaProgramaçãoToolStripMenuItem";
            this.inserirMúsicaVáriasVezesNaProgramaçãoToolStripMenuItem.Size = new System.Drawing.Size(301, 22);
            this.inserirMúsicaVáriasVezesNaProgramaçãoToolStripMenuItem.Text = "Inserir música várias vezes na programação";
            this.inserirMúsicaVáriasVezesNaProgramaçãoToolStripMenuItem.Click += new System.EventHandler(this.menuAddMusicManyTimes_Click);
            // 
            // contextMenuPlaylist
            // 
            this.contextMenuPlaylist.Items.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.menuRemoveMusic,
            this.menuMoveAbove,
            this.menuMoveBellow,
            this.menuSwap,
            this.mnuMerge,
            this.toolStripSeparator1,
            this.menuGoToCurrent,
            this.toolStripSeparator2,
            this.menuPlaySelectedMusic});
            this.contextMenuPlaylist.Name = "contextMenuPlaylist";
            this.contextMenuPlaylist.Size = new System.Drawing.Size(348, 170);
            this.contextMenuPlaylist.Opening += new System.ComponentModel.CancelEventHandler(this.contextMenuPlaylist_Opening);
            // 
            // menuRemoveMusic
            // 
            this.menuRemoveMusic.Name = "menuRemoveMusic";
            this.menuRemoveMusic.ShortcutKeys = ((System.Windows.Forms.Keys)((System.Windows.Forms.Keys.Control | System.Windows.Forms.Keys.Delete)));
            this.menuRemoveMusic.Size = new System.Drawing.Size(347, 22);
            this.menuRemoveMusic.Text = "&Remover musicas selecionadas";
            this.menuRemoveMusic.Click += new System.EventHandler(this.menuRemoveMusic_Click);
            // 
            // menuMoveAbove
            // 
            this.menuMoveAbove.Name = "menuMoveAbove";
            this.menuMoveAbove.ShortcutKeys = ((System.Windows.Forms.Keys)((System.Windows.Forms.Keys.Control | System.Windows.Forms.Keys.Up)));
            this.menuMoveAbove.Size = new System.Drawing.Size(347, 22);
            this.menuMoveAbove.Text = "&Mover para cima músicas selecionadas";
            this.menuMoveAbove.Click += new System.EventHandler(this.menuMoveAbove_Click);
            // 
            // menuMoveBellow
            // 
            this.menuMoveBellow.Name = "menuMoveBellow";
            this.menuMoveBellow.ShortcutKeys = ((System.Windows.Forms.Keys)((System.Windows.Forms.Keys.Control | System.Windows.Forms.Keys.Down)));
            this.menuMoveBellow.Size = new System.Drawing.Size(347, 22);
            this.menuMoveBellow.Text = "Mover para &baixo músicas selecionadas";
            this.menuMoveBellow.Click += new System.EventHandler(this.menuMoveBellow_Click);
            // 
            // menuSwap
            // 
            this.menuSwap.Name = "menuSwap";
            this.menuSwap.ShortcutKeys = ((System.Windows.Forms.Keys)((System.Windows.Forms.Keys.Control | System.Windows.Forms.Keys.Right)));
            this.menuSwap.Size = new System.Drawing.Size(347, 22);
            this.menuSwap.Text = "Trocar de Lugar";
            this.menuSwap.Click += new System.EventHandler(this.menuSwap_Click);
            // 
            // mnuMerge
            // 
            this.mnuMerge.Name = "mnuMerge";
            this.mnuMerge.ShortcutKeys = ((System.Windows.Forms.Keys)((System.Windows.Forms.Keys.Control | System.Windows.Forms.Keys.D0)));
            this.mnuMerge.Size = new System.Drawing.Size(347, 22);
            this.mnuMerge.Text = "&Transformar em Program com Vinhetas";
            this.mnuMerge.Click += new System.EventHandler(this.mnuMerge_Click);
            // 
            // toolStripSeparator1
            // 
            this.toolStripSeparator1.Name = "toolStripSeparator1";
            this.toolStripSeparator1.Size = new System.Drawing.Size(344, 6);
            // 
            // menuGoToCurrent
            // 
            this.menuGoToCurrent.Name = "menuGoToCurrent";
            this.menuGoToCurrent.ShortcutKeys = ((System.Windows.Forms.Keys)((System.Windows.Forms.Keys.Control | System.Windows.Forms.Keys.N)));
            this.menuGoToCurrent.Size = new System.Drawing.Size(347, 22);
            this.menuGoToCurrent.Text = "Ir para música atual";
            this.menuGoToCurrent.Click += new System.EventHandler(this.menuGoToCurrent_Click);
            // 
            // toolStripSeparator2
            // 
            this.toolStripSeparator2.Name = "toolStripSeparator2";
            this.toolStripSeparator2.Size = new System.Drawing.Size(344, 6);
            // 
            // menuPlaySelectedMusic
            // 
            this.menuPlaySelectedMusic.Name = "menuPlaySelectedMusic";
            this.menuPlaySelectedMusic.Size = new System.Drawing.Size(347, 22);
            this.menuPlaySelectedMusic.Text = "Tocar música selecionada";
            this.menuPlaySelectedMusic.Click += new System.EventHandler(this.menuPlaySelectedMusic_Click);
            // 
            // frmMain
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 13F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(1325, 778);
            this.Controls.Add(this.tabMain);
            this.Controls.Add(this.panelMusicStatus);
            this.Name = "frmMain";
            this.Style = MetroFramework.MetroColorStyle.Red;
            this.Text = "BatRadio by Morcegão FM ($version)";
            this.WindowState = System.Windows.Forms.FormWindowState.Maximized;
            this.Load += new System.EventHandler(this.frmMain_Load);
            ((System.ComponentModel.ISupportInitialize)(this.styleManager)).EndInit();
            this.panelMusicStatus.ResumeLayout(false);
            this.tabMain.ResumeLayout(false);
            this.tabPlaying.ResumeLayout(false);
            this.metroPanel1.ResumeLayout(false);
            this.metroPanel1.PerformLayout();
            ((System.ComponentModel.ISupportInitialize)(this.gridPlaylist)).EndInit();
            this.panSearch.ResumeLayout(false);
            this.panSearch.PerformLayout();
            ((System.ComponentModel.ISupportInitialize)(this.gridSearch)).EndInit();
            this.metroTabPage2.ResumeLayout(false);
            this.panelPlaylists.ResumeLayout(false);
            this.panelPlaylists.PerformLayout();
            ((System.ComponentModel.ISupportInitialize)(this.metroGrid1)).EndInit();
            this.tabConfiguration.ResumeLayout(false);
            this.tabConfiguration.PerformLayout();
            this.metroPanel2.ResumeLayout(false);
            this.metroPanel2.PerformLayout();
            this.contextSearchGrid.ResumeLayout(false);
            this.contextMenuPlaylist.ResumeLayout(false);
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
        private MetroFramework.Controls.MetroPanel panSearch;
        private MetroFramework.Controls.MetroTextBox textboxSearch;
        private MetroFramework.Controls.MetroLabel labelSearch;
        private MetroFramework.Controls.MetroGrid gridSearch;
        private MetroFramework.Controls.MetroPanel metroPanel1;
        private MetroFramework.Controls.MetroTextBox textboxSearchPlaylist;
        private MetroFramework.Controls.MetroLabel labelSearchPlaylist;
        private MetroFramework.Controls.MetroGrid gridPlaylist;
        private MetroFramework.Controls.MetroButton metroButton2;
        private MetroFramework.Controls.MetroButton metroButton1;
        private MetroFramework.Controls.MetroButton buttonCurrentMusic;
        private System.Windows.Forms.ContextMenuStrip contextSearchGrid;
        private System.Windows.Forms.ToolStripMenuItem menuAddAbovePlaylist;
        private System.Windows.Forms.ToolStripMenuItem menuAddBellowPlaylist;
        private System.Windows.Forms.ContextMenuStrip contextMenuPlaylist;
        private System.Windows.Forms.ToolStripMenuItem menuRemoveMusic;
        private System.Windows.Forms.ToolStripMenuItem menuMoveAbove;
        private System.Windows.Forms.ToolStripMenuItem menuMoveBellow;
        private System.Windows.Forms.ToolStripMenuItem menuSwap;
        private System.Windows.Forms.ToolStripMenuItem mnuMerge;
        private System.Windows.Forms.ToolStripSeparator toolStripSeparator1;
        private System.Windows.Forms.ToolStripMenuItem menuGoToCurrent;
        private System.Windows.Forms.DataGridViewTextBoxColumn gridSearchFile;
        private System.Windows.Forms.DataGridViewTextBoxColumn gridSearchArtist;
        private System.Windows.Forms.DataGridViewTextBoxColumn gridSearchTitle;
        private System.Windows.Forms.DataGridViewTextBoxColumn gridSearchDuration;
        private System.Windows.Forms.DataGridViewTextBoxColumn gridSearchAlbum;
        private MetroFramework.Controls.MetroPanel metroPanel2;
        private MetroFramework.Controls.MetroLabel labelMusicFade;
        private MetroFramework.Controls.MetroToggle toggleFadeIn;
        private MetroFramework.Controls.MetroLabel labelRepeat;
        private MetroFramework.Controls.MetroToggle toggleRepeatPlaylist;
        private MetroFramework.Controls.MetroLabel labelShuffle;
        private MetroFramework.Controls.MetroToggle toggleShuffle;
        private MetroFramework.Controls.MetroLabel labelIsPlaying;
        private MetroFramework.Controls.MetroToggle toggleIsPlaying;
        private System.Windows.Forms.ToolStripSeparator toolStripSeparator2;
        private System.Windows.Forms.ToolStripMenuItem menuPlaySelectedMusic;
        private System.Windows.Forms.DataGridViewTextBoxColumn datagridPlaylistTime;
        private System.Windows.Forms.DataGridViewTextBoxColumn datagridPlaylistFile;
        private System.Windows.Forms.DataGridViewTextBoxColumn datagridPlaylistArtist;
        private System.Windows.Forms.DataGridViewTextBoxColumn datagridPlaylistTitle;
        private System.Windows.Forms.DataGridViewTextBoxColumn datagridPlaylistTotalDuration;
        private System.Windows.Forms.DataGridViewTextBoxColumn datagridPlaylistAlbum;
        private System.Windows.Forms.DataGridViewTextBoxColumn gridPlaylistPosition;
        private MetroFramework.Controls.MetroPanel panelPlaylists;
        private MetroFramework.Controls.MetroComboBox metroComboBox1;
        private MetroFramework.Controls.MetroLabel labelPlaylists;
        private System.Windows.Forms.ToolStripSeparator toolStripSeparator3;
        private System.Windows.Forms.ToolStripMenuItem inserirMúsicaVáriasVezesNaProgramaçãoToolStripMenuItem;
        private MetroFramework.Controls.MetroLabel labelInfo;
        private MetroFramework.Controls.MetroGrid metroGrid1;
        private System.Windows.Forms.DataGridViewTextBoxColumn dataGridViewTextBoxColumn1;
        private System.Windows.Forms.DataGridViewTextBoxColumn dataGridViewTextBoxColumn2;
        private System.Windows.Forms.DataGridViewTextBoxColumn dataGridViewTextBoxColumn3;
        private System.Windows.Forms.DataGridViewTextBoxColumn dataGridViewTextBoxColumn4;
        private System.Windows.Forms.DataGridViewTextBoxColumn dataGridViewTextBoxColumn5;
        private System.Windows.Forms.DataGridViewTextBoxColumn dataGridViewTextBoxColumn6;
        private System.Windows.Forms.DataGridViewTextBoxColumn dataGridViewTextBoxColumn7;
    }
}