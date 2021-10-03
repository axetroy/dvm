module command

import os

pub fn command_use(v string) ? {
	mut version := v
	if !v.starts_with('v') {
		version = 'v' + v
	}

	root_dir := os.join_path(os.home_dir(), '.dvm')
	exe_file_path := os.join_path(root_dir, 'releases', version, 'deno')
	target_file_path := os.join_path(os.home_dir(), '.deno', 'bin', 'deno')

	if os.is_executable(exe_file_path) {
		if os.exists(target_file_path) {
			os.rm(target_file_path) ?
		}
		os.link(exe_file_path, target_file_path) ?
	} else {
		panic(error('can not found version $version'))
	}
}
