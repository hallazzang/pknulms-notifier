package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/hallazzang/pknulms"
	ini "gopkg.in/ini.v1"
)

type config struct {
	id              string
	pw              string
	slackWebhookURL string
	interval        int
}

func loadConfig(path string) (*config, error) {
	cfg, err := ini.InsensitiveLoad(path)
	if err != nil {
		return nil, err
	}

	id, err := cfg.Section("lms").GetKey("id")
	if err != nil {
		return nil, err
	}
	pw, err := cfg.Section("lms").GetKey("pw")
	if err != nil {
		return nil, err
	}

	webhookURL, err := cfg.Section("slack").GetKey("webhook-url")
	if err != nil {
		return nil, err
	}

	interval, err := cfg.Section("").GetKey("interval")
	if err != nil {
		return nil, err
	}

	return &config{
		id:              id.String(),
		pw:              pw.String(),
		slackWebhookURL: webhookURL.String(),
		interval:        interval.MustInt(60),
	}, nil
}

func printError(format string, args ...interface{}) {
	fmt.Fprintf(color.Output, "[%s] %v %s\n",
		time.Now().Format("01/02 15:04"), color.RedString("ERROR:"), fmt.Sprintf(format, args...))
}

func printInfo(format string, args ...interface{}) {
	fmt.Fprintf(color.Output, "[%s] %v %s\n",
		time.Now().Format("01/02 15:04"), color.CyanString("INFO:"), fmt.Sprintf(format, args...))
}

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "config.ini", "config file path")
	flag.Parse()

	cfg, err := loadConfig(configPath)
	if err != nil {
		printError("Cannot load config file\n  [!] %v", err)
		os.Exit(1)
	}

	client := pknulms.MustNewClient()

	if !client.MustLogin(cfg.id, cfg.pw) {
		printError("Login failed. Check ID and password")
		os.Exit(0)
	} else {
		printInfo("Successfully logged in")
	}

	printInfo("Started monitoring (interval: %d sec(s))", cfg.interval)

	isFirstRun := true
	maxID := -1
	for {
		nts, err := client.GetNotificationsByPage(1)
		if err != nil {
			goto sleepAndContinue
		}

		for i := len(nts) - 1; i >= 0; i-- {
			if nts[i].ID > maxID {
				if !isFirstRun {
					printInfo("New notification has uploaded: %s", nts[i].Title)
					sent, err := sendSlackMessage(cfg.slackWebhookURL,
						fmt.Sprintf("LMS: %s", nts[i].Title))
					if err != nil {
						printError("Couldn't send slack message\n  [!] %v", err)
					} else if !sent {
						printError("Couldn't send slack message\n  [!] Unknown")
					} else {
						printInfo("Sent slack message")
					}
				}
				maxID = nts[i].ID
			}
		}

		isFirstRun = false

	sleepAndContinue:
		time.Sleep(time.Duration(cfg.interval) * time.Second)
	}
}
