using System;
using System.Collections.Generic;
using System.Diagnostics;

// Требование: [DebuggerDisplay] над классом Point
[DebuggerDisplay("x={X}, f(x)={Y}, разрыв={IsDiscontinuous}")]
class Point
{
    public double X { get; set; }
    public double Y { get; set; }
    public bool IsDiscontinuous { get; set; }

    public Point(double x, double y, bool isDiscontinuous)
    {
        X = x;
        Y = y;
        IsDiscontinuous = isDiscontinuous;
    }
}

class ProgramHyperbola
{
    // Требование: TraceSource("Hyperbola")
    private static TraceSource trace = new TraceSource("Hyperbola");

    static void Main(string[] args)
    {
        Console.WriteLine("=== Программа: Гипербола f(x) = 1/x ===");

        // 1. Одиночное вычисление
        Console.Write("Введите x для расчета f(x): ");
        if (double.TryParse(Console.ReadLine(), out double x))
        {
            CalculateSingle(x);
        }
        else
        {
            Console.WriteLine("Вы ввели не число.");
        }

        // 2. Построение таблицы на интервале [a, b] с шагом h
        Console.WriteLine("\n--- Построение таблицы значений ---");
        Console.Write("Введите начало интервала (a): ");
        double.TryParse(Console.ReadLine(), out double a);

        Console.Write("Введите конец интервала (b): ");
        double.TryParse(Console.ReadLine(), out double b);

        Console.Write("Введите шаг (h): ");
        double.TryParse(Console.ReadLine(), out double h);

        if (h <= 0)
        {
            Console.WriteLine("Шаг должен быть больше 0!");
            return;
        }

        BuildTable(a, b, h);
    }

    static void CalculateSingle(double x)
    {
        const double epsilon = 1e-6;

        // Требование: Debug.Assert(x != 0, ...)
        Debug.Assert(x != 0, "Точка разрыва: деление на 0");

        // Требование: Debugger.Break() при |x| < 1e-9, если отладчик подключен
        if (Math.Abs(x) < 1e-9 && Debugger.IsAttached)
        {
            Debugger.Break();
        }

        if (Math.Abs(x) < epsilon)
        {
            if (x == 0)
            {
                trace.TraceEvent(TraceEventType.Error, 1, "Ошибка: x = 0 (деление на ноль)");
            }
            else
            {
                trace.TraceEvent(TraceEventType.Warning, 2, "Предупреждение: x в окрестности точки разрыва 0");
            }

            Console.WriteLine("Разрыв: значение не вычисляется (слишком близко к 0).");
            return;
        }

        double y = 1.0 / x;
        Debug.WriteLine($"x={x}, f={y}");
        Console.WriteLine($"Результат: f({x}) = {y}");
    }

    static void BuildTable(double a, double b, double h)
    {
        // Требование: Stopwatch для замера построения таблицы
        var stopwatch = new Stopwatch();
        var points = new List<Point>();
        const double epsilon = 1e-6;

        stopwatch.Start();

        for (double x = a; x <= b; x += h)
        {
            // Избегаем проблем с плавающей точкой, если попали ровно в 0
            if (Math.Abs(x) < 1e-12) x = 0;

            bool isDiscontinuous = Math.Abs(x) < epsilon;

            if (isDiscontinuous || x == 0)
            {
                if (x == 0)
                {
                    trace.TraceEvent(TraceEventType.Error, 1, $"Ошибка в таблице: x = 0 в точке {x}");
                    // В режиме Debug здесь сработает Assert при точном попадании в 0
                    Debug.Assert(x != 0, "Точка разрыва: деление на 0");
                }
                else
                {
                    trace.TraceEvent(TraceEventType.Warning, 2, $"Предупреждение в таблице: окрестность 0, x = {x}");
                }

                Console.WriteLine($"x={x}: разрыв");
                points.Add(new Point(x, double.NaN, true));
                continue;
            }

            double y = 1.0 / x;
            // Требование: Debug.WriteLine($"x={x}, f={y}")
            Debug.WriteLine($"x={x}, f={y}");
            trace.TraceEvent(TraceEventType.Verbose, 3, $"Точка построена: x={x}, f={y}");
            points.Add(new Point(x, y, false));
        }

        stopwatch.Stop();
        Console.WriteLine($"\nТаблица построена за {stopwatch.ElapsedMilliseconds} мс.");

        // Вывод итоговой таблицы на экран
        Console.WriteLine("\nСводная таблица:");
        foreach (var p in points)
        {
            string yStr = p.IsDiscontinuous ? "разрыв" : p.Y.ToString("F4");
            Console.WriteLine($"x = {p.X,8:F4} | f(x) = {yStr}");
        }
    }
}
