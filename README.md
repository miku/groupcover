groupcover
==========

Layered deduplication.

Input:

```
id, group, attribute, keys
```

Items from different groups may share an attribute. Depending on a order
relation over groups, given for each key, a number of ids may be dropped for a
key.

Example:

```
$ cat fixture/sample.tsv
1, G1, A1, K1, K2
2, G1, A2, K1, K2
3, G2, A1, K1, K2
4, G2, A2, K1, K2
```

Various deduplication strategies:

* single attribute, exact match
* single attribute, fuzzy match
* multiple attributes, various stages, fuzzy matches
