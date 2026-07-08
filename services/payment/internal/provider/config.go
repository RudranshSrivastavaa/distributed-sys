package provider

import "time"

type SimulatorConfig struct {
	Delay time.Duration

	SuccessRate int

	WebhookSecret string
}