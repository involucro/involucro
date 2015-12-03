
inv.task('all')
  .runTask('touch')
  .runTask('clean')

inv.task('touch')
  .using('busybox').run('touch', 'testfile')

inv.task('clean')
  .using('busybox').run('rm', '-f', 'testfile')
