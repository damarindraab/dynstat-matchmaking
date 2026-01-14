# Testing

You can test the app locally and against AGS.

## Local Development Testing

> :warning: **Set `PLUGIN_GRPC_SERVER_AUTH_ENABLED=false` for local tests**: Otherwise, gRPC requests will be rejected by the server.

1. Run the app:

   ```shell
   docker compose up --build
   ```

2. Open Postman, create a new gRPC request, and enter `localhost:6565` as the server URL.

   > :warning: **If you are running [grpc-plugin-dependencies](https://github.com/AccelByte/grpc-plugin-dependencies) alongside this project**: Use `localhost:10000` instead of `localhost:6565`. This routes the request through Envoy in the dependencies stack.

### Testing GetStatCodes

3. Select the `GetStatCodes` method and send:

   ```json
   {
       "rules": {
           "json": "{\"statistics_config\":{\"statistics\":[\"mmr_ryu\",\"mmr_ken\",\"mmr_chun-li\"]}}"
       }
   }
   ```

4. Expected response:

   ```json
   {
       "codes": ["mmr_ryu", "mmr_ken", "mmr_chun-li"]
   }
   ```

### Testing ValidateTicket

5. Select `ValidateTicket` and send:

   ```json
   {
       "ticket": {
           "ticket_id": "ticket-123",
           "players": [
               {
                   "player_id": "playerA",
                   "attributes": {
                       "mmr_ryu": 1850
                   }
               }
           ],
           "ticket_attributes": {
               "playerA": "mmr_ryu"
           }
       },
       "rules": {
           "json": "{\"statistics_config\":{\"statistics\":[\"mmr_ryu\",\"mmr_ken\",\"mmr_chun-li\"]}}"
       }
   }
   ```

6. Expected response:

   ```json
   {
       "validTicket": true
   }
   ```

### Testing EnrichTicket

7. Select `EnrichTicket` and send:

   ```json
   {
       "ticket": {
           "ticket_id": "ticket-123",
           "players": [
               {
                   "player_id": "playerA",
                   "attributes": {
                       "mmr_ryu": 1850
                   }
               }
           ],
           "ticket_attributes": {
               "playerA": "mmr_ryu"
           }
       },
       "rules": {
           "json": "{\"statistics_config\":{\"statistics\":[\"mmr_ryu\",\"mmr_ken\",\"mmr_chun-li\"],\"enriched_key\":\"mmr\"}}"
       }
   }
   ```

8. Expected response (ticket enriched with standard `mmr`):

   ```json
   {
       "ticket": {
           "ticket_id": "ticket-123",
           "ticket_attributes": {
               "mmr": 1850
           },
           "players": [...]
       }
   }
   ```

### Testing MakeMatches (Returns UNIMPLEMENTED)

9. Select `MakeMatches` (stream) and click Invoke. Expected error:

   ```
   Error: 12 UNIMPLEMENTED: MakeMatches not implemented - using AGS default matching
   ```

## Testing With AGS

To test against AGS, expose the local gRPC server using a TCP tunnel.

1. Run the app:

   ```shell
   docker compose up --build
   ```

2. Expose port 6565 using a tunnel provider:

   - Using ngrok:

      ```bash
      ngrok tcp 6565
      ```

   - Using pinggy:

      ```bash
      ssh -p 443 -o StrictHostKeyChecking=no -o ServerAliveInterval=30 -R0:127.0.0.1:6565 tcp@a.pinggy.io
      ```

   Keep the forwarding URL (e.g., `tcp://xxxxx-xxx-xxx-xxx-xxx.a.free.pinggy.link:xxxxx`).

   > :warning: **If you are running grpc-plugin-dependencies**: Run the tunnel from the dependencies stack and forward local port `10000` (Envoy) instead of `6565`.

3. [Create an OAuth client](https://docs.accelbyte.io/gaming-services/services/access/authorization/manage-access-control-for-applications/#create-an-iam-client) with `confidential` client type and the following permissions:

   - For AGS Private Cloud:
      - `NAMESPACE:{namespace}:MATCHMAKING:RULES [CREATE,READ,UPDATE,DELETE]`
      - `NAMESPACE:{namespace}:MATCHMAKING:FUNCTIONS [CREATE,READ,UPDATE,DELETE]`
      - `NAMESPACE:{namespace}:MATCHMAKING:POOL [CREATE,READ,UPDATE,DELETE]`
      - `NAMESPACE:{namespace}:MATCHMAKING:TICKET [CREATE,READ,UPDATE,DELETE]`
      - `ADMIN:NAMESPACE:{namespace}:INFORMATION:USER:* [CREATE,READ,UPDATE,DELETE]`
      - `ADMIN:NAMESPACE:{namespace}:SESSION:CONFIGURATION:* [CREATE,READ,UPDATE,DELETE]`
   - For AGS Shared Cloud:
      - Matchmaking -> Rule Sets (Create, Read, Update, Delete)
      - Matchmaking -> Match Functions (Create, Read, Update, Delete)
      - Matchmaking -> Match Pools (Create, Read, Update, Delete)
      - Matchmaking -> Match Tickets (Create, Read, Update, Delete)
      - IAM -> Users (Create, Read, Update, Delete)
      - Session -> Configuration Template (Create, Read, Update, Delete)

   > :warning: **This OAuth client is different from the one in prerequisites**: It is used by the Postman collection to register the gRPC server URL and to manage test users.

4. Import the Postman collection in `demo/matchmaking-function-grpc-plugin-server.postman_collection.json` and follow its overview instructions to run the flow.

## Test Observability

1. Uncomment the Loki logging driver in `docker-compose.yaml`:

   ```
    # logging:
    #   driver: loki
    #   options:
    #     loki-url: http://host.docker.internal:3100/loki/api/v1/push
    #     mode: non-blocking
    #     max-buffer-size: 4m
    #     loki-retries: "3"
   ```

   > :warning: **Install the Docker Loki plugin first**: `docker plugin install grafana/loki-docker-driver:latest --alias loki --grant-all-permissions`.

2. Clone and run [grpc-plugin-dependencies](https://github.com/AccelByte/grpc-plugin-dependencies) alongside this project. Grafana will be available at http://localhost:3000.

   ```bash
   git clone https://github.com/AccelByte/grpc-plugin-dependencies.git
   cd grpc-plugin-dependencies
   docker-compose up
   ```

3. Run tests in either local or AGS mode above.
