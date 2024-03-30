using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.Extensions.Logging;

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
            // Save the login credentials to a database or any other storage mechanism.
            // For demonstration, let's just log the credentials.
            _logger.LogInformation("Login submitted with Username: {Username}, Password: {Password}", Username, Password);

            // Here you can add logic to validate the credentials, authenticate the user, etc.
            // For demonstration purposes, let's assume validation fails and display an error message.

            if(Username != "admin" || Password != "admin")
            {
                ViewData["AlertMessage"] = "Invalid username or password. Please try again.";;
            }else{
                ViewData["AlertMessage"] = "Good work! Keep it up";
            }

            // Refresh the page to display the alert message
            return Page();
        }
    }
}
