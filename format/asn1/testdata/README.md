tc* files from http://www.strozhevsky.com/free_docs/TEST_SUITE.zip
Files were created using:
for i in tc*.ber; do echo "\$ fq -d asn1_ber v $i" > $i.fqtest ; done
rename 's/transformed_//' transformed_tc*

laymans_guide_examples.json extracted from https://luca.ntop.org/Teaching/Appunti/asn1.html

From https://lapo.it/asn1js/ released under ISC license:
sig-p256-der.p7m
sig-p256-ber.p7m
sig-rsa1024-sha1.p7s
letsencrypt-x3.cer
ed25519.cer
