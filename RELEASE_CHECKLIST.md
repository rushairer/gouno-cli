# Release Checklist

- [ ] Go 1.25.x and 1.26.x quality gates pass.
- [ ] Default template integration follows the template's default branch without SSH credentials; pinned `--template-ref` values use published tags only.
- [ ] `CHANGELOG.md` has a dated release entry; target tag and release do not exist.
- [ ] The worktree is clean and `git diff --check` passes.
