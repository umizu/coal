package coal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	address  = "127.0.0.1:25575"
	password = "secretpassword"
)

func TestConnect(t *testing.T) {
	_, err := Connect(address, password)
	assert.Nil(t, err)
}

func TestConnectFail(t *testing.T) {
	_, err := Connect("127.0.0.1:1", password)
	assert.NotNil(t, err)
}

func TestConnectFailAuth(t *testing.T) {
	_, err := Connect(address, "badpass")
	assert.ErrorIs(t, err, ErrUnauthenticated)
}


