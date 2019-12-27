//+build unit

package service_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)


// ginkgo wrapper
type GinkgoTestReporter struct {}

func (g GinkgoTestReporter) Errorf(format string, args ...interface{}) {
	Fail(fmt.Sprintf(format, args...))
}

func (g GinkgoTestReporter) Fatalf(format string, args ...interface{}) {
	Fail(fmt.Sprintf(format, args...))
}
//ended


func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cart Service Suite")
}
