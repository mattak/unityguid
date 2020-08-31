# unityguid

list up guid from meta files

## Install

```
make install
```

## Usage

### List up guid & file names

```
unityguid list <root_asset_path>
```

```
$ unityguid list Assets
0d3e014d4fe4741d3bb198eeaf4037a8	Assets/Scripts/Script1.cs.meta
9fb931026eb674fe390758d25a38b61a	Assets/Scripts/Script2.cs.meta
...
```

### Conflict detection

Compare guids between base file and target file.
If some guid is conflicted, output conflicted tsv info.

```
unityguid conflict <base_guid_filename_tsv> <target_guid_filename_tsv>+
```

```
$ unityguid conflict project1.tsv project2.tsv project3.tsv
0d3e014d4fe4741d3bb198eeaf4037a8	project1	project2	Assets/Scripts/Script1.cs.meta	Assets/Scripts/Script2.cs.meta
0d3e014d4fe4741d3bb198eeaf4037a8	project1	project3	Assets/Scripts/Script1.cs.meta	Assets/Scripts/Script3.cs.meta
...
```

### Replace specific guids

```
unityguid replace <root_assets_path> <guids>+
```

```
$ unityguid replace Assets 0d3e014d4fe4741d3bb198eeaf4037a8
0d3e014d4fe4741d3bb198eeaf4037a8 => 143c024dcfe474193bb298eeaf405ca9	Assets/Scripts/Script1.cs.meta
...
```

## LICENSE

[MIT](./LICENSE.md)

