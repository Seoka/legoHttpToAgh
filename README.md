# DNS-Resolver with AdGuard Home Integration

This program can be used as a DNS resolver in Traefik and forwards the corresponding rules to AdGuard Home for filtering. The program listens for HTTP requests on the endpoints `/present` and `/cleanup`.
## Installation

To install the program, follow these steps:

1. Clone the repository to your local machine.
2. Set the required environment variables:
  - `ADGUARD_URL`: The URL of the AdGuard Home instance.
  - `ADGUARD_USER`: The username to access the AdGuard Home API.
  - `ADGUARD_PASS`: The password to access the AdGuard Home API.
3. Build the program with the `go build` command. 
4. Run the program with the `./<program-name>` command.

## Usage

The program listens for HTTP requests on the following endpoints:`/present`

To add a new filtering rule, send a POST request to the `/present` endpoint with the following JSON payload:

```json
{
"fqdn": "example.com",
"value": "Some value"
}
```
The program will generate a filter rule and send it to AdGuard Home for filtering. If the filter rule already exists, it will be added.
`/cleanup`

To remove a filtering rule, send a POST request to the `/cleanup` endpoint with the following JSON payload:

```json
{
"fqdn": "example.com",
"value": "Some value"
}
```
The program will generate a filter rule and send it to AdGuard Home for removal. If the filter rule does not exist, no action will be taken. If it exists twice, both will be removed.
## License

This program is licensed under the MIT license. See the LICENSE file for details.