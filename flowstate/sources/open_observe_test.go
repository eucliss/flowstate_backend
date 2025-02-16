package sources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenObserve(t *testing.T) {
	openObserve := &OpenObserveInitializer{}
	openObserveSource, err := openObserve.Initialize()
	OpenObserveStruct := openObserveSource.(*OpenObserve)
	assert.Nil(t, err)
	assert.NotNil(t, openObserveSource)
	assert.Equal(t, "http://localhost:5080", OpenObserveStruct.URL)
}

func TestOpenObserveQuery(t *testing.T) {
	openObserve := &OpenObserveInitializer{}
	openObserveSource, err := openObserve.Initialize()
	assert.Nil(t, err)
	query := Query{
		Query:      "SELECT * FROM \"test_torq3\"",
		SourceType: "flowstate",
		Start:      1738040675782000,
		End:        1738041575782000,
		Limit:      2,
	}
	results := openObserveSource.Query(query)
	assert.NotNil(t, results)
}
