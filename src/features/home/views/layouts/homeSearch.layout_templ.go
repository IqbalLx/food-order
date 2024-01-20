// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.513
package home

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	components "github.com/IqbalLx/food-order/src/features/home/views/components"
	sharedComponents "github.com/IqbalLx/food-order/src/shared/views/components"
)

func HomeSearch(isWithInitialQuery bool, initialQuery string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"flex flex-col overflow-y-auto no-scrollbar\"><div class=\"sticky top-0 bg-white z-50\"><div class=\"px-2\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if isWithInitialQuery {
			templ_7745c5c3_Err = components.Searchbar(initialQuery).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = components.Searchbar("").Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><sl-progress-bar style=\"--height: 4px;\" class=\"htmx-indicator w-full\" indeterminate></sl-progress-bar></div><div id=\"store-search-container\" class=\"flex flex-col gap-y-2\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if isWithInitialQuery {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" hx-post=\"/stores/search\" hx-trigger=\"load\" hx-indicator=\".htmx-indicator\" hx-target=\"#store-search-container\" hx-swap=\"innerHTML\" hx-include=\"#home-searchbar\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if isWithInitialQuery {
			for i := 0; i < 2; i++ {
				templ_7745c5c3_Err = sharedComponents.StoreCardWithMenuIndicator().Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
