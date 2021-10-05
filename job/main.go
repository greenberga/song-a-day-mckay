package main

import (
	"encoding/csv"
	_ "embed"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

//go:embed songs.csv
var content string

// seconds
var maxBackoff float64 = 64

var maxMsgLength = 1600

func chunk(msg string) (ret []string) {
	if len(msg) < maxMsgLength {
		return []string{msg}
	}

	var b strings.Builder
	for _, p := range strings.Split(msg, "\n") {
		p = strings.TrimSpace(p) + "\n\n"
		if b.Len() + len(p) > maxMsgLength {
			ret = append(ret, strings.TrimSpace(b.String()))
			b.Reset()

			b.WriteString(p)
		} else {
			b.WriteString(p)
		}
	}

	if b.Len() > 0 {
		ret = append(ret, strings.TrimSpace(b.String()))
	}

	return
}

func sendMessage(msg string) error {
	sid := os.Getenv("TWILIO_ACCOUNT_SID")
	if sid == "" {
		return fmt.Errorf("missing TWILIO_ACCOUNT_SID")
	}

	token := os.Getenv("TWILIO_AUTH_TOKEN")
	if token == "" {
		return fmt.Errorf("missing TWILIO_AUTH_TOKEN")
	}

	from := os.Getenv("FROM")
	if from == "" {
		return fmt.Errorf("missing FROM")
	}

	to := os.Getenv("TO")
	if to == "" {
		return fmt.Errorf("missing TO")
	}

	v := url.Values{}
	v.Set("Body", msg)
	v.Set("From", from)
	v.Set("To", to)

	s := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", sid)
	u, _ := url.Parse(s)
	u.User = url.UserPassword(sid, token)

	timeout := time.Now().Add(10*time.Minute)
	i := 0

	for {
		resp, err := http.PostForm(u.String(), v)
		if err != nil {
			return fmt.Errorf("doing request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode < 400 {
			return nil
		}

		b, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))

		if resp.StatusCode < 500 {
			return fmt.Errorf("bad request [%d]: %s", resp.StatusCode, string(b))
		}

		if time.Now().After(timeout) {
			return fmt.Errorf("internal server [%d]: %s", resp.StatusCode, string(b))
		}

		seconds := math.Exp2(float64(i))
		millis := float64(rand.Intn(1000)) / 1000
		backoff := math.Min(seconds + millis, maxBackoff)
		time.Sleep(time.Second*time.Duration(backoff))

		i++
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	start := os.Getenv("START_DATE")
	if start == "" {
		log.Fatalf("missing START_DATE")
	}

	t, err := time.Parse("2006-01-02", start)
	if err != nil {
		log.Fatalf("failed to parse start date: %v", err)
	}

	index := int(math.Floor(time.Since(t).Hours() / 24))

	log.Printf("today's index: %d", int(index))

	rows, err := csv.NewReader(strings.NewReader(content)).ReadAll()
	if err != nil {
		log.Fatalf("failed to read csv: %v", err)
	}

	log.Printf("read %d rows", len(rows))

	if index > len(rows) {
		log.Printf("no more songs to send")
		return
	}

	if index == len(rows) {
		if err := sendMessage("That's all folks!"); err != nil {
			log.Fatalf("failed to send message: %v", err)
		}
	}

	row := rows[499]
	msgs := chunk(fmt.Sprintf("🎶 %s\n%s", row[0], strings.TrimSpace(row[1])))

	for _, msg := range msgs {
		if err := sendMessage(msg); err != nil {
			log.Fatalf("failed to send message: %v", err)
		}
	}
}
