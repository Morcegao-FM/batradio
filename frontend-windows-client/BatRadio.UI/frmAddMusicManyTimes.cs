using MetroFramework.Forms;
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
    public partial class frmAddMusicManyTimes : MetroForm
    {
        private string FileName;
        public frmAddMusicManyTimes(string song)
        {
            InitializeComponent();
            this.FileName = song;
            labelMusicName.Text = song;
        }

        private void frmAddMusicManyTimes_Load(object sender, EventArgs e)
        {
            datetimeStart.Value = DateTime.Now;
            datetimeEnd.Value = DateTime.Now.AddDays(1);
        }
    }
}
