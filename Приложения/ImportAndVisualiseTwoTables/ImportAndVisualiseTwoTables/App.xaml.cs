using System.Configuration;
using System.Data;
using System.Windows;
using System.Text;

namespace ImportAndVisualiseTwoTables;

public partial class App : Application
{
    protected override void OnStartup(StartupEventArgs e)
    {
        base.OnStartup(e);
        Console.OutputEncoding = Encoding.UTF8;
    }
}