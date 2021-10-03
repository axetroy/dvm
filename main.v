module main

import os
import flag
import command { command_current, command_install, command_ls, command_ls_remote, command_uninstall, command_use }

fn main() {
	version := '2.0.0'

	mut fp := flag.new_flag_parser(os.args)
	fp.application('dvm')
	fp.version(version)
	fp.description('This tool is only designed to show how the flag lib is working')
	fp.skip_executable()
	fp.limit_free_args_to_at_least(1) or { panic(err) }

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
		'current' {
			if additional_args.len > 1 {
				println(fp.usage())
				exit(1)
			}
			command_current() or { panic(err) }
		}
		'install' {
			if additional_args.len == 1 || additional_args.len > 2 {
				println(fp.usage())
				exit(1)
			}

			command_install(additional_args[1]) or { panic(err) }
		}
		'use' {
			if additional_args.len != 2 {
				println(fp.usage())
				exit(1)
			}
			command_use(additional_args[1]) or { panic(err) }
		}
		'uninstall' {
			if additional_args.len < 2 {
				println(fp.usage())
				exit(1)
			}
			command_uninstall(additional_args[1..]) or { panic(err) }
		}
		'ls', 'list' {
			if additional_args.len > 1 {
				println(fp.usage())
				exit(1)
			}
			command_ls() or { panic(err) }
		}
		'ls-remote', 'list-remote' {
			if additional_args.len > 1 {
				println(fp.usage())
				exit(1)
			}
			command_ls_remote() or { panic(err) }
		}
		else {
			println(fp.usage())
		}
	}
}
