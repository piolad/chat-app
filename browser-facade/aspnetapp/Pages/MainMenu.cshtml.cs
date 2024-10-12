using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.Extensions.Logging;
using Grpc.Net.Client;
using Grpc.Core;
using BrowserFacade;
using System.Security.Claims; // For Claim, ClaimsIdentity, ClaimsPrincipal
using Microsoft.AspNetCore.Authentication; // For AuthenticationProperties and SignInAsync
using Microsoft.AspNetCore.Authentication.Cookies; // For CookieAuthenticationDefaults

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
                    // Create user claims
                    var claims = new List<Claim>
                    {
                        new Claim(ClaimTypes.Name, Username),
                        new Claim(ClaimTypes.Role, "User"), // You can also set roles dynamically
                    };

                    var claimsIdentity = new ClaimsIdentity(claims, CookieAuthenticationDefaults.AuthenticationScheme);

                    var authProperties = new AuthenticationProperties
                    {
                        IsPersistent = true, // Keep the user logged in even after closing the browser (optional)
                        ExpiresUtc = DateTime.UtcNow.AddMinutes(30) // Set session timeout
                    };
                    HttpContext.SignInAsync(CookieAuthenticationDefaults.AuthenticationScheme,
                                            new ClaimsPrincipal(claimsIdentity),
                                            authProperties).Wait();
                    ViewData["AlertMessage"] = "Login successful!";
                }
                else 
                {
                    ViewData["AlertMessage"] = "Invalid username or password. Please try again.";
                }
            }
            catch (RpcException ex)
            {
                _logger.LogError($"Error code: {ex.StatusCode}. Message: {ex.Status.Detail}");
                ViewData["AlertMessage"] = $"Error code: {ex.StatusCode}. Message: {ex.Status.Detail}";
            }       
            catch (Exception ex)
            {
                _logger.LogError(ex, "An unexpected error occurred.");
                ViewData["AlertMessage"] = $"An unexpected error occurred. Please try again later.";
            }

            // Refresh the page to display the alert message
            return Page();
        }
    }
}
