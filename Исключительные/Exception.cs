using System;

namespace Исключительные;

public class CsvParseException : Exception
{
    public int LineNumber { get; }
    public string RawLine { get; }

    public CsvParseException(string message, int lineNumber, string rawLine)
        : base(message)
    {
        LineNumber = lineNumber;
        RawLine = rawLine;
    }

    public CsvParseException(string message, int lineNumber, string rawLine, Exception innerException)
        : base(message, innerException) // Добавлен base
    {
        LineNumber = lineNumber;
        RawLine = rawLine;
    }
}
