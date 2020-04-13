[![License](https://img.shields.io/github/license/kerti/cloudflare-ddns?style=for-the-badge)](https://github.com/kerti/cloudflare-ddns/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/kerti/cloudflare-ddns?style=for-the-badge)](https://goreportcard.com/report/github.com/kerti/cloudflare-ddns)
[![Maintainability](https://img.shields.io/codeclimate/maintainability-percentage/kerti/cloudflare-ddns?style=for-the-badge)](https://codeclimate.com/github/kerti/cloudflare-ddns/maintainability)
[![Build Status](https://img.shields.io/travis/kerti/cloudflare-ddns/master?style=for-the-badge)](https://travis-ci.org/kerti/cloudflare-ddns)
[![Coverage Status](https://img.shields.io/coveralls/github/kerti/cloudflare-ddns?style=for-the-badge)](https://coveralls.io/github/kerti/cloudflare-ddns?branch=master)

# Cloudflare Dynamic DNS

Simple standalone DDNS updater using Cloudflare.

# Providers

Here's a list of providers that enable us to check our external IP address. To prevent bogging down a single server,
we hit each one in a round-robin fashion and set the interval accordingly.

| Service              | URL                           | Type |
-----------------------|-------------------------------|------|
| Big Data Cloud       | https://www.bigdatacloud.com  | JSON |
| I Can Haz IP         | http://icanhazip.com          | Text |
| Ifconfig Me          | https://ifconfig.me           | Text |
| IP API               | https://ipapi.co              | Text |
| IPify                | https://www.ipify.org         | Text |
| My External IP       | https://myexternalip.com      | Text |
| My IP                | https://www.myip.com          | JSON |
| What's My IP Address | https://whatismyipaddress.com | Text |
| WTF Is My IP         | https://wtfismyip.com         | Text |

# Notifications

## IFTTT

You can use IFTTT to hook up this DNS updater to basically anything that IFTTT supports.

* Sign in to [IFTTT](https://ifttt.com)
* Create a new applet
* Use Webhooks as the triggering service
* Choose **Receive a web request**
* Enter an event name
* Choose the action service
* Complete your setup
* Set your IFTTT maker key in the config file

# Roadmap

- [ ] Simplify codebase
  - [ ] Do away with cloudflare wrapper and just do it in worker
- [ ] Do I need to do anything asynchronously?
- [ ] Notifiers
  - [ ] Email
  - [ ] and more...
- [ ] Unit tests
- [ ] Code coverage