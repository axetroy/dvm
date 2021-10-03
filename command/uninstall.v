module command

import os

fn map_versions_helper(version string) string {
	if !version.starts_with('v') {
		return 'v' + version
	} else {
		return version
	}
}

pub fn command_uninstall(v []string) ? {
	versions := v.map(map_versions_helper)

	root_dir := os.join_path(os.home_dir(), '.dvm')
	current_deno_file_path := os.join_path(os.home_dir(), '.deno', 'bin', 'deno')

	for version in versions {
		version_dir := os.join_path(root_dir, 'releases', version)

		if os.is_dir(version_dir) {
			os.rmdir_all(version_dir) ?

			println('$version has been uninstalled')

			if !os.exists(os.real_path(current_deno_file_path)) {
				os.rm(current_deno_file_path) ?
			}
		} else {
			panic(error('can not found version $version'))
		}
	}
}
