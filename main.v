module main

import os
import flag
import command { command_install, command_ls }

fn main() {
	version := '2.0.0'

	mut fp := flag.new_flag_parser(os.args)
	fp.application('dvm')
	fp.version(version)
	fp.description('This tool is only designed to show how the flag lib is working')
	fp.skip_executable()
	fp.limit_free_args_to_at_least(1)

	is_help := fp.bool('help', 0, false, 'print help information')
	is_version := fp.bool('version', 0, false, 'print version information')

	additional_args := fp.finalize() or {
		eprintln(err)
		println(fp.usage())
		return
	}

	if is_help {
		println(fp.usage())
		return
	}

	if is_version {
		println(version)
		return
	}

	cmd := additional_args.first()

	match cmd {
		'install' {
			if additional_args.len == 1 || additional_args.len > 2 {
				println(fp.usage())
				exit(1)
			}

			command_install(additional_args[1]) or { panic(err) }
		}
		'uninstall' {}
		'ls', 'list' {}
		'ls-remote', 'list-remote' {
			if additional_args.len > 1 {
				println(fp.usage())
				exit(1)
			}
			command_ls() or { panic(err) }
		}
		else {
			println(fp.usage())
		}
	}
}
