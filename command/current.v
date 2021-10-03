module command

import os

pub fn command_current() ? {
	deno_exe_path := os.find_abs_path_of_executable('deno') ?

	// output: deno 1.13.0
	result := os.execute('$deno_exe_path -V')

	lines := result.output.trim_space().split_into_lines()

	version := 'v' + lines.first().split(' ')[1]

	println(version)
}
