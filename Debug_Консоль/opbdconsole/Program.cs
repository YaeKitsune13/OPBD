using System;
using System.Diagnostics;

[DebuggerDisplay("x={X}, ε={Epsilon}, N={Iterations}")]
class Program
{
    static TraceSource trace = new TraceSource("TaylorSin");
    static double X, Epsilon;
    static int Iterations;

    static void Main()
    {
        trace.Switch.Level = SourceLevels.All;
        trace.Listeners.Add(new ConsoleTraceListener());

        Console.Write("Введите x: ");
        X = double.Parse(Console.ReadLine());
        Console.Write("Введите epsilon: ");
        Epsilon = double.Parse(Console.ReadLine());

        Debug.Assert(Math.Abs(X) <= 2 * Math.PI, "x вне разумного предела");
        Debug.Assert(Epsilon > 0 && Epsilon < 1, "epsilon вне разумного предела");

        double term = X, sum = X;
        var timer = Stopwatch.StartNew();

        Debug.WriteLine("Начало итераций:");
        Debug.Indent();
        while (Math.Abs(term) >= Epsilon)
        {
            Iterations++;
            term = -term * X * X / (2 * Iterations) / (2 * Iterations + 1);
            sum += term;
            Debug.WriteLine($"Итерация={Iterations}, term={term}, sum={sum}");
            trace.TraceEvent(TraceEventType.Verbose, 1, $"Итерация={Iterations}, term={term}");
            if (Iterations == 100)
                trace.TraceEvent(TraceEventType.Warning, 2, "Больше 100 итераций");
        }
        Debug.Unindent();

        timer.Stop();
        Debug.WriteLine($"Конец итераций, время: {timer.ElapsedMilliseconds} мс");
        Console.WriteLine($"Ряд: {sum}, Math.Sin: {Math.Sin(X)}, итераций: {Iterations}");
    }
}
