package main

import "fmt"

type Authorizer interface {
	Authorize() string
}

type Capturer interface {
	Capture() string
}

type PaymentProvider interface {
	Authorizer
	Capturer
}

type StripeProvider struct {
	MerchantID string
}

func (s StripeProvider) Authorize() string {
	return fmt.Sprintf("Stripe: авторизация платежа для %s", s.MerchantID)
}

func (s StripeProvider) Capture() string {
	return fmt.Sprintf("Stripe: списание средств для %s", s.MerchantID)
}

func (s StripeProvider) String() string {
	return fmt.Sprintf("StripeProvider{MerchantID: %s}", s.MerchantID)
}

func main() {
	var provider PaymentProvider = StripeProvider{MerchantID: "acct_123"}

	fmt.Println(provider.Authorize())
	fmt.Println(provider.Capture())
	fmt.Println(provider)
}
