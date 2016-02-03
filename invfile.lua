
local repo = 'involucro/tool'

inv.task('wrap-yourself')
	.using('busybox:latest')
		.run('mkdir', '-p', 'dist/tmp/')
		.run('cp', 'involucro', 'dist/')
	.wrap('dist').at('/')
		.withConfig({entrypoint = {'/involucro'}})
		.as(repo .. ':latest')
	.using('busybox:latest')
		.run('rm', '-rf', 'dist')


if ENV.TRAVIS_PULL_REQUEST == "false" then
	local tag = ENV.TRAVIS_TAG

	if tag == "" then
		tag = ENV.TRAVIS_BRANCH
	end

	if tag == "master" then
		tag = "latest"
	end

	inv.task('upload-to-hub')
		.tag(repo .. ':latest')
			.as(repo .. ':' .. tag)
		.push(repo .. ':' .. tag)
end
