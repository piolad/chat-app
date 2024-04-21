using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.Extensions.Logging;
using Grpc.Net.Client;
using Greeter; // Import the correct namespace

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
                _logger.LogInformation("Login submitted with Username: {Username}, Password: {Password}", Username, Password);
                
                if(Username != "admin" || Password != "admin")
                {
                    ViewData["AlertMessage"] = "Invalid username or password. Please try again.";
                }
                else
                {
                    // Create an insecure gRPC channel
                    using var channel = GrpcChannel.ForAddress("http://main-service:50050");
                    
                    // Create a gRPC client
                    var client = new Greeter.Greeter.GreeterClient(channel);
                    
                    // Create a request message
                    var request = new HelloRequest { Name = Username };
                    
                    // Call the gRPC service
                    var response = client.SayHello(request);
                    
                    // Display the response message
                    ViewData["AlertMessage"] = response.Message;
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
