# Release Checklist

- [ ] Go 1.26.x and 1.27.x quality gates pass.
- [ ] Default template integration follows the template's default branch without SSH credentials; pinned `--template-ref` values use published tags only.
- [ ] The CLI's real default-template integration workflow passes on Go 1.26.x and 1.27.x, including generated-project Codegen, race tests, and vet.
- [ ] `CHANGELOG.md` has a dated release entry; target tag and release do not exist.
- [ ] The worktree is clean and `git diff --check` passes.
