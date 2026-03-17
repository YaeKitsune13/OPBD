using System;
using System.Collections.Generic;

namespace ImportAndVisualiseTwoTablesAvalonia.Models;

public partial class Employee
{
    public int EmployeeId { get; set; }

    public string? FullName { get; set; }

    public string? Passport { get; set; }

    public string? Position { get; set; }

    public int? CountPolicyholders { get; set; }

    public virtual ICollection<Policyholder> Policyholders { get; set; } = new List<Policyholder>();
}
