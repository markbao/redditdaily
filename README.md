# redditdaily

Daily email notification for Reddit, written in go.

### Why

So I can block Reddit on my computer and still get content from great subreddits.

### What it looks like

What you get in your inbox: http://output.jsbin.com/digukeredu/

### How to use it

1. Move config.yaml.example to config.yaml
2. Edit the config.yaml with your preferred subreddits, email to/from, user agent (to identify yourself to the Reddit API), and your Sendgrid credentials (free account).
3. Run it to test it. Set up a cron job at your specified time.
