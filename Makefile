default:
	make macos
	make linux
	make windows

format:
	v fmt -w *.v **/*.v

windows:
	v -prod -os windows -m64 -o ./bin/cross-env_windows_amd64 main.v
	# v -prod -os windows -m32 -o ./bin/cross-env_windows_386 main.v

macos:
	v -prod -os macos -m64 -o ./bin/cross-env_darwin_amd64 main.v

linux:
	v -prod -os linux -m64 -o ./bin/cross-env_linux_amd64 main.v
	# v -prod -os linux -m32 -o ./bin/cross-env_linux_386 main.v