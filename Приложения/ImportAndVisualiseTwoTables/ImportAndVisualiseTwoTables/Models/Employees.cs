using System;
using System.Collections.Generic;

namespace ImportAndVisualiseTwoTables.Models;

public partial class Employees
{
    public int EmployeeId { get; set; }

    public string? FullName { get; set; }

    public string? Passport { get; set; }

    public string? Position { get; set; }

    public int CountPolicyholders { get; set; }

    public virtual ICollection<Policyholders> Policyholders { get; set; } = new List<Policyholders>();
}
