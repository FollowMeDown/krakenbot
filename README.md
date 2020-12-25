# Krakenbot
_Krakenbot is a telegram bot that fetches prices for cryptocurrency_

## Dependencies
* [go-coinbasepro](github.com/preichenberger/go-coinbasepro)   
`go get github.com/preichenberger/go-coinbasepro`
* [telegram-bot-api](github.com/go-telegram-bot-api/telegram-bot-api)   
`go get github.com/go-telegram-bot-api/telegram-bot-api`

## Install
0. Install Go   
[Official installer](https://golang.org/doc/install)
```brew install go```
1. Clone the repo.   
```git clone https://github.com/sjagoori/krakenbot```
2. Install dependencies.   
```go get github.com/preichenberger/go-coinbasepro```
```go get github.com/go-telegram-bot-api/telegram-bot-api```
3. Get an Telegram bot API key from [@Botfather](t.me/BotFather)
4. Insert your API key in `config.yaml`, see [example](https://github.com/sjagoori/krakenbot/blob/master/example_conf.yaml)

## Future developement plans
- [ ] Add support for pricing in sats
- [ ] Add support for price-spike notifications 
