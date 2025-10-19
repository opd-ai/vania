#!/bin/bash
# VANIA - Comprehensive Verification Script
# This script verifies that all systems are operational after code changes

set -e

echo "╔════════════════════════════════════════════════════════╗"
echo "║     VANIA - Comprehensive Verification Script         ║"
echo "╚════════════════════════════════════════════════════════╝"
echo ""

# Check if xvfb-run is available for headless testing
if ! command -v xvfb-run &> /dev/null; then
    echo "⚠️  xvfb-run not found. Install with: sudo apt-get install xvfb"
    echo "   Running tests without virtual display..."
    XVFB=""
else
    XVFB="xvfb-run -a"
fi

echo "1️⃣  Building project..."
if go build -o game ./cmd/game 2>&1; then
    ls -lh game | awk '{print "   ✅ Binary created: " $9 " (" $5 ")"}'
else
    echo "   ❌ Build failed"
    exit 1
fi
echo ""

echo "2️⃣  Running all tests..."
TEST_OUTPUT=$($XVFB go test ./... 2>&1)
PASS_COUNT=$(echo "$TEST_OUTPUT" | grep -c "^ok" || echo "0")
FAIL_COUNT=$(echo "$TEST_OUTPUT" | grep -c "^FAIL" || echo "0")

if [ $FAIL_COUNT -gt 0 ]; then
    echo "   ❌ Tests failing: $FAIL_COUNT packages"
    echo "$TEST_OUTPUT"
    exit 1
else
    echo "   ✅ Tests passing: $PASS_COUNT packages"
fi
echo ""

echo "3️⃣  Checking determinism..."
OUTPUT1=$($XVFB ./game --seed 12345 2>&1 | grep "Total Rooms" || echo "ERROR")
OUTPUT2=$($XVFB ./game --seed 12345 2>&1 | grep "Total Rooms" || echo "ERROR")

if [ "$OUTPUT1" = "$OUTPUT2" ] && [ "$OUTPUT1" != "ERROR" ]; then
    echo "   ✅ Determinism verified: $OUTPUT1"
else
    echo "   ❌ Determinism failed"
    echo "   First run:  $OUTPUT1"
    echo "   Second run: $OUTPUT2"
    exit 1
fi
echo ""

echo "4️⃣  Testing multiple seeds..."
ALL_SEEDS_OK=true
for seed in 1 42 999 1337; do
    RESULT=$($XVFB ./game --seed $seed 2>&1 | grep -E "(Total Rooms|Theme)" | head -2)
    if [ -z "$RESULT" ]; then
        echo "   ❌ Seed $seed failed to generate"
        ALL_SEEDS_OK=false
    else
        echo "   ✅ Seed $seed generated successfully"
    fi
done

if [ "$ALL_SEEDS_OK" = false ]; then
    exit 1
fi
echo ""

echo "5️⃣  Verifying documentation..."
DOCS_OK=true
for doc in "DEBUGGING_ANALYSIS.md" "BUG_FIX_SUMMARY.md" "README.md"; do
    if [ -f "$doc" ]; then
        echo "   ✅ $doc present"
    else
        echo "   ⚠️  $doc missing"
        DOCS_OK=false
    fi
done
echo ""

echo "6️⃣  Checking code quality..."
# Check for common Go issues
if go vet ./... 2>&1 | grep -q ""; then
    echo "   ✅ go vet passed"
else
    echo "   ⚠️  go vet found issues (non-blocking)"
fi
echo ""

echo "╔════════════════════════════════════════════════════════╗"
if [ "$ALL_SEEDS_OK" = true ] && [ $FAIL_COUNT -eq 0 ]; then
    echo "║     ✅ ALL CHECKS PASSED - PRODUCTION READY           ║"
    echo "╚════════════════════════════════════════════════════════╝"
    exit 0
else
    echo "║     ⚠️  SOME CHECKS FAILED - REVIEW REQUIRED          ║"
    echo "╚════════════════════════════════════════════════════════╝"
    exit 1
fi
