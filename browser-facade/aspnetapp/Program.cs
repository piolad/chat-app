using System.Text.Json.Serialization;
using Microsoft.AspNetCore.Authentication.Cookies;
using BrowserFacade;
using Grpc.Net.Client;

var builder = WebApplication.CreateBuilder(args);

// 1) Framework services
builder.Services.AddRazorPages();
builder.Services.AddHealthChecks();
builder.Services.AddSession(o => o.IdleTimeout = TimeSpan.FromMinutes(30));
builder.Services.AddAuthentication(CookieAuthenticationDefaults.AuthenticationScheme)
    .AddCookie(o =>
    {
        o.LoginPath = "/Login";
        o.LogoutPath = "/Logout";
    });

// 2) gRPC client(s) + app services (BEFORE Build)
var addr = builder.Configuration["MainService:BaseAddress"] ?? "http://main-service:50050";
if (addr.StartsWith("http://", StringComparison.OrdinalIgnoreCase))
    AppContext.SetSwitch("System.Net.Http.SocketsHttpHandler.Http2UnencryptedSupport", true);

builder.Services.AddGrpcClient<BrowserFacade.BrowserFacade.BrowserFacadeClient>(o => o.Address = new Uri(addr));
builder.Services.AddGrpcClient<BrowserFacade.MessageService.MessageServiceClient>(o => o.Address = new Uri(addr));

builder.Services.AddScoped<IMainServiceService, MainServiceService>();

var app = builder.Build();

// 3) Middleware & endpoints
app.MapHealthChecks("/healthz");

if (!app.Environment.IsDevelopment())
{
    app.UseExceptionHandler("/Error");
    app.UseHsts();
}

app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseSession();
app.UseAuthentication();
app.UseAuthorization();

app.MapRazorPages();

app.MapGet("/", () => Results.Redirect("/MainMenu"));

CancellationTokenSource cancellation = new();
app.Lifetime.ApplicationStopping.Register(() => cancellation.Cancel());

app.MapGet("/Environment", () => new EnvironmentInfo());
app.MapGet("/Delay/{value}", async (int value) =>
{
    try { await Task.Delay(value, cancellation.Token); }
    catch (TaskCanceledException) { }
    return new Operation(value);
});

app.Run();

// ----- app services -----
public interface IMainServiceService
{
    Task<Response> SendMessageAsync(string sender, string receiver, string message, string timestamp, CancellationToken ct);
    Task<LoginStatus> LoginAsync(string username, string password, CancellationToken ct);
    Task<FetchLastXConversationsResponse> FetchLastXConversationsAsync(string conversationMember, int startIndex, int count, CancellationToken ct);
    Task<FetchLastXMessagesResponse> FetchLastXMessagesAsync(string sender, string receiver, int startingPoint, int count, CancellationToken ct);

}

public sealed class MainServiceService : IMainServiceService
{
    private readonly BrowserFacade.BrowserFacade.BrowserFacadeClient _client;
    public MainServiceService(BrowserFacade.BrowserFacade.BrowserFacadeClient client) => _client = client;

    public async Task<Response> SendMessageAsync(string sender, string receiver, string message, string timestamp, CancellationToken ct)
    {
        var req = new Message { Sender = sender, Receiver = receiver, Message_ = message, Timestamp = timestamp };
        return await _client.SendMessageAsync(req, cancellationToken: ct);
    }

    public async Task<LoginStatus> LoginAsync(string username, string password, CancellationToken ct){
        var req = new LoginCreds { Username = username, Password = password};
        return await _client.LoginAsync(req, cancellationToken: ct);
    }

    public async Task<FetchLastXConversationsResponse> FetchLastXConversationsAsync(string conversationMember, int startIndex, int count, CancellationToken ct){
        var req = new FetchLastXConversationsRequest { ConversationMember = conversationMember, Count = count, StartIndex = startIndex };
        return await _client.FetchLastXConversationsAsync(req, cancellationToken: ct);
    }

    public async Task<FetchLastXMessagesResponse> FetchLastXMessagesAsync(string sender, string receiver, int startingPoint, int count, CancellationToken ct){
        var req = new FetchLastXMessagesRequest { Sender = sender, Receiver = receiver, StartingPoint = startingPoint, Count = count};
        return await _client.FetchLastXMessagesAsync(req, cancellationToken: ct);
    }

}

[JsonSerializable(typeof(EnvironmentInfo))]
[JsonSerializable(typeof(Operation))]
internal partial class AppJsonSerializerContext : JsonSerializerContext { }

public record struct Operation(int Delay);
