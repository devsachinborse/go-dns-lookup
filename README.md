# DNS Lookup Tool

A simple DNS Lookup Tool built using Golang and HTMX. This tool allows users to input one or multiple domain names (comma-separated) and retrieve detailed DNS records, including A, AAAA, CNAME, MX, and NS records, displayed in a scrollable table format with styled IP addresses. The project also includes a footer with copyright information and is ready to be deployed.

## Features

- Input multiple domain names separated by commas.
- Displays various DNS record types: A, AAAA, CNAME, NS, MX.
- IP addresses are displayed with distinct background colors.
- Footer with copyright information.

## Run this project
#### Local
```
go mod tidy
make run
```
#### Docker
```
# Build the Docker image
docker build -t my-dns-lookup-tool .

# Run the Docker container
docker run -p 9001:9001 my-dns-lookup-tool

```
![Screenshot 2024-08-23 112658](https://github.com/user-attachments/assets/519eb826-6d9a-4fad-8ad2-68c5d4b3b3af)




