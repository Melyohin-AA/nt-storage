package manager

import (
	"fmt"
	"net/http"
	"nt-storage/storage"
	"strconv"
	"strings"
)

type IntRange struct {
	set   bool
	from  int
	count int
}

func (ir *IntRange) Read(r *http.Request) {
	var err error
	if ir.from, err = strconv.Atoi(r.FormValue("from")); (err != nil) || (ir.from < 0) {
		return
	}
	if ir.count, err = strconv.Atoi(r.FormValue("count")); (err != nil) || (ir.count < 0) {
		return
	}
	ir.set = true
}

type Filter struct {
	fid, name, ts, comment, tags bool
}

func (f *Filter) Read(r *http.Request) {
	f.fid = r.FormValue("f_fid") == "true"
	f.name = r.FormValue("f_name") == "true"
	f.ts = r.FormValue("f_mtime") == "true"
	f.comment = r.FormValue("f_comment") == "true"
	f.tags = r.FormValue("f_tags") == "true"
}

func (f *Filter) apply(sb *strings.Builder, blob *storage.Blob, df *storage.DataFile) {
	if f.fid {
		sb.WriteString(fmt.Sprintf(" %s", blob.Fid))
	}
	if f.name {
		sb.WriteString(fmt.Sprintf(" '%s'", blob.Name))
	}
	if f.ts {
		sb.WriteString(fmt.Sprintf(" %d", blob.ModTime))
	}
	if f.comment {
		sb.WriteString(fmt.Sprintf(" {%s}", blob.Comment))
	}
	if f.tags {
		tagVals := make([]string, len(blob.TagValues))
		for i, tv := range blob.TagValues {
			if tag, tagExists := df.GetTag(tv.TagId); tagExists {
				tagVals[i] = fmt.Sprintf(" %s:%v", tag.Label, tv.Value)
			}
		}
		sb.WriteString(" [")
		sb.WriteString(strings.Join(tagVals, " "))
		sb.WriteRune(']')
	}
}
