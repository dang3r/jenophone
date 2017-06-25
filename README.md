# Jenophone

### Overview

Jenophone acts as an SMS forwarder between two specified numbers using Twilio. It is a server that is hit by Twilio everytime a Twilio phone number receives a text. Based on the sender of that text, it will forward it to the other party (ie. user1 -> TwilioNum -> user2 or in the other direction)

### Why

For my use-case, one party has limited data and does not have unlimited text to one country. The other party has both. Thus, getting a Twilio number in the first party's country allows the first party to use their domestic unlimited text, and the other user to text them back.

### How to use

##### Twilio

- Requires a **Twilio** account with a phone number, and a configured **TwiML** app that points toward the server running Jenophone for SMS via HTTP GET Requests

##### Infra

- Requires a server running Jenophone that is reachable on the public internet (Jenophone also uses TLS)

#### Jenophone

```
Usage of ./jenophone:
  -accountsid string
    Twilio Account Sid
  -num1 string
    First phone number
  -num2 string
    Second phone number
```
