def fq_tmpl_adoc($replace_svg):
	# ```jq-eval
	# ...
	# ```
	gsub(
		"////jq-eval\n(?<expr>[\\s\\S]*?)\n////jq-eval";
		[_eval(.expr; {})] | join("\n")
	) |
	# image::display_decode_value_d.svg[]
	if $replace_svg then
		gsub(
			"image::(?<image>.*?)\\[\\]";
			( .image
			| gsub(".svg"; ".txt")
			| open
			| tostring
			| "[source,shell]\n----\n\(.)----\n"
			)
		)
	end |
	gsub(
		"\\$FQ_VERSION";
		_main_input.version
	);
