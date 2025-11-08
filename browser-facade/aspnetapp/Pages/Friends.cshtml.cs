using Grpc.Core;
using Grpc.Net.Client;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.AspNetCore.Authorization;
using System.Security.Claims;
using BrowserFacade;

[Authorize]
public class FriendsModel : PageModel
{
    public List<SenderReceiverPair> Conversations { get; set; } = new List<SenderReceiverPair>();

    private readonly ILogger<FriendsModel> _logger;
    private readonly IMainServiceService _mainsvcsvc;

    public FriendsModel(ILogger<FriendsModel> logger, IMainServiceService mainsvcsvc)
    {
        _logger = logger;
        _mainsvcsvc = mainsvcsvc;
    }

    public async Task<IActionResult> OnGet(CancellationToken ct)
    {
        var user = HttpContext.User;
            var username = user.FindFirst(ClaimTypes.Name)?.Value;

            // Fetch last conversations from the gRPC service
            try
            {
                
                _logger.LogInformation("Fetching messages for username: {username}", username);
                var response = await _mainsvcsvc.FetchLastXConversationsAsync(username, 0, 10, ct);

                _logger.LogInformation("Message response: {Response}", response);

                
                Conversations = response?.Pairs?.ToList() ?? new();
            }
            catch (RpcException ex)
            {
                // Handle gRPC error
                ModelState.AddModelError(string.Empty, $"Error fetching conversations: {ex.Status.Detail}");
            }
            catch (Exception ex)
            {
                // Handle general errors
                ModelState.AddModelError(string.Empty, "An unexpected error occurred while fetching conversations.");
            }

            return Page();
    }
}
