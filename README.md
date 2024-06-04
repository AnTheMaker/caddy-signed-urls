# caddy-signed-urls
A tiny Caddy Server Middleware adding support for Signed URLs. (under development, not ready for production use yet)

## How to use
In your Caddyfile:
```
@allowed {
  signed "supersecretsigningkey"
}
respond @allowed "Yeyy! The URL signature is valid!" 200
respond "Error: Invalid Signature" 403
```
## How to sign a URL
To sign a URL, simply take
- the secret key
- the whole URL you want to request
mesh them together, and SHA256-Hash it. Then append it as a URL parameter with the name ?token=

Example: `SHA256(supersecretsigningkey + url)`

That's it!



> Todo: More Documentation