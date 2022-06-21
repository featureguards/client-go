package featureguards_test

import (
	"context"
	"fmt"

	featureguards "github.com/featureguards/featureguards-go/v2"
)

func ExampleFeatureGuards_IsOn() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ft := featureguards.New(ctx, featureguards.WithApiKey("API_KEY"),
		featureguards.WithDefaults(map[string]bool{"TEST": true}))
	on, _ := ft.IsOn("TEST")
	fmt.Printf("%v\n", on)
	// Output: true
}

func ExampleFeatureGuards_IsOn_attributes() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ft := featureguards.New(ctx, featureguards.WithApiKey("API_KEY"),
		featureguards.WithDefaults(map[string]bool{"TEST": true}))
	on, _ := ft.IsOn("FOO", featureguards.WithAttributes(
		featureguards.Attributes{}.Int64("user_id", 123).String("company_slug", "acme")))
	fmt.Printf("%v\n", on)
	// Output: false

}
