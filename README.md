# pknulms-notifier

Pukyong National University LMS notifier for Slack

## Installation

```bash
$ go get github.com/hallazzang/pknulms-notifier
```

This command will place `pknulms-notifier` executable into `$GOPATH/bin`.

## Configuration

You should provide configuration file to run the program.
Open your favorite editor and write one:

```ini
interval = 30 ; Crawling interval in seconds

[lms]
id = YOUR_STUDENT_NO ; Your LMS ID(=student number)
pw = YOUR_PASSWORD ; Your LMS password

[slack]
webhook-url = SLACK_WEBHOOK_URL ; "https://hooks.slack.com/services/.../.../..."
```

Save it as `config.ini` or whatever you want.

## Run

If you've added `$GOPATH/bin` to your `$PATH`, you can simply run it:
```bash
$ pknulms-notifier -config=/path/to/config.ini
```

Otherwise, you would do like:
```bash
$ $GOPATH/bin/pknulms-notifier -config=/path/to/config.ini
```
