using System.Threading.Tasks;
using Microsoft.AspNetCore.Authentication;
using Microsoft.AspNetCore.Authentication.Cookies;
using Microsoft.AspNetCore.Authentication.OpenIdConnect;
using Microsoft.AspNetCore.HttpOverrides;
using Microsoft.IdentityModel.Tokens;
using Microsoft.IdentityModel.Protocols.OpenIdConnect;

var builder = WebApplication.CreateBuilder(args);

// 1. Configuration – ensure User Secrets are loaded
builder.Configuration.AddUserSecrets<Program>(optional: true);

// 2. MVC
builder.Services.AddControllersWithViews();

// 3. Authentication – explicit cookie + OIDC named "Auth0"
builder.Services
    .AddAuthentication(options =>
    {
        options.DefaultAuthenticateScheme = CookieAuthenticationDefaults.AuthenticationScheme;
        options.DefaultSignInScheme = CookieAuthenticationDefaults.AuthenticationScheme;
        options.DefaultChallengeScheme = "Auth0";
    })
    .AddCookie(cookieOptions =>
    {
        // Make the auth cookie work nicely in Codespaces (HTTPS, proxy, etc.)
        cookieOptions.Cookie.SameSite = SameSiteMode.Lax; // or SameSiteMode.None if needed
        cookieOptions.Cookie.SecurePolicy = CookieSecurePolicy.Always;
    })
    .AddOpenIdConnect("Auth0", options =>
    {
        var domain = builder.Configuration["Auth0:Domain"];
        var clientId = builder.Configuration["Auth0:ClientId"];
        var clientSecret = builder.Configuration["Auth0:ClientSecret"];

        options.Authority = $"https://{domain}";
        options.ClientId = clientId;
        options.ClientSecret = clientSecret;

        options.ResponseType = "code";
        // IMPORTANT: align with what we see in the HAR
        options.ResponseMode = OpenIdConnectResponseMode.Query;

        options.SaveTokens = true;
        options.CallbackPath = "/callback";
        options.ClaimsIssuer = "Auth0";

        options.Scope.Clear();
        options.Scope.Add("openid");
        options.Scope.Add("profile");
        options.Scope.Add("email");

        options.TokenValidationParameters = new TokenValidationParameters
        {
            NameClaimType = "name"
        };

        // Make sure correlation & nonce cookies survive cross-site POST from Auth0
        options.CorrelationCookie.SameSite = SameSiteMode.None;
        options.CorrelationCookie.SecurePolicy = CookieSecurePolicy.Always;
        options.NonceCookie.SameSite = SameSiteMode.None;
        options.NonceCookie.SecurePolicy = CookieSecurePolicy.Always;

        options.Events = new OpenIdConnectEvents
        {
            OnRedirectToIdentityProvider = context =>
            {
                Console.WriteLine("OnRedirectToIdentityProvider:");
                Console.WriteLine("  IssuerAddress: " + context.ProtocolMessage.IssuerAddress);
                Console.WriteLine("  redirect_uri: " + context.ProtocolMessage.RedirectUri);
                return Task.CompletedTask;
            },
            OnTokenValidated = context =>
            {
                Console.WriteLine("OnTokenValidated: user authenticated.");
                return Task.CompletedTask;
            },
            OnRemoteFailure = context =>
            {
                Console.WriteLine("OnRemoteFailure: " + context.Failure?.Message);
                if (context.Failure?.InnerException != null)
                {
                    Console.WriteLine("Inner: " + context.Failure.InnerException.Message);
                }
                context.HandleResponse();
                context.Response.Redirect("/?error=remote_failure");
                return Task.CompletedTask;
            },
            OnAuthenticationFailed = context =>
            {
                Console.WriteLine("OnAuthenticationFailed: " + context.Exception.Message);
                context.HandleResponse();
                context.Response.Redirect("/?error=auth_failed");
                return Task.CompletedTask;
            }
        };
    });

var app = builder.Build();

app.Use(async (context, next) =>
{
    Console.WriteLine($"Request: {context.Request.Method} {context.Request.Path}{context.Request.QueryString} " +
                      $"Auth={context.User?.Identity?.IsAuthenticated}");

    await next();

    Console.WriteLine($"Response: {context.Response.StatusCode} for {context.Request.Path}");
});

// 4. Make ASP.NET Core respect Codespaces proxy headers
var fwd = new ForwardedHeadersOptions
{
    ForwardedHeaders = ForwardedHeaders.XForwardedFor |
                       ForwardedHeaders.XForwardedProto |
                       ForwardedHeaders.XForwardedHost
};
fwd.KnownNetworks.Clear();
fwd.KnownProxies.Clear();
app.UseForwardedHeaders(fwd);

// 5. Standard pipeline
if (!app.Environment.IsDevelopment())
{
    app.UseExceptionHandler("/Home/Error");
    app.UseHsts();
}

app.UseHttpsRedirection();
app.UseRouting();

app.UseAuthentication();
app.UseAuthorization();

// 6. Single default route: Home/Index
app.MapControllerRoute(
    name: "default",
    pattern: "{controller=Home}/{action=Index}/{id?}");

app.Run();
