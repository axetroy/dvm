module command

import json
import regex
import net.http

struct Tag {
pub:
	ref     string
	node_id string
	url     string
}

pub fn command_ls_remote() ? {
	resp := http.get('https://api.github.com/repos/denoland/deno/git/refs/tags') ?

	tags := json.decode([]Tag, resp.text) ?

	query := r'^v\d+\.\d+\.\d+(-.*)?$'

	mut re := regex.regex_opt(query) or { panic(err) }

	for tag in tags {
		name := tag.ref.replace('refs/tags/', '')

		start, end := re.match_string(name)

		if start < 0 || end < 0 {
			continue
		}

		println(name)
	}
}
