using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Net.Http;
using Newtonsoft.Json;

namespace BatRadio.UI
{
    class BatRadioClient
    {
        public static HttpClient Client = new HttpClient();
        public static Config RadioConfig = new Config();
        private string GetString(string URL)
        {
            return this.GetString(URL, new Dictionary<string, string>());
        }


        public static void TimeCalculation(List<StatusSong> Songs, Status CurrentStatus)
        {
            
            var initialIndex = Songs.FindIndex(t => { return t.Pos == CurrentStatus.song; });
            if (initialIndex == -1) return;
            var song = Songs[initialIndex];

            song.NextPresentation = DateTime.Now.AddSeconds((int)CurrentStatus.elapsed * -1);
            DateTime referenteDateTime = song.NextPresentation;
            StatusSong lastSong = song;

            for (int i = initialIndex + 1; i < Songs.Count; i++)
            {
                Songs[i].NextPresentation = referenteDateTime.AddSeconds((int)lastSong.Time);
                referenteDateTime = Songs[i].NextPresentation;
                lastSong = Songs[i];
            }
            
            // Sem repetição, calcula pra trás
            if (CurrentStatus.repeat == 0)
            {
                referenteDateTime = song.NextPresentation;                
                for (int i = initialIndex - 1 ; i > -1; i--)
                {
                    Songs[i].NextPresentation = referenteDateTime.Subtract(new TimeSpan(0, 0, (int)Songs[i].Time));
                    referenteDateTime = Songs[i].NextPresentation;
                }
            }
            else
            {
                for(int i = 0; i < initialIndex; i++)
                {
                    Songs[i].NextPresentation = referenteDateTime.AddSeconds((int)Songs[i].Time);
                    referenteDateTime = Songs[i].NextPresentation;
                }
            }
        }

        private string GetString(string URL, Dictionary<string,string> AdditionalHeaders)
        {                
            Client.DefaultRequestHeaders.Clear();
            Client.DefaultRequestHeaders.Add("batradio-apikey", RadioConfig.APIKey);

            foreach (KeyValuePair<string, string> item in AdditionalHeaders)
                Client.DefaultRequestHeaders.Add(item.Key, item.Value);

            var task = Client.GetAsync(new Uri(URL));
            task.Wait();
            var taskdata = task.Result.Content.ReadAsStringAsync();
            taskdata.Wait();
            if (task.Result.StatusCode != System.Net.HttpStatusCode.OK)
            {
                throw new Exception(JsonConvert.DeserializeObject<Dictionary<string, string>>(taskdata.Result)["message"]); ;
            }
                return taskdata.Result;
            
        }

        private string PostString(string URL)
        {
            return this.PostString(URL, new Dictionary<string, string>());
        }
        private string PostString(string URL, Dictionary<string, string> AdditionalHeaders)
        {
            Client.DefaultRequestHeaders.Clear();
            Client.DefaultRequestHeaders.Add("batradio-apikey", RadioConfig.APIKey);

            foreach (KeyValuePair<string, string> item in AdditionalHeaders)
                Client.DefaultRequestHeaders.Add(item.Key, item.Value);

            var task = Client.PostAsync(new Uri(URL), null);
            task.Wait();
            var taskdata = task.Result.Content.ReadAsStringAsync();
            taskdata.Wait();
            if (task.Result.StatusCode != System.Net.HttpStatusCode.OK)
            {
                throw new Exception(JsonConvert.DeserializeObject<Dictionary<string, string>>(taskdata.Result)["message"]); ;
            }
            return taskdata.Result;

        }

        private string PrepareUrl(string method)
        {
            return string.Format("http://{0}:{1}/{2}",RadioConfig.Server, RadioConfig.Port,method);
        }
        public Status GetStatus()
        {
            
            var status = GetString(PrepareUrl("status"));
            return JsonConvert.DeserializeObject<Status>(status);

        }

        public Status PlayOrPause()
        {
            var status = PostString(PrepareUrl("playorpause"));
            return JsonConvert.DeserializeObject<Status>(status);
        }

        public Status FadeIn()
        {
            var status = PostString(PrepareUrl("fadein"));
            return JsonConvert.DeserializeObject<Status>(status);
        }


        public List<StatusSong> GetPlayList(bool updateDatabase = false)
        {
            var parameters = new Dictionary<string, string>();
            parameters.Add("update", updateDatabase?"true":"false");
            var status = PostString(PrepareUrl("playlist"), parameters);
            return JsonConvert.DeserializeObject<List<StatusSong>>(status);
        }


        public List<StatusSong> GetSongs(bool updateDatabase = false)
        {
            var parameters = new Dictionary<string, string>();
            parameters.Add("update", updateDatabase ? "true" : "false");
            parameters.Add("type", "file");
            var status = GetString(PrepareUrl("list"), parameters);
            return JsonConvert.DeserializeObject<List<StatusSong>>(status);

        }

        internal Status Repeat()
        {
            var status = PostString(PrepareUrl("repeat"));
            return JsonConvert.DeserializeObject<Status>(status);

        }

        internal Status Shuffle()
        {
            var status = PostString(PrepareUrl("shuffle"));
            return JsonConvert.DeserializeObject<Status>(status);

        }

        internal List<StatusSong> AddMusic(List<string> files, int position)
        {
            var parameters = new Dictionary<string, string>();
            parameters.Add("files",  JsonConvert.SerializeObject(files));
            parameters.Add("position", position.ToString());

            return JsonConvert.DeserializeObject<List<StatusSong>>(PostString(PrepareUrl("addtoplaylist"), parameters));
        }

        public List<StatusSong> RemoveMusic(List<string> files)
        {
            var parameters = new Dictionary<string, string>();
            parameters.Add("files", JsonConvert.SerializeObject(files));            

            return JsonConvert.DeserializeObject<List<StatusSong>>(PostString(PrepareUrl("delete"), parameters));
        }

        public List<StatusSong> Move(int currentPos, int newPos)
        {
            var parameters = new Dictionary<string, string>();
            parameters.Add("from-pos", currentPos.ToString());
            parameters.Add("to-pos", newPos.ToString());
            return JsonConvert.DeserializeObject<List<StatusSong>>(PostString(PrepareUrl("move"), parameters));

        }
    }
}
