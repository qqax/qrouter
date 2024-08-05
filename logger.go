package qrouter

import (
	"github.com/fatih/color"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func requestLogger(next http.Handler) http.Handler {
	h := hlog.NewHandler(log.Logger)

	accessHandler := hlog.AccessHandler(
		func(r *http.Request, status, size int, duration time.Duration) {
			var method, statusCode string
			white := color.New(color.FgWhite)
			black := color.New(color.FgBlack)
			green := color.New(color.FgGreen)
			red := color.New(color.FgRed)
			cyan := color.New(color.FgCyan)
			greenBackground := white.Add(color.BgGreen)
			whiteBackground := black.Add(color.BgWhite)
			blueBackground := white.Add(color.BgBlue)
			cyanBackground := white.Add(color.BgCyan)
			redBackground := white.Add(color.BgRed)
			yellowBackground := white.Add(color.BgYellow)
			switch r.Method {
			case "GET":
				method = greenBackground.Sprint("GET")
			case "POST":
				method = whiteBackground.Sprint("POST")
			case "PATCH":
				method = blueBackground.Sprint("PATCH")
			case "PUT":
				method = cyanBackground.Sprint("PUT")
			case "DELETE":
				method = redBackground.Sprint("DELETE")
			case "HEAD":
				method = yellowBackground.Sprint("HEAD")
			default:
				method = r.Method
			}

			switch status {
			case 200:
				statusCode = green.Sprint(status)
			default:
				statusCode = red.Sprint(status)
			}

			hlog.FromRequest(r).Info().
				Msgf("%s %s %s %s%v %s%v",
					method, r.URL, statusCode, cyan.Sprint("elapsed="), duration, cyan.Sprint("size_byte="), size)
		},
	)

	userAgentHandler := hlog.UserAgentHandler("http_user_agent")

	return h(accessHandler(userAgentHandler(next)))
}
