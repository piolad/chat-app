using System;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.Extensions.Logging;
using Grpc.Net.Client;
using Grpc.Core;
using BrowserFacade;

namespace aspnetapp.Pages
{
    public class MainLoginModel : PageModel
    {
        private readonly ILogger<MainLoginModel> _logger;

        public MainLoginModel(ILogger<MainLoginModel> logger)
        {
            _logger = logger;
        }
            public string sender = "kuba";
            public string receiver = "mati";
            public string message = "Ja jestem najlepszy";
            public string timestamp = DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss");

        public IActionResult OnPost()
        {
            try
            {             
                using var channel = GrpcChannel.ForAddress("http://message-data-centre:50052"); 
                var client = new BrowserFacade.BrowserFacade.BrowserFacadeClient(channel);
                var request = new Message { Sender = sender, Receiver = receiver, Message_ = message, Timestamp = timestamp}; //zmiana

                var response = client.SendMessage(request);

                _logger.LogInformation("Message response: {Response}", response);
            }
            catch (RpcException ex)
            {
                _logger.LogError($"Error code: {ex.StatusCode}. Message: {ex.Status.Detail}");
                ViewData["AlertMessage"] = $"Error code: {ex.StatusCode}. Message: {ex.Status.Detail}";
            }       
            return Page();
        }
    }
}