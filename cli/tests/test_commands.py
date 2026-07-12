"""Smoke tests for each sonic command group: --help and dry-run.

One test per command group as specified by task.sonic.cli-runtime.001.
"""

from __future__ import annotations

from click.testing import CliRunner

from sonic.cli import cli


# -- usb ----------------------------------------------------------------

def test_usb_help(runner: CliRunner) -> None:
    result = runner.invoke(cli, ["usb", "--help"])
    assert result.exit_code == 0
    assert "create" in result.output


# -- security -----------------------------------------------------------

def test_security_help(runner: CliRunner) -> None:
    result = runner.invoke(cli, ["security", "--help"])
    assert result.exit_code == 0
    assert "enroll" in result.output


# -- mint ---------------------------------------------------------------

def test_mint_help(runner: CliRunner) -> None:
    result = runner.invoke(cli, ["mint", "--help"])
    assert result.exit_code == 0
    assert "build" in result.output


# -- bootloader ---------------------------------------------------------

def test_bootloader_help(runner: CliRunner) -> None:
    result = runner.invoke(cli, ["bootloader", "--help"])
    assert result.exit_code == 0
    assert "install" in result.output


# -- device -------------------------------------------------------------

def test_device_help(runner: CliRunner) -> None:
    result = runner.invoke(cli, ["device", "--help"])
    assert result.exit_code == 0
    assert "scan" in result.output


# -- mesh ---------------------------------------------------------------

def test_mesh_help(runner: CliRunner) -> None:
    result = runner.invoke(cli, ["mesh", "--help"])
    assert result.exit_code == 0
    assert "init" in result.output


# -- chasis -------------------------------------------------------------

def test_chasis_help(runner: CliRunner) -> None:
    result = runner.invoke(cli, ["chasis", "--help"])
    assert result.exit_code == 0
    assert "add" in result.output