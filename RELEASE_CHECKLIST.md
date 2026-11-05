# Release Checklist

- [x] Go 1.26.x and 1.27.x quality gates pass with latest-patch selection (`check-latest: true`).
- [x] Default-template integration follows the live default branch without SSH credentials.
- [x] The official `gouno-template v1.4.0` annotated tag is public and resolves to `9f0cba9c8a6be6331f48596c83663efbec5490a2`.
- [x] The integration workflow pins `--template-ref v1.4.0`; both Go 1.26.x and 1.27.x jobs must pass automatic tidy, verbatim resource preservation, individual/composition generation, race tests, and vet.
- [x] The CLI's real default-template integration workflow passes on Go 1.26.x and 1.27.x, including automatic tidy, Codegen v1 preservation, individual and composition generators, generated-project race tests, and vet.
- [x] `CHANGELOG.md` has the dated v1.3.0 release entry; target tag and release do not exist.
- [x] Confirm the release commit's module can be resolved through the public Go module proxy after tagging.
- [x] Publish the GitHub Release notes against the immutable v1.3.0 tag and verify the release page points to the same reviewed commit.
- [x] The worktree is clean and `git diff --check` passes.


## v1.3.0 evidence (2026-11-05)

- [Final release-commit CI](https://github.com/rushairer/gouno-cli/actions/runs/42010113946) passed quality and live/pinned Template v1.4.0 integration on Go 1.26.10 / 1.27.3, including automatic tidy, Codegen resource preservation, individual/composition generation, race and vet.
- Annotated tag object `b898d5e921fb9a001e03a65b1b880293794fc606` peels to `087c00a523a6ec8c512d5eb9c65f657f6be8d5fb`.
- [GitHub Release](https://github.com/rushairer/gouno-cli/releases/tag/v1.3.0) exists; public proxy info resolves that exact SHA. Downloaded module sum: `h1:dL7sChra+i9zgpNxmdHFc/tYlBGhDHqfhHC0YWsyLD0=`; go.mod sum: `h1:UVjoYQk7M5sC7ziidDTT3zrG/KqzYBL8pSdM1kH1Pls=`.

This evidence update follows publication; it does not move the release tag.
