package provider

type CreateIntentResult struct {
	ProviderReference string
	PaymentURL        string
}

type CreateRefundResult struct {

	ProviderReference string
}