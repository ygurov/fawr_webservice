package api

import (
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

	http.ListenAndServe(addr, corsMiddleware(mux))
}
