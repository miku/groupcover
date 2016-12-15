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

The attribute value *A1* is shared by records with ID *1* and *3*, which belong to
groups *G1* and *G2*. For each *K1* and *K2*, one or the other is preferred.
E.g. *K1* prefers *G1* over *G2* and *K2* *G2* over *G1*.

```
$ cat fixture/sample.tsv
1   G1  A1  K1,K2
2   G1  A2  K1,K2
3   G2  A1  K1,K2
```

Given the above input, we would learn, that for record *1* the key *K2* can be
dropped and for record *2* the *K1* key can be dropped. The corresponding output would be:

```
$ cat fixture/sample.tsv | groupcover
1   G1  A1  K1
2   G1  A2  K1,K2
3   G2  A1  K2
```

Various deduplication strategies:

* single attribute, exact match
* single attribute, fuzzy match
* multiple attributes, various stages, fuzzy matches

----

```
$ go run cmd/groupcover/main.go < fixtures/sample.tsv
K1 prefers [G1, G2]
K2 prefers [G2, G1]
K3 prefers []

1   G1  A1  [K1, K2]
2   G1  A2  [K1, K2]
3   G2  A1  [K1, K2, K3]

1   G1  A1  [K1]
2   G1  A2  [K1]
3   G2  A1  [K2, K3]
```

----

Notes.

Say, we already reduced the cadidates to a smaller set:

```
1   G1  A1  K1,K2
3   G2  A1  K1,K2,K3
```

Possible interface:

    type Cleaner interface {
        Clean(entries []Entry) []Entry
    }

How to group?

* K1 -> G1, G2; K2 -> G1, G2; K3 -> G2
