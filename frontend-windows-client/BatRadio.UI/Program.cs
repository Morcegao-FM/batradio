using System;
using System.Collections.Generic;
using System.Configuration;
using System.Linq;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace BatRadio.UI
{
    static class Program
    {
        /// <summary>
        /// Ponto de entrada principal para o aplicativo.
        /// </summary>
        [STAThread]
        static void Main()
        {
            string server = ConfigurationManager.AppSettings["Server"]??"localhost";
            int port = int.Parse( ConfigurationManager.AppSettings["Port"]??"9320");
            string apiKey = ConfigurationManager.AppSettings["APIKey"]??"";

            BatRadioClient.RadioConfig.APIKey = apiKey;
            BatRadioClient.RadioConfig.Port = port;
            BatRadioClient.RadioConfig.Server = server;
            
            Application.EnableVisualStyles();
            Application.SetCompatibleTextRenderingDefault(false);
            Application.Run(new frmMain());
        }
    }
}
