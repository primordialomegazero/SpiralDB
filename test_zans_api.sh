#!/bin/bash
echo "========================================="
echo "  ZANS TEST — SpiralDB API"
echo "========================================="

# PUT test value
echo ""
echo "--- PUT test value (42) ---"
curl -s -X POST http://localhost:8094/ -d '{"action":"put","key":"zans_test","value":"42"}'

echo ""
echo "--- GET (verify) ---"
curl -s -X POST http://localhost:8094/ -d '{"action":"get","key":"zans_test"}'

echo ""
echo "--- FHE ADD (42 + 0) blind ---"
curl -s -X POST http://localhost:8094/ -d '{"action":"put","key":"zero","value":"0"}' > /dev/null
curl -s -X POST http://localhost:8094/ -d '{"action":"fhe_compute","key1":"zans_test","key2":"zero","op":"add"}'

echo ""
echo "========================================="
echo "  ZANS API test complete!"
echo "========================================="
