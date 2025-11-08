using System;
using System.ComponentModel.DataAnnotations;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.Extensions.Logging;

namespace aspnetapp.Pages
{
    public class MainLoginModel : PageModel
    {
        private readonly ILogger<MainLoginModel> _logger;
        private readonly IChatService _chat;

        public MainLoginModel(ILogger<MainLoginModel> logger, IChatService chat)
        {
            _logger = logger;
            _chat = chat;
        }

        // Bind form fields
        [BindProperty, Required] public string Sender { get; set; } = "kuba";
        [BindProperty, Required] public string Receiver { get; set; } = "mati";
        [BindProperty, Required, StringLength(2000)] public string Body { get; set; } = "Ja jestem najlepszy";

        // Expose a message to the view
        public string? Alert { get; private set; }

        [ValidateAntiForgeryToken]
        public async Task<IActionResult> OnPostAsync()
        {
            if (!ModelState.IsValid)
            {
                Alert = "Form has validation errors.";
                return Page();
            }

            var timestamp = DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss");

            try
            {
                var resp = await _chat.SendMessageAsync(Sender, Receiver, Body, timestamp, HttpContext.RequestAborted);
                _logger.LogInformation("Message response: {Response}", resp?.Message);
                Alert = $"✅ Sent. Server says: {resp?.Message}";
            }
            catch (Grpc.Core.RpcException ex)
            {
                _logger.LogError(ex, "gRPC error {Code}: {Detail}", ex.StatusCode, ex.Status.Detail);
                Alert = $"❌ gRPC error {ex.StatusCode}: {ex.Status.Detail}";
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Unexpected error");
                Alert = "❌ Unexpected error. Check logs for details.";
            }

            return Page();
        }
    }
}
