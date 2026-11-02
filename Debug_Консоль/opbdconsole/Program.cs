using System;
using System.Diagnostics;

[DebuggerDisplay("x={X}, ε={Epsilon}, N={Iterations}")]
class TaylorResult
{
    public double X, Epsilon, Sum;
    public int Iterations;
}

class Program
{
    static void Main()
    {
        var trace = new TraceSource("SinTaylor");
        trace.Switch.Level = SourceLevels.All;
        trace.Listeners.Add(new ConsoleTraceListener());

        Console.Write("Введите x (радианы): ");
        double x = double.Parse(Console.ReadLine());
        Console.Write("Введите epsilon: ");
        double epsilon = double.Parse(Console.ReadLine());

        if (double.IsNaN(x))
        {
            trace.TraceEvent(TraceEventType.Error, 3, "x = NaN");
            Console.WriteLine("x не является числом (NaN).");
            return;
        }

        Debug.Assert(Math.Abs(x) <= 2 * Math.PI, "x вне разумного диапазона");
        Debug.Assert(epsilon > 0 && epsilon < 1, "epsilon вне разумного диапазона");

        // Вложенная функция: факториал (в Debug-логе с отступом)
        double Factorial(int k)
        {
            Debug.Indent();
            double f = 1;
            for (int i = 2; i <= k; i++) f *= i;
            Debug.WriteLine($"{k}! = {f}");
            Debug.Unindent();
            return f;
        }

        var r = new TaylorResult();
        r.X = x;
        r.Epsilon = epsilon;

        var timer = Stopwatch.StartNew();

        int n = 0;
        double term = x;
        while (Math.Abs(term) > epsilon && n < 100)
        {
            r.Sum += term;
            Debug.WriteLine($"n={n}, term={term}, sum={r.Sum}");
            trace.TraceEvent(TraceEventType.Verbose, 1, $"n={n}, term={term}");

            n++;
            term = Math.Pow(-1, n) * Math.Pow(x, 2 * n + 1) / Factorial(2 * n + 1);
        }

        if (Math.Abs(term) > epsilon)
            trace.TraceEvent(TraceEventType.Warning, 2, "100 итераций, а точность не достигнута");

        r.Iterations = n;
        timer.Stop();

        double real = Math.Sin(x);
        Console.WriteLine($"Ряд Тейлора: {r.Sum}");
        Console.WriteLine($"Math.Sin:    {real}");
        Console.WriteLine($"Разница:     {Math.Abs(r.Sum - real)}");
        Console.WriteLine($"Итераций: {r.Iterations}, время: {timer.ElapsedMilliseconds} мс");
    }
}
