using Grpc.Core;
using Grpc.Net.Client;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using System.Security.Claims;
using BrowserFacade; // Ensure this is included

public class MessagesModel : PageModel
{
    public string Receiver { get; set; } // The user you are conversing with
    public List<Message> Messages { get; set; } = new List<Message>();
    [BindProperty]
    public string NewMessage { get; set; } // For the message input

    public async Task<IActionResult> OnGet(string receiver)
    {
        var user = HttpContext.User;

        if (!user.Identity.IsAuthenticated)
        {
            return RedirectToPage("/Login");
        }

        var sender = user.FindFirst(ClaimTypes.Name)?.Value;
        Receiver = receiver;

        // Fetch last 10 messages from the gRPC service
        try
        {
            using var channel = GrpcChannel.ForAddress("http://main-service:50050");
            var client = new BrowserFacade.BrowserFacade.BrowserFacadeClient(channel);

            var request = new FetchLastXMessagesRequest
            {
                Sender = sender,
                Receiver = receiver,
                StartingPoint = 0, // Fetch from the beginning (latest messages)
                Count = 10 // Fetch the last 10 messages
            };

            var response = await client.FetchLastXMessagesAsync(request);

            if (response != null && response.Messages != null)
            {
                Messages = response.Messages.ToList();
            }
        }
        catch (RpcException ex)
        {
            ModelState.AddModelError(string.Empty, $"Error fetching messages: {ex.Status.Detail}");
        }
        catch (Exception ex)
        {
            ModelState.AddModelError(string.Empty, "An unexpected error occurred while fetching messages.");
        }

        return Page();
    }

    public async Task<IActionResult> OnPostSend(string receiver)
    {
        var user = HttpContext.User;

        if (!user.Identity.IsAuthenticated)
        {
            return RedirectToPage("/Login");
        }

        var sender = user.FindFirst(ClaimTypes.Name)?.Value;
        Receiver = receiver;

        // Send a new message via gRPC service
        try
        {
            using var channel = GrpcChannel.ForAddress("http://main-service:50050");
            var client = new BrowserFacade.BrowserFacade.BrowserFacadeClient(channel);

            var message = new Message
            {
                Sender = sender,
                Receiver = receiver,
                Message_ = NewMessage, // The text message from the form
                Timestamp = DateTime.UtcNow.ToString("o") // Use ISO 8601 format for timestamp
            };

            var response = await client.SendMessageAsync(message);

            if (response != null)
            {
                // Optionally handle the response message
            }

            // After sending, fetch the updated message list again
            return RedirectToPage(new { receiver });
        }
        catch (RpcException ex)
        {
            ModelState.AddModelError(string.Empty, $"Error sending message: {ex.Status.Detail}");
        }
        catch (Exception ex)
        {
            ModelState.AddModelError(string.Empty, "An unexpected error occurred while sending the message.");
        }

        return Page();
    }
}
