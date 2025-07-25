#!/usr/bin/env python3
import csv
import json
import time
from urllib.parse import urlencode
import requests

API_URL     = "http://localhost:8000/api/logs"
PROJECT_IDS = [
    "a63d3d4f-656f-49f6-be2c-7cab588092f7",
    "2000cc4f-64b3-4760-806f-fdc2910f338a",
    "342a666f-2873-4355-963c-88dd72081c4a",
    "a7ad2ae3-9b29-4e05-b3d0-ca224d746051"
]

SCENARIOS = [
    ("no_filter", {}),
    ("by_event", {"event_name": "login"}),
    ("by_key", {"searchable_keys[k1]": "v1_0"}),
    ("combo", {"searchable_keys[k2]": "v2_0"}),
]

PAGES = range(1, 6)
OUTPUT_CSV = "metrics.csv"

def safe_json(resp):
    """Return parsed JSON or empty dict on failure."""
    try:
        return resp.json() or {}
    except (ValueError, json.JSONDecodeError):
        return {}

def run_test():
    with open(OUTPUT_CSV, "w", newline="") as fout:
        writer = csv.writer(fout)
        writer.writerow([
            "project_id",
            "scenario",
            "searchable_keys",
            "page",
            "status_code",
            "duration_ms",
            "total_count",
            "returned_items",
        ])

        for pid in PROJECT_IDS:
            for name, extra_q in SCENARIOS:
                for page in PAGES:
                    # build query params
                    params = {
                        "project_id": pid,
                        "page": page,
                        **extra_q,
                    }
                    url = f"{API_URL}?{urlencode(params)}"

                    start = time.perf_counter()
                    resp  = requests.get(url)
                    delta = (time.perf_counter() - start) * 1000

                    body = safe_json(resp)
                    # ensure data is a list
                    data = body.get("data")
                    if not isinstance(data, list):
                        data = []

                    tc = body.get("total_count")
                    # if total_count missing or not an int, fallback to len(data)
                    total_count = tc if isinstance(tc, int) else len(data)
                    returned    = len(data)

                    writer.writerow([
                        pid,
                        name,
                        json.dumps(extra_q, ensure_ascii=False),
                        page,
                        resp.status_code,
                        round(delta, 2),
                        total_count,
                        returned,
                    ])
                    print(f"{pid} | {name:9} | page={page} | status={resp.status_code} | {delta:.1f}ms")

    print("Done. Results in", OUTPUT_CSV)

if __name__ == "__main__":
    run_test()
