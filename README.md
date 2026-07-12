# Project CLI

Easily manage your projects with Project CLI. Use commands such as `add`, `list`, `remove`, and `open`.

> [!NOTE]
> This tool is developed on and for MacOs machines. Some features might not work properly on other platforms because of platform specific terminal commands.

## Content:

- [**How do I install it?**](#how-do-i-install-it)
- [**How do I get started?**](#how-do-i-get-started)
- [**How do I use it?**](#how-do-i-use-it)

## How do I install it?

There are 2 ways to install the cli.

**1. Using the `go install` command:**

For this step you need to have [The Go Programming Language](https://go.dev/) installed so that you can use the `install` command.

Then just run the command in your terminal:

```
go install github.com/Revi-Studios/project-cli
```

This will get the code from this repository and install it to your path.

**2. Download a prebuilt binary and adding it to your path manually:**

For this step you only need to download the prebuilt binary and find where you path lives on you computer. (Where new commands are saved and found)

1. Download the latest stable prebuilt binary from [releases](https://github.com/Revi-Studios/project-cli/releases).
2. Add it to you path
3. Try it out!

## How do I get started?

After installing the cli to your terminal path, it's time to set up you project path/folder.

```
mkdir Projects
cd Projects
project path set $(pwd)
```

These commands:

1. Create a folder called "Projects"
2. Cd into it
3. Configures the cli to use the current path as the project path using `path set <path>`.

> [!NOTE]
> `pwd` returns the current path the terminal is at.
> You don't need to create a new folder, name it _Projects_ nore use it as the project folder.
>
> You can use which ever path you want, just replace `$(pwd)` with your own path.

Now you´re all set. You can use the cli to create new projects with tags, list all of them, and open them.

## How do I use it?

> [!TIP]
> Check out the [docs](docs/md/project.md)
