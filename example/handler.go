package main

import (
	"fmt"
	"net/http"
)

func (r *userRepo) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(`
<div>
<a href="/json-form/form.html?title=Create%20user&amp;schemaName=` + r.schemaName + `&amp;submitUrl=/users&amp;submitMethod=POST">Create user with dynamic form</a>
<br />
<a href="/create-user">Create user with static form</a>
</div>

<ul>
`))

	for i, u := range r.list() {
		_, _ = w.Write([]byte(fmt.Sprintf(`

<li> %s %s
<a href="/json-form/form.html?title=Edit%%20user&amp;schemaName=`+r.schemaName+`&amp;valueUrl=/user/%d.json&amp;submitUrl=/user/%d.json&amp;submitMethod=PUT">Edit with dynamic form</a>
<a href="/edit-user/%d">Edit with static form</a><br />
</li>
`, u.FirstName, u.LastName, i+1, i+1, i+1)))
	}

	_, _ = w.Write([]byte(`
</ul>
`))
}
