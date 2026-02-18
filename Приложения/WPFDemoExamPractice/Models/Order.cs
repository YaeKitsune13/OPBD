using System;
using System.Collections.Generic;

namespace WPFDemoExamPractice.Models;

public partial class Order
{
    public int Id { get; set; }

    public DateOnly DateOrder { get; set; }

    public DateOnly DateDelivary { get; set; }

    public int? AdressId { get; set; }

    public int UserId { get; set; }

    public int Code { get; set; }

    public string State { get; set; } = null!;

    public virtual Address? Adress { get; set; }

    public virtual User User { get; set; } = null!;
}
