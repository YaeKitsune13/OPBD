using System;
using System.Collections.Generic;

namespace WPFDemoExamPractice.Models;

public partial class OrderItem
{
    public int? OrderId { get; set; }

    public string Arcticle { get; set; } = null!;

    public int? Count { get; set; }

    public virtual Tovar ArcticleNavigation { get; set; } = null!;

    public virtual Order? Order { get; set; }
}
