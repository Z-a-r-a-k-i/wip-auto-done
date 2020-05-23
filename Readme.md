# wip-auto-done
This project is written in golang and will simply check if you did any task today and then automatically add one if it wasn't the case

## Installation
All you need to use this project is the latest golang version installed.
```
git clone https://github.com/Z-a-r-a-k-i/wip-auto-done
cd wip-auto-done
go install ./...
```
You need your wip.chat API key as well, you can find it [here](https://wip.chat/api)

## Usage
```
USAGE
  wip-auto-done [global flags]

FLAGS
  -api-key ...              your wip.chat api key (https://wip.chat/api)
  -message #tryhardinglife  the message body of the completed todo
  -wip-user ...             your username on wip.chat

```
**Example**
```
wip-auto-done -wip-user username -api-key ************************ -message "One more day tryharding life"
```

## Crontab
If you want to run the script every day 11pm (timezone is UTC)
```
0 23 * * * /path/to/go/bin/wip-auto-done -wip-user username -api-key ************************ -message "One more day tryharding life"
```

## Random / Contribute
This project was made for fun while I was bored and my friends were playing Valorant without me.

Main purpose was to learn graphQL basics and because I don't want to loose my streak on wip.chat when I don't do anything interesting enough to be mentioned publicly for 24h #iwnl. 

If you don't like this code, feel free to pm and discuss I love to make new friends !

If you like this code, feel free to pm and discuss I love to make new friends !