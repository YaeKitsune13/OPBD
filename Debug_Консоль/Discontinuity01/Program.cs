using System;
using System.Collections.Generic;
using System.Diagnostics;

[DebuggerDisplay("x={X}, f(x)={Y}, разрыв={IsDiscontinuous}")]
class Point
{
    public double X, Y;
    public bool IsDiscontinuous;
}

class Program
{
    static void Main()
    {
        var trace = new TraceSource("Hyperbola");
        trace.Switch.Level = SourceLevels.All;
        trace.Listeners.Add(new ConsoleTraceListener());
        double eps = 1e-6;

        Console.Write("Введите x: ");
        double x = double.Parse(Console.ReadLine());

        Debug.Assert(x != 0, "Точка разрыва: деление на 0");

        if (Math.Abs(x) < eps)
        {
            if (Math.Abs(x) < 1e-9 && Debugger.IsAttached) Debugger.Break();

            if (x == 0) trace.TraceEvent(TraceEventType.Error, 1, "Error x = 0");
            else trace.TraceEvent(TraceEventType.Warning, 2, "Warning окрестность 0, x = " + x);

            Console.WriteLine("значение не вычисляется.");
        }
        else
        {
            Console.WriteLine($"f({x}) = {1 / x}");
        }

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
        {
            var p = new Point();
            p.X = t;

            if (Math.Abs(t) < eps)
            {
                if (Math.Abs(t) < 1e-9 && Debugger.IsAttached) Debugger.Break();

                if (t == 0) trace.TraceEvent(TraceEventType.Error, 1, "Error x = 0");
                else trace.TraceEvent(TraceEventType.Warning, 2, "Warring окрестность 0, x = " + t);

                p.Y = double.NaN;
                p.IsDiscontinuous = true;
                Console.WriteLine($"x={t}: разрыв, точка пропущена");
            }
            else
            {
                p.Y = 1 / t;
                Debug.WriteLine($"x={t}, f={p.Y}");
            }

            points.Add(p);
        }

        timer.Stop();
        Console.WriteLine($"Таблица построена за {timer.ElapsedMilliseconds} мс");

        foreach (var p in points)
        {
            string y = p.IsDiscontinuous ? "разрыв" : p.Y.ToString("F4");
            Console.WriteLine($"x = {p.X,8:F4} | f(x) = {y}");
        }
    }
}
