syntax = "v1"

info (
	title: // : add title
	desc: // : add description
	author: "{{.gitUser}}"
	email: "{{.gitEmail}}"
)

type request {
	// : add members here and delete this comment
}

type response {
	// : add members here and delete this comment
}

service {{.serviceName}} {
	@handler GetUser // : set handler name and delete this comment
	get /users/id/:userId(request) returns(response)

	@handler CreateUser // : set handler name and delete this comment
	post /users/create(request)
}
