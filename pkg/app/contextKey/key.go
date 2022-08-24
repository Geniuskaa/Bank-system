package contextKey

var AuthenticationContextKey = &contextKey{"authentication context"}

type contextKey struct {
	name string
}

func (c *contextKey) String() string {
	return c.name
}
