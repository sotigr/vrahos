package vrahos

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"text/template"
)

func minifyHtml(s string) string {
	if os.Getenv("MINIFY") == "true" {
		var re = regexp.MustCompile(`(?m)<!--(.|\n)*?-->`)
		s = re.ReplaceAllString(s, "")
		return strings.ReplaceAll(strings.ReplaceAll(string(strings.Join(strings.Fields(s), " ")), "\t", ""), "\n", "")
	}

	return s
}

func Vrahos(mux *http.ServeMux, components []Component, sse *Sse, meta *MetaData, handler func(next http.Handler) http.Handler) {

	fs := http.FileServer(http.Dir("./static"))

	mux.Handle("/static/", http.StripPrefix("/static", fs))

	templates := template.Must(template.New(""), nil)

	templateFuncs := map[string]any{
		"makeprops": MakeProps,
	}

	for _, component := range components {
		c := component
		name := c.Name()
		funcs := c.Functions()
		if funcs != nil {
			for k, v := range *funcs {
				templateFuncs[fmt.Sprintf("%s%s", name, k)] = v
			}
			templates.Funcs(templateFuncs)
		}

	}

	templates.Funcs(templateFuncs)

	validRoutes := []string{}
	meta.Template = templates

	for _, component := range components {
		c := component
		name := c.Name()
		tmp := templates.New(name)

		tmpStr := c.Template()
		var err error
		if strings.HasPrefix(tmpStr, "@file:") {
			filePath := strings.Replace(tmpStr, "@file:", "", 1)
			f, err := os.ReadFile(filePath)
			if err != nil {
				panic(err)
			}
			tmp, err = tmp.Parse(minifyHtml(string(f)))

		} else {
			tmp, err = tmp.Parse(minifyHtml(tmpStr))
		}

		if err != nil {
			panic(err)
		}

		funcs := c.Functions()
		if funcs != nil {
			templateFuncs := map[string]any{}
			for k, v := range *funcs {
				templateFuncs[fmt.Sprintf("%s%s", name, k)] = v
			}
			templates.Funcs(templateFuncs)
		}

		url := c.URL()
		if url != "" {
			if url[len(url)-1:] != "/" {
				panic("All urls must end with \"/\" Error in " + url + " " + name)
			}
			validRoutes = append(validRoutes, url)
			mux.Handle(url, handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !slices.Contains(validRoutes, r.URL.Path) {
					http.Error(w, "Not Found", http.StatusNotFound)
					return
				}
				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusOK)
					return
				}

				w.Header().Set("Content-Type", "text/html")
				w.Header().Set("CharSet", "utf-8")

				if err != nil {
					io.WriteString(w, err.Error())
					return
				}

				props, headers := c.Props(r, meta)

				if headers != nil {
					for key, value := range headers {
						if strings.ToLower(key) == "redirect" {
							// http.Redirect(w, r, value, http.StatusSeeOther)
						} else if strings.ToLower(key) == "status" {
							h, err := strconv.Atoi(value)
							if err != nil {
								fmt.Println(err)
								h = 500
							}
							w.WriteHeader(h)
						} else {
							w.Header().Set(key, value)
						}
					}
				}

				data := TemplateData{
					Props:   props,
					Page:    name,
					Request: r,
					Meta:    meta,
				}

				tmp.Execute(w, &data)
			})))
		}
	}

	if sse != nil {
		fmt.Println("Configuring sse server at", sse.GetPath())
		sse.Init(mux)
	}

}
