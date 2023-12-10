import os
from flask import Flask

app = Flask(__name__)

@app.route("/")
def root():
    return "App is running"

@app.route("/health")
def health():
    return "OK"

if __name__ == "__main__":
    # Get the port from the environment variable or use 5000 as default
    port_env = os.environ.get("PORT", None)
    if port_env is None:
        print("No valid port specified")
        exit(1)
    
    port = int(port_env)
    
    # Run the app on the specified port
    app.run(port=port)