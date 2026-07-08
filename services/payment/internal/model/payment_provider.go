package model

type PaymentProvider string

const (
	ProviderSimulator PaymentProvider = "SIMULATOR"

	ProviderRazorpay PaymentProvider = "RAZORPAY"

	ProviderStripe PaymentProvider = "STRIPE"
)

func (p PaymentProvider) IsSupported() bool {

	switch p {

	case ProviderSimulator,
		ProviderStripe,
		ProviderRazorpay:

		return true

	default:

		return false

	}

}