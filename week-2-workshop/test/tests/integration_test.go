package tests

import (
	"github.com/stretchr/testify/suite"
	//test_suite "gitlab.ozon.dev/14/week-2-workshop/test/suite"
	tc_suite "gitlab.ozon.dev/14/week-2-workshop/test/suite"
	"testing"
)

func TestIntegrationSuite(t *testing.T) {
	//suite.Run(t, new(test_suite.ItemS))
	suite.Run(t, new(tc_suite.TCSuite))
}
