GROUPCOVER 1 "JANUAR 2017" "Leipzig University Library" "Manuals"
=================================================================

NAME
----

groupcover - identifier-based deduplication with groups.

SYNOPSIS
--------

`groupcover` [`-prefs` *preferences*] < *file*

DESCRIPTION
-----------

Deduplication on identifiers is simple, however it can be a more tedious, when
there are different key (e.g. ISIL) attached to an identifier, and the keys
need to be considered for a set of duplicates, with potentiall different
preferences per key.

The groupcover tool tries to by lazy in different ways:

* Work on data, that only contains the necessary fields: identifier, group
  (e.g. source-id), value (e.g. DOI) and keys (e.g. ISIL).

* Input data needs to be sorted by value (e.g. DOI) - similar to uniq(1).

* Output only the rows, that contain changes.

With these simplifications, deduplication becomes a fast operation, taking
about 10min for a file with 150M rows and around 6M changes.

The more time-consuming part will be the preparation of the input CSV
(https://git.io/vprui) and the application of the changes
(https://git.io/vpru1).

OPTIONS
-------

`-cpuprofile` *filename*
  Write cpu profile to given filename.

`-f` *N*
  Column number to use for grouping. One-based, defaults to column 3. 

`-lower`
  Lowercase the input for case-insensitive values.

`-prefs` *string*
  A space separated list of strings, each string names a group (e.g. a source id).

`-version`
  Program version.

`-verbose`
  Show progress.

EXAMPLES
--------

As of early 2018, for AI (finc.info) the deduplication is realized with the following command:

  `groupcover -lower -prefs '85 55 89 60 50 105 34 101 53 49 28 48 121' < a > b`

Where *a* and *b* are CSV files: `id,source_id,doi,isil,isil,isil,...` - with doi and isil being optional.

BUGS
----

Please report bugs to https://github.com/miku/groupcover/issues.

AUTHORS
------

Martin Czygan <https://github.com/miku>, <martin.czygan@uni-leipzig.de>

SEE ALSO
--------

[FINC](https://finc.info), [AMSL](http://amsl.technology/)

