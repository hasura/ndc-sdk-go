# TLS trust requirements

Build SDK consumers, the generator, and generated connectors with Go 1.26.10+
on the 1.26 line or Go 1.27.3+ on the 1.27 line. The module and template minimum
is 1.26.10. The reference image and generated Dockerfile pin Go 1.26.10 and
select it explicitly with GOTOOLCHAIN, so updating the base-image tag alone
does not accidentally select an older compiler.

The standard library includes the eMudhra certificate distrust policy.
Update OS CA stores separately, particularly for non-Go executables in images.
Existing generated connectors need their own module/build-input update and
rebuild; installing a new generator does not rewrite old projects.

Incident and patched-version reference: https://emudhra.we-have-a-problem.com/

Release validation must show certificate-verification rejection of all 12
published TLS chains (including cross-signs), plus acceptance of an unrelated
valid chain. Timeouts and connection failures are not evidence of distrust.
Rebuild consumers and update immutable deployed image references; an SDK source
change does not modify existing binaries, images, or hibernated instances.
Custom truststores, CA overrides, or custom verification require separate review.
