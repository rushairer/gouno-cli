# Release Checklist

- [ ] Go 1.26.x and 1.27.x quality gates pass with latest-patch selection (`check-latest: true`).
- [ ] Default-template integration follows the live default branch without SSH credentials.
- [x] The official `gouno-template v1.4.0` annotated tag is public and resolves to `9f0cba9c8a6be6331f48596c83663efbec5490a2`.
- [ ] The integration workflow pins `--template-ref v1.4.0`; both Go 1.26.x and 1.27.x jobs must pass automatic tidy, verbatim resource preservation, individual/composition generation, race tests, and vet.
- [ ] The CLI's real default-template integration workflow passes on Go 1.26.x and 1.27.x, including automatic tidy, Codegen v1 preservation, individual and composition generators, generated-project race tests, and vet.
- [ ] `CHANGELOG.md` has the dated v1.3.0 release entry; target tag and release do not exist.
- [ ] Confirm the release commit's module can be resolved through the public Go module proxy after tagging.
- [ ] Publish the GitHub Release notes against the immutable v1.3.0 tag and verify the release page points to the same reviewed commit.
- [ ] The worktree is clean and `git diff --check` passes.
