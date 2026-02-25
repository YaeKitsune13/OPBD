using System;
using System.Collections.Generic;

namespace ImportAndVisualiseTwoTables.Models;

public partial class InsuranceTypes
{
    public short InsuranceTypeId { get; set; }

    public string? Name { get; set; }

    public string? Description { get; set; }

    public decimal? AnnualCost { get; set; }

    public virtual ICollection<Policyholders> Policyholders { get; set; } = new List<Policyholders>();
}
