using System.Windows;
using System.Windows.Controls;
using ImportAndVisualiseTwoTables.Models;
using ImportAndVisualiseTwoTables.ModelsPostgresSQL;
using Microsoft.EntityFrameworkCore;

namespace ImportAndVisualiseTwoTables;

public partial class MainWindow : Window
{
    private readonly Models.InsuranceContext dbContext;
    private readonly ModelsPostgresSQL.InsuranceContext dbContext2;

    public MainWindow()
    {
        InitializeComponent();
        dbContext = new Models.InsuranceContext();
        dbContext2 = new ModelsPostgresSQL.InsuranceContext();
        LoadEmployeesDataGridMySQL();
        LoadEmployeesDataGridPostgreSQL();
    }

    private void LoadEmployeesDataGridMySQL()
    {
        EmployeesDataGrid.ItemsSource = dbContext.Employees
            .FromSqlRaw("SELECT * FROM employees")
            .ToList();
        EmployeesDataGrid.SelectionChanged += EmployeesDataGrid_SelectionChanged;
    }

    private void LoadEmployeesDataGridPostgreSQL()
    {
        EmployeesPostgreDataGrid.ItemsSource = dbContext2.Employees
            .FromSqlRaw("SELECT * FROM employees")
            .ToList();
        EmployeesPostgreDataGrid.SelectionChanged += EmployeesPostgreDataGrid_SelectionChanged;
    }

    private void EmployeesDataGrid_SelectionChanged(object sender, SelectionChangedEventArgs e)
    {
        if (EmployeesDataGrid.SelectedItem is Models.Employees selectedEmployee)
            ClaimsIncludedDataGrid.ItemsSource = dbContext.Claims
                .Include(c => c.PolicyNumberNavigation)
                .Where(c => c.PolicyNumberNavigation!.EmployeeId == selectedEmployee.EmployeeId)
                .AsNoTracking()
                .ToList();
    }

    private void EmployeesPostgreDataGrid_SelectionChanged(object sender, SelectionChangedEventArgs e)
    {
        if (EmployeesPostgreDataGrid.SelectedItem is ModelsPostgresSQL.Employee selectedEmployee)
            ClaimsPostgreIncludedDataGrid.ItemsSource = dbContext2.Claims
                .Include(c => c.PolicyNumberNavigation)
                .Where(c => c.PolicyNumberNavigation!.EmployeeId == selectedEmployee.EmployeeId)
                .AsNoTracking()
                .ToList();
    }
}