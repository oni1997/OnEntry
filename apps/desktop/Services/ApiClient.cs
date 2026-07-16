using System.Net.Http.Json;
using OnEntryDesktop.Models;

namespace OnEntryDesktop.Services;

public class ApiClient
{
    private readonly HttpClient _http;
    private string? _accessToken;

    public ApiClient(string baseUrl)
    {
        _http = new HttpClient { BaseAddress = new Uri(baseUrl) };
    }

    public void SetToken(string token)
    {
        _accessToken = token;
        if (!string.IsNullOrEmpty(token))
        {
            _http.DefaultRequestHeaders.Authorization = new System.Net.Http.Headers.AuthenticationHeaderValue("Bearer", token);
        }
    }

    public async Task<AuthResponse?> LoginAsync(string email, string password)
    {
        var response = await _http.PostAsJsonAsync("/login", new { email, password });
        response.EnsureSuccessStatusCode();
        return await response.Content.ReadFromJsonAsync<AuthResponse>();
    }

    public async Task SyncAsync()
    {
        if (string.IsNullOrEmpty(_accessToken)) return;
        var response = await _http.GetAsync("/vault");
        response.EnsureSuccessStatusCode();
    }
}