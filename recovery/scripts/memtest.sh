#!/bin/bash
# =============================================================================
# SonicScrewdriver — Memory Test Tool
# =============================================================================
# Runs memory diagnostics to detect faulty RAM.
# Part of the Sonic recovery toolkit.
# =============================================================================

set -euo pipefail

echo ""
echo "=============================================="
echo "  SonicScrewdriver — Memory Test Tool"
echo "=============================================="
echo ""

echo "Memory diagnostics can be run via:"
echo ""
echo "  1) memtest86+ — Boot from USB (included in SonicScrewloader menu)"
echo "  2) memtester — Run from within Linux"
echo ""

if command -v memtester &>/dev/null; then
    echo "memtester is available."
    echo ""
    echo "Available RAM: $(free -m | awk '/Mem:/ {print $7}') MB free"
    echo ""
    read -rp "Amount of RAM to test in MB (e.g., 1024): " TEST_SIZE
    read -rp "Number of test passes (e.g., 3): " TEST_PASSES

    echo ""
    echo "Running memtester on ${TEST_SIZE}MB for ${TEST_PASSES} passes..."
    echo "Press Ctrl+C to stop."
    echo ""

    memtester "${TEST_SIZE}M" "${TEST_PASSES}"
else
    echo "memtester not installed."
    echo "Install with: sudo apt-get install memtester"
    echo ""
    echo "For full diagnostics, boot memtest86+ from the SonicScrewloader menu."
fi

echo ""
echo "Memory test complete."
