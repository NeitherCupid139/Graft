from __future__ import annotations

import subprocess
import unittest
from pathlib import Path
from tempfile import TemporaryDirectory
from unittest.mock import patch

from check_migration_versions import candidate_dirs, validate


class CandidateDirsTest(unittest.TestCase):
    def test_changed_mode_includes_all_default_chain_dirs_for_global_conflict_checks(self) -> None:
        with TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            httpx_dir = root / "server" / "internal" / "httpx" / "migrations"
            user_dir = root / "server" / "modules" / "user" / "migrations"
            rbac_dir = root / "server" / "modules" / "rbac" / "migrations"
            httpx_dir.mkdir(parents=True)
            user_dir.mkdir(parents=True)
            rbac_dir.mkdir(parents=True)

            with patch(
                "check_migration_versions.subprocess.check_output",
                return_value="server/modules/user/migrations/202605280001_user.sql\nserver/modules/user/README.md\n",
            ):
                dirs = candidate_dirs(root, "changed")

            self.assertEqual(dirs, [httpx_dir, rbac_dir, user_dir])

    def test_all_mode_excludes_historical_shared_migration_dir(self) -> None:
        with TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            httpx_dir = root / "server" / "internal" / "httpx" / "migrations"
            historical_dir = root / "server" / "internal" / "ent" / "migrate" / "migrations"
            user_dir = root / "server" / "modules" / "user" / "migrations"
            httpx_dir.mkdir(parents=True)
            historical_dir.mkdir(parents=True)
            user_dir.mkdir(parents=True)

            dirs = candidate_dirs(root, "all")

            self.assertEqual(dirs, [httpx_dir, user_dir])


class ValidateTest(unittest.TestCase):
    def test_validate_reports_duplicate_versions(self) -> None:
        with TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            httpx_dir = root / "server" / "internal" / "httpx" / "migrations"
            user_dir = root / "server" / "modules" / "user" / "migrations"
            httpx_dir.mkdir(parents=True)
            user_dir.mkdir(parents=True)
            (httpx_dir / "202605280001_access_log.sql").write_text("SELECT 1;\n", encoding="utf-8")
            (user_dir / "202605280001_user.sql").write_text("SELECT 1;\n", encoding="utf-8")

            errors = validate([httpx_dir, user_dir], root)

            self.assertEqual(
                errors,
                [
                    "default migration chain version conflict: 202605280001 appears in "
                    "server/internal/httpx/migrations/202605280001_access_log.sql, "
                    "server/modules/user/migrations/202605280001_user.sql"
                ],
            )

    def test_validate_allows_unique_versions(self) -> None:
        with TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            user_dir = root / "server" / "modules" / "user" / "migrations"
            rbac_dir = root / "server" / "modules" / "rbac" / "migrations"
            user_dir.mkdir(parents=True)
            rbac_dir.mkdir(parents=True)
            (user_dir / "202605280001_user.sql").write_text("SELECT 1;\n", encoding="utf-8")
            (rbac_dir / "202605280002_rbac.sql").write_text("SELECT 1;\n", encoding="utf-8")

            errors = validate([user_dir, rbac_dir], root)

            self.assertEqual(errors, [])


if __name__ == "__main__":
    unittest.main()
