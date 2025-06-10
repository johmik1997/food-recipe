package uploadimage

import (
	"fmt"
	"io"

	"github.com/google/uuid"
)

type GraphQLUUID string
func (u GraphQLUUID) MarshalGQL(w io.Writer) {
	fmt.Fprintf(w, "%q", string(u))
}


func (u *GraphQLUUID) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("GraphQLUUID must be a string")
	}
	_, err := uuid.Parse(s)
	if err != nil {
		return err
	}
	*u = GraphQLUUID(s)
	return nil
}
