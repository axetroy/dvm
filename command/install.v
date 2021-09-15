module command

import net.http
import os
import szip

fn get_download_url(version string) string {
	platform := $if windows {
		'pc-windows-msvc'
	} $else $if darwin {
		'apple-darwin'
	} $else $if macos {
		'apple-darwin'
	} $else {
		'unknown-linux-gnu'
	}

	arch := $if arm64 { 'aarch64' } $else { 'x86_64' }

	filename := 'deno-$arch-$platform' + '.zip'

	return 'https://github.com/denoland/deno/releases/download/$version/$filename'
}

fn on_download(progress int) {
	println('$progress')
}

fn on_finish() {
	println('finish')
}

pub fn command_install(version string) ? {
	download_url := get_download_url(version)

	println(download_url)

	target := '/Users/axetroy/go/src/github.com/axetroy/dvm/' + os.file_name(download_url)

	println(target)

	eprintln("Downloading...")

	// http.download_file_with_progress(download_url, target, on_download, on_finish)
	http.download_file(download_url, target) or {
		return err
	}

	eprintln("Download done!")

	szip.extract_zip_to_dir(target, './dist') or {
		return err
	}
}
