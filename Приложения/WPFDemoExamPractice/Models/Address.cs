using System;
using System.Collections.Generic;

namespace WPFDemoExamPractice.Models;

public partial class Address
{
    public int Id { get; set; }

    public string City { get; set; } = null!;

    public string Street { get; set; } = null!;

    public int? House { get; set; }

    public virtual ICollection<Order> Orders { get; set; } = new List<Order>();
}
