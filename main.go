package main

import (
	"bytes"
	"flag"
	"log"
	"net/smtp"
	"os"
	"strings"
)

func SendMail(addr string, from string, to string, msg []byte) error {
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	if err := c.Mail(from); err != nil {
		return err
	}
	if err := c.Rcpt(to); err != nil {
		return err
	}
	wc, err := c.Data()
	if err != nil {
		return err
	}
	_, err = wc.Write(msg)
	if err != nil {
		return err
	}
	err = wc.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func flagOrEnv(flagName string, envName string) string {
	f := flag.Lookup(flagName)
	val := f.Value.String()
	if val == f.DefValue {
		if env, env_ok := os.LookupEnv(envName); env_ok {
			val = env
		}
	}
	return val
}

func main() {
	var headers []string
	addHeader := func(header string) error {
		headers = append(headers, header)
		return nil
	}

	if sep := os.Getenv("SMTPCLI_HEADERS_SEP"); len(sep) > 0 {
		for _, header := range strings.Split(os.Getenv("SMTPCLI_HEADERS"), sep) {
			if len(header) > 0 {
				addHeader(header)
			}
		}
	}

	flag.String("s", "localhost:25", "server")
	flag.String("f", "from@smtpcli", "from")
	flag.String("t", "to@smtpcli", "to")
	flag.String("u", "", "subject")
	flag.String("b", "", "body")
	flag.Func("h", "header, multiple", addHeader)
	debugPtr := flag.Bool("d", false, "debug")
	flag.Parse()

	server := flagOrEnv("s", "SMTPCLI_SERVER")
	from := flagOrEnv("f", "SMTPCLI_FROM")
	to := flagOrEnv("t", "SMTPCLI_TO")
	subject := flagOrEnv("u", "SMTPCLI_SUBJECT")
	body := flagOrEnv("b", "SMTPCLI_BODY")
	debug := *debugPtr

	if len(subject) > 0 {
		addHeader("Subject: " + subject)
	}
	if len(headers) > 0 {
		addHeader("")
	}
	if debug {
		log.Println("Headers:", headers)
	}

	var msg bytes.Buffer
	for _, header := range headers {
		if _, err := msg.WriteString(header + "\n"); err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}
	if _, err := msg.WriteString(body); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	if debug {
		log.Printf("Msg:\n%s\n", msg.String())
	}
	if err := SendMail(server, from, to, msg.Bytes()); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
