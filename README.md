# gitrecent

CLI tool for selecting recent branches to checkout.

## Install

```
git clone git@github.com:jonatanblue/gitrecent.git
cd gitrecent
go build && cp gitrecent ~/your-preferred-bin-path/gr
```

## Usage

```
$ gr c
? Pick branch to checkout  [Use arrows to move, type to filter]
>   other-branch
    + branch-in-other-worktree
    another-branch
    * main
```


