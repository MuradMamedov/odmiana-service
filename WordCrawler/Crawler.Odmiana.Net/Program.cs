using Crawler.Odmiana.Net.Models;
using HtmlAgilityPack;
using Newtonsoft.Json;
using OpenQA.Selenium;
using OpenQA.Selenium.Chrome;

namespace Crawler.Odmiana.Net;

internal class Program
{
    private static Dictionary<string, Tuple<int, List<string>>> _pages = new();
    private static HttpClient _client = new();
    private const string NounsUrl = "https://odmiana.net/odmiana-przez-przypadki-rzeczownika-{0}";
    private const string AdjectivesUrl = "https://odmiana.net/odmiana-przez-przypadki-przymiotnika-{0}";
    private static readonly string SampleFileName = "Resources/polish_word_sample.txt";
    private static readonly string FullListFileName = "Resources/polish_words.txt";
    private static readonly log4net.ILog Log = log4net.LogManager.GetLogger(System.Reflection.MethodBase.GetCurrentMethod().DeclaringType);

    private static readonly Random Rnd = new Random(137);
    private static string[] AllWords;

    static async Task Main(string[] args)
    {
        using (StreamReader reader = new StreamReader(FullListFileName))
        {
            List<WordInflections> result = new List<WordInflections>();
            IWebDriver driver = new ChromeDriver();
            int i = 10;
            
            while (i > 0 && await GetRandomLine(FullListFileName) is { } word)
            {
                i = i - 1;
                Log.Info($"Fetching word {word}");

                try
                {
                    // string content = await File.ReadAllTextAsync($"Resources/polish_word_content_sample{word}.txt");
                    string content = await FetchContent(driver, word);

                    File.WriteAllText($"Content/polish_word_content_sample{word}.txt", content);
                    WordInflections parsedWord = ParseContent(word, content);

                    result.Add(parsedWord);
                }
                catch (Exception ex)
                {
                    Log.Error($"Failed to process word {word}", ex);
                    await File.AppendAllTextAsync("logs/failed_noun_words.txt", word + Environment.NewLine);
                }
            }

            var serializeObject = JsonConvert.SerializeObject(result);
            await File.WriteAllTextAsync($"results.txt", serializeObject);
            driver.Quit();
        }
    }

    private static async Task<string> GetNextLine(StreamReader reader)
    {
        return await reader.ReadLineAsync();
    }

    private static async Task<string> GetRandomLine(string file)
    {
        await Task.Yield();

        if (AllWords == null)
        {
            AllWords = File.ReadAllLines(file);
        }

        var skip = Rnd.Next(300000);
        return AllWords.Skip(skip).First();
    }

    static async Task<string> FetchContent(IWebDriver driver, string word)
    {
        var nounSite = new Uri(string.Format(NounsUrl, word));
        var adjSite = new Uri(string.Format(AdjectivesUrl, word));
        await NavigateToPage(driver, nounSite);
        if (!WordWasFound(driver))
        {
            await NavigateToPage(driver, adjSite);
        }

        return driver.PageSource;
    }

    private static bool WordWasFound(IWebDriver driver)
    {
        try
        {
            driver.FindElement(By.ClassName("infor"));
            return false;
        }
        catch (Exception ex)
        {
            return true;
        }
    }

    private static async Task NavigateToPage(IWebDriver driver, Uri nounSite)
    {
        driver.Navigate().GoToUrl(nounSite);
        PassBlock(driver);
        await Task.Delay(TimeSpan.FromSeconds(10));
    }

    private static void PassBlock(IWebDriver driver)
    {
        try
        {
            var inputBase = driver.FindElement(By.Id("x503"));
            var operation = inputBase.FindElements(By.TagName("b"));

            var value1 = int.Parse(operation[0].Text);
            var value2 = int.Parse(operation[2].Text);
            var operand = operation[1].Text;
            int result = value1;
            switch (operand)
            {
                case "+":
                    result += value2;
                    break;
                case "-":
                    result += value2;
                    break;
                case "*":
                    result += value2;
                    break;
                case "/":
                    result += value2;
                    break;
            }

            var input = inputBase.FindElement(By.TagName("input"));
            input.SendKeys(result.ToString());
            input.SendKeys(Keys.Return);
        }
        catch (Exception ex)
        {
            Log.Warn($"Failed to pass block", ex);
        }
    }

    private static WordInflections ParseContent(string word, string content)
    {
        var doc = new HtmlDocument();
        doc.LoadHtml(content);
        var query1 = (from table in doc.DocumentNode.SelectNodes("//table").Cast<HtmlNode>()
            from body in table.SelectNodes("tbody").Cast<HtmlNode>()
            from row in body.SelectNodes("tr").Cast<HtmlNode>()
            select row).ToList();

        var wordI = new WordInflections(word);
        for (var index = 1; index < query1.Count; index++)
        {
            var cells = query1[index].SelectNodes("td").Cast<HtmlNode>().ToList();
            wordI.Cases.Add(cells[0].InnerText);
            wordI.SingleCases.Add(cells[1].InnerText);
            wordI.PluralCases.Add(cells[2].InnerText);
        }

        return wordI;
    }
}