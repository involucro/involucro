
inv.task('compile')
	.using('sojournlabs/gcc').run('gcc', '-o', 'dist/add', 'add.c', '-static')

inv.task('package')
	.wrap('dist').inImage('busybox').at('/usr/local/bin').as('test/showaddition:v1')

inv.task('run')
  .using('test/showaddition:v1')
  .withExpectation({code = 0, stdout = "5 \\+ 10 = 15\n"})
	.run('/usr/local/bin/add')
