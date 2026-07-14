package provider

import (
	//"fmt"
	"fmt"
	"math/rand"
	"time"

	"github.com/rudransh/distributed-commerce/notification/internal/errors"
	"github.com/rudransh/distributed-commerce/notification/internal/model"
	"github.com/rudransh/distributed-commerce/pkg/retry"
)

type ConsoleEmailProvider struct {
    successRate int
    minLatency time.Duration
    maxLatency time.Duration
}

func NewConsoleEmailProvider(successRate int,minLatency time.Duration,maxLatency time.Duration) *ConsoleEmailProvider {

	return &ConsoleEmailProvider{
		successRate: successRate,
		minLatency: minLatency,
		maxLatency: maxLatency,
	}

}

func (p *ConsoleEmailProvider) simulateLatency() {

    if p.maxLatency <= p.minLatency {
        time.Sleep(p.minLatency)

        return
    }
    diff := p.maxLatency - p.minLatency
    delay := p.minLatency +
        time.Duration(rand.Int63n(int64(diff)))
    time.Sleep(delay)
}

func (p *ConsoleEmailProvider) shouldSucceed() bool {

    return rand.Intn(100) < p.successRate

}

func (p *ConsoleEmailProvider) Send(
    notification *model.Notification,
) error {

    p.simulateLatency()

    fmt.Println("====================================")
    fmt.Printf("To      : %s\n", notification.Recipient)
    fmt.Printf("Subject : %s\n", notification.Subject)
    fmt.Println("------------------------------------")
    fmt.Println(notification.Body)
    fmt.Println("------------------------------------")

    if !p.shouldSucceed() {

        return retry.NewRetryable(
            errors.ErrTemporaryFailure,
        )
    }

    fmt.Println("Email Sent Successfully")
    fmt.Println("====================================")

    return nil
}