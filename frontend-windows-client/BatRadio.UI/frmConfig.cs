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
    public partial class frmConfig : Form
    {
        public frmConfig()
        {
            InitializeComponent();
        }

        private void frmConfig_Load(object sender, EventArgs e)
        {
            this.textServidor.Text = BatRadioClient.RadioConfig.Server;
            this.textPorta.Text = BatRadioClient.RadioConfig.Port.ToString();
            this.textAPIKey.Text = BatRadioClient.RadioConfig.APIKey;
        }

        private void buttonSalvar_Click(object sender, EventArgs e)
        {
            Configuration config = ConfigurationManager.OpenExeConfiguration(ConfigurationUserLevel.None);
            
            config.AppSettings.Settings["Server"].Value = textServidor.Text;
            config.AppSettings.Settings["Port"].Value = textPorta.Text;
            config.AppSettings.Settings["APIKey"].Value = textAPIKey.Text;

            config.Save(ConfigurationSaveMode.Modified);
            ConfigurationManager.RefreshSection("appSettings");

            BatRadioClient.RadioConfig.Server = textServidor.Text;
            BatRadioClient.RadioConfig.Port = int.Parse(textPorta.Text);
            BatRadioClient.RadioConfig.APIKey = textAPIKey.Text;



            this.Close();
        }

        private void buttonCancelar_Click(object sender, EventArgs e)
        {
            this.Close();
        }
    }
}
