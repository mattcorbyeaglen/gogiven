package gogiven

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type HtmlGenerator interface {
	Generate(context *TestContext) (html string)
}

type HtmlFileGenerator struct {
	HtmlGenerator *HtmlFileGenerator
}

func (generator *HtmlFileGenerator) Generate(context *TestContext) string {
	html := "<html><title>Given Test</title>" +
		"<body>"
	safeMap := context.someTests
	for _, key := range safeMap.Keys() {
		if some, ok := safeMap.Load(key); ok{
			tests := some.(*Some)
			html += tests.globalTestingT.Name()
		}
	}
	return html
}

var Generator HtmlGenerator = new(HtmlFileGenerator)

func GenerateTestOutput() {
	for _, key := range globalTestContextMap.Keys() {
		value, _ := globalTestContextMap.Load(key)
		currentTestContext := value.(*TestContext)

		html := Generator.Generate(currentTestContext)
		htmlFileName := fmt.Sprintf("%s%c%s", os.TempDir(),
			os.PathSeparator,
			strings.Replace(filepath.Base(currentTestContext.fileName), ".go", ".html", 1))

		errWritingHtml := ioutil.WriteFile(htmlFileName, []byte(html), 0644)
		if errWritingHtml != nil {
			panic("error writing html:" + errWritingHtml.Error())
		}
		fmt.Printf("generated test output: file:///%s\n", strings.Replace(htmlFileName, "\\", "/", -1))
	}
}