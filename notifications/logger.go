package notifications

import (
	"github.com/AreaHQ/logging"
)

var logger *logging.Logger

func init() {
	logger = logging.New(nil, nil, new(logging.ColouredFormatter))
}
