using System;
using System.Collections.Generic;
using Microsoft.EntityFrameworkCore;
using Pomelo.EntityFrameworkCore.MySql.Scaffolding.Internal;

namespace WPFDemoExamPractice.Models;

public partial class AppDbContext : DbContext
{
    public AppDbContext()
    {
    }

    public AppDbContext(DbContextOptions<AppDbContext> options)
        : base(options)
    {
    }

    public virtual DbSet<Address> Addresses { get; set; }

    public virtual DbSet<Order> Orders { get; set; }

    public virtual DbSet<OrderItem> OrderItems { get; set; }

    public virtual DbSet<Tovar> Tovars { get; set; }

    public virtual DbSet<User> Users { get; set; }

    protected override void OnConfiguring(DbContextOptionsBuilder optionsBuilder)
#warning To protect potentially sensitive information in your connection string, you should move it out of source code. You can avoid scaffolding the connection string by using the Name= syntax to read it from configuration - see https://go.microsoft.com/fwlink/?linkid=2131148. For more guidance on storing connection strings, see https://go.microsoft.com/fwlink/?LinkId=723263.
        => optionsBuilder.UseMySql("server=localhost;user=root;password=1234;database=app_db", Microsoft.EntityFrameworkCore.ServerVersion.Parse("9.5.0-mysql"));

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder
            .UseCollation("utf8mb4_0900_ai_ci")
            .HasCharSet("utf8mb4");

        modelBuilder.Entity<Address>(entity =>
        {
            entity.HasKey(e => e.Id).HasName("PRIMARY");

            entity.ToTable("address");

            entity.Property(e => e.Id).HasColumnName("id");
            entity.Property(e => e.City)
                .HasMaxLength(40)
                .HasColumnName("city");
            entity.Property(e => e.House).HasColumnName("house");
            entity.Property(e => e.Street)
                .HasMaxLength(50)
                .HasColumnName("street");
        });

        modelBuilder.Entity<Order>(entity =>
        {
            entity.HasKey(e => e.Id).HasName("PRIMARY");

            entity.ToTable("orders");

            entity.HasIndex(e => e.AdressId, "orders_address_id_fk");

            entity.HasIndex(e => e.UserId, "user_id_fk");

            entity.Property(e => e.Id).HasColumnName("id");
            entity.Property(e => e.AdressId).HasColumnName("adress_id");
            entity.Property(e => e.Code).HasColumnName("code");
            entity.Property(e => e.DateDelivary).HasColumnName("date_delivary");
            entity.Property(e => e.DateOrder).HasColumnName("date_order");
            entity.Property(e => e.State)
                .HasColumnType("enum('Завершен','Новый')")
                .HasColumnName("state");
            entity.Property(e => e.UserId).HasColumnName("user_id");

            entity.HasOne(d => d.Adress).WithMany(p => p.Orders)
                .HasForeignKey(d => d.AdressId)
                .HasConstraintName("orders_address_id_fk");

            entity.HasOne(d => d.User).WithMany(p => p.Orders)
                .HasForeignKey(d => d.UserId)
                .OnDelete(DeleteBehavior.ClientSetNull)
                .HasConstraintName("user_id_fk");
        });

        modelBuilder.Entity<OrderItem>(entity =>
        {
            entity
                .HasNoKey()
                .ToTable("order_items");

            entity.HasIndex(e => e.OrderId, "order_items_orders_id_fk");

            entity.HasIndex(e => e.Arcticle, "order_items_tovar_article_fk");

            entity.Property(e => e.Arcticle)
                .HasMaxLength(6)
                .HasColumnName("arcticle");
            entity.Property(e => e.Count).HasColumnName("count");
            entity.Property(e => e.OrderId).HasColumnName("order_id");

            entity.HasOne(d => d.ArcticleNavigation).WithMany()
                .HasForeignKey(d => d.Arcticle)
                .OnDelete(DeleteBehavior.ClientSetNull)
                .HasConstraintName("order_items_tovar_article_fk");

            entity.HasOne(d => d.Order).WithMany()
                .HasForeignKey(d => d.OrderId)
                .HasConstraintName("order_items_orders_id_fk");
        });

        modelBuilder.Entity<Tovar>(entity =>
        {
            entity.HasKey(e => e.Article).HasName("PRIMARY");

            entity.ToTable("tovar");

            entity.Property(e => e.Article)
                .HasMaxLength(6)
                .HasColumnName("article");
            entity.Property(e => e.Category)
                .HasColumnType("enum('Женская обувь','Мужская обувь')")
                .HasColumnName("category");
            entity.Property(e => e.Count).HasColumnName("count");
            entity.Property(e => e.Creator)
                .HasMaxLength(40)
                .HasColumnName("creator");
            entity.Property(e => e.DeliveryPerson)
                .HasMaxLength(40)
                .HasColumnName("delivery_person");
            entity.Property(e => e.Description)
                .HasMaxLength(200)
                .HasColumnName("description");
            entity.Property(e => e.Discount).HasColumnName("discount");
            entity.Property(e => e.Name)
                .HasMaxLength(40)
                .HasColumnName("name");
            entity.Property(e => e.Photo)
                .HasMaxLength(30)
                .HasColumnName("photo");
            entity.Property(e => e.Price).HasColumnName("price");
        });

        modelBuilder.Entity<User>(entity =>
        {
            entity.HasKey(e => e.Id).HasName("PRIMARY");

            entity.ToTable("users");

            entity.Property(e => e.Id).HasColumnName("id");
            entity.Property(e => e.FirstName)
                .HasMaxLength(30)
                .HasColumnName("first_name");
            entity.Property(e => e.LastName)
                .HasMaxLength(30)
                .HasColumnName("last_name");
            entity.Property(e => e.Login)
                .HasMaxLength(40)
                .HasColumnName("login");
            entity.Property(e => e.Name)
                .HasMaxLength(30)
                .HasColumnName("name");
            entity.Property(e => e.Password)
                .HasMaxLength(20)
                .HasColumnName("password");
            entity.Property(e => e.Rule)
                .HasColumnType("enum('Администратор','Менеджер','Авторизированный клиент')")
                .HasColumnName("rule");
        });

        OnModelCreatingPartial(modelBuilder);
    }

    partial void OnModelCreatingPartial(ModelBuilder modelBuilder);
}
