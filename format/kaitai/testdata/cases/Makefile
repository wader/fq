
%.fqtest: %.ksy %.input
	echo "$$ fq -d kaitai -o source=@$(@:.fqtest=.ksy) 'dv,tovalue' $(@:.fqtest=.input)" > $@
	docker run \
		-v "${PWD}:/share" \
		-it --entrypoint=ksdump kaitai/ksv \
		-f json \
		$(@:.fqtest=.input) \
		$(@:.fqtest=.ksy) \
		| tail +5 | jq >> $@
