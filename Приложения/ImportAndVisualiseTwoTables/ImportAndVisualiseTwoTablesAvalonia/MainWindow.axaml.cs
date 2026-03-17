using System;
using System.Linq;
using Avalonia.Controls;
using ImportAndVisualiseTwoTablesAvalonia.Models;
using ImportAndVisualiseTwoTablesAvalonia.ModelsPosgresSQL;
using Microsoft.EntityFrameworkCore;

namespace ImportAndVisualiseTwoTablesAvalonia;

public partial class MainWindow : Window
{
    private readonly Models.InsuranceContext dbContext;
    private readonly ModelsPosgresSQL.InsuranceContext dbContext2;
    
    public MainWindow()
    {
        InitializeComponent();
        dbContext = new Models.InsuranceContext();
        dbContext2 = new ModelsPosgresSQL.InsuranceContext();
        
        this.Opened += (s,e) => {
            LoadData();
        };
    }
    
    private void LoadData()
    {
        // MySQL Employees
        var mysqlEmp = dbContext.Employees
            .FromSqlRaw("SELECT * FROM employees")
            .ToList();
        MySQLEmployeesListBox.ItemsSource = mysqlEmp;
        MySQLEmployeesListBox.SelectionChanged += MySQLEmployees_SelectionChanged;
        
        // PostgreSQL Employees
        var pgEmp = dbContext2.Employees
            .FromSqlRaw("SELECT * FROM employees")
            .ToList();
        PostgreSQLEmployeesListBox.ItemsSource = pgEmp;
        PostgreSQLEmployeesListBox.SelectionChanged += PostgreSQLEmployees_SelectionChanged;
    }
    
    private void MySQLEmployees_SelectionChanged(object? sender, SelectionChangedEventArgs e)
    {
        if (MySQLEmployeesListBox.SelectedItem is Models.Employee selectedEmployee)
        {
            // LINQ запрос
            var claimsLINQ = dbContext.Claims
                .Include(c => c.PolicyNumberNavigation)
                .Where(c => c.PolicyNumberNavigation!.EmployeeId == selectedEmployee.EmployeeId)
                .AsNoTracking()
                .ToList();
            MySQLClaimsLINQ.ItemsSource = claimsLINQ;
            
            // SQL запрос
            var claimsSQL = dbContext.Claims
                .FromSqlRaw("SELECT c.* FROM claims c JOIN policyholders p ON c.policy_number = p.policy_number WHERE p.employee_id = {0}", selectedEmployee.EmployeeId)
                .Include(c => c.PolicyNumberNavigation)
                .ToList();
            MySQLClaimsSQL.ItemsSource = claimsSQL;
        }
    }
    
    private void PostgreSQLEmployees_SelectionChanged(object? sender, SelectionChangedEventArgs e)
    {
        if (PostgreSQLEmployeesListBox.SelectedItem is ModelsPosgresSQL.Employee selectedEmployee)
        {
            // LINQ запрос
            var claimsLINQ = dbContext2.Claims
                .Include(c => c.PolicyNumberNavigation)
                .Where(c => c.PolicyNumberNavigation!.EmployeeId == selectedEmployee.EmployeeId)
                .AsNoTracking()
                .ToList();
            PostgreSQLClaimsLINQ.ItemsSource = claimsLINQ;
            
            // SQL запрос
            var claimsSQL = dbContext2.Claims
                .FromSqlRaw("SELECT c.* FROM claims c JOIN policyholders p ON c.policy_number = p.policy_number WHERE p.employee_id = {0}", selectedEmployee.EmployeeId)
                .Include(c => c.PolicyNumberNavigation)
                .ToList();
            PostgreSQLClaimsSQL.ItemsSource = claimsSQL;
        }
    }
}