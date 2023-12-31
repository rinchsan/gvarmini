package gvarmini_test

import (
	"testing"

	"github.com/gostaticanalysis/testutil"
	"github.com/rinchsan/gvarmini"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.RunWithSuggestedFixes(t, testdata, gvarmini.Analyzer, "a")
}
