# project list

&nbsp;&nbsp;&nbsp;&nbsp;`$ project list [flags]`

Filter different types of projects and list them quickly.

This allows you to view your projects in different ways
	

## Options


 `-b`,  `--browse`

&nbsp;&nbsp;&nbsp;&nbsp;--browse &lt;string&gt;

 `-e`,  `--exclude`

&nbsp;&nbsp;&nbsp;&nbsp;--exclude &lt;tags&gt;

 `-f`,  `--filter`

&nbsp;&nbsp;&nbsp;&nbsp;--filter &lt;tags&gt;

 `-v`,  `--view <list|category|grid>`

&nbsp;&nbsp;&nbsp;&nbsp;Specifies which layout you want to view your projects in.


## Examples
```
project list --view grid --exclude none

project list -v c --filter go,rust,c#

project list --browse cli
```



## See also
- [`project`](project.md) - Base command

	- [`add`](project-add.md) - Add a new project

	- [`list`](project-list.md) - List projects ***&lt;-- current*** 

	- [`open`](project-open.md) - Open a project

	- [`path`](project-path.md) - Show the path to the project folder

	- [`quick`](project-quick.md) - Open the Quick opener

	- [`remove`](project-remove.md) - Delete a project

	- [`tag`](project-tag.md) - tag project


