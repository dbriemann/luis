# Plain Text Relations

This document describes the relations of models as plain text. It ignores all basic attributes of the models.

## Collection

* is owned by 1 owner(User),
* has 1 cover(File),
* has 0..n files(File),
* has 0..n permissions(Permission),

## File

* is owned by 1 owner(User),
* is used as cover by 0..n collections,
* is in 0..n collections,
* has 0..n tags(Tag),
* has 0..n stars(Star),
* has 0..n comments(Comment),
* has 0..n permissions(Permission),

## User

* owns 0..n files(File)
* owns 0..n collections(Collection)
* has 0..n permissions(Permission)
 
## Permission



## Tag

* is attached to 0..n files(File),

## Star

* belongs to 1 File

## Comment

* belongs to 1 file
* belongs to 1 collection ?