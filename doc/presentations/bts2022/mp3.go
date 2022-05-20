func decode(d *decode.D, in any) any {
	d.FieldArray("headers", func(d *decode.D) {
		for !d.End() {
			d.TryFieldFormat("header", headerGroup)
		}
	})

	d.FieldArray("frames", func(d *decode.D) {
		for !d.End() {
			d.TryFieldFormat("frame", mp3Group)
		}
	})

	d.FieldArray("footers", func(d *decode.D) {
		for !d.End() {
			d.TryFieldFormat("footer", footerGroup)
		}
	})

	return nil
}
