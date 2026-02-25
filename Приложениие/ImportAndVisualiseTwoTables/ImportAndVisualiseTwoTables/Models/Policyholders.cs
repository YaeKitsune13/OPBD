using System;
using System.Collections.Generic;

namespace ImportAndVisualiseTwoTables.Models;

public partial class Policyholders
{
    public string PolicyNumber { get; set; } = null!;

    public string? Passport { get; set; }

    public string? FullName { get; set; }

    public DateOnly? BirthDate { get; set; }

    public short? InsuranceTypeId { get; set; }

    public int? EmployeeId { get; set; }

    public DateOnly? ContractDate { get; set; }

    public DateOnly? EndDate { get; set; }

    public decimal? PremiumAmount { get; set; }

    public decimal? PolicyCost { get; set; }

    public virtual ICollection<Claims> Claims { get; set; } = new List<Claims>();

    public virtual Employees? Employee { get; set; }

    public virtual InsuranceTypes? InsuranceType { get; set; }
}