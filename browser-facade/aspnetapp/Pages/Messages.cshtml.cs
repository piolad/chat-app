using Grpc.Core;
using Grpc.Net.Client;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.AspNetCore.Authorization;
using System.Security.Claims;
using BrowserFacade; // Ensure this is included

[Authorize]
public class MessagesModel : PageModel
{

    private readonly ILogger<MessagesModel> _logger;
    private readonly IMainServiceService _mainsvcsvc;

    public MessagesModel(ILogger<MessagesModel> logger, IMainServiceService mainsvcsvc)
    {
        _logger = logger;
        _mainsvcsvc = mainsvcsvc;
    }

    public string Receiver { get; set; } // The user you are conversing with

    public List<Message> Messages { get; set; } = new List<Message>();

    [BindProperty]
    public string NewMessage { get; set; } // For the message input

    public async Task<IActionResult> OnGet(string receiver, CancellationToken ct)
    {

        var user = HttpContext.User;

        var sender = user.FindFirst(ClaimTypes.Name)?.Value;
        Receiver = receiver;

        // Fetch last 10 messages from the gRPC service
        try
        {
            _logger.LogInformation("Fetching last 10 messages between {sender} and {receiver}", sender, receiver);

            var response = await _mainsvcsvc.FetchLastXMessagesAsync(sender, receiver, 0, 10, ct);

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

    public async Task<IActionResult> OnPostSend(string receiver, CancellationToken ct)
    {
        var user = HttpContext.User;

        if(string.IsNullOrEmpty(NewMessage)){
            //show error
            _logger.LogInformation("Empty Message");
            return RedirectToPage(new { receiver });
        }

        var sender = user.FindFirst(ClaimTypes.Name)?.Value;
        Receiver = receiver;

        // Send a new message via gRPC service
        try
        {

            var message = new Message
            {
                Sender = sender,
                Receiver = receiver,
                Message_ = NewMessage, // The text message from the form
                Timestamp = DateTime.UtcNow.ToString("yyyy-MM-dd HH:mm:ss")
            };

            var response = await _mainsvcsvc.SendMessageAsync(sender, receiver, NewMessage, DateTime.UtcNow.ToString("yyyy-MM-dd HH:mm:ss"), ct);

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
