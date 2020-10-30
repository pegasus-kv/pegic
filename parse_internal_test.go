package pegic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseScan(t *testing.T) {
	assert := assert.New(t)

	c := &scanCommand{}
	err := c.parse(`count hashkey "key" sortkey suffix "suffix" novalue`)
	assert.Nil(err)
	assert.Equal(&scanCommand{
		count: true,
		delete: false,
		hashKey: "key",
		sortKey: &scanSortKey{
			suffix: "suffix",
		},
		noValue: true,
	}, c)

	err = c.parse(`count delete hashkey "k"`)
	assert.NotNil(err)

	err = c.parse(`hashkey "k" novalue xxxx`)
	assert.NotNil(err)

	err = c.parse(`hashkey "with space" novalue`)
	assert.Nil(err)
	assert.Equal(&scanCommand{
		hashKey: "with space",
		noValue: true,
	}, c)

	err = c.parse(`delete hashkey "key " sortkey between "start " and "stop" novalue`)
	assert.Nil(err)
	assert.Equal(&scanCommand{
		delete: true,
		hashKey: "key ",
		sortKey: &scanSortKey{
			start: "start ",
			stop: "stop",
		},
		noValue: true,
	}, c)

	err = c.parse(`hashkey "k  " sortkey between and "stop" novalue`)
	assert.NotNil(err)
}

func TestParseFullScan(t *testing.T) {
	assert := assert.New(t)

	c := &fullScanCommand{}
	err := c.parse(`delete hashkey contains "contains" novalue`)
	assert.Nil(err)
	assert.Equal(&fullScanCommand{
		delete: true,
		hashKey: &fullScanHashKey{
			contains: "contains",
		},
		noValue: true,
	}, c)

	err = c.parse(`count hashkey suffix "suffix" sortkey is "is"`)
	assert.Nil(err)
	assert.Equal(&fullScanCommand{
		count: true,
		hashKey: &fullScanHashKey{
			suffix: "suffix",
		},
		sortKey: &fullScanSortKey{
			is: "is",
		},
	}, c)

	err = c.parse(`count hashkey suffix "suffix" sortkey is "is" contains "c"`)
	assert.NotNil(err)

	err = c.parse(`sortkey between "1" and "2" novalue`)
	assert.Nil(err)
	assert.Equal(&fullScanCommand{
		sortKey: &fullScanSortKey{
			start: "1",
			stop: "2",
		},
		noValue: true,
	}, c)

	err = c.parse(``)
	assert.Nil(err)
	assert.Equal(&fullScanCommand{}, c)

	err = c.parse(`count delete hashkey suffix "suffix"`)
	assert.NotNil(err)
}

func TestParseEncoding(t *testing.T) {
	assert := assert.New(t)

	c := &encodingCommand{}
	err := c.parse(``)
	assert.Nil(err)
	assert.Equal(&encodingCommand{}, c)

	err = c.parse(`hashkey utf-8`)
	assert.Nil(err)
	assert.Equal(&encodingCommand{
		keyType: "hashkey",
		encoding: "utf-8",
	}, c)

	err = c.parse(`reset`)
	assert.Nil(err)
	assert.Equal(&encodingCommand{reset: true}, c)

	err = c.parse(`hashkey`)
	assert.NotNil(err)
}
