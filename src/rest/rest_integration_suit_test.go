// +build integration

package rest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestRestAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rest API Suite")
}
