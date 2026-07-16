using System.Windows;
using System.Windows.Threading;

namespace OnEntryDesktop;

public partial class MainWindow : Window
{
    private readonly DispatcherTimer _autoLockTimer;
    private readonly DispatcherTimer _clipboardTimer;
    private readonly DispatcherTimer _syncTimer;

    public MainWindow()
    {
        InitializeComponent();
        SetupTimers();
        SetupTrayIcon();
    }

    private void SetupTimers()
    {
        _autoLockTimer = new DispatcherTimer { Interval = TimeSpan.FromMinutes(5) };
        _autoLockTimer.Tick += (s, e) => Lock();

        _clipboardTimer = new DispatcherTimer { Interval = TimeSpan.FromSeconds(30) };
        _clipboardTimer.Tick += (s, e) => ClearClipboard();

        _syncTimer = new DispatcherTimer { Interval = TimeSpan.FromMinutes(1) };
        _syncTimer.Tick += async (s, e) => await SyncAsync();
    }

    private void SetupTrayIcon()
    {
    }

    private void Lock()
    {
        _autoLockTimer.Stop();
        _clipboardTimer.Stop();
        MessageBox.Show("Locked", "OnEntry", MessageBoxButton.OK, MessageBoxImage.Information);
    }

    private void ClearClipboard()
    {
        try
        {
            Clipboard.Clear();
        }
        catch { }
    }

    private async System.Threading.Tasks.Task SyncAsync()
    {
    }
}