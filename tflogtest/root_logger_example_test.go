// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tflogtest

import (
	"bytes"
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func ExampleRootLogger() {
	var output bytes.Buffer

	ctx := RootLogger(context.Background(), &output)

	// Root provider logger is now available for usage, such as writing
	// entries, calling SetField(), or calling NewSubsystem().
	tflog.Trace(ctx, "hello, world", map[string]interface{}{
		"foo":    123,
		"colors": []string{"red", "blue", "green"},
	})

	fmt.Println(output.String())

	// Output:
	// {"@level":"trace","@message":"hello, world","@module":"provider","colors":["red","blue","green"],"foo":123}
}
