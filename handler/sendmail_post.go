package handler

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type user struct {
	Username string `json:"username"`
	Mail     string `json:"mailadress"`
	MailType string `json:"mailtype"`
}

type mailContent struct {
	Name        string
	HeaderLogo  string
	HeaderImage string
}

// use viper package to read .env file
// return the value of the key
func viperEnvVariable(key string, whichEnv string) string {

	// SetConfigFile explicitly defines the path, name and extension of the config file.
	// Viper will use this and not check any of the config paths.
	// .env - It will search for the .env file in the current directory
	viper.SetConfigFile("./config/" + whichEnv)

	// Find and read the config file
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	// viper.Get() returns an empty interface{}
	// to get the underlying type of the key,
	// we have to do the type assertion, we know the underlying value is string
	// if we type assert to other type it will throw an error
	value, ok := viper.Get(key).(string)

	// If the type is a string then ok will be true
	// ok will make sure the program not break
	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	return value
}

func (c mailContent) sendMail(userMailAdress string) error {
	host := viperEnvVariable("SMTP_HOST", "smtp.env")
	stringPort := viperEnvVariable("SMTP_PORT", "smtp.env")
	username := viperEnvVariable("SMTP_USERNAME", "smtp.env")
	password := viperEnvVariable("SMTP_PASSWORD", "smtp.env")

	// While sending email, port variable have to be integer type.
	// While getting value of an environment variable it returns string.
	// So we have to convert to integer it.
	port, err := strconv.Atoi(stringPort)
	if err != nil {
		log.Fatalf("Error while converting port value to integer %s", err)
	}

	t := template.New("subscribe.html")

	t, err = t.ParseFiles("templates/subscribe.html")
	if err != nil {
		log.Fatalf("Error while parsing template %s", err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, c); err != nil {
		log.Println(err)
	}

	result := tpl.String()

	m := gomail.NewMessage()
	m.SetHeader("From", "itsbaseline.testing.dev@gmail.com")
	m.SetHeader("To", userMailAdress)
	m.SetHeader("Subject", "Welcome to Line.")
	m.SetBody("text/html", result)
	d := gomail.NewDialer(host, port, username, password)

	if err := d.DialAndSend(m); err != nil {
		log.Fatalf("Error while sending email %s", err)
	}

	return err
}

//SendMailPost function
func SendMailPost(c *gin.Context) {

	var u user
	c.BindJSON(&u)

	d := mailContent{
		u.Username,
		viperEnvVariable("HEADER_LOGO", u.MailType+".env"),
		viperEnvVariable("HEADER_IMG", u.MailType+".env"),
	}

	if err := d.sendMail(u.Mail); err != nil {
		log.Fatalf("Error while sending email %s", err)
	}

	c.JSON(http.StatusOK, map[string]string{
		"message": "200",
	})

}
