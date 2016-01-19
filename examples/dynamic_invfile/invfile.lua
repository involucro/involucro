
function defineMoreTasks()
  for i = 1, 3 do
    inv.task('do' .. tostring(i)).using('busybox').run('echo', 'Task ' .. tostring(i))
  end
end

inv.task('prep')
  .hook(defineMoreTasks)
  .runTask('do2')
