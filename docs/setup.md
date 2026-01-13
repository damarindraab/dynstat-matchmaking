# Setup

## Prerequisites

1. Windows 11 WSL2 or Linux Ubuntu 22.04 or macOS 14+ with the following tools installed:

   a. Bash

      - On Windows WSL2 or Linux Ubuntu:

         ```
         bash --version

         GNU bash, version 5.1.16(1)-release (x86_64-pc-linux-gnu)
         ...
         ```

      - On macOS:

         ```
         bash --version

         GNU bash, version 3.2.57(1)-release (arm64-apple-darwin23)
         ...
         ```

   b. Make

      - On Windows WSL2 or Linux Ubuntu:

         To install from the Ubuntu repository, run `sudo apt update && sudo apt install make`.

         ```
         make --version

         GNU Make 4.3
         ...
         ```

      - On macOS:

         ```
         make --version

         GNU Make 3.81
         ...
         ```

   c. Docker (Docker Desktop 4.30+ or Docker Engine v23.0+)

      - On Linux Ubuntu:

         1. To install from the Ubuntu repository, run `sudo apt update && sudo apt install docker.io docker-buildx docker-compose-v2`.
         2. Add your user to the `docker` group: `sudo usermod -aG docker $USER`.
         3. Log out and log back in to allow the changes to take effect.

      - On Windows or macOS:

         Follow Docker's documentation on installing Docker Desktop on [Windows](https://docs.docker.com/desktop/install/windows-install/) or [macOS](https://docs.docker.com/desktop/install/mac-install/).

         ```
         docker version

         ...
         Server: Docker Desktop
            Engine:
            Version:          24.0.5
         ...
         ```

   d. Go v1.24

      - Follow [Go's installation guide](https://go.dev/doc/install) to install Go.

      ```
      go version

      go version go1.24.0 ...
      ```

   e. [Postman](https://www.postman.com/)

      - Use the available binary from [Postman](https://www.postman.com/downloads/).

   f. [extend-helper-cli](https://github.com/AccelByte/extend-helper-cli)

      - Use the available binary from [extend-helper-cli](https://github.com/AccelByte/extend-helper-cli/releases).

   g. Local tunnel service that has TCP forwarding capability, such as:

      - [Ngrok](https://ngrok.com/)

         Need registration for free tier. Please refer to [ngrok documentation](https://ngrok.com/docs/getting-started/) for a quick start.

      - [Pinggy](https://pinggy.io/)

         Free to try without registration. Please refer to [pinggy documentation](https://pinggy.io/docs/) for a quick start.

   > :exclamation: In macOS, you may use [Homebrew](https://brew.sh/) to easily install some of the tools above.

2. Access to an AGS environment.

   a. Base URL:

      - For Shared Cloud tier, e.g., https://spaceshooter.prod.gamingservices.accelbyte.io
      - For Private Cloud tier, e.g., https://dev.accelbyte.io

   b. [Create a game namespace](https://docs.accelbyte.io/gaming-services/services/access/reference/namespaces/manage-your-namespaces/) if you do not have one yet. Keep the `Namespace ID`.

   c. [Create an OAuth client](https://docs.accelbyte.io/gaming-services/services/access/authorization/manage-access-control-for-applications/#create-an-iam-client) with confidential client type. Keep the `Client ID` and `Client Secret`.

## Configure The App

1. Create a Docker compose `.env` file by copying the content of `.env.template`.

   > :warning: **Host OS environment variables take precedence over `.env` file values**: If variables in `.env` do not seem to take effect properly, check for host OS environment variables with the same names. See Docker's documentation about [Docker compose environment variables precedence](https://docs.docker.com/compose/how-tos/environment-variables/envvars-precedence/).

2. Fill in the required environment variables in `.env`:

   ```
   AB_BASE_URL=https://test.accelbyte.io     # Base URL of AGS environment
   AB_CLIENT_ID='xxxxxxxxxx'                 # Client ID from the Prerequisites section
   AB_CLIENT_SECRET='xxxxxxxxxx'             # Client Secret from the Prerequisites section
   AB_NAMESPACE='xxxxxxxxxx'                 # Namespace ID from the Prerequisites section
   PLUGIN_GRPC_SERVER_AUTH_ENABLED=true     # Enable or disable access token validation
   ```

   > :exclamation: **PLUGIN_GRPC_SERVER_AUTH_ENABLED is `true` by default**: If it is set to `false`, the gRPC server can be invoked without an AGS access token. This option is provided for development purposes only. It is recommended to enable gRPC server access token validation in production.

## Build

```shell
make build
```

## Run

```shell
docker compose up --build
```
