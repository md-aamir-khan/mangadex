# commands
## initdir
    --name
    --id
    --parentdir

- creates .cache yaml file

```yaml
    downloadedWebsiteIndices:
        - 1
    offset: 25
    nextIdxToCreatePDF: 35
    id: "slfjalfjsalfjaslf-asfaljjfl-a;sfja"
    name: "Berserk"
    dir: "manga/Berserk"
    url: "alsfjlasjf.com/alfjlasffljj"
```

## downloadmangachapters
    --parentdir
    --name
    --limit

error out if .cache does not exist, suggest the initdir command

### Later
- if id is not given and files in .cache does not exist or are not valid, then we return error 
