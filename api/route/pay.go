package route

import (
	"fmt"
	"net/http"
	"os"

	"github.com/fawrwebservice/model"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"gorm.io/gorm"
)

type Pay struct {
	DB *gorm.DB
}

func (route *Pay) Register(parent *http.ServeMux) {
	parent.HandleFunc("/pay", route.handler)
}

func (route *Pay) handler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		route.redirectHandler(w, req)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (route *Pay) redirectHandler(w http.ResponseWriter, req *http.Request) {
	commentid := req.URL.Query().Get("commentid")

	var cnt int64
	err := route.DB.Model(&model.Comment{}).Where("id = ?", commentid).Count(&cnt).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if cnt != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String("https://nohate.femantiwar.org/"),
		CancelURL:  stripe.String("https://nohate.femantiwar.org/"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String("price_1Ro6LaFAfeQ2yJkxlmz9c2Ep"),
				Quantity: stripe.Int64(1),
			},
		},
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Metadata: map[string]string{
				"commentid": commentid,
			},
		},
		Mode: stripe.String(stripe.CheckoutSessionModePayment),
	}

	result, err := session.New(params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating checkout session: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, result.URL, http.StatusTemporaryRedirect)
}
