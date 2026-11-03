using Avalonia.Controls;
using Avalonia.Interactivity;
using Avalonia.Platform.Storage;
using System;
using System.Collections.Generic;
using System.Globalization;
using System.IO;
using System.Threading.Tasks;

namespace Исключительные;

public partial class MainWindow : Window
{
    public MainWindow()
    {
        InitializeComponent();
    }

    private async void BtnOpenCsv_Click(object sender, RoutedEventArgs e)
    {
        try
        {
            if (StorageProvider == null)
            {
                await ShowFallbackDialog();
                return;
            }

            var files = await StorageProvider.OpenFilePickerAsync(new FilePickerOpenOptions
            {
                Title = "Открыть CSV файл",
                AllowMultiple = false,
                FileTypeFilter = new[]
                {
                    new FilePickerFileType("CSV файлы") { Patterns = new[] { "*.csv" } },
                    new FilePickerFileType("Все файлы") { Patterns = new[] { "*.*" } }
                }
            });

            if (files != null && files.Count > 0)
            {
                string filePath = files[0].Path.LocalPath;
                await ImportCsvFileAsync(filePath);
            }
            else
            {
                TxtStatus.Text = "Выбор файла был отменен";
            }
        }
        catch (Exception ex)
        {
            TxtStatus.Text = "Ошибка системного диалога. Пробуем загрузить data.csv...";
            Console.WriteLine($"[StorageProvider Error]: {ex.Message}");
            await ShowFallbackDialog();
        }
    }

    private async Task ShowFallbackDialog()
    {
        string defaultPath = "data.csv";
        if (File.Exists(defaultPath))
        {
            TxtStatus.Text = "Используется файл по умолчанию (data.csv)";
            await ImportCsvFileAsync(defaultPath);
        }
        else
        {
            TxtStatus.Text = "Поместите data.csv в папку проекта или почините xdg-desktop-portal.";
        }
    }

    private async Task ImportCsvFileAsync(string filePath)
    {
        StreamReader? reader = null;
        var products = new List<Product>();

        try
        {
            ProgressBarImport.IsVisible = true;
            TxtStatus.Text = "Импорт файла...";

            await Task.Delay(400);

            reader = new StreamReader(filePath);
            int lineNumber = 0;

            while (!reader.EndOfStream)
            {
                string? line = await reader.ReadLineAsync();
                lineNumber++;

                if (string.IsNullOrWhiteSpace(line))
                    continue;

                string[] parts = line.Split(';');

                if (parts.Length != 3)
                {
                    throw new CsvParseException("Неверный формат строки (ожидается 3 поля через ';')", lineNumber, line);
                }

                if (!decimal.TryParse(parts[1].Trim(), NumberStyles.Any, CultureInfo.InvariantCulture, out decimal price) &&
                    !decimal.TryParse(parts[1].Trim(), out price))
                {
                    throw new CsvParseException($"Ошибка преобразования цены: '{parts[1]}'", lineNumber, line);
                }

                if (!int.TryParse(parts[2].Trim(), out int quantity))
                {
                    throw new CsvParseException($"Ошибка преобразования количества: '{parts[2]}'", lineNumber, line);
                }

                products.Add(new Product
                {
                    Name = parts[0].Trim(),
                    Price = price,
                    Quantity = quantity
                });
            }

            DataGridProducts.ItemsSource = products;
            TxtStatus.Text = $"Успешно импортировано строк: {products.Count}";
        }
        catch (FileNotFoundException ex)
        {
            TxtStatus.Text = $"Ошибка: файл не найден ({ex.FileName ?? ex.Message})";
        }
        catch (DirectoryNotFoundException ex)
        {
            TxtStatus.Text = $"Ошибка: указанный путь не найден ({ex.Message})";
        }
        catch (UnauthorizedAccessException ex)
        {
            TxtStatus.Text = $"Ошибка: нет прав доступа к файлу ({ex.Message})";
        }
        catch (IOException ex)
        {
            TxtStatus.Text = $"Ошибка ввода-вывода: {ex.Message}";
        }
        catch (CsvParseException ex) when (ex.LineNumber > 0)
        {
            TxtStatus.Text = $"Ошибка парсинга в строке {ex.LineNumber}: {ex.Message}";
            await ShowMessageBoxAsync(
                $"Ошибка в строке №{ex.LineNumber}\nСодержимое: \"{ex.RawLine}\"\n\nПричина: {ex.Message}",
                "Ошибка чтения CSV");
        }
        catch (Exception ex)
        {
            TxtStatus.Text = $"Ошибка: {ex.Message}";
        }
        finally
        {
            reader?.Dispose();
            ProgressBarImport.IsVisible = false;
        }
    }

    private async Task ShowMessageBoxAsync(string message, string title)
    {
        var btnOk = new Button
        {
            Content = "OK",
            HorizontalAlignment = Avalonia.Layout.HorizontalAlignment.Right,
            Width = 80
        };

        var dialog = new Window
        {
            Title = title,
            Width = 420,
            Height = 190,
            WindowStartupLocation = WindowStartupLocation.CenterOwner,
            CanResize = false,
            Content = new StackPanel
            {
                Margin = new Avalonia.Thickness(20),
                Spacing = 15,
                Children =
                {
                    new TextBlock { Text = message, TextWrapping = Avalonia.Media.TextWrapping.Wrap },
                    btnOk
                }
            }
        };

        btnOk.Click += (_, _) => dialog.Close();
        await dialog.ShowDialog(this);
    }

    private void BtnExport_Click(object sender, RoutedEventArgs e)
    {
        StreamWriter? writer = null;
        try
        {
            if (DataGridProducts.ItemsSource is not IEnumerable<Product> products)
            {
                throw new InvalidOperationException("Сначала загрузите CSV-файл!");
            }

            var list = new List<Product>(products);
            if (list.Count == 0)
            {
                throw new InvalidOperationException("Таблица пуста — нет данных для отчёта.");
            }

            string reportPath = "report.txt";
            writer = new StreamWriter(reportPath);

            writer.WriteLine("ОТЧЁТ ПО СКЛАДУ");
            writer.WriteLine($"Дата: {DateTime.Now:dd.MM.yyyy HH:mm}");
            writer.WriteLine("--------------------------------------------");

            decimal totalSum = 0;
            foreach (var p in list)
            {
                decimal sum = p.Price * p.Quantity;
                totalSum += sum;
                writer.WriteLine($"{p.Name} — {p.Price} руб. x {p.Quantity} шт. = {sum} руб.");
            }

            writer.WriteLine("--------------------------------------------");
            writer.WriteLine($"Всего позиций: {list.Count}");
            writer.WriteLine($"Итоговая стоимость: {totalSum} руб.");

            TxtStatus.Text = $"Отчёт успешно сохранён в {reportPath}";
        }
        catch (InvalidOperationException ex)
        {
            TxtStatus.Text = $"Ошибка экспорта: {ex.Message}";
        }
        catch (UnauthorizedAccessException)
        {
            TxtStatus.Text = "Ошибка: нет прав на запись файла report.txt";
        }
        catch (IOException ex)
        {
            TxtStatus.Text = $"Ошибка записи файла: {ex.Message}";
        }
        finally
        {
            writer?.Dispose();
        }
    }
}
