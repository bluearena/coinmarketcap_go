package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"net/http"
)

var (
	baseURL    = "https://api.coinmarketcap.com/v1"
	url        string
	coin_id    string
	coin_names = map[string]string{"BTC": "bitcoin", "ETH": "ethereum", "BCH": "bitcoin-cash", "XRP": "ripple", "DASH": "dash", "LTC": "litecoin", "NEM": "nem", "MIOTA": "iota", "XMR": "monero", "ETC": "ethereum-classic", "NEO": "neo", "OMG": "omisego", "BCC": "bitconnect", "LSK": "lisk", "USDT": "tether", "QTUM": "qtum", "ZEC": "zcash", "STRAT": "stratis", "WAVES": "waves", "ARK": "ark", "STEEM": "steem", "MAID": "maidsafecoin", "BCN": "bytecoin-bcn", "EOS": "eos", "GNT": "golem-network-tokens", "BAT": "basic-attention-token", "DCR": "decred", "REP": "augur", "BTS": "bitshares", "XLM": "stellar", "PAY": "tenx", "HSR": "hshare", "KMD": "komodo", "VERI": "veritaseum", "MTL": "metal", "PIVX": "pivx", "ICN": "iconomi", "FCT": "factom", "NXS": "nexus", "DGD": "digixdao", "GBYTE": "byteball", "ARDR": "ardor", "CVC": "civic", "SC": "siacoin", "PPT": "populous", "DGB": "digibyte", "SNGLS": "singulardtv", "GAS": "gas", "GNO": "gnosis-gno", "BTCD": "bitcoindark", "GAME": "gamecredits", "GXS": "gxshares", "LKK": "lykke", "ZRX": "0x", "BLOCK": "blocknet", "DOGE": "dogecoin", "BNT": "bancor", "FUN": "funfair", "AE": "aeternity", "DCN": "dentacoin", "SNT": "status", "XVG": "verge", "SYS": "syscoin", "MCO": "monaco", "BNB": "binance-coin", "BTM": "bytom", "BQX": "bitquence", "FRST": "firstcoin", "NXT": "nxt", "IOC": "iocoin", "EDG": "edgeless", "LINK": "chainlink", "ANT": "aragon", "UBQ": "ubiq", "PART": "particl", "WINGS": "wings", "NAV": "nav-coin", "RISE": "rise", "MGO": "mobilego", "VTC": "vertcoin", "STORJ": "storj", "CFI": "cofound-it", "TNT": "tierion", "BDL": "bitdeal", "RLC": "rlc", "ETP": "metaverse", "NLG": "gulden", "XZC": "zcoin", "FAIR": "faircoin", "CLOAK": "cloakcoin", "PLR": "pillar", "MLN": "melon", "XEL": "elastic", "TRIG": "triggers", "NLC2": "nolimitcoin", "MTH": "monetha", "WTC": "walton", "PPC": "peercoin", "XRL": "rialto", "LRC": "loopring"}
)

//Coin struct
type Coin struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Symbol           string  `json:"symbol"`
	Rank             int     `json:"rank,string"`
	PriceUsd         float64 `json:"price_usd,string"`
	PriceCny         float64 `json:"price_cny,string"`
	PriceBtc         float64 `json:"price_btc,string"`
	Two4HVolumeUsd   float64 `json:"24h_volume_usd,string"`
	MarketCapUsd     float64 `json:"market_cap_usd,string"`
	Two4HVolumeCny   float64 `json:"24h_volume_cny,string"`
	MarketCapCny     float64 `json:"market_cap_cny,string"`
	AvailableSupply  float64 `json:"available_supply,string"`
	TotalSupply      float64 `json:"total_supply,string"`
	PercentChange1H  float64 `json:"percent_change_1h,string"`
	PercentChange24H float64 `json:"percent_change_24h,string"`
	PercentChange7D  float64 `json:"percent_change_7d,string"`
	LastUpdated      string  `json:"last_updated"`
}

//GlobalMarketData struct
type GlobalMarketData struct {
	TotalMarketCapUsd float64 `json:"total_market_cap_usd"`
	Total24HVolumeUsd float64 `json:"total_24h_volume_usd"`
	BitcoinDominance  float64 `json:"bitcoin_percentage_of_market_cap"`
	ActiveCurrencies  int     `json:"active_currencies"`
	ActiveAssets      int     `json:"active_assets"`
	ActiveMarkets     int     `json:"active_markets"`
}

//Client
func doReq(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}

	return body, nil
}

func makeReq(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	resp, err := doReq(req)
	if err != nil {
		log.Println(err)
	}

	return resp, err
}

func GetCoinInfoFromCMC(coin string, convert string) (Coin, error) {
	if len(convert) == 0 {
		convert = "CNY"
	}
	url = fmt.Sprintf("%s/ticker/%s?convert=%s", baseURL, coin, convert)
	resp, err := makeReq(url)
	if err != nil {
		log.Println(err)
		return Coin{}, err
	}
	var data []Coin
	err = json.Unmarshal(resp, &data)
	if err != nil {
		log.Println(err)
		return Coin{}, err
	}

	return data[0], nil
}

//GetCoinPrice - Get price information from Coin Market Cap and make a message
//func GetCoinPrice(s string) (msg string, err error) {
//	coin, err := cmcAPI.GetCoinInfo(s)
//	if err != nil {
//		msg = "Can't find anything. Type the full name of the coin"
//		return
//	}
//
//	name := fmt.Sprintf("Name: %s | %s | #%d\n", coin.Name, coin.Symbol, coin.Rank)
//	price := fmt.Sprintf("PriceBTC: %f\nPriceUSD: %.2f\n", coin.PriceBtc, coin.PriceUsd)
//	change := fmt.Sprintf("Change 1H/24H/7d: %.2f | %.2f | %.2f\n", coin.PercentChange1H, coin.PercentChange24H, coin.PercentChange7D)
//
//	msg = name + price + change
//
//	return
//}

//GetCoinInfo - Get full infomartion from Coin Market Cap and make a message
func GetCoinInfo(s string) (msg string, err error) {
	s = strings.ToUpper(s)
	coin_id = coin_names[s]
	log.Printf("[Coin Name] %s", coin_id)
	if len(coin_id) > 0 {

		coin, err1 := GetCoinInfoFromCMC(coin_id,"CNY")
		if err1 != nil {
			msg = "未找到数字货币 " + s + " 。请输入数字货币全称 如 bitcoin hshare"
			err = err1
			return
		}

		name := fmt.Sprintf("名称: %s | %s | #%d\n", coin.Name, coin.Symbol, coin.Rank)
		supply := fmt.Sprintf("现有总量: %d\n", int(coin.AvailableSupply))
		mcap := fmt.Sprintf("市值(美金): %d\n", int(coin.MarketCapUsd))
		mcap_cny := fmt.Sprintf("市值(人民币): %d\n", int(coin.MarketCapCny))
		volume := fmt.Sprintf("24小时成交额(美金): %d\n", int(coin.Two4HVolumeUsd))
		volume_cny := fmt.Sprintf("24小时成交额(人民币): %d\n", int(coin.Two4HVolumeCny))
		price := fmt.Sprintf("BTC 价格: %f\n美元价格: %.2f\n人民币价格: %.2f\n", coin.PriceBtc, coin.PriceUsd, coin.PriceCny)
		change := fmt.Sprintf("涨幅 1小时/24小时/7天: %.2f | %.2f | %.2f\n", coin.PercentChange1H, coin.PercentChange24H, coin.PercentChange7D)

		msg = name + supply + mcap + mcap_cny + volume + volume_cny + price + change
	}else{
		msg = "未找到数字货币 " + s + " 。请输入数字货币符号 如 btc hsr"
		return
	}

	return
}

func main() {
	//GetCoinInfo("BTC")
	bot, err := tgbotapi.NewBotAPI("392886056:AAFJI2_snuukiF--XnnFGsxdC-o5eziK9vI")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			commandArgs := update.Message.CommandArguments()

			switch update.Message.Command() {
			case "help":
				//msg.Text = "there is no help for the people here\nbut you can try /status /info /price"
				msg.Text = "/price 币种符号(如 btc hsr bch(比特币现金) 等) 显示数字货币价格\n/status 查看机器人状态(是否活着:P)"
			case "status":
				msg.Text = "我很好，谢谢检查我的健康状态:P"
			case "price":
				msg.Text, err = GetCoinInfo(commandArgs)
				if err != nil {
					log.Println(err)
				}
			case "price_ltc":
				msg.Text, err = GetCoinInfo("ltc")
				if err != nil {
					log.Println(err)
				}
			case "price_eth":
				msg.Text, err = GetCoinInfo("eth")
				if err != nil {
					log.Println(err)
				}
			case "price_bch":
				msg.Text, err = GetCoinInfo("bch")
				if err != nil {
					log.Println(err)
				}
			case "price_hsr":
				msg.Text, err = GetCoinInfo("hsr")
				if err != nil {
					log.Println(err)
				}
			default:
				msg.Text = "无效的指令，请输入/help查询可用指令"
			}
			bot.Send(msg)
		}
	}
}
