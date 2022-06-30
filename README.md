# Cypher `MERGE` does not guarantee uniqueness

Spoiler alert: you need to define a [uniqueness constraint](https://neo4j.com/docs/cypher-manual/current/constraints/syntax/#administration-constraints-syntax-create-unique) to make sure the 
combination of a node labels and property is unique.

## How to run

Prereqs:

 - Docker
 - Go 1.18

```shell
go run ./cmd/merge-gotcha
```
