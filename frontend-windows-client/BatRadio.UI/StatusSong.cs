using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace BatRadio.UI
{
    class StatusSong
    {
        public string file { get; set; }
        
        [Newtonsoft.Json.JsonProperty(propertyName:"Last-Modified")]
        public DateTime LastModified { get; set; }

        public string Artist { get; set; }
        public string Title { get; set; }
        public int Time { get; set; }
        public float duration { get; set; }
        public int Pos { get; set; }
        public int Id { get; set; }
        public string Album { get; set; }
        public string Genre { get; set; }
    }
}
