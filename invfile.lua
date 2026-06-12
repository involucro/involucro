-- local repo = "involucro/tool"
local repo = "quay.io/bgruening/involucro"

local tag = ENV.IMAGE_TAG or "latest"
local arch = ENV.TARGETARCH or "amd64"
local platforms = {
	amd64 = "linux/amd64",
	arm64 = "linux/arm64",
	armv7 = "linux/arm/v7",
	riscv64 = "linux/riscv64",
}
local platform = platforms[arch]

if platform == nil then
	error("unsupported TARGETARCH: " .. tostring(arch))
end

inv.task("wrap")
	.wrap("docker-context/" .. arch)
	.at("/")
	.withConfig({ entrypoint = { "/involucro" } })
	.withPlatform(platform)
	.as(repo .. ":" .. tag .. "-" .. arch)

inv.task("push").push(repo .. ":" .. tag .. "-" .. arch)

-- Keep wrap-yourself for local dev/integration testing on the host arch
inv.task("wrap-yourself")
	.using("busybox:latest")
	.run("mkdir", "-p", "dist/tmp/")
	.run("cp", "involucro", "dist/")
	.wrap("dist")
	.at("/")
	.withConfig({ entrypoint = { "/involucro" } })
	.as(repo .. ":latest")
	.using("busybox:latest")
	.run("rm", "-rf", "dist")
