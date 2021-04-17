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
    }
}
