using System;
using System.Collections.Generic;
using Microsoft.EntityFrameworkCore;

namespace ImportAndVisualiseTwoTables.ModelsPostgresSQL;

public partial class InsuranceContext : DbContext
{
    public InsuranceContext()
    {
    }

    public InsuranceContext(DbContextOptions<InsuranceContext> options)
        : base(options)
    {
    }

    public virtual DbSet<Claim> Claims { get; set; }

    public virtual DbSet<Employee> Employees { get; set; }

    public virtual DbSet<InsuranceType> InsuranceTypes { get; set; }

    public virtual DbSet<Policyholder> Policyholders { get; set; }

    protected override void OnConfiguring(DbContextOptionsBuilder optionsBuilder)
#warning To protect potentially sensitive information in your connection string, you should move it out of source code. You can avoid scaffolding the connection string by using the Name= syntax to read it from configuration - see https://go.microsoft.com/fwlink/?linkid=2131148. For more guidance on storing connection strings, see https://go.microsoft.com/fwlink/?LinkId=723263.
        => optionsBuilder.UseNpgsql("Server=localhost;username=postgres;password=1234;database=insurance");

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.Entity<Claim>(entity =>
        {
            entity.HasKey(e => e.ClaimId).HasName("claims_pkey");

            entity.ToTable("claims");

            entity.Property(e => e.ClaimId).HasColumnName("claim_id");
            entity.Property(e => e.Description).HasColumnName("description");
            entity.Property(e => e.EventDate).HasColumnName("event_date");
            entity.Property(e => e.Payout)
                .HasPrecision(12, 2)
                .HasColumnName("payout");
            entity.Property(e => e.PolicyNumber)
                .HasMaxLength(10)
                .IsFixedLength()
                .HasColumnName("policy_number");

            entity.HasOne(d => d.PolicyNumberNavigation).WithMany(p => p.Claims)
                .HasForeignKey(d => d.PolicyNumber)
                .HasConstraintName("claims_policy_number_fkey");
        });

        modelBuilder.Entity<Employee>(entity =>
        {
            entity.HasKey(e => e.EmployeeId).HasName("employees_pkey");

            entity.ToTable("employees");

            entity.Property(e => e.EmployeeId)
                .ValueGeneratedNever()
                .HasColumnName("employee_id");
            entity.Property(e => e.FullName)
                .HasMaxLength(100)
                .HasColumnName("full_name");
            entity.Property(e => e.Passport)
                .HasMaxLength(50)
                .HasColumnName("passport");
            entity.Property(e => e.Position)
                .HasMaxLength(60)
                .HasColumnName("position");
        });

        modelBuilder.Entity<InsuranceType>(entity =>
        {
            entity.HasKey(e => e.InsuranceTypeId).HasName("insurance_types_pkey");

            entity.ToTable("insurance_types");

            entity.Property(e => e.InsuranceTypeId)
                .ValueGeneratedNever()
                .HasColumnName("insurance_type_id");
            entity.Property(e => e.AnnualCost)
                .HasPrecision(12, 2)
                .HasColumnName("annual_cost");
            entity.Property(e => e.Description).HasColumnName("description");
            entity.Property(e => e.Name)
                .HasMaxLength(100)
                .HasColumnName("name");
        });

        modelBuilder.Entity<Policyholder>(entity =>
        {
            entity.HasKey(e => e.PolicyNumber).HasName("policyholders_pkey");

            entity.ToTable("policyholders");

            entity.Property(e => e.PolicyNumber)
                .HasMaxLength(10)
                .IsFixedLength()
                .HasColumnName("policy_number");
            entity.Property(e => e.BirthDate).HasColumnName("birth_date");
            entity.Property(e => e.ContractDate).HasColumnName("contract_date");
            entity.Property(e => e.EmployeeId).HasColumnName("employee_id");
            entity.Property(e => e.EndDate).HasColumnName("end_date");
            entity.Property(e => e.FullName)
                .HasMaxLength(40)
                .HasColumnName("full_name");
            entity.Property(e => e.InsuranceTypeId).HasColumnName("insurance_type_id");
            entity.Property(e => e.Passport)
                .HasMaxLength(50)
                .HasColumnName("passport");
            entity.Property(e => e.PolicyCost)
                .HasPrecision(8)
                .HasColumnName("policy_cost");
            entity.Property(e => e.PremiumAmount)
                .HasPrecision(8)
                .HasColumnName("premium_amount");

            entity.HasOne(d => d.Employee).WithMany(p => p.Policyholders)
                .HasForeignKey(d => d.EmployeeId)
                .HasConstraintName("policyholders_employee_id_fkey");

            entity.HasOne(d => d.InsuranceType).WithMany(p => p.Policyholders)
                .HasForeignKey(d => d.InsuranceTypeId)
                .HasConstraintName("policyholders_insurance_type_id_fkey");
        });

        OnModelCreatingPartial(modelBuilder);
    }

    partial void OnModelCreatingPartial(ModelBuilder modelBuilder);
}
