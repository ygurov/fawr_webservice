package route

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/fawrwebservice/model"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
	"gorm.io/gorm"
)

var stripeWebhookKey = os.Getenv("STRIPE_WEBHOOK_KEY")

type StripeWebhook struct {
	DB *gorm.DB
}

func (route *StripeWebhook) Register(parent *http.ServeMux) {
	parent.HandleFunc("/stripe_webhook", route.stripeWebhookHandle)
}

func (route *StripeWebhook) setCommentBought(id int) {
	var comment model.Comment
	err := route.DB.Where("id = ?", id).First(&comment).Error
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting comment from DB: %v\n", err)
		return
	}

	comment.Bought = true

	err = route.DB.Updates(&comment).Error
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating comment in DB: %v\n", err)
	}
}

func (route *StripeWebhook) stripeWebhookHandle(w http.ResponseWriter, req *http.Request) {
	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	event, err := webhook.ConstructEventWithOptions(payload, req.Header.Get("Stripe-Signature"), stripeWebhookKey,
		webhook.ConstructEventOptions{IgnoreAPIVersionMismatch: true})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if event.Type != "payment_intent.succeeded" {
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var paymentIntent stripe.PaymentIntent
	err = json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	commentid, err := strconv.Atoi(paymentIntent.Metadata["commentid"])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing comment id: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go route.setCommentBought(commentid)
	w.WriteHeader(http.StatusOK)
}
