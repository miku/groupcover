solrdups
========

SOLR based deduplication, very hard coded.

```
$ solrdups -fq "source_id:10 OR source_id:11" -by source_id -by doi -by institution
```

This would look at all records that match the filter query. For each record,
we would extract the doi and the institutions.

    <SID> <ID>               <ISIL>
    14    10.1201/21212.abc. [DE-15, DE-14, DE-Brt1]

Or
--

Extract a list of fields from SOLR, like estab. Already supported:

* http://localhost:8983/solr/select?q=ipod&fl=id,cat,name,popularity,price,score&wt=csv
