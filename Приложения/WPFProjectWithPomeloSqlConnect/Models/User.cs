using System;
using System.Collections.Generic;

namespace WPFProjectWithPomeloSqlConnect.Models;

public partial class User
{
    public int Id { get; set; }

    public string Login { get; set; } = null!;

    public string Password { get; set; } = null!;

    public string? Email { get; set; }

    public DateTime? CreatedAt { get; set; }
}
