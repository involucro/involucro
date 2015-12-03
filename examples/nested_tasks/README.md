
Nested Tasks
============

Tasks can call other tasks as well (they do not have to be declared at that point):

    inv.task('all')
      .runTask('touch')
      .runTask('clean')

    inv.task('touch')
      .using('busybox').run('touch', 'testfile')

    inv.task('clean')
      .using('busybox').run('rm', '-f', 'testfile')
