# cloudflare-ddns

[![Go Report Card](https://goreportcard.com/badge/github.com/kerti/cloudflare-ddns)](https://goreportcard.com/report/github.com/kerti/cloudflare-ddns)
[![Maintainability](https://api.codeclimate.com/v1/badges/77883b0508313dc1ba32/maintainability)](https://codeclimate.com/github/kerti/cloudflare-ddns/maintainability)

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
  - [ ] Use single class for simple IP lookup provider
  - [ ] Do away with cloudflare wrapper and just do it in worker
- [ ] Do I need to do anything asynchronously?
- [ ] Optimize binary executable size
- [ ] Notifiers
  - [ ] Email
  - [ ] Telegram
  - [ ] and more...
- [ ] Code linter/vetter/checker
- [ ] Unit tests
- [ ] Code coverage
- [ ] CI integration