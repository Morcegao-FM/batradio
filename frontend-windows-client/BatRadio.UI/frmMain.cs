using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Configuration;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace BatRadio.UI
{
    public partial class frmMain : MetroFramework.Forms.MetroForm
    {
        private BatRadioClient client = new BatRadioClient();
        Status status = new Status();
        List<StatusSong> playlist = new List<StatusSong>();
        List<StatusSong> files = new List<StatusSong>();
        public frmMain()
        {
            InitializeComponent();

        }
        private void frmMain_Load(object sender, EventArgs e)
        {
            LoadConfiguration();
            try
            {
                UpdateDatabaseInformation();                
                ShowStatus();
                timerUpdate.Enabled = true;
                tabMain.SelectTab(tabPlaying);

            }
            catch(Exception err)
            {
                tabMain.SelectTab(tabConfiguration);
                ShowError("Não foi possível conectar-se ao servidor da rádio. Por favor, verifique se o servidor está funcionando ou confirme as configurações da aplicação");
                return;
            }

            ShowMessage(string.Format("{0} - {1} : [{2}]", status.currentSong.Artist, status.currentSong.Title, status.currentSong.file));
        }

        private void ShowStatus()
        {
            labelNowPlaying.Text = string.Format("{0} - {1}   {2}", status.currentSong.Artist, status.currentSong.Title, status.GetDurationStatus());
            trackCurrentMusic.Maximum = (int)status.duration;
            trackCurrentMusic.Value = (int)status.elapsed;
            toggleIsPlaying.Checked = status.state == "play";
            toggleShuffle.Checked = status.random == 1;
            toggleRepeatPlaylist.Checked = status.repeat == 1;
        }
        private void ShowMessage(string message)
        {
        }
        private void ShowError(string message, string title = "Erro", bool error = true)
        {
            MetroFramework.MetroMessageBox.Show(this, message, title, MessageBoxButtons.OK, error?MessageBoxIcon.Error:MessageBoxIcon.Warning);
        }

        
        private void timeUpdate_Tick(object sender, EventArgs e)
        {
            try 
            {
                status = client.GetStatus();
                ShowStatus();
            }
            catch(Exception error)
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

            if(port < 1024 || port > 65535)
            {
                throw new Exception("A porta do sevidor deve estar entre 1024 e 65535");
            }

            Configuration config = ConfigurationManager.OpenExeConfiguration(ConfigurationUserLevel.None);

            config.AppSettings.Settings["Server"].Value = textboxServerAddress.Text;
            config.AppSettings.Settings["Port"].Value = port.ToString();
            config.AppSettings.Settings["APIKey"].Value = textboxServerAPIKey.Text;

            config.Save(ConfigurationSaveMode.Modified);
            ConfigurationManager.RefreshSection("appSettings");

            BatRadioClient.RadioConfig.Server = textboxServerAddress.Text;
            BatRadioClient.RadioConfig.Port = port;
            BatRadioClient.RadioConfig.APIKey = textboxServerAPIKey.Text;
        }

        private void UpdateDatabaseInformation()
        {
            status = client.GetStatus();
            playlist = client.GetPlayList(true);
            files = client.GetSongs(false);            
        }
        private void buttonConfigurationSave_Click(object sender, EventArgs e)
        {
            try
            {
                SaveConfiguration();
                if (timerUpdate.Enabled == false) timerUpdate.Enabled = true;
            }
            catch(Exception error)
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

    }
}
