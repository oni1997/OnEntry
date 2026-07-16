using System.Windows;
using System.Windows.Input;

namespace OnEntryDesktop;

public partial class LoginWindow : Window
{
    public LoginWindow()
    {
        InitializeComponent();
    }

    private void Login_Click(object sender, RoutedEventArgs e)
    {
        var email = EmailBox.Text;
        var password = PasswordBox.Password;

        if (string.IsNullOrEmpty(email) || string.IsNullOrEmpty(password))
        {
            MessageBox.Show("Please enter email and password", "OnEntry", MessageBoxButton.OK, MessageBoxImage.Warning);
            return;
        }

        var main = new MainWindow();
        main.Show();
        Close();
    }
}