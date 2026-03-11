# Graph Projection Export

Internal tool for exporting the `graph-projection-v1` payload from canonical Aion entities.

## Usage

```bash
go run ./hack/tools/graph-projection-export --user-id 999
go run ./hack/tools/graph-projection-export --user-id 999 --window WINDOW_90D --output ./tmp/graph.json
go run ./hack/tools/graph-projection-export --user-id 999 --category-id 3 --tag-ids 14,15
```

## Notes

- This is a dev/lab tool, not a public runtime endpoint.
- Graph rules remain in `internal/record/core/...`.
- The tool only bootstraps config, DB and repositories, then delegates to canonical services/mappers.
