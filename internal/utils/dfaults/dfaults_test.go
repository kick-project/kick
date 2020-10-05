package dfaults

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func TestString(t *testing.T) {
	tz := faker.Address().TimeZone()
	assert.Equal(t, String("UTC", ""), "UTC")
	assert.Equal(t, String("UTC", tz), tz)
}

func TestInterface(t *testing.T) {
	type Contact struct {
		name   string
		number string
	}

	companyName := faker.Company().Name()
	companyNumber := faker.PhoneNumber().PhoneNumber()
	dflt := &Contact{
		name:   companyName,
		number: companyNumber,
	}

	personName := faker.Name().Name()
	personNumber := faker.PhoneNumber().PhoneNumber()
	val := &Contact{
		name:   personName,
		number: personNumber,
	}

	assert.Equal(t, Interface(dflt, nil), dflt)
	assert.Equal(t, Interface(dflt, val), val)
}
