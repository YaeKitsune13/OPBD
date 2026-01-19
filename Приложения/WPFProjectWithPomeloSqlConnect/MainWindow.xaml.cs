using System.Text;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Navigation;
using System.Windows.Shapes;
using WPFProjectWithPomeloSqlConnect.Models;

namespace WPFProjectWithPomeloSqlConnect;

/// <summary>
/// Interaction logic for MainWindow.xaml
/// </summary>
public partial class MainWindow : Window
{
    public AppDbContext appDb;
    public MainWindow()
    {
        InitializeComponent();
        appDb = new AppDbContext();
        UpdateUsersDataGrid();
    }

    public void UpdateUsersDataGrid()
    {
        var users = appDb.Users.ToList();
        UsersDataGrid.ItemsSource = users;
        MessageBox.Show(users.Count.ToString());
    }
}