using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace BatRadio.UI
{
    public partial class frmMain : Form
    {
        private BatRadioClient client = new BatRadioClient();
        public frmMain()
        {
            InitializeComponent();

        }

        private void configuraçãoToolStripMenuItem_Click(object sender, EventArgs e)
        {
            frmConfig config = new frmConfig();
            config.ShowDialog(this);
        }

        private void frmMain_Load(object sender, EventArgs e)
        {
            Status status = new Status();
            try
            {
                status = client.GetStatus();
                var playlist = client.GetPlayList();
                gridPlaylist.DataSource = playlist;
                gridPlaylist.Update();
            }
            catch(Exception err)
            {
                ShowError(err.Message);
                return;
            }

            ShowMessage(string.Format("{0} - {1} : [{2}]", status.currentSong.Artist, status.currentSong.Title, status.currentSong.file));
        }


        private void ShowMessage(string message)
        {
            labelMusic.ForeColor = Color.Black;
            labelMusic.Text = message;
        }
        private void ShowError(string message)
        {
            labelMusic.ForeColor = Color.Red;
            labelMusic.Text = message;
        }
    }
}
