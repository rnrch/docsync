# docsync

a tool to write folder structure to markdown files

## Usage

e.g.

```bash
./docsync -ea .git -ea docsync -ea .gitignore -o example-out.md -f example.tmpl
```

[example output](example-out.md)

| flag | meaning                                                                                               |
| ---- | ----------------------------------------------------------------------------------------------------- |
| f    | name of template file(s), make sure the first one specified by this flag is the final output template |
| e    | absolute path of folders/files to be excluded                                                         |
| ea   | name of  folders/files to be excluded                                                                 |
| o    | name of the output file                                                                               |
| d    | absolute path of the directory to process, current directory by default                               |
