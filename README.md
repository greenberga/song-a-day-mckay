A "bot" that texts you the title and description of one song per day from
Rolling Stone's ["The 500 Greatest Songs of All Time"][rs]. That article is too
long to handle, plus it loads really slowly on my computer. But I wanted to
listen to all the songs! So I made this. Built with Twilio.

You can build it with

```sh
go build -o song-a-day-mckay job/main.go
```

Create a file (mine's named `.song-a-day-mckay-env`) that exports some required
environment variables. Its content should look like this:

```
export START_DATE='YYYY-MM-DD' # A future date at which to start texting.
export FROM='+XXXXXXXXXXX'     # A Twilio phone number that exists in your account.
export TO='+XXXXXXXXXXX'       # Recipient phone number.
export TWILIO_ACCOUNT_SID='your sid'
export TWILIO_AUTH_TOKEN='your token'
```

And install by running `crontab -e` and entering

```
0 10 * * * . /path/to/.song-a-day-mckay-env; /path/to/song-a-day-mckay
```

at the bottom. That will run it every day at 10 AM.

[rs]: https://www.rollingstone.com/music/music-lists/best-songs-of-all-time-1224767
