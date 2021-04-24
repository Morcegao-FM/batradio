using Newtonsoft.Json;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace BatRadio.UI
{
    class Config
    {       
        public string Server { get; set; }
        public int Port { get; set; }
        public string APIKey { get; set; }

        public Config()
        {
            this.Server = "localhost";
            this.Port = 9320;
            this.APIKey = "no-api-key";              
        }

        public void Load()
        {
            string fileName = GetFileName();

            if (System.IO.File.Exists(fileName))
            {
                try
                {
                    string content = System.IO.File.ReadAllText(fileName);
                    Config config = JsonConvert.DeserializeObject<Config>(content);
                    this.Server = config.Server;
                    this.Port = config.Port;
                    this.APIKey = config.APIKey;
                    return;
                }
                catch (Exception e)
                {

                }
            }

        }
        public void Save()
        {
            string filename = GetFileName();
            System.IO.File.WriteAllText(filename, JsonConvert.SerializeObject(this));
        }

        private string GetFileName()
        {
#if DEBUG
            return System.IO.Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.UserProfile), "batradio.debug.config");
#endif
            return System.IO.Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.UserProfile), "batradio.config");
        }
    }

}
