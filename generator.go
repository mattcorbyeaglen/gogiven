package gogiven

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"github.com/corbym/gogiven/generator"
	"github.com/corbym/htmlspec"
)

// Generator is a global variable that holds the GoGivensOutputGenerator.
// You can replace the generator with your own if you match the interface here
// and set Generator = new(myFooGenerator) in a method (usually TestMain or init).
// Don't forget to add the call to the generator function in a "func TestMain(testing.M)" method
// in your test package.
// One file per test file will be generated containing output.
var Generator generator.GoGivensOutputGenerator = new(htmlspec.TestOutputGenerator)

// TransformFileNameToHeader takes a test filename e.g. /foo/bar/my_test.go and returns a header e.g. "My Test".
// Strips off the file path and removes the extension.
func TransformFileNameToHeader(fileName string) (header string) {
	return strings.Title(strings.Replace(strings.TrimSuffix(filepath.Base(fileName), ".go"), "_", " ", -1))
}

// GenerateTestOutput generates the test output. Call this method from TestMain.
func GenerateTestOutput() {
	for _, key := range globalTestContextMap.Keys() {
		value, _ := globalTestContextMap.Load(key)
		currentTestContext := value.(*TestContext)
		pageData := &generator.PageData{
			Title:   TransformFileNameToHeader(currentTestContext.FileName()),
			SomeMap: currentTestContext.SomeTests().AsMapOfSome(),
		}
		output := Generator.Generate(pageData)
		extension := Generator.FileExtension()

		outputFileName := fmt.Sprintf("%s%c%s", outputDirectory(),
			os.PathSeparator,
			strings.Replace(filepath.Base(currentTestContext.fileName), ".go", extension, 1))

		err := ioutil.WriteFile(outputFileName, []byte(output), 0644)
		if err != nil {
			panic("error generating gogiven output:" + err.Error())
		}
		fmt.Printf("\ngenerated test output: file://%s\n", strings.Replace(outputFileName, "\\", "/", -1))
	}
}

func outputDirectory() string {
	outputDir := os.Getenv("GOGIVENS_OUTPUT_DIR")
	if outputDir == "" {
		os.Stdout.WriteString("env var GOGIVENS_OUTPUT_DIR was not found, using TempDir " + os.TempDir())
		outputDir = os.TempDir()
	}
	if _, err := os.Stat(outputDir); err == nil {
		return outputDir
	}
	os.Stderr.WriteString("output dir not found:" + outputDir + ", defaulting to ./")
	return "."
}
