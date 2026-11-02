using System;
using System.Collections.Generic;
using System.Diagnostics;

[DebuggerDisplay("x={X}, f(x)={Y}, разрыв={IsDiscontinuous}")]
class Point
{
    public double X, Y;
    public bool IsDiscontinuous => double.IsNaN(Y);

    public Point(double x, double y)
    {
        X = x;
        Y = y;
    }
}

class ProgramHyperbola
{
    const double Epsilon = 1e-6;
    static TraceSource trace = new TraceSource("Hyperbola");

    static void Main()
    {
        trace.Switch.Level = SourceLevels.All;
        trace.Listeners.Add(new ConsoleTraceListener());

        Console.Write("Введите x: ");
        double x = double.Parse(Console.ReadLine());
        double y = F(x);
        Console.WriteLine(double.IsNaN(y) ? "Разрыв" : $"f({x}) = {y}");

        Console.Write("Введите a: ");
        double a = double.Parse(Console.ReadLine());
        Console.Write("Введите b: ");
        double b = double.Parse(Console.ReadLine());
        Console.Write("Введите шаг h: ");
        double h = double.Parse(Console.ReadLine());

        if (h <= 0)
        {
            Console.WriteLine("Шаг должен быть больше 0!");
            return;
        }

        var points = new List<Point>();
        var timer = Stopwatch.StartNew();
        for (double t = a; t <= b; t += h)
            points.Add(new Point(t, F(t)));
        timer.Stop();

        Console.WriteLine($"Таблица построена за {timer.ElapsedMilliseconds} мс");
        foreach (var p in points)
        {
            string yStr = p.IsDiscontinuous ? "разрыв" : p.Y.ToString("F4");
            Console.WriteLine($"x = {p.X,8:F4} | f(x) = {yStr}");
        }
    }

    static double F(double x)
    {
        Debug.Assert(x != 0, "Точка разрыва: деление на 0");

        if (Math.Abs(x) < Epsilon)
        {
            if (Debugger.IsAttached) Debugger.Break();
            var level = x == 0 ? TraceEventType.Error : TraceEventType.Warning;
            trace.TraceEvent(level, 1, $"Разрыв: x = {x}");
            return double.NaN;
        }

        double y = 1.0 / x;
        Debug.WriteLine($"x={x}, f={y}");
        trace.TraceEvent(TraceEventType.Verbose, 2, $"x={x}, f={y}");
        return y;
    }
}
