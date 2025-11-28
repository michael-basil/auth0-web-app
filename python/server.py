import json
from os import environ as env
from urllib.parse import quote_plus, urlencode

from authlib.integrations.base_client import OAuthError
from authlib.integrations.flask_client import OAuth
from dotenv import find_dotenv, load_dotenv
from flask import Flask, redirect, render_template, session, url_for

# Load .env file
ENV_FILE = find_dotenv()
if ENV_FILE:
    load_dotenv(ENV_FILE)

# Flask app setup
app = Flask(__name__)
app.secret_key = env.get("APP_SECRET_KEY")

app.config.update(
    SERVER_NAME=env.get("SERVER_NAME"),
    APPLICATION_ROOT=env.get("APPLICATION_ROOT", "/"),
    PREFERRED_URL_SCHEME=env.get("PREFERRED_URL_SCHEME", "https"),
)

# Auth0 / Authlib setup
oauth = OAuth(app)
oauth.register(
    "auth0",
    client_id=env.get("AUTH0_CLIENT_ID"),
    client_secret=env.get("AUTH0_CLIENT_SECRET"),
    client_kwargs={
        "scope": "openid profile email",
    },
    server_metadata_url=f'https://{env.get("AUTH0_DOMAIN")}/.well-known/openid-configuration',
)

# Routes

@app.route("/login")
def login():
    scheme = env.get("PREFERRED_URL_SCHEME", "https")
    return oauth.auth0.authorize_redirect(
        redirect_uri = url_for("callback", _external=True, _scheme=scheme)
    )

@app.route("/callback", methods=["GET", "POST"])
def callback():
    try:
        token = oauth.auth0.authorize_access_token()
    except OAuthError as e:
        print("OAuthError:", e.error, e.description)
        return f"OAuth error: {e.error} – {e.description}", 400

    session["user"] = token
    return redirect(url_for("home"))

@app.route("/logout")
def logout():
    session.clear()
    scheme = env.get("PREFERRED_URL_SCHEME", "https")
    return redirect(
        "https://" + env.get("AUTH0_DOMAIN")
        + "/v2/logout?"
        + urlencode(
            {
                "returnTo": url_for("home", _external=True, _scheme=scheme),
                "client_id": env.get("AUTH0_CLIENT_ID"),
            },
            quote_via=quote_plus,
        )
    )

@app.route("/")
def home():
    user = session.get("user")
    pretty = json.dumps(user, indent=4) if user else None
    return render_template("home.html", session=user, pretty=pretty)

if __name__ == "__main__":
    # PORT is configurable via env; default 3000 (matches quickstart)
    app.run(host="0.0.0.0", port=int(env.get("PORT", 3000)))
