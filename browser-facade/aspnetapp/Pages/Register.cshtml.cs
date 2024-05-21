using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.Extensions.Logging;
using Grpc.Net.Client;
using Grpc.Core;
using BrowserFacade;

namespace aspnetapp.Pages
{
    public class RegisterModel : PageModel
    {
        private readonly ILogger<RegisterModel> _logger;

        public RegisterModel(ILogger<RegisterModel> logger)
        {
            _logger = logger;
        }

        [BindProperty]
        public string FirstName { get; set; }
        [BindProperty]
        public string LastName { get; set; }
        [BindProperty]
        public string BirthDate { get; set; }
        [BindProperty]
        public string Email { get; set; }
        [BindProperty]
        public string Username { get; set; }
        [BindProperty]
        public string Password { get; set; }

        public IActionResult  OnPost(){
            try
            {
                //TODO - need to add a validation for the date


                //_logger.LogInformation("Starts with: {FirstName}, LastName: {LastName}, BirthDate: {BirthDate}, Email: {Email}, Username: {Username}, Password: {Password}", FirstName, LastName, BirthDate, Email, Username, Password);
                if(string.IsNullOrWhiteSpace(FirstName) || string.IsNullOrWhiteSpace(LastName) || string.IsNullOrWhiteSpace(BirthDate) || string.IsNullOrWhiteSpace(Email) || string.IsNullOrWhiteSpace(Username) || string.IsNullOrWhiteSpace(Password))
                {
                    ViewData["AlertMessage"] = "Please fill out all fields.";
                    return Page();
                }

                ViewData["AlertMessage"] = "Good done!";
                _logger.LogInformation("Register submitted with FirstName: {FirstName}, LastName: {LastName}, BirthDate: {BirthDate}, Email: {Email}, Username: {Username}, Password: {Password}", FirstName, LastName, BirthDate, Email, Username, Password);
                
                using var channel = GrpcChannel.ForAddress("http://main-service:50050");
                var client = new BrowserFacade.BrowserFacade.BrowserFacadeClient(channel);
                var request = new RegisterCreds { Firstname = FirstName, Lastname = LastName, Birthdate = BirthDate, Email = Email, Username = Username, Password = Password };

                var response = client.Register(request);

                _logger.LogInformation("Register response: {Response}", response);

                if(response.Success)
                {
                    ViewData["AlertMessage"] = "Registration successful!";
                }
                else
                {
                     ViewData["AlertMessage"] = "Registration failed. Please try again.";
                }
            }
            catch (RpcException ex)
            {
                _logger.LogError($"Error code: {ex.StatusCode}. Message: {ex.Status.Detail}");
                ViewData["AlertMessage"] = $"Error code: {ex.StatusCode}. Message: {ex.Status.Detail}";
            }       
            catch (Exception ex)
            {
                _logger.LogError($"Error: {ex.Message}");
                ViewData["AlertMessage"] = $"Error: {ex.Message}";
            }
            return Page();
        }
    }
}