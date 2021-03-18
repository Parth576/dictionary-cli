# dictionary-cli
CLI wrapper for [dictionaryapi.dev](https://dictionaryapi.dev/). Currently can search for only english definitions.

## Installation 

#### For Linux / macOS

- Download the binary file from github releases (darwin-amd64 is the macOS release)
- Latest version can be found on the [releases](https://github.com/Parth576/dictionary-cli/releases/latest) page
```bash
wget https://github.com/Parth576/dictionary-cli/releases/download/<version>/gowords-linux-amd64
```
- Make the binary executable
```bash
chmod +x /path/to/binary/gowords-linux-amd64
./path/to/binary/gowords-linux-amd64 search
```
- Using alias in ~/.bashrc
```bash
alias gowords = /path/to/binary/gowords-linux-amd64
```

#### For Windows

- Download the windows release (.exe extension)
- In Powershell / cmd, navigate to the folder where the binary is downloaded
```bash
.\gowords-windows-amd64.exe help
```


## Features

- Search for definitions by word
```bash
gowords search <word>
gowords search obsequious
```
- Print random 'n' words and their definitions from cache

```bash
gowords random --number 5  
gowords random -n 3
```

- Definitions are cached at ~/gowords.json
