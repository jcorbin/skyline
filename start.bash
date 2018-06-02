#!/usr/bin/env bash
set -e

base=go_setup
branch="go_solution_$(date +%Y%m%d)"

head=$(git symbolic-ref HEAD)
head=${head#refs/heads/}

prior=
if [[ $head = go_solution_* ]]; then
    prior=$head
elif [ -n "$1" ]; then
    prior=$1
fi

rev=$(git rev-parse --verify "$base")
prior_desc=
if [ -n "$prior" ]; then
    if rev=$(git rev-parse --verify -q "$prior^2"); then
        prior_desc="$prior^2 $(git show "$prior" --pretty='%h (%p) %s')"
    else
        rev=$(git rev-parse "$prior")
        prior_desc="$prior $(git show "$prior" --pretty='%h %s')"
    fi
fi

git branch -f "$branch" "$rev"
git branch -u "$base" "$branch"
if [ -n "$prior_desc" ]; then
    echo "Starting from prior $prior_desc:"
    git log --pretty='- %h %ad %s' "$base".."$branch"
fi
git checkout "$branch"
git commit --allow-empty -m mark
