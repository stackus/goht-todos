// Code generated by GoHT - DO NOT EDIT.
// https://github.com/stackus/goht

package partials

import "context"
import "io"
import "github.com/stackus/goht"

func AddTodoForm() goht.Template {
	return goht.TemplateFunc(func(ctx context.Context, __w io.Writer) (__err error) {
		__buf, __isBuf := __w.(goht.Buffer)
		if !__isBuf {
			__buf = goht.GetBuffer()
			defer goht.ReleaseBuffer(__buf)
		}
		var __children goht.Template
		ctx, __children = goht.PopChildren(ctx)
		_ = __children
		if _, __err = __buf.WriteString("<form class=\"inline\" method=\"POST\" action=\"/todos\" hx-post=\"/todos\" hx-target=\"#no-todos\" hx-swap=\"beforebegin\">\n<label class=\"flex items-center\">\n<span class=\"text-lg font-bold\">Add Todo</span>\n<input class=\"ml-2 grow\" type=\"text\" name=\"description\" _=\"on keyup if the event&#39;s key is &#39;Enter&#39; set my value to &#39;&#39; trigger keyup\"></label>\n</form>\n"); __err != nil {
			return
		}
		if !__isBuf {
			_, __err = __w.Write(__buf.Bytes())
		}
		return
	})
}
