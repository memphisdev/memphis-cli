# Hello world (Node.js)
## Description
Hello world function in Node.js that add a new field named as the `field_to_ingest` input.
## Supported event formats
JSON 
## Inputs:
Input name | Description | Type
|---|---|---|
| field_to_ingest | The name of the field that will be added to the event | string |
## Test event 

### Inputs
Input name | Value
|---|---|
| field_to_ingest | hello_world_enrichment

### Event:

```json
{
    "id": 1
}
```

## Output to the test event

### Modified Event:
```json
{
    "id": 1,
    "hello_world_enrichment": "Hello from Memphis!"
}
```