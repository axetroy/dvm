module command

import os

fn find_exe_path(exepath string) ?string {
	if exepath == '' {
		return error('expected non empty `exepath`')
	}
	if os.is_abs_path(exepath) {
		return os.real_path(exepath)
	}
	mut res := ''
	paths := os.getenv('PATH').split(os.path_delimiter)
	for p in paths {
		found_abs_path := os.join_path(p, exepath)
		if os.exists(found_abs_path) && os.is_executable(found_abs_path) {
			res = found_abs_path
			break
		}
	}
	if res.len > 0 {
		return res
	}
	return error('failed to find executable')
}

pub fn command_unused() ? {
	deno_exe_path := find_exe_path('deno') or {
		eprintln('Deno are not been used yet!')
		return
	}

	if os.is_link(deno_exe_path) {
		os.rm(deno_exe_path) ?
		eprintln('Deno has been remove!')
	} else {
		eprintln('Deno are not been used yet!')
	}
}
