package failure

import (
	"errors"
	"fmt"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/fault"
	"github.com/semi-technologies/weaviate/entities/models"
)

// CombineGraphQLErrors combines and returns a slice of
// *models.GraphQLError as a single error
func CombineGraphQLErrors(errs []*models.GraphQLError) error {
	message := "weaviate gql errors"

	for _, e := range errs {
		message = fmt.Sprintf("%s: %+v", message, e)
	}

	return errors.New(message)
}

func WeaviateError(srcErr error) (dstErr error) {
	if werr, ok := srcErr.(*fault.WeaviateClientError); ok &&
		werr != nil && werr.DerivedFromError != nil {
		dstErr = fmt.Errorf("%s: %s", srcErr, werr.DerivedFromError.Error())
		return
	}

	dstErr = srcErr
	return
}
