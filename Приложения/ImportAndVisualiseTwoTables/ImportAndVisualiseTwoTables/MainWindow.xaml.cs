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

    public MainWindow()
    {
        InitializeComponent();
        dbContext = new InsuranceContext();
        LoadEmployeesDataGrid();
        LoadClaimsDataGrid();
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
}