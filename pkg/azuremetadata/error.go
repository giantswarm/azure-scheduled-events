package azuremetadata

import "github.com/giantswarm/microerror"

var unexpectedStatusCodeError = &microerror.Error{
	Kind: "unexpectedStatusCodeError",
}
