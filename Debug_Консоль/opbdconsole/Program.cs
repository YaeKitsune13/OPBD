using System.Diagnostics;
[DebuggerDisplay("x={X}, ε={Epsilon}, N={Iterations}")]
class Program
{
    public static TraceSource trace = new TraceSource("TeilorSin");

    static void Main()
    {
        bool success;

        Console.Write("Введите x: ");
        success = double.TryParse(Console.ReadLine(), out double x);
        if (!success)
        {
            Console.WriteLine("\nВы ввели не число");
        }
        Console.Write("Введите epsilon: ");
        success = double.TryParse(Console.ReadLine(), out double epsilon);
        if (!success)
        {
            Console.WriteLine("\nВы ввели не число");
        }

        double teilor = TeilorSin(x, epsilon);
        if (Math.Sin(x) >= teilor + epsilon || Math.Sin(x) <= teilor - epsilon)
            Console.WriteLine($"Точность получилась такой себе: {teilor} когда sin {Math.Sin(x)}");
        else Console.WriteLine($"Точность получилась такой: {teilor} когда sin {Math.Sin(x)}");
    }


    static double TeilorSin(double x, double epsilon = 1e-6)
    {
        double term = x, sum = x;
        var timer = new Stopwatch();
        bool success = false;
        int iterations = 0;
        do
        {
            Console.WriteLine("Введите кол-во итераций");
            success = int.TryParse(Console.ReadLine(), out iterations);
            if (!success)
            {
                Console.WriteLine("Вы ввели не число");
            }
        }
        while (!success);
        Debug.WriteLine("Список ошибок");
        Debug.Assert(Math.Abs(x) <= 2 * Math.PI, "x вне разумного предела");
        Debug.Assert(epsilon > 0 && epsilon < 1, "epsilon вне разумного предела");
        Debug.WriteLine("Конец списка");

        Debug.WriteLine("Начало итераций:");
        Debug.Indent();

        timer.Start();
        for (int iter = 0; iter < iterations; iter++)
        {
            term = -term * x * x / (2 * iter + 2) / (2 * iter + 3);
            Debug.WriteLine($"Итерация={iter}, term={term}, sum={sum}");
            sum += term;
            trace.TraceEvent(TraceEventType.Verbose, 2, $"Итерация={iter}, term={term}, sum={sum}");
            if (iter >= 100 && (Math.Sin(x) >= sum + epsilon || Math.Sin(x) <= sum - epsilon))
            {
                trace.TraceEvent(TraceEventType.Warning, 2, "Больше 100 итераций всё ещё точно");
            }
            if (term == 0 || term == -0)
            {
                Debug.WriteLine($"На итераций: {iter} мы пришли к 0");
                break;
            }
            ;
        }
        timer.Stop();
        TimeSpan timeTaken = timer.Elapsed;
        Debug.Unindent();
        Debug.WriteLine($"Конец итераций, выполнено за: {timeTaken.ToString(@"m\:ss\.ff")}");
        trace.Close();
        return sum;
    }
}
