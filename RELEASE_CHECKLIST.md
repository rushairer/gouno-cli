# Release Checklist

- [ ] Go 1.26.x and 1.27.x quality gates pass with latest-patch selection (`check-latest: true`).
- [ ] Default-template integration follows the live default branch without SSH credentials.
- [ ] The official `gouno-template v1.4.0` tag exists and resolves to the reviewed release commit.
- [ ] Update the CLI integration workflow's pinned `--template-ref` from v1.3.0 to v1.4.0; both Go 1.26.x and 1.27.x pinned-tag bootstrap jobs must pass.
- [ ] The CLI's real default-template integration workflow passes on Go 1.26.x and 1.27.x, including automatic tidy, Codegen v1 preservation, individual and composition generators, generated-project race tests, and vet.
- [ ] `CHANGELOG.md` has the dated v1.3.0 release entry; target tag and release do not exist.
- [ ] Confirm the release commit's module can be resolved through the public Go module proxy after tagging.
- [ ] The worktree is clean and `git diff --check` passes.
