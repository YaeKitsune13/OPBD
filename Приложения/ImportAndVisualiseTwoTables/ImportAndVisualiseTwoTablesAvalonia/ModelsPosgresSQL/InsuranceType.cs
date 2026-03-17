using System;
using System.Collections.Generic;

namespace ImportAndVisualiseTwoTablesAvalonia.ModelsPosgresSQL;

public partial class InsuranceType
{
    public short InsuranceTypeId { get; set; }

    public string? Name { get; set; }

    public string? Description { get; set; }

    public decimal? AnnualCost { get; set; }

    public virtual ICollection<Policyholder> Policyholders { get; set; } = new List<Policyholder>();
}
