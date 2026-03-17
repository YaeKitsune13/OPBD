using System;
using System.Collections.Generic;

namespace ImportAndVisualiseTwoTablesAvalonia.Models;

public partial class Claim
{
    public long ClaimId { get; set; }

    public string? PolicyNumber { get; set; }

    public string? Description { get; set; }

    public DateOnly? EventDate { get; set; }

    public decimal? Payout { get; set; }

    public virtual Policyholder? PolicyNumberNavigation { get; set; }
}
