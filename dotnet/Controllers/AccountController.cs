using System.Threading.Tasks;
using Microsoft.AspNetCore.Authentication;
using Microsoft.AspNetCore.Authentication.Cookies;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;

namespace Auth0Minimal.Controllers
{
    public class AccountController : Controller
    {
        public async Task Login(string returnUrl = "/Home/Portal")
        {
            var props = new AuthenticationProperties { RedirectUri = returnUrl };
            await HttpContext.ChallengeAsync("Auth0", props);
        }

        [Authorize]
        public IActionResult Profile()
        {
            return View();
        }

        [Authorize]
        public async Task Logout()
        {
            var props = new AuthenticationProperties
            {
                RedirectUri = Url.Action("Index", "Home")!
            };

            await HttpContext.SignOutAsync("Auth0", props);
            await HttpContext.SignOutAsync(CookieAuthenticationDefaults.AuthenticationScheme);
        }
    }
}
