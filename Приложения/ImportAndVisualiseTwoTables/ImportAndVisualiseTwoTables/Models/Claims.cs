using System;
using System.Collections.Generic;

namespace ImportAndVisualiseTwoTables.Models;

public partial class Claims
{
    public long ClaimId { get; set; }

    public string? PolicyNumber { get; set; }

    public string? Description { get; set; }

    public DateOnly? EventDate { get; set; }

    public decimal? Payout { get; set; }

    public virtual Policyholders? PolicyNumberNavigation { get; set; }
}
