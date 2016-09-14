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
0 ignored paths
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
    "language": "SVG",
    "percent": 38.07195038817389
  },
  {
    "language": "HTML",
    "percent": 20.51148776433667
  },
  {
    "language": "JavaScript",
    "percent": 18.27789994256565
  }
]
```

### -json-with-colors

> Output Results in JSON format, including any HTML color codes defined for associated languages.

```
tso@chopstick ~/sirupeuse (master) $ l -json-with-colors -limit 3
[
  {
    "language": "SVG",
    "percent": 38.07195038817389,
    "color": ""
  },
  {
    "language": "HTML",
    "percent": 20.51148776433667,
    "color": "#e44b23"
  },
  {
    "language": "JavaScript",
    "percent": 18.27789994256565,
    "color": "#f1e05a"
  }
]
```

Please note that `Color` will be the empty string `""` if no color is associated with the language.

### -limit n

> Limit number of languages to `n` results, where `n` is a number `> 0`.

> Default is 10.

> An `n` of 0 or less indicates unlimited result set, which may result in lots of erroneous "noise".

### -unignore-contents

### -unignore-filenames

#### (NOT RECOMMENDED)

> By default, this program will ignore certain types of files, such as documentation,

> configuration files, binary data (images, audio, video, executables, etc...)

> which can skew results in undesirable ways, since these files tend to be much

> larger than source code files. 

> This ignoring behavior is based on filename or file contents, and one or both

> can be disabled with these flags, thus including them in the result set.

> This can be useful if too many files are being ignored.

```
tso@chopstick /tmp/react-boilerplate $ l
              JavaScript: 75.1445%
                    JSON: 11.1835%
                Markdown: 04.0704%
              Handlebars: 04.0655%
                     CSS: 02.4746%
                    HTML: 00.8763%
                    YAML: 00.7576%
                     PHP: 00.7551%
                   Nginx: 00.6230%
                   OCaml: 00.0495%

10 languages detected in 183 files
tso@chopstick /tmp/react-boilerplate $ l -unignore-filenames
              JavaScript: 50.8235%
                    JSON: 35.8719%
                Markdown: 05.0065%
                       C: 03.5231%
                    HTML: 02.6561%
                    ABAP: 00.4643%
                     XML: 00.2748%
                   OCaml: 00.2368%
                       M: 00.1317%
                    Diff: 00.1269%

90 languages detected in 28660 files
tso@chopstick /tmp/react-boilerplate $ l -unignore-filenames -unignore-contents
              JavaScript: 46.6284%
                    JSON: 32.7377%
                       C: 09.8128%
                Markdown: 04.5691%
                    HTML: 02.4241%
                    Hack: 01.2408%
                     XML: 00.9649%
                    ABAP: 00.4237%
                   OCaml: 00.2161%
                       M: 00.1202%

93 languages detected in 28723 files
```
