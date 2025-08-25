package ticket

import (
	"ticket"

	u "github.com/lemoras/goutils/api"
)

func Main(in ticket.Request) (*u.Response, error) {
	return ticket.Invoke(in)
}
