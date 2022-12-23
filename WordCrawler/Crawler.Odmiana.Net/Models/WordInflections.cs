namespace Crawler.Odmiana.Net.Models;

public class WordInflections
{
    public WordInflections(string word)
    {
        Word = word;
    }

    public string Word { get; }

    public List<string> Cases { get; set; } = new List<string>();

    public List<string> SingleCases { get; set; } = new List<string>();

    public List<string> PluralCases { get; set; } = new List<string>();
}