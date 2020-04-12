package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func Test_parser_realPage(t *testing.T) {
	var (
		doc   *html.Node
		value string
	)
	file, err := os.Open("./sample.html")
	assert.NoError(t, err)
	doc, err = html.Parse(file)
	assert.NoError(t, err)
	value, err = crawl(doc)
	assert.NoError(t, err)
	assert.Equal(t, "€518.910", value)
}

func Test_parser_notFound(t *testing.T) {
	var doc *html.Node
	file, err := os.Open("./example.html")
	assert.NoError(t, err)
	doc, err = html.Parse(file)
	assert.NoError(t, err)
	_, err = crawl(doc)
	assert.Error(t, err)
}

func Test_parseAmount(t *testing.T) {
	v, err := parseAmount("€518.910")
	assert.NoError(t, err)
	assert.EqualValues(t, 518.910, v)
}

func Test_getLatestValue(t *testing.T) {
	v, err := getLatestValue()
	assert.NoError(t, err)
	assert.True(t, v > 0)
}
