# Inventory Reference

Use this reference when the main workflow needs a stable alert schema or an offline review artifact.

## Live Fetch

From the repository root:

```bash
python3 .agents/skills/graft-security-remediation/scripts/inventory_open_alerts.py \
  --repo GeWuYou/Graft \
  --kind all \
  --pretty
```

The helper is read-only. It calls `gh api --paginate` for:

- `/repos/<repo>/code-scanning/alerts?state=open&per_page=100`
- `/repos/<repo>/dependabot/alerts?state=open&per_page=100`

## Offline Import

If live GitHub access is unavailable, save raw JSON arrays first and then normalize them:

```bash
python3 .agents/skills/graft-security-remediation/scripts/inventory_open_alerts.py \
  --kind all \
  --input-json code-scanning=/tmp/code-scanning.json \
  --input-json dependabot=/tmp/dependabot.json \
  --pretty
```

Each input file must be one top-level JSON array returned by `gh api`.

## Normalized Output

The helper writes one JSON object:

```json
{
  "repo": "GeWuYou/Graft",
  "requested_kind": "all",
  "sources": {
    "code-scanning": "gh-api",
    "dependabot": "gh-api"
  },
  "summary": {
    "total": 3,
    "by_kind": {
      "code-scanning": 1,
      "dependabot": 2
    },
    "by_severity": {
      "high": 2,
      "moderate": 1
    }
  },
  "alerts": []
}
```

### Code Scanning Fields

- `kind`
- `number`
- `rule_id`
- `severity`
- `state`
- `tool`
- `path`
- `line`
- `message`
- `html_url`
- `created_at`
- `dismissed_at`

### Dependabot Fields

- `kind`
- `number`
- `package`
- `ecosystem`
- `manifest_path`
- `scope`
- `severity`
- `state`
- `vulnerable_version_range`
- `first_patched_version`
- `advisory_ids`
- `summary`
- `html_url`
- `created_at`
- `dismissed_at`

## Review Use

Use the normalized file to:

- prove the full alert set was inventoried before remediation
- batch classification notes without re-querying GitHub each time
- compare local remediation scope against the remaining open-alert set
