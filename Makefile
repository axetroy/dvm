default:
	make macos
	make linux
	make windows

format:
	v fmt -w *.v **/*.v

windows:
	v -prod -os windows -m64 -o bin/dvm_windows_amd64 main.v

macos:
	v -prod -os macos -m64 -o bin/dvm_darwin_amd64 main.v

linux:
	v -prod -os linux -m64 -o bin/dvm_linux_amd64 main.v