using System.Windows;
using System.Windows.Controls;
using ImportAndVisualiseTwoTables.Models;
using Microsoft.EntityFrameworkCore;
using ImportAndVisualiseTwoTables.ModelsPostgresSQL;

namespace ImportAndVisualiseTwoTables;

/// <summary>
///     Interaction logic for MainWindow.xaml
/// </summary>
public partial class MainWindow : Window
{
    private readonly Models.InsuranceContext dbContext;
    private readonly ModelsPostgresSQL.InsuranceContext dbContext2;
    private List<int> employeeIds = new List<int>();
    public MainWindow()
    {
        InitializeComponent();
        dbContext = new Models.InsuranceContext();
        dbContext2 = new ModelsPostgresSQL.InsuranceContext();
        LoadEmployeesDataGrid();
        LoadClaimsDataGrid();
        // LoadEmployersComboBox();
        // LoadPolicyHoldersDataGrid();
        // LoadInsuranceTypes();
        // LoadInsuranceClaimsByType();
    }

    private void LoadEmployeesDataGrid()
    {
        var dbEmployees = dbContext.Employees
            .FromSqlRaw("SELECT * FROM employees")
            .ToList();

        EmployeesDataGrid.ItemsSource = dbEmployees;
        EmployeesDataGrid.SelectionChanged += EmployeesDataGrid_SelectionChanged;
    }

    private void EmployeesDataGrid_SelectionChanged(object sender, SelectionChangedEventArgs e)
    {
        if (EmployeesDataGrid.SelectedItem is Models.Employees selectedEmployee)
        {
            ClaimsIncludedDataGrid.ItemsSource = dbContext.Claims
                .Include(c => c.PolicyNumberNavigation)
                .Where(c => c.PolicyNumberNavigation!.EmployeeId == selectedEmployee.EmployeeId)
                .AsNoTracking()
                .ToList();
        }
    }

    private void LoadClaimsDataGrid()
    {
        ClaimsIncludedDataGrid.ItemsSource = dbContext.Claims
            .Include(c => c.PolicyNumberNavigation)
            .AsNoTracking()
            .ToList();
    }

   
}