namespace BatRadio.UI
{
    partial class frmAddMusicManyTimes
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
            this.datetimeStart = new MetroFramework.Controls.MetroDateTime();
            this.checkPeriodo = new MetroFramework.Controls.MetroCheckBox();
            this.labelMusicName = new MetroFramework.Controls.MetroLabel();
            this.labelAt = new MetroFramework.Controls.MetroLabel();
            this.datetimeEnd = new MetroFramework.Controls.MetroDateTime();
            this.SuspendLayout();
            // 
            // datetimeStart
            // 
            this.datetimeStart.CustomFormat = "dd/MM/yyyy HH:mm";
            this.datetimeStart.DisplayFocus = true;
            this.datetimeStart.Format = System.Windows.Forms.DateTimePickerFormat.Custom;
            this.datetimeStart.Location = new System.Drawing.Point(27, 131);
            this.datetimeStart.MinimumSize = new System.Drawing.Size(0, 29);
            this.datetimeStart.Name = "datetimeStart";
            this.datetimeStart.Size = new System.Drawing.Size(200, 29);
            this.datetimeStart.Style = MetroFramework.MetroColorStyle.Green;
            this.datetimeStart.TabIndex = 1;
            // 
            // checkPeriodo
            // 
            this.checkPeriodo.AutoSize = true;
            this.checkPeriodo.Location = new System.Drawing.Point(27, 98);
            this.checkPeriodo.Name = "checkPeriodo";
            this.checkPeriodo.Size = new System.Drawing.Size(153, 15);
            this.checkPeriodo.TabIndex = 2;
            this.checkPeriodo.Text = "Apenas por um período?";
            this.checkPeriodo.UseSelectable = true;
            // 
            // labelMusicName
            // 
            this.labelMusicName.AutoSize = true;
            this.labelMusicName.Location = new System.Drawing.Point(27, 64);
            this.labelMusicName.Name = "labelMusicName";
            this.labelMusicName.Size = new System.Drawing.Size(78, 19);
            this.labelMusicName.TabIndex = 3;
            this.labelMusicName.Text = "musicName";
            // 
            // labelAt
            // 
            this.labelAt.AutoSize = true;
            this.labelAt.Location = new System.Drawing.Point(249, 141);
            this.labelAt.Name = "labelAt";
            this.labelAt.Size = new System.Drawing.Size(27, 19);
            this.labelAt.TabIndex = 4;
            this.labelAt.Text = "até";
            // 
            // datetimeEnd
            // 
            this.datetimeEnd.CustomFormat = "dd/MM/yyyy HH:mm";
            this.datetimeEnd.Format = System.Windows.Forms.DateTimePickerFormat.Custom;
            this.datetimeEnd.Location = new System.Drawing.Point(300, 131);
            this.datetimeEnd.MinimumSize = new System.Drawing.Size(0, 29);
            this.datetimeEnd.Name = "datetimeEnd";
            this.datetimeEnd.Size = new System.Drawing.Size(200, 29);
            this.datetimeEnd.Style = MetroFramework.MetroColorStyle.Green;
            this.datetimeEnd.TabIndex = 5;
            // 
            // frmAddMusicManyTimes
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 13F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(1053, 450);
            this.Controls.Add(this.datetimeEnd);
            this.Controls.Add(this.labelAt);
            this.Controls.Add(this.labelMusicName);
            this.Controls.Add(this.checkPeriodo);
            this.Controls.Add(this.datetimeStart);
            this.Name = "frmAddMusicManyTimes";
            this.Style = MetroFramework.MetroColorStyle.Black;
            this.Text = "Adicionar música na programação";
            this.Load += new System.EventHandler(this.frmAddMusicManyTimes_Load);
            this.ResumeLayout(false);
            this.PerformLayout();

        }

        #endregion
        private MetroFramework.Controls.MetroDateTime datetimeStart;
        private MetroFramework.Controls.MetroCheckBox checkPeriodo;
        private MetroFramework.Controls.MetroLabel labelMusicName;
        private MetroFramework.Controls.MetroLabel labelAt;
        private MetroFramework.Controls.MetroDateTime datetimeEnd;
    }
}