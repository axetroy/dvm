module command

import os

pub fn command_ls() ? {
	root_dir := os.join_path(os.home_dir(), '.dvm')
	releases_dir := os.join_path(root_dir, 'releases')

	mut files := os.ls(releases_dir) ?

	files.sort_ignore_case()

	for file in files {
		deno_exe_file_path := os.join_path(releases_dir, file, 'deno')

		if os.is_executable(deno_exe_file_path) {
			println(file)
		}
	}
}
