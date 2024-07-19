# simple-share

A simple personal file sharing service.

## Features

- Completely backed by Alibaba Cloud OSS without a DB
- OIDC authentication
- Generate directory / file / text / url shares
  - Customizable share name
  - Customizable generated share link length
  - Password protection
  - Expiration date

## Configuration

### General

- `DEBUG`: debug mode (defaults to `false`)
- `EMBED_DISABLE`: disable web assets embedding (defaults to `true` in debug mode)
- `HOST`, `PORT`: address to listen
- `BASEURL`: base URL of the server (required unless in debug mode)

### OIDC Authentication

The application currently only supports OIDC authentication.
A user will be allowed to create shares if a valid OIDC token is present.
Note that we do not authorize the user. Instead, the OIDC provider is expected to do so.
If a user is not allowed to create shares, the OIDC provider should not return a valid token at all.

If the following entries are not properly configured, authentication will be disabled, and the application will be read-only.

Remember to set `BASEURL` correctly for callback url to work.

- `OIDC_ISSUER`: OIDC issuer
- `OIDC_CLIENT_ID`: OIDC client_id
- `OIDC_CLIENT_SECRET`: OIDC client_secret
- `OIDC_NAME_CLAIM`: name of the username claim (default: `username`)

### Storage

The application uses Alibaba Cloud OSS service for storage.

- `ALIBABA_CLOUD_ACCESS_KEY_ID`: Alibaba Cloud AccessKey ID
- `ALIBABA_CLOUD_ACCESS_KEY_SECRET`: Alibaba Cloud AccessKey Secret
- `OSS_ENDPOINT`: datacenter endpoint to use (example: `region-internal.aliyuncs.com`)
- `OSS_ENDPOINT_PUBLIC`: public bucket endpoint to use (for custom domain; defaults to `OSS_ENDPOINT` if not set)
- `OSS_BUCKET`: bucket name
