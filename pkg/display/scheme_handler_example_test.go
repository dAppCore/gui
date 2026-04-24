package display

import (
	"fmt"
	"net/url"

	core "dappco.re/go/core"
)

func ExampleNewCoreSchemeHandler() {
	c := core.New()
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		name, ok := q.(string)
		if !ok || name != "core.settings" {
			return core.Result{}
		}
		return core.Result{Value: "settings-query", OK: true}
	})

	parsedURL, _ := url.Parse("core://settings")
	result := NewCoreSchemeHandler(c).Handle(parsedURL)

	fmt.Println(result.OK, result.Value)
	// Output: true settings-query
}
