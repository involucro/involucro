
inv.task('compile')
	.using('sojournlabs/gcc').run('gcc', '-o', 'dist/add', 'add.c', '-static')

inv.task('package')
	.wrap('dist').inImage('busybox').at('/usr/local/bin').as('test/showaddition')

inv.task('run')
  .using('test/showaddition')
  .withExpectation({code = 0, stdout = "5 \\+ 10 = 15\n"})
	.run('/usr/local/bin/add')
