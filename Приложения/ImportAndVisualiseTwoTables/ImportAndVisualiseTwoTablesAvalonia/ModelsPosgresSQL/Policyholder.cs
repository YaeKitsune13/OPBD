using System;
using System.Collections.Generic;

namespace ImportAndVisualiseTwoTablesAvalonia.ModelsPosgresSQL;

public partial class Policyholder
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

    public virtual ICollection<Claim> Claims { get; set; } = new List<Claim>();

    public virtual Employee? Employee { get; set; }

    public virtual InsuranceType? InsuranceType { get; set; }
}
