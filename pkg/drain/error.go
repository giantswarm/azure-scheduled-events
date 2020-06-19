package drain

import (
	"strings"

	"github.com/giantswarm/microerror"
)

var cannotEvictPodError = &microerror.Error{
	Kind: "cannotEvictPodError",
}

func IsCannotEvictPod(err error) bool {
	c := microerror.Cause(err)

	if err == nil {
		return false
	}

	if strings.Contains(c.Error(), "Cannot evict pod") {
		return true
	}

	return c == cannotEvictPodError
}
