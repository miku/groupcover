# GROUPCOVER 1 "JANUAR 2017" "Leipzig University Library" "Manuals"

## NAME

groupcover - identifier-based deduplication with groups.

## SYNOPSIS

`groupcover` [`-prefs` *preferences*] < *file*

## DESCRIPTION

Deduplication on identifiers is simple, however it can be a more tedious, when
there are different keys (e.g. ISIL) attached to an identifier, and the keys
need to be considered for a set of duplicates, with potentially different
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

## OPTIONS

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

## EXAMPLES

As of early 2018, for AI (finc.info) the deduplication is realized with the following command:

  `groupcover -lower -prefs '85 55 89 60 50 105 34 101 53 49 28 48 121' < a > b`

Where *a* and *b* are CSV files: `id,source_id,doi,isil,isil,isil,...` - with doi and isil being optional.

## EXAMPLE INPUT AND OUTPUT

The DOI 10.7557/13.2301 appears twice (in source 49 and 28), with various attachments.

```
ai-49-aHR0cDovL2R4LmRvaS5vcmcvMTAuNzU1Ny8xMy4yMzAx,49,10.7557/13.2301,DE-105,DE-15,DE-Ch1,DE-Brt1,DE-14,DE-82,DE-Gla1,DE-D275,DE-Zwi2,DE-Zi4
ai-28-8c9842c492474445a79c2bfb51fd41d1,28,10.7557/13.2301,DE-15,DE-L242,DE-14,DE-82,DE-540,DE-D161,DE-D275,DE-Zwi2,DE-D117,DE-Bn3
```

The resulting adjustment, after groupcover (preferring 49 over 28):

```
ai-28-8c9842c492474445a79c2bfb51fd41d1,28,10.7557/13.2301,DE-540,DE-Bn3,DE-D117,DE-D161,DE-L242
```

## BUGS

Please report bugs to https://github.com/miku/groupcover/issues.

## AUTHORS

Martin Czygan <https://github.com/miku>, <martin.czygan@uni-leipzig.de>

## SEE ALSO

[FINC](https://finc.info), [AMSL](http://amsl.technology/)

