#!/bin/bash
echo "========================================="
echo "  SPIRALDB FHE BENCHMARK"
echo "========================================="

echo ""
echo "--- PUT Benchmark (100 operations) ---"
time for i in $(seq 1 100); do
  curl -s -X POST http://localhost:8094/ -d "{\"action\":\"put\",\"key\":\"bench_$i\",\"value\":\"$i\"}" > /dev/null
done
echo "✅ 100 PUTs complete"

echo ""
echo "--- FHE ADD Benchmark (50 blind additions) ---"
# Setup
curl -s -X POST http://localhost:8094/ -d '{"action":"put","key":"ba","value":"10"}' > /dev/null
curl -s -X POST http://localhost:8094/ -d '{"action":"put","key":"bb","value":"20"}' > /dev/null

time for i in $(seq 1 50); do
  curl -s -X POST http://localhost:8094/ -d '{"action":"fhe_compute","key1":"ba","key2":"bb","op":"add"}' > /dev/null
done
echo "✅ 50 FHE additions complete"

echo ""
echo "--- FHE MUL Benchmark (50 blind multiplications) ---"
time for i in $(seq 1 50); do
  curl -s -X POST http://localhost:8094/ -d '{"action":"fhe_compute","key1":"ba","key2":"bb","op":"mul"}' > /dev/null
done
echo "✅ 50 FHE multiplications complete"

echo ""
echo "--- GET Benchmark (100 reads) ---"
time for i in $(seq 1 100); do
  curl -s -X POST http://localhost:8094/ -d "{\"action\":\"get\",\"key\":\"bench_$i\"}" > /dev/null
done
echo "✅ 100 GETs complete"

echo ""
echo "========================================="
echo "  BENCHMARK COMPLETE"
echo "========================================="
