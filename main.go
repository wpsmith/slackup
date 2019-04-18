package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	//for extracting service credentials from VCAP_SERVICES
	//"github.com/cloudfoundry-community/go-cfenv"

	"github.com/wpsmith/slash"
)

const (
	DEFAULT_PORT = "8080"
)

func main() {
	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = DEFAULT_PORT
	}

	//h := slash.HandlerFunc(Handle)
	//s := slash.NewServer(h)

	r := slash.NewMux()
	r.Command("/weather", "secrettoken", slash.HandlerFunc(Weather))
	r.Command("/test", "secrettoken", slash.HandlerFunc(Handle))

	s := slash.NewServer(r)

	http.ListenAndServe(":" + port, s)
}

func Handle(ctx context.Context, r slash.Responder, command slash.Command) error {
	if err := r.Respond(slash.Reply("Cool beans")); err != nil {
		return err
	}

	for i := 0; i < 4; i++ {
		<-time.After(time.Second)
		if err := r.Respond(slash.Reply(fmt.Sprintf("Async response %d", i))); err != nil {
			return err
		}
	}

	return nil
}

func printErrors(h slash.Handler) slash.Handler {
	return slash.HandlerFunc(func(ctx context.Context, r slash.Responder, command slash.Command) error {
		if err := h.ServeCommand(ctx, r, command); err != nil {
			fmt.Printf("error: %v\n", err)
		}
		return nil
	})
}

// Weather is the primary slash handler for the /weather command.
func Weather(ctx context.Context, r slash.Responder, command slash.Command) error {
	h := slash.NewMux()

	var zipcodeRegex = regexp.MustCompile(`(?P<zip>[0-9])`)
	h.MatchText(zipcodeRegex, slash.HandlerFunc(Zipcode))

	return h.ServeCommand(ctx, r, command)
}

// Zipcode is a slash handler that returns the weather for a zip code.
func Zipcode(ctx context.Context, r slash.Responder, command slash.Command) error {
	params := slash.Params(ctx)
	zip := params["zip"]
	return r.Respond(slash.Reply(zip))
}