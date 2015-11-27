
inv.task('compile')
	.using('gcc').run('gcc', '-o', 'dist/add', 'add.c', '-static')

inv.task('package')
	.wrap('dist').inImage('busybox').at('/usr/local/bin').as('test/showaddition')

inv.task('run')
	.using('test/showaddition').run('/usr/local/bin/add')
