
inv.task('build')
	.using('golang:latest')
		.withConfig({
			env = {"CGO_ENABLED=0"}
		})
		.run('go', 'build', 'hello.go')

inv.task('package')
	.wrap('hello')
		.at('/hello')
		.withConfig({
			cmd = {"/hello"},
			expose = {"8080/tcp"}
		})
		.as('involucro/demo')
