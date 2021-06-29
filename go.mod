module InstagramStatistic

// +heroku goVersion go1.16
go 1.16

require (
	github.com/TheForgotten69/goinsta/v2 v2.7.0
	github.com/go-redis/redis/v8 v8.10.0
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
)

//replace github.com/TheForgotten69/goinsta/v2 v2.7.0 => github.com/MatviecKirill/goinsta/v2 v2.7.0
