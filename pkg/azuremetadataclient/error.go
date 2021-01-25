package azuremetadataclient

import "github.com/giantswarm/microerror"

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

var unexpectedStatusCodeError = &microerror.Error{
	Kind: "unexpectedStatusCodeError",
}
