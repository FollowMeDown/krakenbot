package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "math"
  "time"

  tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
  "github.com/preichenberger/go-coinbasepro"
  "gopkg.in/yaml.v2"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
  tgbotapi.NewKeyboardButtonRow(
    tgbotapi.NewKeyboardButton("BTC"),
    tgbotapi.NewKeyboardButton("ETH"),
    tgbotapi.NewKeyboardButton("XRP"),
  ),
  tgbotapi.NewKeyboardButtonRow(
    tgbotapi.NewKeyboardButton("EOS"),
    tgbotapi.NewKeyboardButton("SNX"),
    tgbotapi.NewKeyboardButton("LINK"),
  ),
)

var client = coinbasepro.NewClient()

type conf struct {
  APIKey string `yaml:"API_KEY"`
}

func main() {
  c, err := readConf("conf.yaml")
  if err != nil {
    log.Fatal(err)
  }

  client.UpdateConfig(&coinbasepro.ClientConfig{
    BaseURL: "https://api.pro.coinbase.com",
  })

  bot, err := tgbotapi.NewBotAPI(c.APIKey)
  if err != nil {
    log.Panic(err)
  }

  log.Printf("Authorized on account %s", bot.Self.UserName)

  u := tgbotapi.NewUpdate(0)
  u.Timeout = 60
  updates, _ := bot.GetUpdatesChan(u)

  for update := range updates {
    if update.Message == nil {
      continue
    }

    if update.Message.Text == "/start" {
      reply := "Welcome, my name is Kraken! Send me a ticker (e.g. `BTC`) or pick one from the buttons below. You can find my source code [here](www.github.com/sjagoori/krakenbot)."
      message := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
      message.ParseMode = "markdown"
      message.ReplyMarkup = numericKeyboard
      bot.Send(message)
    } else {
      send(bot, update.Message.Chat.ID, update.Message.Text, getCoinPrice(update.Message.Text))
    }
  }
}

func getCoinPrice(coinType string) string {
  bookUSD, err := client.GetBook(coinType+"-USD", 1)
  if err != nil {
    return "Coin not found"
  }

  priceUSD := bookUSD.Asks[0].Price

  stamp := time.Now()
  location, _ := time.LoadLocation("Europe/Amsterdam")

  startOfTheDay := time.Date(stamp.Year(), stamp.Month(), stamp.Day(), int(0), int(0), int(0), int(0), location)
  now := time.Date(stamp.Year(), stamp.Month(), stamp.Day(), stamp.Hour(), stamp.Minute(), stamp.Second(), stamp.Nanosecond(), location)

  data24h := coinbasepro.GetHistoricRatesParams{Start: startOfTheDay, End: now, Granularity: 86400}

  product24h, err := client.GetHistoricRates(coinType+"-USD", data24h)
  if err != nil {
    return fmt.Sprintf("%s", err)
  }

  table := "\n\n*24H:*\t\t\t\t\t\t\t`" + fmt.Sprintf("%v", (math.Round(percentageChange(product24h[0].Open, product24h[0].Close)*100)/100)) + "%" +
    "`\n*Lowest:*\t `" + "$" + fmt.Sprintf("%v", product24h[0].Low) +
    "`\n*Highest:* `" + "$" + fmt.Sprintf("%v", product24h[0].High) +
    "`\n\n_" + stamp.Format("15:04:05 02-01") + "_"

  return "*Price:*\t\t\t\t\t`$" + priceUSD + "`" + table
}

func readConf(filename string) (*conf, error) {
  buf, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  c := &conf{}
  err = yaml.Unmarshal(buf, c)
  if err != nil {
    return nil, fmt.Errorf("in file %q: %v", filename, err)
  }
  return c, nil
}

func percentageChange(old, new float64) float64 {
  diff := float64(new - old)
  return (diff / float64(old)) * 100

}

func send(bot *tgbotapi.BotAPI, chatid int64, coin string, price string) {
  reply := "`" + coin + "`\n\n" + price
  message := tgbotapi.NewMessage(chatid, reply)
  message.ParseMode = "markdown"
  message.ReplyMarkup = numericKeyboard
  bot.Send(message)
}
