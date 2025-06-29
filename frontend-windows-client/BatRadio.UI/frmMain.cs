using MetroFramework.Forms;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Configuration;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Reflection;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace BatRadio.UI
{
    public partial class frmMain : MetroFramework.Forms.MetroForm
    {
        private BatRadioClient client = new BatRadioClient();
        Status status = new Status();
        private List<StatusSong> playlist { get { return _playlist; } set { _playlist = value; BatRadioClient.TimeCalculation(_playlist, status); } }
        List<StatusSong> _playlist = new List<StatusSong>();
        List<StatusSong> files = new List<StatusSong>();
        private bool CheckChanging = false;
        private bool selectCurrent = true;
        private int current_playlist_row = 0;
        private List<PlaylistList> playlistlist = new List<PlaylistList>();
        public frmMain()
        {
            InitializeComponent();
            gridSearch.AutoGenerateColumns = false;
            gridPlaylist.AutoGenerateColumns = false;
            gridPlaylist.RowPrePaint += GridPlaylist_RowPrePaint;
            client.StartUpdate += Client_StartUpdate;
            client.EndUpdate += Client_EndUpdate;

            try
            {
                UpdateDatabaseInformation();
                ShowStatus();

                BatRadioClient.TimeCalculation(playlist, status);
                gridPlaylist.DataSource = playlist;
                gridPlaylist.DataBindingComplete += GridPlaylist_DataBindingComplete;
                timerUpdate.Enabled = true;
                tabMain.SelectTab(tabPlaying);
                this.Text = this.Text.Replace("$version", VersionLabel);
            }
            catch (Exception err)
            {
                tabMain.SelectTab(tabConfiguration);
                ShowError("Não foi possível conectar-se ao servidor da rádio. Por favor, verifique se o servidor está funcionando ou confirme as configurações da aplicação");
                return;
            }
        }

        private void Client_EndUpdate(object sender, EventArgs e)
        {
            Cursor.Current = Cursors.Default;
        }

        private void Client_StartUpdate(object sender, EventArgs e)
        {
            Cursor.Current = Cursors.WaitCursor;
        }

        private void GridPlaylist_RowPrePaint(object sender, DataGridViewRowPrePaintEventArgs e)
        {
            current_playlist_row = gridPlaylist.FirstDisplayedScrollingRowIndex;
            //Console.WriteLine(current_playlist_row);
        }

        private void GridPlaylist_DataBindingComplete(object sender, DataGridViewBindingCompleteEventArgs e)
        {
            if (selectCurrent)
                SelectCurrentMusic();
            selectCurrent = false;
        }

        private void frmMain_Load(object sender, EventArgs e)
        {
            LoadConfiguration();
            gridPlaylist.RowHeadersWidthSizeMode = DataGridViewRowHeadersWidthSizeMode.DisableResizing;
            gridSearch.RowHeadersWidthSizeMode = DataGridViewRowHeadersWidthSizeMode.DisableResizing;
            // or even better, use .DisableResizing. Most time consuming enum is DataGridViewRowHeadersWidthSizeMode.AutoSizeToAllHeaders

            // set it to false if not needed
            gridPlaylist.RowHeadersVisible = false;
            if (status.currentSong != null)
                ShowMessage(string.Format("{0} - {1} : [{2}]", status.currentSong.Artist, status.currentSong.Title, status.currentSong.file));
        }

        private void SelectCurrentMusic(int currentIndex = -1, bool scroll = false)
        {

            gridPlaylist.ClearSelection();
            if (currentIndex == -1) scroll = true;
            if (currentIndex == -1)
                currentIndex = playlist.FindIndex(t => t.Pos == status.song);
            if (currentIndex >= gridPlaylist.Rows.Count) currentIndex = gridPlaylist.Rows.Count - 1;
            gridPlaylist.Rows[currentIndex].Selected = true;
            gridPlaylist.Rows[currentIndex].Cells[0].Selected = true;
            if (currentIndex < 3) currentIndex = 3;
            gridPlaylist.CurrentCell = gridPlaylist.Rows[currentIndex].Cells[0];
            if (scroll)
                gridPlaylist.FirstDisplayedScrollingRowIndex = currentIndex - 3;
            else
                gridPlaylist.FirstDisplayedScrollingRowIndex = current_playlist_row;

        }
        private void SelectSearchMusic(int currentIndex)
        {
            if (currentIndex == -1) return;
            gridSearch.ClearSelection();
            gridSearch.Rows[currentIndex].Selected = true;
            gridSearch.Rows[currentIndex].Cells[0].Selected = true;
            gridSearch.CurrentCell = gridSearch.Rows[currentIndex].Cells[0];
        }


        private void ShowStatus()
        {
            this.CheckChanging = true;
            labelNowPlaying.Text = string.Format("[{3}] {0} - {1}   {2}", status.currentSong.Artist, status.currentSong.Title, status.GetDurationStatus(), status.songid - 1);
            if (status.currentSong.Time < status.elapsed)
                status.currentSong.Time = (int)status.elapsed;
            trackCurrentMusic.Maximum = (int)status.currentSong.Time;
            trackCurrentMusic.Value = (int)status.elapsed;
            toggleIsPlaying.Checked = status.state == "play";
            toggleShuffle.Checked = status.random == 1;
            toggleRepeatPlaylist.Checked = status.repeat == 1;
            toggleFadeIn.Checked = status.xfade > 0;
            this.CheckChanging = false;
        }
        private void ShowMessage(string message)
        {
        }
        private void ShowError(string message, string title = "Erro", bool error = true)
        {
            MetroFramework.MetroMessageBox.Show(this, message, title, MessageBoxButtons.OK, error ? MessageBoxIcon.Error : MessageBoxIcon.Warning);
        }


        private void timeUpdate_Tick(object sender, EventArgs e)
        {
            try
            {
                status = client.GetStatus();
                ShowStatus();
            }
            catch (Exception error)
            {
                timerUpdate.Enabled = false;
                ShowError("Falha ao atualizar o status da rádio, por favor, verifique sua conexão com a internet ou se o servidor de streaming está ativo." + Environment.NewLine + error.Message, "Falha ao atualizar Status", false);
                tabMain.SelectTab(tabConfiguration);
            }
        }



        #region Configuration Features

        public void LoadConfiguration()
        {
            this.textboxServerAddress.Text = BatRadioClient.RadioConfig.Server;
            this.textboxServerPort.Text = BatRadioClient.RadioConfig.Port.ToString();
            this.textboxServerAPIKey.Text = BatRadioClient.RadioConfig.APIKey;
        }

        public void SaveConfiguration()
        {
            int port = 0;
            if (!int.TryParse(textboxServerPort.Text, out port))
                throw new Exception("A porta do servidor deve ser um valor numérico.");

            if (port < 1024 || port > 65535)
            {
                throw new Exception("A porta do servidor deve estar entre 1024 e 65535");
            }
            BatRadioClient.RadioConfig.Server = textboxServerAddress.Text;
            BatRadioClient.RadioConfig.Port = port;
            BatRadioClient.RadioConfig.APIKey = textboxServerAPIKey.Text;
            BatRadioClient.RadioConfig.Save();
            UpdateDatabaseInformation();
        }

        private void UpdateDatabaseInformation()
        {
            status = client.GetStatus();
            playlist = client.GetPlayList(false);
            timeUpdate_Tick(this, null);
            files = client.GetSongs(true);
            playlistlist = client.GetPlayListList(false);

            gridPlaylist.DataSource = playlist;
            cmbPlaylist.DataSource = playlistlist;
            tabMain.SelectTab(tabPlaying);

        }
        private void buttonConfigurationSave_Click(object sender, EventArgs e)
        {
            try
            {
                SaveConfiguration();
                if (timerUpdate.Enabled == false) timerUpdate.Enabled = true;
                tabMain.SelectTab(tabPlaying);
            }
            catch (Exception error)
            {
                MetroFramework.MetroMessageBox.Show(this, error.Message, "Falha ao salvar configurações", MessageBoxButtons.OK, MessageBoxIcon.Error);
            }
        }

        private void buttonConfigurationCancel_Click(object sender, EventArgs e)
        {
            LoadConfiguration();
        }

        private void buttonUpdateMusicDatabase_Click(object sender, EventArgs e)
        {
            UpdateDatabaseInformation();
        }

        #endregion


        #region playback functions
        private void toggleIsPlaying_CheckedChanged(object sender, EventArgs e)
        {
            if (this.CheckChanging)
                return;

            status = client.PlayOrPause();
            ShowStatus();
        }



        private void toggleFadeIn_CheckedChanged(object sender, EventArgs e)
        {
            if (this.CheckChanging) return;

            status = client.FadeIn();
            ShowStatus();
        }

        private void toggleRepeatPlaylist_CheckedChanged(object sender, EventArgs e)
        {
            if (this.CheckChanging) return;
            status = client.Repeat();
            ShowStatus();
        }

        private void toggleShuffle_CheckedChanged(object sender, EventArgs e)
        {
            if (this.CheckChanging) return;
            status = client.Shuffle();
            ShowStatus();
        }

        private void TextboxSearch_TextChanged(object sender, System.EventArgs e)
        {
            if (textboxSearch.Text.Length < 3) return;
            var list = files.FindAll(t => { return t.Matches(textboxSearch.Text.Split(' ')); }).ToList();
            gridSearch.DataSource = list.OrderBy(t => t.file).ToList();
            gridSearch.Update();
        }


        #endregion



        private void buttonCurrentMusic_Click_1(object sender, EventArgs e)
        {
            SelectCurrentMusic();
        }

        private void textboxSearchPlaylist_KeyPress(object sender, KeyPressEventArgs e)
        {
            if (e.KeyChar == (char)Keys.Enter)
            {
                var searchValue = textboxSearchPlaylist.Text;
                var selectedRow = gridPlaylist.SelectedRows[0].Index;
                var output = playlist.FirstOrDefault(t => t.Matches(searchValue.Split(' ')) && t.Id > selectedRow + 1);

                if (output == null)
                    output = playlist.FirstOrDefault(t => t.Matches(searchValue.Split(' ')));

                if (output == null) return;

                SelectCurrentMusic(output.Pos, true);


            }
        }

        private void trackCurrentMusic_Scroll(object sender, ScrollEventArgs e)
        {
            if (!CheckChanging)
            {
                Console.WriteLine(e.NewValue);
            }

        }

        private void contextSearchGrid_Opening(object sender, CancelEventArgs e)
        {
            menuAddAbovePlaylist.Enabled = gridSearch.SelectedRows.Count > 0;
            menuAddBellowPlaylist.Enabled = gridSearch.SelectedRows.Count > 0;
        }

        private void menuAddBellowPlaylist_Click(object sender, EventArgs e)
        {
            List<string> files = new List<string>();
            int position = 0;

            foreach (DataGridViewRow item in gridSearch.SelectedRows)
            {
                files.Add(item.Cells["gridSearchFile"].Value.ToString());
            }
            if (gridPlaylist.SelectedRows.Count > 0)
                position = int.Parse(gridPlaylist.SelectedRows[gridPlaylist.SelectedRows.Count - 1].Cells["gridPlaylistPosition"].Value.ToString()) + 1;

            playlist = client.AddMusic(files, position);

            gridPlaylist.DataSource = playlist;
            SelectCurrentMusic(position);
        }

        private void menuAddAbovePlaylist_Click(object sender, EventArgs e)
        {
            List<string> files = new List<string>();
            int position = 0;

            foreach (DataGridViewRow item in gridSearch.SelectedRows)
            {
                files.Add(item.Cells["gridSearchFile"].Value.ToString());
            }
            if (gridPlaylist.SelectedRows.Count > 0)
                position = int.Parse(gridPlaylist.SelectedRows[gridPlaylist.SelectedRows.Count - 1].Cells["gridPlaylistPosition"].Value.ToString());

            playlist = client.AddMusic(files, position);
            gridPlaylist.DataSource = playlist;
            SelectCurrentMusic(position);
        }

        private void buttonAddBellow_Click(object sender, EventArgs e)
        {
            menuAddBellowPlaylist_Click(sender, e);
        }

        private void buttonAddAbove_Click(object sender, EventArgs e)
        {
            menuAddAbovePlaylist_Click(sender, e);
        }

        private void menuRemoveMusic_Click(object sender, EventArgs e)
        {
            List<string> files = new List<string>();
            int position = int.Parse(gridPlaylist.SelectedRows[0].Cells["gridPlaylistPosition"].Value.ToString());
            foreach (DataGridViewRow item in gridPlaylist.SelectedRows)
            {
                files.Add(item.Cells["gridPlaylistPosition"].Value.ToString());
            }

            playlist = client.RemoveMusic(files);
            gridPlaylist.DataSource = playlist;
            SelectCurrentMusic(position);

        }

        private void menuMoveAbove_Click(object sender, EventArgs e)
        {
            if (gridPlaylist.SelectedRows.Count == 0)
            {
                MessageBox.Show("Selecion apenas uma musica para movimentar");
                return;
            }
            int position = int.Parse(gridPlaylist.SelectedRows[0].Cells["gridPlaylistPosition"].Value.ToString());
            if (position == 0) return;
            playlist = client.Move(position, position - 1);
            gridPlaylist.DataSource = playlist;
            SelectCurrentMusic(position);
        }

        private void menuMoveBellow_Click(object sender, EventArgs e)
        {
            if (gridPlaylist.SelectedRows.Count == 0)
            {
                MessageBox.Show("Selecione apenas uma musica para movimentar");
                return;
            }
            int position = int.Parse(gridPlaylist.SelectedRows[0].Cells["gridPlaylistPosition"].Value.ToString());
            if (position == playlist.Max(t => t.Pos)) return;
            playlist = client.Move(position, position + 1);
            gridPlaylist.DataSource = playlist;
            SelectCurrentMusic(position);
        }

        private void menuSwap_Click(object sender, EventArgs e)
        {
            if (gridPlaylist.SelectedRows.Count < 2)
            {
                MessageBox.Show("Selecione duas músicas para trocar de lugar");
                return;
            }
            List<int> values = new List<int>();
            values.Add(int.Parse(gridPlaylist.SelectedRows[0].Cells["gridPlaylistPosition"].Value.ToString()));
            values.Add(int.Parse(gridPlaylist.SelectedRows[1].Cells["gridPlaylistPosition"].Value.ToString()));
            values.Sort();
            client.Move(values[0], values[1]);
            playlist = client.Move(values[1] - 1, values[0]);
            gridPlaylist.DataSource = playlist;
            SelectCurrentMusic(values[0]);
        }

        private void mnuMerge_Click(object sender, EventArgs e)
        {
            if (gridPlaylist.SelectedRows.Count % 2 != 0)
            {
                MessageBox.Show("Por favor, selecione múltiplos de 2");
            }
            List<int> values = new List<int>();
            foreach (DataGridViewRow row in gridPlaylist.SelectedRows)
            {
                values.Add(int.Parse(row.Cells["gridPlaylistPosition"].Value.ToString()));
            }
            values.Sort();
            var initial = values.Count / 2;
            int aux = 1;
            for (int i = initial; i < values.Count; i++)
            {
                playlist = client.Move(values[i], values[aux]);
                aux = aux + 2;
            }
            gridPlaylist.DataSource = playlist;
            SelectCurrentMusic(values[0]);
        }

        private void menuGoToCurrent_Click(object sender, EventArgs e)
        {
            buttonCurrentMusic_Click_1(sender, e);
        }

        private void menuPlaySelectedMusic_Click(object sender, EventArgs e)
        {
            if (gridPlaylist.SelectedRows.Count == 0) return;
            DataGridViewRow currentRow = gridPlaylist.SelectedRows[0];
            var position = int.Parse(currentRow.Cells["gridPlayListPosition"].Value.ToString());
            var filename = currentRow.Cells["datagridPlaylistFile"].Value.ToString();

            DialogResult result = MetroFramework.MetroMessageBox.Show(this, "Você tem certeza de que quer tocar " + filename + "?", "TEM CERTEZA?", MessageBoxButtons.YesNo);
            if (result == DialogResult.Yes)
            {
                result = MetroFramework.MetroMessageBox.Show(this, "Você vai parar de tocar a música " + status.currentSong.file + ", tem CERTEZA???????????", "TEM CERTEZA", MessageBoxButtons.YesNo);
                if (result == DialogResult.No)
                    return;
            }
            else
                return;

            status = client.Play(position);
            ShowStatus();

        }

        private void gridPlaylist_MouseClick(object sender, MouseEventArgs e)
        {
            if (e.Button == MouseButtons.Right && gridPlaylist.SelectedRows.Count < 2)
            {
                var rowIndex = gridPlaylist.HitTest(e.X, e.Y).RowIndex;
                SelectCurrentMusic(rowIndex, false);
            }
        }

        public string VersionLabel
        {
            get
            {
                if (System.Deployment.Application.ApplicationDeployment.IsNetworkDeployed)
                {
                    Version ver = System.Deployment.Application.ApplicationDeployment.CurrentDeployment.CurrentVersion;
                    return string.Format("{0}.{1}.{2}.{3}", ver.Major, ver.Minor, ver.Build, ver.Revision, Assembly.GetEntryAssembly().GetName().Name);
                }
                else
                {
                    var ver = Assembly.GetExecutingAssembly().GetName().Version;
                    return string.Format("{0}.{1}.{2}.{3}", ver.Major, ver.Minor, ver.Build, ver.Revision, Assembly.GetEntryAssembly().GetName().Name);
                }
            }
        }

        private void metroPanel3_Paint(object sender, PaintEventArgs e)
        {

        }

        private void menuAddMusicManyTimes_Click(object sender, EventArgs e)
        {
            var song = gridSearch.SelectedRows[0].Cells["gridSearchFile"].Value.ToString(); ;
            frmAddMusicManyTimes addForm = new frmAddMusicManyTimes(song);
            addForm.FormClosed += AddForm_FormClosed;
            addForm.Show(this);

        }

        private void AddForm_FormClosed(object sender, FormClosedEventArgs e)
        {
            ((MetroForm)sender).Dispose();
        }

        private void gridSearch_MouseClick(object sender, MouseEventArgs e)
        {
            if (e.Button == MouseButtons.Right && gridSearch.SelectedRows.Count < 2)
            {
                var rowIndex = gridSearch.HitTest(e.X, e.Y).RowIndex;
                if (rowIndex == -1) return;
                SelectSearchMusic(rowIndex);
            }
        }

        private void contextMenuPlaylist_Opening(object sender, CancelEventArgs e)
        {

        }

        private void textboxSearch_Click(object sender, EventArgs e)
        {

        }

        private void button1_Click(object sender, EventArgs e)
        {
            playlist = client.LoadPlaylist(cmbPlaylist.Text);
            
            playlist = client.GetPlayList();
            var max = playlist.Count;
            var pos = new Random().Next(1, max - 1);
            status = client.Play(pos);
            ShowStatus();
            this.tabPlaying.Select();
            this.UpdateDatabaseInformation(); 
        }
    }
}
