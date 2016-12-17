groupcover
==========

Staged deduplication.

Input:

```
id, group, attribute, [key, key, ...]
```

Items from different groups may share an attribute. Depending on a order
relation over groups (possibly per key), a number of keys may be dropped for
an entry.

Simple case
-----------

We have a single valued attribute, e.g. a DOI or an ISBN. We sort the input
the attribute and then process the list.

Sort by attribute:

```shell
$ sort -k3 < input.list > sorted.list
```

Process list:

```shell
$ groupcover < sorted.list > cleaned.list
```

Example
-------

```
$ cat fixtures/mini.csv
1,G1,A1,"K1,K2"
2,G1,A2,"K1,K2"
3,G2,A2,"K1,K2,K3"
4,G3,A2,"K2"
5,G1,A3,"K1,K2,K3"
6,G2,A3,"K2,K3"

$ groupcover < fixtures/mini.csv
1,G1,A1,"K1,K2"
2,G1,A2,
3,G2,A2,"K1,K3"
4,G3,A2,K2
5,G1,A3,K1
6,G2,A3,"K2,K3"
```

Finc Index
----------

There is no DOI field in SOLR schema. The licensing information is available
only in *AILicensing*.

```
$ jq -r '[
    .["finc.record_id"],
    .["finc.source_id"],
    .["doi"],
    .["x.labels"][]?] | @csv' < <(unpigz -c /tmp/AILicensing/date-2016-11-28.ldj.gz) > input.csv
```

![](sketch.jpg)
