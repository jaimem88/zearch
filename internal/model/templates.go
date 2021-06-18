package model

import "text/template"

var OrgResultTemplate = template.Must(template.New("orgResultTemplate").Parse(orgResultTemplate))

const orgResultTemplate = `
_id                 {{ index .Organization "_id" }}
url                 {{ index .Organization "url" }}
external_id         {{ index .Organization "external_id" }}
name                {{ index .Organization "name" }}
domain_names        {{ index .Organization "domain_names" }}
created_at          {{ index .Organization "created_at" }}
details             {{ index .Organization "details" }}
shared_tickets      {{ index .Organization "shared_tickets" }}
tags                {{ index .Organization "tags" }}
{{ range $index, $element := .UserNames }}user_{{ $index }}              {{ $element }}
{{ end }}{{ range $index, $element := .TicketSubjects }}ticket_{{ $index }}            {{ $element }}
{{ end }}`
