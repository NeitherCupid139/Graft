import json
import tempfile
import unittest
from pathlib import Path

from inventory_open_alerts import (
    load_json_array,
    normalize_code_scanning,
    normalize_dependabot,
    parse_input_json,
    summarize,
)


class InventoryOpenAlertsTests(unittest.TestCase):
    def test_parse_input_json_rejects_invalid_key(self) -> None:
        with self.assertRaises(SystemExit):
            parse_input_json(["other=file.json"])

    def test_parse_input_json_accepts_existing_file(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            path = Path(tmpdir) / "alerts.json"
            path.write_text("[]", encoding="utf-8")
            parsed = parse_input_json([f"code-scanning={path}"])
        self.assertEqual(parsed["code-scanning"], path)

    def test_load_json_array_requires_array(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            path = Path(tmpdir) / "alerts.json"
            path.write_text("{}", encoding="utf-8")
            with self.assertRaises(SystemExit):
                load_json_array(path)

    def test_normalize_code_scanning(self) -> None:
        raw = {
            "number": 17,
            "rule_severity": "high",
            "state": "open",
            "tool": {"name": "CodeQL"},
            "rule": {"id": "go/log-injection"},
            "most_recent_instance": {
                "location": {"path": "server/main.go", "start_line": 44},
                "message": {"text": "User data flows into log output"},
            },
            "html_url": "https://example.test/17",
            "created_at": "2026-06-01T00:00:00Z",
            "dismissed_at": None,
        }
        normalized = normalize_code_scanning(raw)
        self.assertEqual(normalized["rule_id"], "go/log-injection")
        self.assertEqual(normalized["path"], "server/main.go")
        self.assertEqual(normalized["line"], 44)
        self.assertEqual(normalized["tool"], "CodeQL")

    def test_normalize_dependabot(self) -> None:
        raw = {
            "number": 19,
            "state": "open",
            "dependency": {
                "manifest_path": "web/package.json",
                "scope": "runtime",
                "package": {"name": "vite", "ecosystem": "npm"},
            },
            "security_advisory": {
                "severity": "high",
                "summary": "summary",
                "identifiers": [{"value": "GHSA-xxxx"}, {"value": "CVE-2026-0001"}],
            },
            "security_vulnerability": {
                "vulnerable_version_range": "< 7.3.5",
                "first_patched_version": {"identifier": "7.3.5"},
            },
            "html_url": "https://example.test/19",
            "created_at": "2026-06-01T00:00:00Z",
            "dismissed_at": None,
        }
        normalized = normalize_dependabot(raw)
        self.assertEqual(normalized["package"], "vite")
        self.assertEqual(normalized["ecosystem"], "npm")
        self.assertEqual(normalized["first_patched_version"], "7.3.5")
        self.assertEqual(normalized["advisory_ids"], ["GHSA-xxxx", "CVE-2026-0001"])

    def test_summarize_groups_by_kind_and_severity(self) -> None:
        alerts = [
            {"kind": "code-scanning", "severity": "high"},
            {"kind": "dependabot", "severity": "moderate"},
            {"kind": "dependabot", "severity": "moderate"},
        ]
        summary = summarize(alerts)
        self.assertEqual(summary["total"], 3)
        self.assertEqual(summary["by_kind"]["dependabot"], 2)
        self.assertEqual(summary["by_severity"]["moderate"], 2)

    def test_load_json_array_filters_non_object_items(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            path = Path(tmpdir) / "alerts.json"
            path.write_text(json.dumps([{"ok": True}, 1, "x"]), encoding="utf-8")
            loaded = load_json_array(path)
        self.assertEqual(loaded, [{"ok": True}])


if __name__ == "__main__":
    unittest.main()
