using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace BatRadio.UI
{
    class Status
    {
        public int xfade { get; set; }

        public int repeat { get; set; }
        public int random { get; set; }
        public int single { get; set; }
        public int consume { get; set; }
        public string partition { get; set; }
        public int playlist { get; set; }
        public int playlistlength { get; set; }
        public float mixrampdb { get; set; }
        public string state { get; set; }
        public int song { get; set; }
        public int songid { get; set; }
        public string time { get; set; }
        public float elapsed { get; set; }
        public int bitrate { get; set; }

        public float duration { get; set; }
        public string audio { get; set; }

        public int nextsong { get; set; }
        public int nextsongid { get; set; }

        public StatusSong currentSong { get; set; }
        public StatusSong nextSong { get; set; }


        public string GetDurationStatus()
        {
            return string.Format("{0} de {1} )  ",  RadioHelper.SecondsToString(elapsed), RadioHelper.SecondsToString(duration));
        }


    }

}
