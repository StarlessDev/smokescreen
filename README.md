# Smokescreen
A simple CLI utility to manage Cloudflare email address routing written in go.

## Why?
When I bought my domain I figured out that you could use Cloudflare's email routing feature to generate **temporary email addresses** and route the incoming emails to a real email address. This way, if your email gets into the hands of advertisers, you can just **revoke** the address and **prevent your inbox from being spammed**.

*aaaand I have started learning go recently, so I wanted to create a small project to get used to the language (the first version of smokescreen was written in Java)*

## Download
All the releases are compiled for **linux**. 
1. Download the latest stable release from the [Releases section](https://github.com/StarlessDev/smokescreen/releases/latest) or from Github actions
2. Make the binary executable
3. (Optional) Add the binary directory to your `$PATH` environment variable to use it from anywhere. 

## How to use
To start using smokescreen you need to add your domain to Cloudflare and setup email routing.

After you are done you need:
- A Cloudflare API Key
- Your ZoneID (visible at the bottom-right of your Cloudflare domain dashboard)
- A real email address which you setup for email routing beforehand

Use the command `$ ./smokescreen addidentity <identity>`, where `<identity>` is the name of your identity, and follow the steps. The program will prompt you for your details:
```
Insert your API Token: <token>
Insert your zone ID: <zone_id>
Insert your domain: example.org
Where do you want the emails to be redirected: your@email.com
```

Now you can use `$ ./smokescreen <identity> test` to generate an address like `test-90213@example.org`, which will redirect the incoming emails to `your@email.com`.

## Commands
Use `$ ./smokescreen -h` for help:
```
CLI utility to manage email aliases using Cloudflare's email routing feature.

Usage:
  smokescreen [command]

Available Commands:
  addidentity    Add your Cloudflare token and zone id to start managing your emails
  completion     Generate the autocompletion script for the specified shell
  gen            Generate an email address using a certain identity.
  help           Help about any command
  list           List the emails you created
  removeidentity Remove a previously added identity
  revoke         Revoke a generated email. This will completely delete the address!

Flags:
  -h, --help   help for smokescreen

Use "smokescreen [command] --help" for more information about a command.
```
