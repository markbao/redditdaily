package redditdaily

import(
	"fmt"
	"sync"
	"time"
	"bytes"
	"text/template"
	"github.com/jzelinskie/geddit"
	"github.com/spf13/viper"
	"github.com/sendgrid/sendgrid-go"
)

var content map[string][]*geddit.Submission
var delay int

func Run() {
	// Config - todo: should accept like a real package
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("../")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// TODO: move to config
	subreddits := viper.GetStringSlice("subreddits")

	content = make(map[string][]*geddit.Submission)

	// Delay fetching by 500ms to fetch 2 per second
	delay = 0

	// Wait group for waiting until all have finished processing
	var wg sync.WaitGroup

	for _,subreddit := range subreddits {
		content[subreddit] = nil
		wg.Add(1)
		fmt.Println("About to fetch " + subreddit)
		go fetchSubreddit(subreddit, delay, &wg)
		delay = delay + 500
	}

	wg.Wait()

	// Apply template and send
	emailBody := applyTemplate(content)
	sendEmail(emailBody)
}

func applyTemplate(emailContent map[string][]*geddit.Submission) (string) {
	var emailBuffer bytes.Buffer

	// Templatize
	t := template.New("email.html")
	t, _ = template.ParseFiles("email.html")
	t.Execute(&emailBuffer, emailContent)

	emailBody := emailBuffer.String()

	fmt.Print(emailBody)

	return emailBody
}

func sendEmail(emailBody string) {
	// Format date
	time := time.Now()
	todayDate := time.Format("Monday, January 2, 2006")

	// Format subject
	subject := "Reddit Daily — " + todayDate
	fmt.Println(viper.GetString("sendgrid_key"))

	sg := sendgrid.NewSendGridClient(viper.GetString("sendgrid_user"), viper.GetString("sendgrid_key"))
	message := sendgrid.NewMail()
	message.AddTo(viper.GetString("email_to"))
	message.AddToName(viper.GetString("email_to_name"))
	message.SetFrom(viper.GetString("email_from"))
	message.SetFromName(viper.GetString("email_from_name"))
	message.SetSubject(subject)
	message.SetHTML(emailBody)

	if r := sg.Send(message); r == nil {
		fmt.Println("Email sent!")
	} else {
		fmt.Println(r)
	}
}

func fetchSubreddit(subreddit string, delay int, wg *sync.WaitGroup) {
	fmt.Println(fmt.Sprintf("Sleeping %v before fetching %v...", delay, subreddit))
	time.Sleep(time.Duration(delay) * time.Millisecond)

	fmt.Println("Fetching subreddit " + subreddit)
	submissions, err := GetSubmissions(subreddit)

	if err == nil {
		content[subreddit] = submissions
	} else {
		fmt.Println(err)
	}

	wg.Done()
	return
}

func GetSubmissions(subreddit string) ([]*geddit.Submission, error) {
	// Top in 24 hours, limit 10
	options := geddit.ListingOptions{
		Limit: 10,
		Time: geddit.ThisDay,
	}

	session := geddit.NewSession(viper.GetString("reddit_user_agent"))

	submissions, err := session.SubredditSubmissions(subreddit, geddit.TopSubmissions, options)

	if err != nil {
		return nil, err
	}

	fmt.Println(submissions)

	return submissions, nil
}