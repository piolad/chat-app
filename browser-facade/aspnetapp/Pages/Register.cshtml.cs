using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.Extensions.Logging;
using Grpc.Net.Client;
using Grpc.Core;
using BrowserFacade;
using System.ComponentModel.DataAnnotations;
// using DataAnnotations;

namespace aspnetapp.Pages
{
    public class RegisterVm
    {
        [Required, StringLength(30)]
        public string Username { get; set; } = "";

        [Required, EmailAddress]
        public string Email { get; set; } = "";

        [Required, DataType(DataType.Password), MinLength(6)]
        public string Password { get; set; } = "";

        [Required, StringLength(40)]
        public string FirstName { get; set; } = "";

        [Required, StringLength(40)]
        public string LastName { get; set; } = "";

        [Required, DataType(DataType.Date)]
        public DateOnly BirthDate { get; set; }   // .NET 8 can bind DateOnly
    }


    public class RegisterModel : PageModel
    {
        private readonly ILogger<RegisterModel> _logger;
        private readonly IMainServiceService _mainsvcsvc;

        public RegisterModel(ILogger<RegisterModel> logger, IMainServiceService mainsvcsvc)
        {
            _logger = logger;
            _mainsvcsvc = mainsvcsvc;
        }

        [BindProperty] public RegisterVm Input {get; set;} = new();
        

        public async Task<IActionResult>  OnPost(CancellationToken ct){
             if (!ModelState.IsValid)
                return Page();

            try
            {

                ViewData["AlertMessage"] = "Good done!";
                _logger.LogInformation("Register submitted with FirstName: {Input.FirstName}, LastName: {Input.LastName}, BirthDate: {Input.BirthDate}, Email: {Input.Email}, Username: {Input.Username}, Password: {Input.Password}", Input.FirstName, Input.LastName, Input.BirthDate, Input.Email, Input.Username, Input.Password);
                           
                var response = await _mainsvcsvc.RegisterAsync(Input.FirstName, Input.LastName,  Input.BirthDate.ToString("yyyy-MM-dd"), Input.Email, Input.Username, Input.Password, ct);

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