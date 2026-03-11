using System.Windows;
using System.Windows.Controls;
using ImportAndVisualiseTwoTables.Models;
using Microsoft.EntityFrameworkCore;

namespace ImportAndVisualiseTwoTables;

/// <summary>
///     Interaction logic for MainWindow.xaml
/// </summary>
public partial class MainWindow : Window
{
    private readonly InsuranceContext dbContext;
    private readonly ModelsPostgresSQL.InsuranceContext dbContext2;
    private List<int> employeeIds = new List<int>();
    public MainWindow()
    {
        InitializeComponent();
        dbContext = new InsuranceContext();
        dbContext2 = new ModelsPostgresSQL.InsuranceContext();
        LoadEmployeesDataGrid();
        LoadClaimsDataGrid();
        LoadEmployersComboBox();
        LoadPolicyHoldersDataGrid();
    }

    private void LoadEmployeesDataGrid()
    {
        var dbEmployees = dbContext.Employees
            .FromSqlRaw("SELECT * FROM employees")
            .ToList();
    
        EmployeesDataGrid.ItemsSource = dbEmployees;
    }

    private void LoadClaimsDataGrid()
    {
        ClaimsIncludedDataGrid.ItemsSource = dbContext.Claims
            .Include(c => c.PolicyNumberNavigation)
            .AsNoTracking()
            .ToList();
    }

    private void LoadEmployersComboBox()
    {
        employeeIds = dbContext.Employees.Select(e => e.EmployeeId).ToList();
        DatabaseComboBox.ItemsSource = employeeIds;

        DatabaseComboBox.SelectionChanged += (s, e) =>
        {
            if (DatabaseComboBox.SelectedItem is int selectedEmployeeId)
            {
                PolicyHoldersDataGrid.ItemsSource = dbContext2.Policyholders
                    .Where(p => p.EmployeeId == selectedEmployeeId)
                    .AsNoTracking()
                    .ToList();
            }
        };
    }

    private void LoadPolicyHoldersDataGrid()
    {
        PolicyHoldersDataGrid.ItemsSource = dbContext2.Policyholders
            .AsNoTracking()
            .ToList();
    }

    private void LoadButton_Click(object sender, RoutedEventArgs e)
    {
        LoadPolicyHoldersDataGrid();
    }

}