using System;
using System.Collections.Generic;
using System.Text;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Navigation;
using System.Windows.Shapes;
using WPFDemoExamPractice.Models;

namespace WPFDemoExamPractice.UserControls
{
    /// <summary>
    /// Логика взаимодействия для LoginUC.xaml
    /// </summary>
    public partial class DataUC : UserControl
    {
        public DataUC()
        {
            InitializeComponent();
            Loaded += (s, e) =>
            {
                itemsLoad();
            };
        }
        public async void itemsLoad()
        {
            ItemGrid.ItemsSource = appstat.Context.Addresses.ToList();
            combobox.ItemsSource = appstat.Context.Addresses.ToList();
        }
    }
}
