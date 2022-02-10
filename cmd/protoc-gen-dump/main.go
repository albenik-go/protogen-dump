package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/albenik/protoc-gen-dump/wellknown"
)

func main() {
	fs := new(flag.FlagSet)
	noComments := fs.String("comments", "yes", "")

	protogen.Options{ParamFunc: fs.Set}.Run(func(gen *protogen.Plugin) error {
		dumpln("=== DUMP BEGIN ===")
		defer dumpln("=== DUMP END ===")

		showComments := !(*noComments == "hide")

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			dumpln("File: %s (%s)", f.Desc.Path(), f.Desc.FullName())
			dumpln("Generated: %s", f.GeneratedFilenamePrefix)

			p1 := "  "

			for _, s := range f.Services {
				if showComments {
					dumpLeadingComments(p1, &s.Comments)
				}
				if showComments && len(s.Comments.Trailing) > 0 {
					dumpln("%sService: %s %s", p1, s.Desc.FullName(), trailingComment(s.Comments.Trailing))
				} else {
					dumpln("%sService: %s", p1, s.Desc.FullName())
				}

				p2 := p1 + "  "

				for _, m := range s.Methods {
					if showComments {
						dumpLeadingComments(p2, &m.Comments)
					}
					if showComments && len(m.Comments.Trailing) > 0 {
						dumpln("%sMethod: %s (%s) %s", p2, m.Desc.Name(), m.Desc.FullName(),
							trailingComment(m.Comments.Trailing))
					} else {
						dumpln("%sMethod: %s (%s)", p2, m.Desc.Name(), m.Desc.FullName())
					}
					p3 := p2 + "  "
					in := m.Input.Desc
					out := m.Output.Desc
					dumpln("%s Input: %s (%s)", p3, in.Name(), in.FullName())
					dumpln("%sOutput: %s (%s)", p3, out.Name(), out.FullName())
				}
			}

			for _, m := range f.Messages {
				dumpMessage(p1, m, showComments)
			}
		}

		return nil
	})

}

func dumpMessage(p string, m *protogen.Message, showComments bool) {
	if showComments {
		dumpLeadingComments(p, &m.Comments)
	}
	if showComments && len(m.Comments.Trailing) > 0 {
		dumpln("%sMessage: %s (%s) %s", p, m.Desc.Name(), m.Desc.FullName(), trailingComment(m.Comments.Trailing))
	} else {
		dumpln("%sMessage: %s (%s)", p, m.Desc.Name(), m.Desc.FullName())
	}

	p1 := p + "  "

	for _, f := range m.Fields {
		if showComments {
			dumpLeadingComments(p1, &f.Comments)
		}

		switch f.Desc.Kind() {
		case protoreflect.MessageKind:
			typeName := string(f.Message.Desc.FullName())
			if wellknown.Wellknown(f.Message.Desc.FullName()) {
				typeName += " (wellknown)"
			}

			if showComments && len(f.Comments.Trailing) > 0 {
				dumpln("%sField: %s (%s) <%s %s> %s", p1, f.Desc.Name(), f.Desc.FullName(),
					f.Desc.Cardinality(), typeName, trailingComment(f.Comments.Trailing))
			} else {
				dumpln("%sField: %s (%s) <%s %s>", p1, f.Desc.Name(), f.Desc.FullName(),
					f.Desc.Cardinality(), typeName)
			}

			p2 := p1 + "  "
			if f.Desc.IsMap() {
				k := f.Desc.MapKey()
				v := f.Desc.MapValue()
				dumpln("%sMap <%s> → %s <%s>", p2, k.Kind(), v.FullName(), v.Kind())
			}

		default:
			dumpln("%sField: %s (%s) <%s %s>", p1, f.Desc.Name(), f.Desc.FullName(),
				f.Desc.Cardinality(), f.Desc.Kind())
		}

	}

	for _, mm := range m.Messages {
		dumpMessage(p1, mm, showComments)
	}
}

func dumpDeattachedComments(prefix string, comments []protogen.Comments) {
	strs := make([]string, 0, len(comments))
	for _, c := range comments {
		strs = append(strs, "//!"+strings.TrimPrefix(strings.TrimSpace(c.String()), "//"))
	}
	dumpln("%s%s", prefix, strings.Join(strs, "\n"))
}

func dumpLeadingComments(prefix string, cs *protogen.CommentSet) {
	if len(cs.LeadingDetached) > 0 {
		dumpDeattachedComments(prefix, cs.LeadingDetached)
	}
	if len(cs.Leading) > 0 {
		dumpln("%s%s", prefix, strings.TrimSpace(cs.Leading.String()))
	}
}

func dump(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func dumpln(format string, a ...interface{}) {
	dump(format+"\n", a...)
}

func trailingComment(c protogen.Comments) string {
	s := strings.Split(strings.TrimSuffix(string(c), "\n"), "\n")
	return strings.Join(s, "|")
}
