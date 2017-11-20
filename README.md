# git pif

**git pif** (**Push in Fork**)

git pif is a simple utility to increase collaboration across github. 

It simplifies the process of forking a repo. Cloning it. Pushing commits to it. Only step left for user is to create a PR to original Repo.

The `git pif` does all the next steps:

1. forks the current git repo in path
2. add a new remote `fork` with forked repository
3. push the current branch and set the remote of current branch to `fork` 



## Example:

Check out to a new branch in the repo where you want to open a new PR

```
➜  git-pif git:(master) git co -b my-fix
Switched to a new branch 'my-fix'
```

let git pif to fork the repo, in case it hasn't been. And push the branch to the fork.

```
➜  git-pif git:(my-fix) git pif
Using config file: /Users/zainmalik/.git-pif/config.yaml
Checking if user octagen have a fork of zmalik/git-pif
Forking the repo zmalik/git-pif to octagen/git-pif
Waiting the creation of the fork octagen/git-pif
Fork octagen/git-pif created successfully
Branch my-fix set up to track remote branch my-fix from fork
```

### What has happened here?

- A new fork has been created for the repo. And a new remote`fork` has been added.

```
➜  git-pif git:(my-fix) git remote -v
fork	https://github.com/octagen/git-pif (fetch)
fork	https://github.com/octagen/git-pif (push)
origin	https://github.com/zmalik/git-pif (fetch)
origin	https://github.com/zmalik/git-pif (push)
```

- The branch `my-fix` has been pushed to `octagen/gif-pif` and its remote is the `fork`. Which basically means that from now on `git push` in the `my-fix` will go to the `fork`  

```
➜  git-pif git:(my-fix) cat .git/config
[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
	logallrefupdates = true
	ignorecase = true
	precomposeunicode = true
[remote "origin"]
	url = https://github.com/zmalik/git-pif
	fetch = +refs/heads/*:refs/remotes/origin/*
[branch "master"]
	remote = origin
	merge = refs/heads/master
[remote "fork"]
	url = https://github.com/octagen/git-pif
	fetch = +refs/heads/*:refs/remotes/fork/*
[branch "my-fix"]
	remote = fork
	merge = refs/heads/my-fix
```



## Who can use it?

Everybody willing to do a Pull Request Cycle to a github repository. 

But I made it out specially for **Go** developers. As in case of **Go** the path is also really important. So you cannot just fork a repository and clone it in another path. As it will mess up with the packages. You have to play with the remotes. And this tool intends to simplify that process.



## Configuration

`git pif` needs an environment variable. `GITHUB_TOKEN` The github personal API token is used to fork the repository. 

There are two ways to setting it:

- Set the environment variable `GITHUB_TOKEN` or inline`GITHUB_TOKEN=mygithubtooken git pif`
- Add the `GITHUB_TOKEN` to `$HOME/.git-pif/config.yaml` 

```
cat $HOME/.git-pif/config.yaml
GITHUB_TOKEN : mygithubtoken
```



## Basic rules AKA best practices

- Mixing up the branches in the fork and origin is prohibited. That means:
  - using master branch to push in the fork is prohibited. 
  - pushing a branch that exists in origin is prohibited

