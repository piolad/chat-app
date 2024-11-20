using Grpc.Core;
using Grpc.Net.Client;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using System.Security.Claims;
using BrowserFacade;

public class ChatModel : PageModel
{

    private readonly ILogger<ChatModel> _logger;

    public ChatModel(ILogger<ChatModel> logger)
    {
        _logger = logger;
    }

    public async Task<IActionResult> OnGet() {
        return Page();
    }
}  
