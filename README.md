[![License](https://img.shields.io/github/license/kerti/cloudflare-ddns?style=for-the-badge)](https://github.com/kerti/cloudflare-ddns/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/kerti/cloudflare-ddns?style=for-the-badge)](https://goreportcard.com/report/github.com/kerti/cloudflare-ddns)
[![Maintainability](https://img.shields.io/codeclimate/maintainability-percentage/kerti/cloudflare-ddns?style=for-the-badge)](https://codeclimate.com/github/kerti/cloudflare-ddns/maintainability)
[![Build Status](https://img.shields.io/travis/kerti/cloudflare-ddns/master?style=for-the-badge)](https://travis-ci.org/kerti/cloudflare-ddns)
[![Coverage Status](https://img.shields.io/coveralls/github/kerti/cloudflare-ddns?style=for-the-badge)](https://coveralls.io/github/kerti/cloudflare-ddns?branch=master)

# Cloudflare Dynamic DNS

Simple standalone DDNS updater using Cloudflare.

# Roadmap

- [x] Multiple provider support
  - [x] Big Data Cloud
  - [x] I Can Haz IP
  - [x] Ifconfig Me
  - [x] IP API
  - [x] IPify
  - [x] My External IP
  - [x] My IP
  - [x] What's My IP Address
  - [x] WTF Is My IP
  - [ ] and more...
- [x] Round-robin checking
- [x] Automatically create A records
- [x] Automatically update A records
- [ ] Simplify codebase
  - [x] Use single class for simple IP lookup provider
  - [ ] Do away with cloudflare wrapper and just do it in worker
- [ ] Do I need to do anything asynchronously?
- [x] Optimize binary executable size
- [ ] Notifiers
  - [ ] Email
  - [ ] Telegram
  - [ ] and more...
- [x] Code linter/vetter/checker
- [ ] Unit tests
- [x] Code coverage
- [x] CI integration