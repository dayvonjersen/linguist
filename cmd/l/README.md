*you still shouldn't use this yet*

# l

A command-line utility to report programming language usage in a project.

A reference implementation of [github.com/generaltso/linguist](https://github.com/generaltso/linguist).

## Usage

#### in:

```bash
$ cd /some/project/dir
$ l
```

#### out:

```
      Go: 98.9999%
Markdown: 01.0001%

2 languages detected in 10 files
```

#### flags:

### -debug

> Print debug information.

### -git

> Scan for files using git ls-tree and cat-file, rather than filesystem.

### -git-tree [treeish]

> Use `treeish` as root to scan, default is `HEAD`. From the manual for git(1):


> &lt;tree-ish&gt;
> Indicates a tree, commit or tag object name. A command that takes a 
> &lt;tree-ish &gt; argument ultimately wants to operate on a &lt;tree&gt; object
> but automatically dereferences &lt;commit&gt; and &lt;tag&gt; objects that point at a &lt;tree&gt;.


> Basically anything like `master`, sha1 hash ids of commits, branch names, and sha1 hash ids of directories.

### -fs

> Scan for files using filesystem

---

**NOTE:**

By default, this tool will use `-git` behavior if a `.git` directory exists, otherwise it will use the `-fs` behavior.

---

### -json

> Output Results in JSON format.

```
tso@chopstick ~/sirupeuse (master) $ l -json -limit 3
[
  {
    "Language": "SVG",
    "Percent": 38.07195038817389
  },
  {
    "Language": "HTML",
    "Percent": 20.51148776433667
  },
  {
    "Language": "JavaScript",
    "Percent": 18.27789994256565
  }
]
```

### -json-with-colors

> Output Results in JSON format, including any HTML color codes defined for associated languages.

```
tso@chopstick ~/sirupeuse (master) $ l -json-with-colors -limit 3
[
  {
    "Language": "SVG",
    "Percent": 38.07195038817389,
    "Color": ""
  },
  {
    "Language": "HTML",
    "Percent": 20.51148776433667,
    "Color": "#e44b23"
  },
  {
    "Language": "JavaScript",
    "Percent": 18.27789994256565,
    "Color": "#f1e05a"
  }
]
```

Please note that `Color` will be the empty string `""` if no color is associated with the language.

### -limit n

> Limit result set to `n` results, where `n` is a number `> 0`.

> Default is 10.

> An `n` of 0 or less indicates unlimited result set, which may result in lots of erroneous "noise".
