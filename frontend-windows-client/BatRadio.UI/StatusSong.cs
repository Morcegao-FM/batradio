﻿using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace BatRadio.UI
{
    class PlaylistList
    {
        [Newtonsoft.Json.JsonProperty(propertyName: "Last-Modified")]
        public DateTime LastModified { get; set; }
        public string playlist { get; set; }

        public override string ToString()
        {
            return playlist;
        }
    }
    class StatusSong
    {
        public string file { get; set; }        

        public string InternalID { get { return string.Format("{0}:{1}", this.Id, this.Pos); } }

        [Newtonsoft.Json.JsonProperty(propertyName: "Last-Modified")]
        public DateTime LastModified { get; set; }

        public string Artist { get; set; }
        public string Title { get; set; }
        public int Time { get; set; }
        public double duration { get; set; }

        public string TotalDuration { get { return RadioHelper.SecondsToString(this.Time); } }
        public int Pos { get; set; }
        public int Id { get; set; }
        public string Album { get; set; }
        public string Genre { get; set; }
        public DateTime NextPresentation {get;set;}

        public bool Matches(string[] Information)
        {            
            var key = this.ToString().ToUpper();
            foreach(var item in Information)
            {
                if (!key.Contains(item.ToUpper())) return false;
            }
            return true;
        }
        public override string ToString()
        {
            return string.Format("{0} {1} {2} {3}", file, Title, Album, Artist);
        }
    }
}
