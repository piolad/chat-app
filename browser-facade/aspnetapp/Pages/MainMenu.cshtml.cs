using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.Extensions.Logging;
using Grpc.Net.Client;
using BrowserFacade;

namespace aspnetapp.Pages
{
    public class MainMenuModel : PageModel
    {
        private readonly ILogger<MainMenuModel> _logger;

        public MainMenuModel(ILogger<MainMenuModel> logger)
        {
            _logger = logger;
        }

        [BindProperty]
        public string Username { get; set; }

        [BindProperty]
        public string Password { get; set; }

        public IActionResult OnPost()
        {
            try
            {
                if(string.IsNullOrWhiteSpace(Username) || string.IsNullOrWhiteSpace(Password))
                {
                    ViewData["AlertMessage"] = "Please enter a username and password.";
                    return Page();
                }

                _logger.LogInformation("Login submitted with Username: {Username}, Password: {Password}", Username, Password);
                
                using var channel = GrpcChannel.ForAddress("http://main-service:50050");
                var client = new BrowserFacade.BrowserFacade.BrowserFacadeClient(channel);
                var request = new LoginCreds { Username = Username, Password = Password };

                var response = client.Login(request);

                _logger.LogInformation("Login response: {Response}", response);

                if(response.Success)
                {
                    ViewData["AlertMessage"] = "Login successful!";
                }
                else
                {
                    ViewData["AlertMessage"] = "Invalid username or password. Please try again.";
                }
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "An error occurred while processing the gRPC request.");
                ViewData["AlertMessage"] = "An error occurred while processing your request. Please try again later.";
            }

            // Refresh the page to display the alert message
            return Page();
        }
    }
}
