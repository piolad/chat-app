using Grpc.Core;
using Grpc.Net.Client;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using System.Security.Claims;
using BrowserFacade;  // Ensure this is included

public class FriendsModel : PageModel
{
    public List<SenderReceiverPair> Conversations { get; set; } = new List<SenderReceiverPair>();

    public async Task<IActionResult> OnGet()
    {
        var user = HttpContext.User;

        if (user.Identity.IsAuthenticated)
        {
            var username = user.FindFirst(ClaimTypes.Name)?.Value;

            // Fetch last conversations from the gRPC service
            try
            {
                using var channel = GrpcChannel.ForAddress("http://main-service:50050");
                var client = new BrowserFacade.BrowserFacade.BrowserFacadeClient(channel);

                var request = new FetchLastXConversationsRequest
                {
                    ConversationMember = username, // Use the authenticated user's name
                    Count = 10, // Fetch the last 10 conversations
                    StartIndex = 0
                };

                var response = await client.FetchLastXConversationsAsync(request);

                if (response != null && response.Pairs != null)
                {
                    Conversations = response.Pairs.ToList(); // Populate the Conversations property
                }
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
        else
        {
            return RedirectToPage("/Login");
        }
    }
}
