
inv.task('greet')
	.using('busybox').run('echo', 'Hello, World!')
