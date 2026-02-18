using System;
using System.Collections.Generic;

namespace WPFDemoExamPractice.Models;

public partial class User
{
    public int Id { get; set; }

    public string Rule { get; set; } = null!;

    public string FirstName { get; set; } = null!;

    public string Name { get; set; } = null!;

    public string LastName { get; set; } = null!;

    public string Login { get; set; } = null!;

    public string Password { get; set; } = null!;

    public virtual ICollection<Order> Orders { get; set; } = new List<Order>();
}
