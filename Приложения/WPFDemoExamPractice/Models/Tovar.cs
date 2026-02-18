using System;
using System.Collections.Generic;

namespace WPFDemoExamPractice.Models;

public partial class Tovar
{
    public string Article { get; set; } = null!;

    public string Name { get; set; } = null!;

    public int? Price { get; set; }

    public string DeliveryPerson { get; set; } = null!;

    public string Creator { get; set; } = null!;

    public string Category { get; set; } = null!;

    public int? Discount { get; set; }

    public int? Count { get; set; }

    public string Description { get; set; } = null!;

    public string? Photo { get; set; }
}
