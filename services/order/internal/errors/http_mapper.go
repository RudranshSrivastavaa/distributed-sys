package errors

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	httpresponse "github.com/rudransh/distributed-commerce/pkg/http/response"
)

func Handle(
	c *fiber.Ctx,
	err error,
) error {

	switch {

	case errors.Is(err, ErrOrderNotFound):
		return httpresponse.NotFound(
			c,
			err.Error(),
		)

	case errors.Is(err, ErrCustomerIDRequired):

		return httpresponse.BadRequest(
			c,
			err.Error(),
		)

	case errors.Is(err, ErrEmptyOrder):

		return httpresponse.BadRequest(
			c,
			err.Error(),
		)

	case errors.Is(err, ErrInvalidQuantity):

		return httpresponse.BadRequest(
			c,
			err.Error(),
		)

	case errors.Is(err, ErrInvalidPrice):

		return httpresponse.BadRequest(
			c,
			err.Error(),
		)

	default:

		return httpresponse.InternalServerError(c)

	}

}