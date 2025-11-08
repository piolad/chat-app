using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.Extensions.Logging;
using Grpc.Net.Client;
using Grpc.Core;
using BrowserFacade;
using System.ComponentModel.DataAnnotations;
using System.Security.Claims; // For Claim, ClaimsIdentity, ClaimsPrincipal
using Microsoft.AspNetCore.Authentication; // For AuthenticationProperties and SignInAsync
using Microsoft.AspNetCore.Authentication.Cookies; // For CookieAuthenticationDefaults

namespace aspnetapp.Pages
{
    [ValidateAntiForgeryToken]
    public class MainMenuModel : PageModel
    {
        private readonly ILogger<MainMenuModel> _logger;
        private readonly IMainServiceService _mainssvcsvc;

        public MainMenuModel(ILogger<MainMenuModel> logger, IMainServiceService mainssvcsvc)
        {
            _logger = logger;
            _mainssvcsvc = mainssvcsvc;
        }

        [BindProperty, Required] public string Username { get; set; } = string.Empty;
        [BindProperty, Required] public string Password { get; set; } = string.Empty;

        public async Task<IActionResult> OnPostAsync()
        {
            if(!ModelState.IsValid || string.IsNullOrWhiteSpace(Username) || string.IsNullOrWhiteSpace(Password)){
                ViewData["AlertMessage"] = "Please enter a username and password.";
                return Page();
            }
            
            // loging password on purpose
            _logger.LogInformation("Login submitted with Username: {Username}, Password: {Password}", Username, Password);

            try
            {
                var response = await _mainssvcsvc.LoginAsync(Username, Password, HttpContext.RequestAborted);

                _logger.LogInformation("Login response: {Response}", response);

                if(response.Success)
                {
                    // Create user claims
                    var claims = new List<Claim>
                    {
                        new Claim(ClaimTypes.Name, response.Username ?? Username),
                        new Claim(ClaimTypes.Role, "User"), // You can also set roles dynamically
                    };

                    var claimsIdentity = new ClaimsIdentity(claims, CookieAuthenticationDefaults.AuthenticationScheme);

                    var authProperties = new AuthenticationProperties
                    {
                        IsPersistent = true, // Keep the user logged in even after closing the browser (optional)
                        ExpiresUtc = DateTime.UtcNow.AddMinutes(30) // Set session timeout
                    };
                    HttpContext.SignInAsync(
                        CookieAuthenticationDefaults.AuthenticationScheme,
                        new ClaimsPrincipal(claimsIdentity),
                        authProperties);
                    ViewData["AlertMessage"] = "Login successful!";

                    return RedirectToPage("/Friends");
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
