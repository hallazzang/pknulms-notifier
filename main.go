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
	fmt.Fprintf(color.Output, "%v %s\n", color.RedString("ERROR:"), fmt.Sprintf(format, args...))
}

func printErrorAndExit(format string, args ...interface{}) {
	printError(format, args...)
	os.Exit(1)
}

func printInfo(format string, args ...interface{}) {
	fmt.Fprintf(color.Output, "%v %s\n", color.CyanString("INFO:"), fmt.Sprintf(format, args...))
}

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "config.ini", "config file path")
	flag.Parse()

	cfg, err := loadConfig(configPath)
	if err != nil {
		printErrorAndExit("Cannot load config file\n  [!] %v", err)
	}

	client := pknulms.MustNewClient()

	if !client.MustLogin(cfg.id, cfg.pw) {
		printErrorAndExit("Login failed. Check ID and password")
	} else {
		printInfo("Successfully logged in")
	}

	printInfo("Started monitoring (interval: %d sec(s))", cfg.interval)

	isFirstRun := true
	maxID := -1
	for {
		nts, err := client.GetNotificationsByPage(1)
		if err != nil {
			printError("An error occurred while fetching notifications\n  [!] %v", err)
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
