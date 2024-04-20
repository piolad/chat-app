using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.Extensions.Logging;
using Grpc.Core;

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
            _logger.LogInformation("Login submitted with Username: {Username}, Password: {Password}", Username, Password);
            
            if(Username != "admin" || Password != "admin")
            {
                ViewData["AlertMessage"] = "Invalid username or password. Please try again.";;
            }else{
                ViewData["AlertMessage"] = "Good work! Keep it up";
            }

            var channel = new Channel("localhost:50050", ChannelCredentials.Insecure);

            var client = new BrowserFacade.BrowserFacadeClient(channel);

            //We need to check if this is username or email

            var request = new LoginCreds
            {
                // You can either set username or email based on your oneof definition
                Username = Username,
                Password = Password
            };

            var response = await client.LoginAsync(request);

            if (response.Success)
            {
                Console.WriteLine("Login successful.");
                Console.WriteLine($"Token: {response.Token}");
                Console.WriteLine($"Username: {response.Username}");
            }
            else
            {
                Console.WriteLine($"Login failed: {response.Message}");
            }
            
            // Refresh the page to display the alert message
            return Page();
        }
    }
}
