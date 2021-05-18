package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"log"
	"net/http"
	"encoding/json"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

type Cryptoresponse struct {
	Sourceftm struct {
		Value float64 `json:"usd"`
	} `json:"fantom"`
	Sourcespirit struct {
		Value float64 `json:"usd"`
	} `json:"spiritswap"`
}
 
 //TextOutput is exported,it formats the data to plain text.
 func (c Cryptoresponse) TextOutput(coin string) string {
	var p string
	if coin == "fantom" {
		p= fmt.Sprintf(
		"Fantom(FTM): %f", c.Sourceftm.Value)
	} else if coin == "spiritswap" {
		p= fmt.Sprintf(
		"SpiritSwap(SPIRIT): %f", c.Sourcespirit.Value)
	}
	
	return p
 }

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func FetchCrypto(crypto string) (string, error) {	
		req, reqErr := http.NewRequest("GET", "https://api.coingecko.com/api/v3/simple/price?ids=" + crypto + "&vs_currencies=usd", nil)
		if reqErr != nil {
			log.Fatal("ooopsss an error occurred, please try again", reqErr)
		}
	
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
	
		resp, respErr := client.Do(req)
		if respErr != nil {
			log.Fatal("ooopsss an error occurred, please try again", respErr)
		}
		defer resp.Body.Close()

		var cResp Cryptoresponse
		if decErr := json.NewDecoder(resp.Body).Decode(&cResp); decErr != nil {
			log.Fatal("ooopsss! an error occurred, please try again", decErr)
		}
	    return cResp.TextOutput(crypto), nil
}

func main() {

	// Create a new Discord session using the provided bot token.
	fmt.Println("token", Token)
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)


	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ftm" {
		fetchResp,fetchErr := FetchCrypto("fantom")
		if fetchErr != nil {
			log.Println(fetchErr)
		}
		s.ChannelMessageSend(m.ChannelID, fetchResp)
		// s.UpdateGameStatus(0, fetchResp)
		s.UpdateListeningStatus(fetchResp)
	}

	if m.Content == "spirit" {
		fetchResp,fetchErr := FetchCrypto("spiritswap")
		if fetchErr != nil {
			log.Println(fetchErr)
		}
		s.ChannelMessageSend(m.ChannelID, fetchResp)
		s.UpdateListeningStatus(fetchResp)
	}
}