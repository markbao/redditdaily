# redditdaily

Daily email notification for Reddit, written in go.

### Why

So I can block Reddit on my computer and still get content from great subreddits.

### What it looks like

What you get in your inbox: http://output.jsbin.com/digukeredu/

### How to use it

1. Move config.yaml.example to config.yaml
2. Edit the config.yaml with your preferred subreddits, email to/from, user agent (to identify yourself to the Reddit API), and your Sendgrid credentials (free account). Optionally set your preferred cron run time (24h time format).
3. Run redditdaily.go to test it. Then run redditdaily_cron.go to run a long-running go process that will execute redditdaily when you specify it to.
