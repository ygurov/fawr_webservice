package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/fawrwebservice/api/route"
	"github.com/stripe/stripe-go/v82"
	"gorm.io/gorm"
)

func Register(addr string, db *gorm.DB) {
	stripe.Key = os.Getenv("STRIPE_KEY")

	mux := http.NewServeMux()
	{
		api := http.NewServeMux()
		{
			comment := route.CommentRoute{DB: db}
			comment.Register(api)

			stripewebhook := route.StripeWebhook{DB: db}
			stripewebhook.Register(api)

			pay := route.Pay{DB: db}
			pay.Register(api)
		}
		mux.Handle("/api/", http.StripPrefix("/api", api))

		static := route.StaticRoute{}
		static.Register(mux)
	}

	certPath := os.Getenv("SSL_CERT_PATH")
	keyPath := os.Getenv("SSL_KEY_PATH")

	err := http.ListenAndServeTLS(":443", certPath, keyPath, mux)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while starting webserver: %v\n", err)
	}
}
