package main

import (
	"database/sql"
	_ "embed"
  "fmt"
	"time"
  "log"
  "net/smtp"

	_ "github.com/mattn/go-sqlite3"
	. "github.com/workspace-9/erk"
	. "github.com/workspace-9/ik"
  "git.jordanbonecutter.com/bibleplan/backend/calendar"
  "github.com/BurntSushi/toml"
)

//go:embed config.toml
var configData string

var config struct {
  Email struct {
    User, Password, Addr, Host string
  }
  RenewLink string
}

func main() {
  db := Try(sql.Open("sqlite3", "bibleplan.db"))
  now := float64(time.Now().UnixMilli())/1000.
  Must(loadConfig())
  log.Println(config)

  for scan := range Sql(Try(db.Query(`
    SELECT email, start_time FROM subscribers;
  `))) {
    var email string
    var startTime float64
    Must(scan(&email, &startTime))
    daysSinceSubscription := int((now - startTime)/(3600*24))
    log.Println(email, daysSinceSubscription)
    if daysSinceSubscription < 0 {
      continue
    }
    if daysSinceSubscription >= 365 {
      Must(sendCongratsEmail(email))
      Try(db.Exec(`
        DELETE FROM subscribers WHERE email = $1;
      `, email))
    } else {
      Must(sendDailyEmail(email, calendar.MCheyne[daysSinceSubscription]))
    }
  }
}

func sendCongratsEmail(to string) error {
  return smtp.SendMail(
    config.Email.Addr,
    smtp.PlainAuth("", config.Email.User, config.Email.Password, config.Email.Host),
    config.Email.User, []string{to}, basicEmail(config.Email.User, to, "You've Finished your daily bible reading!", `Congrats on finishing your bible plan!
Visit %s to renew another year of your plan!
  `, config.RenewLink))
}

func sendDailyEmail(to string, plan [4]calendar.Passage) error {
  return smtp.SendMail(
    config.Email.Addr,
    smtp.PlainAuth("", config.Email.User, config.Email.Password, config.Email.Host),
    config.Email.User, []string{to}, basicEmail(config.Email.User, to, "Daily Bible Reading", `
Your daily bible readings are:
%s
%s
%s
%s
    `, plan[0], plan[1], plan[2], plan[3]))
}

func basicEmail(
  from, to, subject, body string,
  args ...any,
) []byte {
  return []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n", from, to, subject) + fmt.Sprintf(body, args...))
}

func loadConfig() error {
  _, err := toml.Decode(configData, &config)
  return err
}
