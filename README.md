# dictionary-cli
CLI wrapper for dictionaryapi.dev. Currently can search for only english definitions.

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
