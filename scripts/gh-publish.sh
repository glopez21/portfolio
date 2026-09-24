#!/usr/bin/env bash
# Publish the 25 project repos to GitHub under glopez21/ with a MINIMAL README,
# keeping Forgejo main (full README) as the source of truth.
#
# Strategy per repo:
#   1. ensure `github` remote -> https://github.com/glopez21/<name>.git
#   2. create GitHub repo (private) if it does not exist
#   3. build a temp worktree at forgejo/main, swap README.md to the minimal
#      version from scripts/gh-readmes/<name>.md, commit, force-push to GitHub `main`
#   4. The `main` branch / full README stay unchanged on Forgejo.
set -euo pipefail

REPOS=(3v4l 4sk augur Bl4ck1c3 hermes jobtracker KRAKEN homelab logsentry m0rpheus-rs n3xus-flow n3xusDB neural-sim ph4nt0m r0ut3r r4g rust-infer s3arc SENTINEL shadowsim SLEEPER termvault threat-pulse transitflow-nyc Vestige)

PROJECTS_DIR="/home/w01f/projects"
READMES_DIR="/home/w01f/projects/portfolio/scripts/gh-readmes"
WORKBASE="/tmp/ghpub"
GITHUB_USER="glopez21"
DRY_RUN="${DRY_RUN:-0}"

mkdir -p "${WORKBASE}"

for name in "${REPOS[@]}"; do
  repo="${PROJECTS_DIR}/${name}"
  readme="${READMES_DIR}/${name}.md"
  rem="https://github.com/${GITHUB_USER}/${name}.git"

  echo "===== ${name} ====="

  if [[ ! -d "${repo}" ]]; then
    echo "  SKIP: no local repo dir"; continue
  fi
  if [[ ! -f "${readme}" ]]; then
    echo "  SKIP: no minimal readme"; continue
  fi

  cd "${repo}"
  git remote get-url github >/dev/null 2>&1 || git remote add github "${rem}"
  actual="$(git remote get-url github)"
  if [[ "${actual}" != "${rem}" ]]; then
    git remote set-url github "${rem}"
    actual="${rem}"
  fi
  echo "  github remote: ${actual}"

  # ensure Forgejo main is current
  git fetch forgejo main >/dev/null 2>&1 || true

  # repo existence + visibility
  if gh repo view "${GITHUB_USER}/${name}" >/dev/null 2>&1; then
    vis="$(gh repo view "${GITHUB_USER}/${name}" --json isPrivate -q .isPrivate)"
    echo "  GitHub repo EXISTS (private=${vis})"
    fpush="--force"
  else
    echo "  GitHub repo MISSING -> create private"
    if [[ "${DRY_RUN}" == "0" ]]; then
      gh repo create "${GITHUB_USER}/${name}" --private >/dev/null
    fi
    fpush="--force"
  fi

  wt="${WORKBASE}/${name}"
  if [[ "${DRY_RUN}" == "1" ]]; then
    echo "  [dry-run] would build worktree from forgejo/main, swap README, push ${fpush} github commit:main"
    continue
  fi

  rm -rf "${wt}"
  git worktree prune
  git worktree add --detach "${wt}" forgejo/main >/dev/null 2>&1
  cp -f "${readme}" "${wt}/README.md"
  git -C "${wt}" add README.md
  git -C "${wt}" -c user.name="glopez21" -c user.email="glopez21@users.noreply.github.com" \
      commit -q -m "docs: GitHub mirror minimal README" || git -C "${wt}" commit -q --allow-empty \
      -m "docs: GitHub mirror minimal README (no change)"
  git -C "${wt}" push "${fpush}" github HEAD:refs/heads/main 2>&1 | sed 's/^/    /'
  git worktree remove "${wt}" --force

  echo "  OK: github main set to minimal-README tree (forgejo main untouched)"
done

echo
echo "DONE"