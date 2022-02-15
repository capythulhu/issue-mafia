# Welcome to issue-mafia!

**issue-mafia** is an out-of-the-box CLI that helps you to easily synchronize Git hooks with a remote repository. You can learn more about Git hooks [here](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks).

## Getting started

### Installing issue-mafia via Go CLI

To install issue-mafia Go binary, run:

```
go install github.com/thzoid/issue-mafia@latest
```

### Using the `help` command

To see a list of available commands, run:

```
issue-mafia --help
```

## Synchronizing Git hooks

### Setting up a configuration repository

To get started with Git hooks synchronization, a GitHub repository containing every hook to be shared with other repositories needs to be created.

> issue-mafia will only be able to fetch hooks if the repository is public.

On a configuration repository, hooks need to be placed at the root directory. Each script file should be named after the hook they correspond, without extension. For example, if the Git hooks `commit-msg` and `post-commit` were to be synchronized across repositories, the directory structure would look like this:

```
my-own-hooks/
├─ commit-msg
├─ post-commit
```

Then, the GitHub repository could be initialized through the following workflow:

```
cd my-own-hooks
git init
git add .
git commit -m "add msg and post-commit hooks"
git remote add origin git@github.com:<username>/my-own-hooks.git
git push -u origin main
```

> If issue-mafia is being used within an organization, it is recommended to create this repository using the organization's own GitHub account.

### Setting up synchronization on repository

To tell issue-mafia that a specific repository needs to have its hooks synchronized with a remote repository, use the following command:

```
issue-mafia init
```

After completing the configuration file creation wizard and providing a valid issue-mafia configuration repository (explained in the previous step), an `.issue-mafia` file is going to be generated, and this repository is ready to have its hooks fetched from the specified remote origin.

### Synchronizing hooks

To synchronize (or update) hooks on a single repository, the directory should have a folder structure similar to the following:

```
my-local-repo/
├─ .issue-mafia
├─ ...
```

In that case, just run the root command, and the hooks will be synchronized with the remote repository:

```
cd my-local-repo
issue-mafia
```

However, when many repositories have to be synchronized at once, the repositories folder structure might be similar to this one:

```
my-repos/
├─ foo-repo/
│  ├─ .issue-mafia
│  ├─ ...
├─ bar-repo/
│  ├─ .issue-mafia
│  ├─ ...
├─ ...
```

In that case, the base command can be executed using the `--recursive` flag, and issue-mafia will automatically search for repositories in sub-folders:

```
cd my-repos
issue-mafia --recursive
```

> Always make sure that the repositories that provide Git hooks (shell scripts) are trustworthy.

## Collaborating

Any person is invited to collaborate with issue-mafia. Just open a Pull Request and we will review it as soon as possible.
